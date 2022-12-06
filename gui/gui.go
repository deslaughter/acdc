package gui

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/deslaughter/acdc/anl"
	"github.com/deslaughter/acdc/input"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

//go:embed index.html static
var staticContent embed.FS

func Run(staticFS fs.FS) error {

	if staticFS == nil {
		staticFS = staticContent
	}

	hub := newHub()
	go hub.run()

	r := mux.NewRouter()
	root := r.PathPrefix("/acdc/")

	api := root.PathPrefix("/api/").Subrouter()
	api.HandleFunc("/schemas", schemaHandler).Methods("GET")
	api.HandleFunc("/analysis", getAnalysisHandler).Methods("GET")
	api.HandleFunc("/analysis", putAnalysisHandler).Methods("PUT")
	api.HandleFunc("/conditions", updateConditionsHandler).Methods("POST")
	api.HandleFunc("/model", importModelHandler).Methods("POST")
	api.HandleFunc("/evaluate", hub.evaluateStartHandler).Methods("POST")
	api.HandleFunc("/evaluate", hub.evaluateCancelHandler).Methods("DELETE")
	api.HandleFunc("/evaluate", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	}).Methods("GET")
	api.HandleFunc("/validate-path", validatePathHandler).Methods("POST")

	// root.PathPrefix("/static/").Handler(http.StripPrefix("/fasted/", http.FileServer(http.FS(staticFS))))
	staticHandler := func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/acdc/", http.FileServer(http.FS(staticFS))).ServeHTTP(w, r)
	}
	r.PathPrefix("/").HandlerFunc(staticHandler)

	return http.ListenAndServe(":8080", r)
}

func schemaHandler(w http.ResponseWriter, r *http.Request) {

	err := json.NewEncoder(w).Encode(input.Schemas)
	if err != nil {
		http.Error(w, fmt.Sprintf("error encoding schemas: %s", err), http.StatusInternalServerError)
	}
}

const AnalysisFile = "analysis.json"

func getAnalysisHandler(w http.ResponseWriter, r *http.Request) {

	analysis, err := anl.Read(AnalysisFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("error reading '%s': %s", AnalysisFile, err),
			http.StatusInternalServerError)
		return
	}

	analysis.ValidatePaths()

	err = json.NewEncoder(w).Encode(analysis)
	if err != nil {
		http.Error(w, fmt.Sprintf("error encoding analysis: %s", err), http.StatusInternalServerError)
	}
}

