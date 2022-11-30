package input_test

import (
	"acdc/input"
	"os"
	"testing"
)

func TestFASTFormat(t *testing.T) {

	text, err := FASTExp.Format()
	if err != nil {
		t.Fatal(err)
	}

	act := input.NewFAST()
	if err := act.Parse(text); err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, FASTExp); err != nil {
		t.Fatal(err)
	}

	os.WriteFile("testdata/test_FAST.dat", text, 0777)
}

func TestFASTParse(t *testing.T) {

	bs, err := os.ReadFile("testdata/AOC_WSt.fst")
	if err != nil {
		t.Fatal(err)
	}

	act := input.NewFAST()
	if err := act.Parse(bs); err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, FASTExp); err != nil {
		t.Fatal(err)
	}
}

var FASTExp = &input.FAST{
	Title:       "FAST Certification Test #06",
	Echo:        true,
	AbortLevel:  "FATAL",
	TMax:        35,
	DT:          0.005,
	InterpOrder: 1,
	NumCrctn:    0,
	DT_UJac:     99999,
	UJacSclFact: 1e+06,
	CompElast:   1,
	CompInflow:  1,
	CompAero:    1,
	CompServo:   1,
	CompHydro:   0,
	CompSub:     0,
	CompMooring: 0,
	CompIce:     0,
	MHK:         0,
	Gravity:     9.80665,
	AirDens:     0.9526,
	WtrDens:     0,
	KinVisc:     1.4639e-05,
	SpdSound:    0,
	Patm:        0,
	Pvap:        0,
	WtrDpth:     0,
	MSL2SWL:     0,
	EDFile:      "AOC_WSt_ElastoDyn.dat",
	BDBldFile1:  "unused",
	BDBldFile2:  "unused",
	BDBldFile3:  "unused",
	InflowFile:  "AOC_WSt_InflowWind.dat",
	AeroFile:    "AOC_WSt_AD.ipt",
	ServoFile:   "AOC_WSt_ServoDyn.dat",
	HydroFile:   "unused",
	SubFile:     "unused",
	MooringFile: "unused",
	IceFile:     "unused",
	SumPrint:    true,
	SttsTime:    5,
	ChkptTime:   99999,
	DT_Out:      0.05,
	TStart:      5,
	OutFileFmt:  0,
	TabDelim:    true,
	OutFmt:      "ES10.3E2",
	Linearize:   false,
	CalcSteady:  false,
	TrimCase:    3,
	TrimTol:     0.001,
	TrimGain:    0.01,
	Twr_Kdmp:    0,
	Bld_Kdmp:    0,
	NLinTimes:   2,
	LinTimes:    []float64{30, 60},
	LinInputs:   1,
	LinOutputs:  1,
	LinOutJac:   false,
	LinOutMod:   false,
	WrVTK:       0,
	VTK_type:    2,
	VTK_fields:  false,
	VTK_fps:     15,
	Defaults:    map[string]struct{}{},
}
