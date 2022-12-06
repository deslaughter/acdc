package input_test

import (
	"os"
	"testing"

	"github.com/deslaughter/acdc/input"
)

func TestAeroDyn15Format(t *testing.T) {

	text, err := AeroDyn15Exp.Format()
	if err != nil {
		t.Fatal(err)
	}

	act := input.NewAeroDyn15()
	if err := act.Parse(text); err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, AeroDyn15Exp); err != nil {
		t.Fatal(err)
	}

	os.WriteFile("testdata/test_AeroDyn15.dat", text, 0777)
}

func TestAeroDyn15Parse(t *testing.T) {

	act, err := input.ReadAeroDyn15("testdata/AOC_YFix_WSt_AD15.ipt")
	if err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, AeroDyn15Exp); err != nil {
		t.Fatal(err)
	}
}

var AeroDyn15Exp = &input.AeroDyn15{
	Title:             "AOC aerodynamic parameters for FAST Certification Test #8.",
	Echo:              false,
	DTAero:            0,
	WakeMod:           1,
	AFAeroMod:         2,
	TwrPotent:         0,
	TwrShadow:         1,
	TwrAero:           false,
	FrozenWake:        false,
	CavitCheck:        false,
	CompAA:            false,
	AA_InputFile:      "unused",
	AirDens:           0,
	KinVisc:           0,
	SpdSound:          0,
	Patm:              0,
	Pvap:              0,
	SkewMod:           2,
	SkewModFactor:     0,
	TipLoss:           true,
	HubLoss:           false,
	TanInd:            true,
	AIDrag:            false,
	TIDrag:            false,
	IndToler:          0,
	MaxIter:           100,
	DBEMT_Mod:         2,
	Tau1_const:        4,
	OLAFInputFileName: "unused",
	UAMod:             3,
	FLookup:           true,
	AFTabMod:          1,
	InCol_Alfa:        1,
	InCol_Cl:          2,
	InCol_Cd:          3,
	InCol_Cm:          0,
	InCol_Cpmin:       0,
	NumAFfiles:        5,
	AFNames: []string{
		"../AOC/Airfoils/S814_1.dat",
		"../AOC/Airfoils/S814_15.dat",
		"../AOC/Airfoils/S812_15.dat",
		"../AOC/Airfoils/S812_2.dat",
		"../AOC/Airfoils/S813_15.dat",
	},
	UseBlCm:   false,
	ADBlFile1: "../AOC/AOC_AeroDyn_blade.dat",
	ADBlFile2: "../AOC/AOC_AeroDyn_blade.dat",
	ADBlFile3: "../AOC/AOC_AeroDyn_blade.dat",
	NumTwrNds: 10,
	TwrNodes: []input.AeroDyn15TwrNodes{
		{TwrElev: 16, TwrDiam: 0.059656972, TwrCd: 2.0115, TwrTI: 0.1},
		{TwrElev: 17, TwrDiam: 0.059656972, TwrCd: 2.0115, TwrTI: 0.1},
		{TwrElev: 18, TwrDiam: 0.059656972, TwrCd: 2.0115, TwrTI: 0.1},
		{TwrElev: 19, TwrDiam: 0.059656972, TwrCd: 2.0115, TwrTI: 0.1},
		{TwrElev: 20, TwrDiam: 0.059656972, TwrCd: 2.0115, TwrTI: 0.1},
		{TwrElev: 21, TwrDiam: 0.059656972, TwrCd: 2.0115, TwrTI: 0.1},
		{TwrElev: 22, TwrDiam: 0.059656972, TwrCd: 2.0115, TwrTI: 0.1},
		{TwrElev: 23, TwrDiam: 0.059656972, TwrCd: 2.0115, TwrTI: 0.1},
		{TwrElev: 24, TwrDiam: 0.059656972, TwrCd: 2.0115, TwrTI: 0.1},
		{TwrElev: 24.4, TwrDiam: 0.059656972, TwrCd: 2.0115, TwrTI: 0.1},
	},
	SumPrint: true,
	NBlOuts:  4,
	BlOutNd:  []int{4, 6, 8, 10},
	NTwOuts:  0,
	TwOutNd:  []int{1, 2, 3, 4, 5},
	OutList: []string{
		"RtSpeed", "RtSkew",
		"B1N1Theta", "B1N2Theta", "B1N3Theta", "B1N4Theta",
		"B1N1AxInd", "B1N2AxInd", "B1N3AxInd", "B1N4AxInd",
		"B1N1Alpha", "B1N2Alpha", "B1N3Alpha", "B1N4Alpha",
		"B1N1Vrel", "B1N2Vrel", "B1N3Vrel", "B1N4Vrel",
		"B1N1Cl", "B1N2Cl", "B1N3Cl", "B1N4Cl",
		"B1N1Cd", "B1N2Cd", "B1N3Cd", "B1N4Cd",
		"B1N1Cm", "B1N2Cm", "B1N3Cm", "B1N4Cm",
	},
	Defaults: map[string]struct{}{
		"DTAero":        {},
		"AirDens":       {},
		"KinVisc":       {},
		"SpdSound":      {},
		"Patm":          {},
		"Pvap":          {},
		"SkewModFactor": {},
		"IndToler":      {},
	},
}
