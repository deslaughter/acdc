package acdc_test

import (
	"acdc"
	"context"
	"path/filepath"
	"testing"
)

func TestEvalConditions(t *testing.T) {

	// model := "5MW_Land_BD_linear"
	model := "StC_test_OC4Semi_Linear_Tow"
	// model := "5MW_OC4Semi_Linear"

	modelPath := filepath.Join("testdata", model, model+".fst")
	runDir := filepath.Join("testdata", model, "turbine")

	conditions := []acdc.Conditions{
		{WindSpeed: 0, BladePitch: 0, RotorSpeed: 0},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	execPath := "/Users/dslaught/projects/openfast-mac/install/bin/openfast"

	err := acdc.EvaluateModel(ctx, modelPath, execPath, runDir, conditions)
	if err != nil {
		t.Fatal(err)
	}
}