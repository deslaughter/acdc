package main

import (
	"acdc/gui"
	"os"
)

type Conditions struct {
	WindSpeed            float64 // Wind speed (m/s)
	BladePitch           float64 // Blade pitch (deg)
	RotorSpeed           float64 // Rotor speed in (rpm)
	TowerTopDispForeAft  float64 // Tower Top Displacement Fore-Aft (m)
	TowerTopDispSideSide float64 // Tower Top Displacement Side-Side (m)
}

func main() {

	gui.Run(os.DirFS("gui"))

	// inp, err := input.ReadFiles("testdata/Model_GE-1.5SLE/GE-1.5SLE.fst")
	// if err != nil {
	// 	log.Fatalf("error reading input files: %s", err)
	// }

	// conditions := []Conditions{
	// 	{WindSpeed: 0, BladePitch: 0, RotorSpeed: 0},
	// }

	// // Create new fast library
	// turbines, err := NewTurbines(len(conditions))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Create turbine options structure
	// turbineOpts := &TurbineOpts{MinOutput: true}

	// // Initialize turbines
	// for i, turbine := range turbines {

	// 	// Get conditions for this turbine
	// 	c := conditions[i]

	// 	// Modify input file with conditions
	// 	inp.ElastoDyn.BlPitch1 = c.BladePitch
	// 	inp.ElastoDyn.BlPitch2 = c.BladePitch
	// 	inp.ElastoDyn.BlPitch3 = c.BladePitch
	// 	inp.ElastoDyn.RotSpeed = c.RotorSpeed
	// 	inp.ElastoDyn.TTDspFA = c.TowerTopDispForeAft
	// 	inp.ElastoDyn.TTDspSS = c.TowerTopDispSideSide

	// 	// Disable DOFs
	// 	inp.ElastoDyn.YawDOF = false
	// 	inp.ElastoDyn.TeetDOF = false

	// 	// Write input files
	// 	inputFilePath := fmt.Sprintf("testdata/stab_%02d.fst", i+1)
	// 	err := inp.WriteFiles(inputFilePath)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	// Initialize turbine
	// 	if err := turbine.Initialize(inputFilePath, turbineOpts); err != nil {
	// 		log.Fatalf("turbine %d failed to initialize with input '%s': %s",
	// 			turbine.Num, inputFilePath, err)
	// 	}

	// 	// Run simulation
	// 	if err := turbine.Run(); err != nil {
	// 		log.Fatalf("turbine %d failed to run: %s", turbine.Num, err)
	// 	}

	// 	// Stop simulation
	// 	turbine.Stop()
	// }

	// if err := turbines.Delete(); err != nil {
	// 	log.Fatalf("error deleting turbines: %s", err)
	// }
}
