package input_test

import (
	"acdc/input"
	"os"
	"testing"
)

func TestBeamDynFormat(t *testing.T) {

	text, err := BeamDynExp.Format()
	if err != nil {
		t.Fatal(err)
	}

	act := input.NewBeamDyn()
	if err := act.Parse(text); err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, BeamDynExp); err != nil {
		t.Fatal(err)
	}

	os.WriteFile("testdata/test_BeamDyn.dat", text, 0777)
}

func TestBeamDynParse(t *testing.T) {

	act, err := input.ReadBeamDyn("testdata/NRELOffshrBsline5MW_BeamDyn.dat")
	if err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, BeamDynExp); err != nil {
		t.Fatal(err)
	}
}

var BeamDynExp = &input.BeamDyn{Title: "NREL 5MW blade",
	Echo:             false,
	QuasiStaticInit:  true,
	Rhoinf:           0,
	Quadrature:       2,
	Refine:           0,
	N_fact:           0,
	DTBeam:           0,
	Load_retries:     0,
	NRMax:            0,
	Stop_tol:         0,
	Tngt_stf_fd:      false,
	Tngt_stf_comp:    false,
	Tngt_stf_pert:    0,
	Tngt_stf_difftol: 0,
	RotStates:        true,
	Member_total:     1,
	Kp_total:         49,
	Members: [][]input.BeamDynMembers{
		{
			{0, 0, 0, 13.308},
			{0, 0, 0.199875, 13.308},
			{0, 0, 1.199865, 13.308},
			{0, 0, 2.199855, 13.308},
			{0, 0, 3.199845, 13.308},
			{0, 0, 4.199835, 13.308},
			{0, 0, 5.199825, 13.308},
			{0, 0, 6.199815, 13.308},
			{0, 0, 7.199805, 13.308},
			{0, 0, 8.201025, 13.308},
			{0, 0, 9.199785, 13.308},
			{0, 0, 10.199775, 13.308},
			{0, 0, 11.199765, 13.181},
			{0, 0, 12.199755, 12.848},
			{0, 0, 13.200975, 12.192},
			{0, 0, 14.199735, 11.561},
			{0, 0, 15.199725, 11.072},
			{0, 0, 16.199715, 10.792},
			{0, 0, 18.200925, 10.232},
			{0, 0, 20.20029, 9.672},
			{0, 0, 22.20027, 9.11},
			{0, 0, 24.20025, 8.534},
			{0, 0, 26.20023, 7.932},
			{0, 0, 28.200825, 7.321},
			{0, 0, 30.20019, 6.711},
			{0, 0, 32.20017, 6.122},
			{0, 0, 34.20015, 5.546},
			{0, 0, 36.20013, 4.971},
			{0, 0, 38.200725, 4.401},
			{0, 0, 40.20009, 3.834},
			{0, 0, 42.20007, 3.332},
			{0, 0, 44.20005, 2.89},
			{0, 0, 46.20003, 2.503},
			{0, 0, 48.20124, 2.116},
			{0, 0, 50.19999, 1.73},
			{0, 0, 52.19997, 1.342},
			{0, 0, 54.19995, 0.954},
			{0, 0, 55.19994, 0.76},
			{0, 0, 56.19993, 0.574},
			{0, 0, 57.19992, 0.404},
			{0, 0, 57.699915, 0.319},
			{0, 0, 58.20114, 0.253},
			{0, 0, 58.699905, 0.216},
			{0, 0, 59.1999, 0.178},
			{0, 0, 59.699895, 0.14},
			{0, 0, 60.19989, 0.101},
			{0, 0, 60.699885, 0.062},
			{0, 0, 61.19988, 0.023},
			{0, 0, 61.5, 0},
		},
	},
	Order_elem:  5,
	BldFile:     "NRELOffshrBsline5MW_BeamDyn_Blade.dat",
	UsePitchAct: false,
	PitchJ:      200,
	PitchK:      2e+07,
	PitchC:      500000,
	SumPrint:    true,
	OutFmt:      "ES10.3E2",
	NNodeOuts:   0,
	OutNd:       []int{1, 2, 3, 4, 5, 6},
	OutList: []string{
		"RootFxr",
		"RootFyr",
		"RootFzr",
		"RootMxr",
		"RootMyr",
		"RootMzr",
		"TipTDxr",
		"TipTDyr",
		"TipTDzr",
		"TipRDxr",
		"TipRDyr",
		"TipRDzr",
	},
	Defaults: map[string]struct{}{
		"DTBeam":           {},
		"NRMax":            {},
		"load_retries":     {},
		"n_fact":           {},
		"refine":           {},
		"stop_tol":         {},
		"tngt_stf_comp":    {},
		"tngt_stf_difftol": {},
		"tngt_stf_fd":      {},
		"tngt_stf_pert":    {},
	},
}
