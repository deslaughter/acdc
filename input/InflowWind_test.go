package input_test

import (
	"acdc/input"
	"os"
	"testing"
)

func TestInflowWindFormat(t *testing.T) {

	text, err := InflowWindExp.Format()
	if err != nil {
		t.Fatal(err)
	}

	act := input.NewInflowWind()
	if err := act.Parse(text); err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, InflowWindExp); err != nil {
		t.Fatal(err)
	}

	os.WriteFile("testdata/test_InflowWind.dat", text, 0777)
}

func TestInflowWindParse(t *testing.T) {

	bs, err := os.ReadFile("testdata/AOC_WSt_InflowWind.dat")
	if err != nil {
		t.Fatal(err)
	}

	act := input.NewInflowWind()
	if err := act.Parse(bs); err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, InflowWindExp); err != nil {
		t.Fatal(err)
	}
}

var InflowWindExp = &input.InflowWind{
	Title:          "Uniform winds for FAST CertTest #06: 12 m/s with no shear or turbulence",
	Echo:           false,
	WindType:       2,
	PropagationDir: 0,
	VFlowAng:       0,
	NWindVel:       1,
	WindVxiList:    []float64{0},
	WindVyiList:    []float64{0},
	WindVziList:    []float64{25},
	HWindSpeed:     0,
	RefHt:          25,
	PLExp:          0.2,
	FileName_Uni:   "../AOC/Wind/NoShr_12.wnd",
	RefHt_Uni:      25,
	RefLength:      14.898,
	FileName_BTS:   "unused",
	FileNameRoot:   "unused",
	TowerFile:      false,
	FileName_u:     "wasp\\Output\\basic_5u.bin",
	FileName_v:     "wasp\\Output\\basic_5v.bin",
	FileName_w:     "wasp\\Output\\basic_5w.bin",
	NX:             64,
	NY:             32,
	NZ:             32,
	DX:             16,
	DY:             3,
	DZ:             3,
	RefHt_Hawc:     25,
	ScaleMethod:    1,
	SFx:            1,
	SFy:            1,
	SFz:            1,
	SigmaFx:        12,
	SigmaFy:        8,
	SigmaFz:        2,
	URef:           5,
	WindProfile:    2,
	PLExp_Hawc:     0.2,
	Z0:             0.03,
	XOffset:        0,
	SumPrint:       false,
	OutList:        []string{"Wind1VelX", "Wind1VelY", "Wind1VelZ"},
	Defaults:       map[string]struct{}{},
}