func putAnalysisHandler(w http.ResponseWriter, r *http.Request) {

	analysis, err := anl.Read(AnalysisFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("error reading '%s': %s", AnalysisFile, err),
			http.StatusInternalServerError)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(analysis); err != nil {
		http.Error(w, fmt.Sprintf("error decoding condition: %s", err),
			http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Validate paths
	analysis.ValidatePaths()

	// Save analysis with turbine data
	if err = analysis.Write(AnalysisFile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(analysis)
	if err != nil {
		http.Error(w, fmt.Sprintf("error encoding analysis: %s", err), http.StatusInternalServerError)
	}
}

func updateConditionsHandler(w http.ResponseWriter, r *http.Request) {

	// Read condition from body
	conditions := []anl.Conditions{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&conditions); err != nil {
		http.Error(w, fmt.Sprintf("error decoding condition: %s", err),
			http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Sort by wind speed and rotor speed
	sort.SliceStable(conditions, func(i, j int) bool {
		if conditions[i].WindSpeed != conditions[j].WindSpeed {
			return conditions[i].WindSpeed < conditions[j].WindSpeed
		}
		return conditions[i].RotorSpeed < conditions[j].RotorSpeed
	})

	analysis, err := anl.Read(AnalysisFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("error reading '%s': %s", AnalysisFile, err),
			http.StatusInternalServerError)
		return
	}

	// Add new condition to list
	analysis.Conditions = conditions

	// Save analysis with conditions data
	if err = analysis.Write(AnalysisFile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Encode turbine as response
	if err = json.NewEncoder(w).Encode(analysis.Conditions); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

const MAX_UPLOAD_SIZE = 20 * 1024 * 1024 // 20MB

func importModelHandler(w http.ResponseWriter, r *http.Request) {

	// Get path to main model file
	path := r.FormValue("path")

	// Read analysis file
	analysis, err := anl.Read(AnalysisFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("error reading '%s': %s", AnalysisFile, err),
			http.StatusInternalServerError)
		return
	}

	// Read turbine from files
	analysis.Model, err = input.ReadFiles(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save model path
	analysis.ModelPath = path

	// Save analysis with turbine data
	if err = analysis.Write(AnalysisFile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Encode turbine as response
	if err = json.NewEncoder(w).Encode(analysis.Model); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (hub *Hub) evaluateStartHandler(w http.ResponseWriter, r *http.Request) {

	analysis, err := anl.Read(AnalysisFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("error reading '%s': %s", AnalysisFile, err),
			http.StatusInternalServerError)
		return
	}

	// Create local copy of conditions
	conditions := analysis.Conditions

	// Send command to reset evaluation
	hub.resetChan <- ResetEval{
		StatusCount: len(conditions),
	}

	var ctx context.Context
	ctx, hub.cancelFunc = context.WithCancel(context.Background())

	g, ctx2 := errgroup.WithContext(ctx)

	// Launch evaluations
	semChan := make(chan struct{}, analysis.NumCPUs)
	for _, conditions := range conditions {
		conditions := conditions
		semChan <- struct{}{}
		g.Go(func() error {
			defer func() { <-semChan }()
			return analysis.Evaluate(ctx2, conditions, hub.statusChan)
		})
	}

	// Wait for evaluations to complete. If error, print
	go func() {
		if err := g.Wait(); err != nil {
			fmt.Println(err)
		}
		hub.cancelFunc()
	}()

	w.WriteHeader(http.StatusNoContent)
}

func (hub *Hub) evaluateCancelHandler(w http.ResponseWriter, r *http.Request) {
	hub.cancelFunc()
	w.WriteHeader(http.StatusNoContent)
}

func validatePathHandler(w http.ResponseWriter, r *http.Request) {

	// Get path to validate
	path := r.FormValue("path")

	// If path exists
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

//------------------------------------------------------------------------------
// Hub
//------------------------------------------------------------------------------

type ResetEval struct {
	StatusCount int
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	cancelFunc  context.CancelFunc
	statusMap   []anl.EvalStatus
	statusChan  chan anl.EvalStatus
	resetChan   chan ResetEval
	clients     map[*Client]struct{}
	register    chan *Client
	unregister  chan *Client
	lastMessage []byte
}

func newHub() *Hub {
	return &Hub{
		cancelFunc:  func() {},
		statusMap:   make([]anl.EvalStatus, 0),
		statusChan:  make(chan anl.EvalStatus, 10),
		resetChan:   make(chan ResetEval),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]struct{}),
		lastMessage: []byte("[]"),
	}
}

func (h *Hub) run() {

	sendChan := make(chan struct{})
	sendTimer := time.AfterFunc(0, func() {
		sendChan <- struct{}{}
	})

	for {
		select {
		case client := <-h.register:
			h.clients[client] = struct{}{}
			client.send <- h.lastMessage

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case reset := <-h.resetChan:
			h.statusMap = make([]anl.EvalStatus, reset.StatusCount)
			for i := range h.statusMap {
				h.statusMap[i].ID = i + 1
				h.statusMap[i].State = "Queued"
			}
			h.statusChan <- anl.EvalStatus{ID: 0}

		case status := <-h.statusChan:

			// Update status if ID in range
			if status.ID >= 1 && status.ID <= len(h.statusMap) {
				h.statusMap[status.ID-1] = status
			}

			sendTimer.Reset(time.Millisecond * 100)

		case <-sendChan:

			// Convert map to json data
			data, err := json.Marshal(h.statusMap)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// Store data as last message
			h.lastMessage = data

			// Broadcast data to clients
			for client := range h.clients {
				select {
				case client.send <- data:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

//------------------------------------------------------------------------------
// Client
//------------------------------------------------------------------------------

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	go client.writePump()
	go client.readPump()
}
