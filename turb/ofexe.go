//go:build !oflib

package turb

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Turbines []*Turbine

func NewTurbines(numTurbines int) (Turbines, error) {

	// Create slice to represent each turbine
	turbines := make(Turbines, numTurbines)
	for i := range turbines {
		turbines[i] = &Turbine{
			Num: i,
		}
	}

	return turbines, nil
}

// Delete deallocates memory associated with the turbines
func (ts *Turbines) Delete() error {

	// Nullify slice of turbines so they can't be reused
	*ts = nil

	return nil
}

type Turbine struct {
	Num       int
	InputPath string
	ExecPath  string
}

type Opts struct {
	MinOutput bool
	ExecPath  string
}

func (t *Turbine) Initialize(inputPath, execPath string, opts *Opts) error {
	t.InputPath = inputPath
	t.ExecPath = execPath
	return nil
}

func (t *Turbine) Run(ctx context.Context) error {

	// Create log file
	logFile, err := os.Create(strings.TrimSuffix(t.InputPath, filepath.Ext(t.InputPath)) + ".log")
	if err != nil {
		return fmt.Errorf("error creating log file '%s': %w", logFile.Name(), err)
	}

	// Create command
	cmd := exec.CommandContext(ctx, t.ExecPath, t.InputPath)
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	// Run command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running '%v': %w", cmd.Args, err)
	}

	// If context was canceled
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("run canceled")
	}

	return nil
}
