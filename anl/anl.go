package anl

import (
	"acdc"
	"acdc/input"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

type Analysis struct {
	Name           string
	ModelPath      string
	ModelPathValid bool
	ExecPath       string
	ExecPathValid  bool
	NumCPUs        int
	Conditions     []acdc.Conditions
	Viz            VizData
	Model          *input.Model
	Campbell       *CampbellData
}

func New() *Analysis {
	return &Analysis{
		NumCPUs: 1,
	}
}

func Read(path string) (*Analysis, error) {

	// Read analysis file
	bs, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse json to structure
	a := &Analysis{}
	err = json.Unmarshal(bs, a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Analysis) Write(path string) error {

	// Sort conditions
	sort.SliceStable(a.Conditions, func(i, j int) bool {
		if a.Conditions[i].WindSpeed != a.Conditions[j].WindSpeed {
			return a.Conditions[i].WindSpeed < a.Conditions[j].WindSpeed
		}
		return a.Conditions[i].RotorSpeed < a.Conditions[j].RotorSpeed
	})

	// Set condition identifiers
	for i := range a.Conditions {
		a.Conditions[i].ID = i + 1
	}

	// Convert analysis to json
	bs, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		return err
	}

	// Write json to file
	if err = os.WriteFile(path, bs, 0777); err != nil {
		return err
	}

	return nil
}

func (a *Analysis) ValidatePaths() {

	if _, err := os.Stat(a.ModelPath); !os.IsNotExist(err) {
		a.ModelPathValid = true
	} else {
		a.ModelPathValid = false
	}

	if _, err := os.Stat(a.ExecPath); !os.IsNotExist(err) {
		a.ExecPathValid = true
	} else {
		a.ExecPathValid = false
	}
}

type CampbellData struct {
}

type VizData struct {
}

type EvalStatus struct {
	ID       int
	State    string
	Progress int
	Error    string
}

func (a *Analysis) Evaluate(ctx context.Context, conditions acdc.Conditions,
	statusChan chan<- EvalStatus) error {

	// Make a copy of the model
	model := &input.Model{}
	bs, err := json.Marshal(a.Model)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bs, model); err != nil {
		return err
	}

	// Modify model for conditions
	model.ElastoDyn.BlPitch1 = conditions.BladePitch
	model.ElastoDyn.BlPitch2 = conditions.BladePitch
	model.ElastoDyn.BlPitch3 = conditions.BladePitch
	model.ElastoDyn.RotSpeed = conditions.RotorSpeed
	model.ElastoDyn.TTDspFA = conditions.TowerTopDispForeAft
	model.ElastoDyn.TTDspSS = conditions.TowerTopDispSideSide

	// If wind speed is zero, disable inflow wind
	if conditions.WindSpeed == 0 {
		model.FAST.CompInflow = 0
		model.FAST.CompAero = 0
	} else {
		model.FAST.CompInflow = 1
		model.InflowWind.WindType = 1
		model.InflowWind.HWindSpeed = conditions.WindSpeed
		model.InflowWind.PLExp = 0
	}

	// Create turbine from model and conditions
	turbine := NewTurbine(conditions, model)

	// Create directory for turbine
	if err := os.MkdirAll(filepath.Dir(turbine.ModelPath), 0777); err != nil {
		return err
	}

	// Run turbine simulation
	if err := turbine.Simulate(ctx, a.ExecPath, statusChan); err != nil {
		return err
	}

	return nil
}
