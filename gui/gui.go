package gui

import (
	"acdc"
	"acdc/anl"
	"acdc/input"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gorilla/mux"
)

//go:embed index.html static
var staticContent embed.FS

func Run(staticFS fs.FS) error {

	if staticFS == nil {
		staticFS = staticContent
	}

	r := mux.NewRouter()
	root := r.PathPrefix("/acdc/")

	api := root.PathPrefix("/api/").Subrouter()
	api.HandleFunc("/schemas", schemaHandler).Methods("GET")
	api.HandleFunc("/analysis", getAnalysisHandler).Methods("GET")
	api.HandleFunc("/conditions", updateConditionsHandler).Methods("POST")
	api.HandleFunc("/model", uploadHandler).Methods("POST")

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

	err = json.NewEncoder(w).Encode(analysis)
	if err != nil {
		http.Error(w, fmt.Sprintf("error encoding analysis: %s", err), http.StatusInternalServerError)
	}
}

func updateConditionsHandler(w http.ResponseWriter, r *http.Request) {

	// Read condition from body
	conditions := []acdc.Conditions{}
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

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	// Create temporary directory for model
	modelDir, err := os.MkdirTemp(".", "model")
	if err != nil {
		http.Error(w, fmt.Sprintf("error creating model temporary directory %s", err), http.StatusBadRequest)
		return
	}
	defer os.RemoveAll(modelDir)

	// 32 MB is the default used by FormFile()
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get list of file paths
	filePaths := r.MultipartForm.Value["paths"]
	for i := range filePaths {
		tmp := filepath.Join(strings.Split(filepath.Clean(filePaths[i]), string(filepath.Separator))[1:]...)
		filePaths[i] = filepath.Join(modelDir, tmp)
	}

	// Get a reference to the fileHeaders.
	// They are accessible only after ParseMultipartForm is called
	fileHeaders := r.MultipartForm.File["files"]

	// Loop through file headers
	for i, fileHeader := range fileHeaders {

		if fileHeader.Size > MAX_UPLOAD_SIZE {
			http.Error(w, fmt.Sprintf("The uploaded image is too big: %s. Please use an image less than 1MB in size", fileHeader.Filename), http.StatusBadRequest)
			return
		}

		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer file.Close()

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = os.MkdirAll(filepath.Dir(filePaths[i]), os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f, err := os.Create(filePaths[i])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	matches, err := filepath.Glob(filepath.Join(modelDir, "*.fst"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(matches) == 0 {
		http.Error(w, "please upload main file with '.fst' extension", http.StatusBadRequest)
		return
	}

	// Read analysis file
	analysis, err := anl.Read(AnalysisFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("error reading '%s': %s", AnalysisFile, err),
			http.StatusInternalServerError)
		return
	}

	// Read turbine from files
	analysis.Turbine, err = input.ReadFiles(matches[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save analysis with turbine data
	if err = analysis.Write(AnalysisFile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Encode turbine as response
	if err = json.NewEncoder(w).Encode(analysis.Turbine); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
