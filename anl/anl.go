package anl

import (
	"acdc"
	"acdc/input"
	"encoding/json"
	"os"
)

type Analysis struct {
	Name       string
	ExecPath   string
	Conditions []acdc.Conditions
	Viz        VizData
	Turbine    *input.Turbine
	Campbell   *CampbellData
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

type CampbellData struct {
}

type VizData struct {
}
