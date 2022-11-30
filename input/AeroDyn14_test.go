package input_test

import (
	"acdc/input"
	"os"
	"testing"
)

func TestAeroDyn14Format(t *testing.T) {

	text, err := AeroDyn14Exp.Format()
	if err != nil {
		t.Fatal(err)
	}

	act := input.NewAeroDyn14()
	if err := act.Parse(text); err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, AeroDyn14Exp); err != nil {
		t.Fatal(err)
	}

	os.WriteFile("testdata/test_AeroDyn14.dat", text, 0777)
}

func TestAeroDyn14Parse(t *testing.T) {

	act, err := input.ReadAeroDyn14("testdata/AOC_WSt_AD.ipt")
	if err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, AeroDyn14Exp); err != nil {
		t.Fatal(err)
	}
}

var AeroDyn14Exp = &input.AeroDyn14{
	Title:        "AOC aerodynamic parameters for FAST Certification Test #6.",
	StallMod:     "STEADY",
	UseCm:        "NO_CM",
	InfModel:     "EQUIL",
	IndModel:     "SWIRL",
	AToler:       0.005,
	TLModel:      "PRANDtl",
	HLModel:      "PRANDtl",
	TwrShad:      0.3,
	ShadHWid:     0.2,
	T_Shad_Refpt: 1.341,
	AirDens:      0.9526,
	KinVisc:      1.4639e-05,
	DTAero:       0.005,
	NumFoil:      5,
	FoilNm: []string{
		"../AOC/AeroData/S814_1.dat",
		"../AOC/AeroData/S814_15.dat",
		"../AOC/AeroData/S812_15.dat",
		"../AOC/AeroData/S812_2.dat",
		"../AOC/AeroData/S813_15.dat",
	},
	BldNodes: 10,
	BlNd: []input.AeroDyn14BlNd{
		{0.515, 7.69, 0.47, 0.494, 1, "NOPRINT"},
		{1.124, 5.04, 0.748, 0.579, 1, "NOPRINT"},
		{1.874, 4.6, 0.752, 0.68, 1, "NOPRINT"},
		{2.624, 4.26, 0.748, 0.744, 1, "NOPRINT"},
		{3.374, 3.85, 0.752, 0.738, 2, "NOPRINT"},
		{4.12, 3.15, 0.74, 0.677, 2, "NOPRINT"},
		{4.87, 2.45, 0.76, 0.616, 3, "NOPRINT"},
		{5.62, 1.75, 0.74, 0.558, 3, "NOPRINT"},
		{6.37, 1.05, 0.76, 0.497, 5, "NOPRINT"},
		{7.12, 0.35, 0.74, 0.436, 5, "NOPRINT"},
	},
	Defaults: map[string]struct{}{},
}
