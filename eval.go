package acdc

import (
	"acdc/input"
	"acdc/turb"
	"context"
	"fmt"
	"log"
	"path/filepath"
)

func EvaluateModel(ctx context.Context, modelPath, execPath, runDir string, conditions []Conditions) error {

	inp, err := input.ReadFiles(modelPath)
	if err != nil {
		return fmt.Errorf("error reading input files: %s", err)
	}

	// Create new fast library
	turbines, err := turb.NewTurbines(len(conditions))
	if err != nil {
		log.Fatal(err)
	}

	// Create turbine options structure
	turbineOpts := &turb.Opts{MinOutput: true}

	// Initialize turbines
	for i, turbine := range turbines {

		// Get conditions for this turbine
		c := conditions[i]

		// Modify input file with conditions
		inp.ElastoDyn.BlPitch1 = c.BladePitch
		inp.ElastoDyn.BlPitch2 = c.BladePitch
		inp.ElastoDyn.BlPitch3 = c.BladePitch
		inp.ElastoDyn.RotSpeed = c.RotorSpeed
		inp.ElastoDyn.TTDspFA = c.TowerTopDispForeAft
		inp.ElastoDyn.TTDspSS = c.TowerTopDispSideSide

		// Disable DOFs
		inp.ElastoDyn.YawDOF = false
		inp.ElastoDyn.TeetDOF = false

		// Write input files
		inputFilePath := filepath.Join(runDir, fmt.Sprintf("turb_%02d.fst", i+1))
		err := inp.WriteFiles(inputFilePath)
		if err != nil {
			log.Fatal(err)
		}

		// Initialize turbine
		if err := turbine.Initialize(inputFilePath, execPath, turbineOpts); err != nil {
			return fmt.Errorf("turbine %d failed to initialize with input '%s': %s",
				turbine.Num, inputFilePath, err)
		}

		// Run simulation
		if err := turbine.Run(ctx); err != nil {
			return fmt.Errorf("turbine %d: %s", turbine.Num, err)
		}
	}

	if err := turbines.Delete(); err != nil {
		return fmt.Errorf("error deleting turbines: %s", err)
	}

	return nil
}
