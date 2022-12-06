package anl

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/deslaughter/acdc/input"
)

type Turbine struct {
	ID             int
	Name           string
	OperatingPoint Conditions
	Dir            string
	ModelPath      string
	LogPath        string
	Model          *input.Model
}

func NewTurbine(c Conditions, model *input.Model) *Turbine {

	// Create turbine instance
	turb := &Turbine{
		ID:             c.ID,
		Name:           fmt.Sprintf("turb_%02d", c.ID),
		OperatingPoint: c,
		Model:          model,
	}

	// Make turbine directory the same as name
	turb.Dir = turb.Name

	// Build path to turbine model file and log file
	turb.ModelPath = filepath.Join(turb.Dir, turb.Name+".fst")
	turb.LogPath = filepath.Join(turb.Dir, turb.Name+".log")

	return turb
}

func (turb *Turbine) Simulate(ctx context.Context, execPath string, statusChan chan<- EvalStatus) error {

	// Write model input files
	if err := turb.Model.WriteFiles(turb.ModelPath); err != nil {
		log.Fatal(err)
	}

	// Create log file
	logFile, err := os.Create(turb.LogPath)
	if err != nil {
		return fmt.Errorf("error creating log file '%s': %w", logFile.Name(), err)
	}
	defer logFile.Close()

	// Create command
	cmd := exec.CommandContext(ctx, execPath, turb.ModelPath)

	// Get pipe to read stdout and stderr
	outputReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = cmd.Stdout

	// Start command
	cmd.Start()

	// Get progress
	scanner := bufio.NewScanner(outputReader)
	for scanner.Scan() {
		line := scanner.Text()
		logFile.WriteString(line + "\n")
		if strings.Contains(line, "Time: ") {
			fields := strings.Fields(line)
			currentTime, err := strconv.ParseFloat(fields[1], 32)
			if err != nil {
				continue
			}
			totalTime, err := strconv.ParseFloat(fields[3], 32)
			if err != nil {
				continue
			}
			statusChan <- EvalStatus{
				ID:       turb.ID,
				State:    "Simulation",
				Progress: int(100 * currentTime / totalTime),
			}
		} else if strings.Contains(line, "Performing linearization") {
			fields := strings.Fields(line)
			linNumber, err := strconv.ParseFloat(fields[2], 32)
			if err != nil {
				continue
			}
			statusChan <- EvalStatus{
				ID:       turb.ID,
				State:    "Linearization",
				Progress: int(100 * linNumber / float64(turb.Model.FAST.NLinTimes)),
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	// Wait for command to exit
	if err := cmd.Wait(); err != nil {
		statusChan <- EvalStatus{
			ID:       turb.ID,
			State:    "Error",
			Progress: 100,
			Error:    err.Error(),
		}
		return err
	}

	// If context was canceled
	if err := ctx.Err(); err != nil {
		statusChan <- EvalStatus{
			ID:       turb.ID,
			State:    "Error",
			Progress: 100,
			Error:    "Canceled: " + err.Error(),
		}
		return fmt.Errorf("run canceled")
	}

	// Send complete status
	statusChan <- EvalStatus{
		ID:       turb.ID,
		State:    "Complete",
		Progress: 100,
	}

	return nil
}

func (turb *Turbine) PerformMBC() (*MBC, error) {

	// Red linearization files produced by this turbine
	linFiles, err := filepath.Glob(filepath.Join(turb.Dir, turb.Name+"*.lin"))
	if err != nil {
		return nil, err
	}
	linData := make([]*LinData, len(linFiles))
	for i, f := range linFiles {
		if linData[i], err = ReadLinData(f); err != nil {
			return nil, err
		}
	}

	// Combine linearization data into matrix data
	matData, err := collectMatrixData(linData)
	if err != nil {
		return nil, err
	}

	_ = matData

	return nil, nil
}
