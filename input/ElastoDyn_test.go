package input_test

import (
	"os"
	"testing"

	"github.com/deslaughter/acdc/input"
)

func TestElastoDynFormat(t *testing.T) {

	text, err := ElastoDynExp.Format()
	if err != nil {
		t.Fatal(err)
	}

	act := input.NewElastoDyn()
	if err := act.Parse(text); err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, ElastoDynExp); err != nil {
		t.Fatal(err)
	}

	os.WriteFile("testdata/test_ElastoDyn.dat", text, 0777)
}

func TestElastoDynParse(t *testing.T) {

	bs, err := os.ReadFile("testdata/AOC_WSt_ElastoDyn.dat")
	if err != nil {
		t.Fatal(err)
	}

	act := input.NewElastoDyn()
	if err = act.Parse(bs); err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, ElastoDynExp); err != nil {
		t.Fatal(err)
	}
}

var ElastoDynExp = &input.ElastoDyn{
	Title:    "FAST certification Test #06",
	Echo:     false,
	Method:   3,
	DT:       0.005,
	FlapDOF1: true, FlapDOF2: true,
	EdgeDOF:  true,
	TeetDOF:  false,
	DrTrDOF:  true,
	GenDOF:   true,
	YawDOF:   false,
	TwFADOF1: false, TwFADOF2: false,
	TwSSDOF1: false, TwSSDOF2: false,
	PtfmSgDOF: false, PtfmSwDOF: false, PtfmHvDOF: false,
	PtfmRDOF: false, PtfmPDOF: false, PtfmYDOF: false,
	OoPDefl:  0,
	IPDefl:   0,
	BlPitch1: 1.54, BlPitch2: 1.54, BlPitch3: 1.54, TeetDefl: 0,
	Azimuth:   0,
	RotSpeed:  0,
	NacYaw:    0,
	TTDspFA:   0,
	TTDspSS:   0,
	PtfmSurge: 0,
	PtfmSway:  0,
	PtfmHeave: 0,
	PtfmRoll:  0,
	PtfmPitch: 0,
	PtfmYaw:   0,
	NumBl:     3,
	TipRad:    7.49,
	HubRad:    0.28,
	PreCone1:  6, PreCone2: 6, PreCone3: 6,
	HubCM:    0,
	UndSling: 0,
	Delta3:   0,
	AzimB1Up: 0,
	OverHang: 1.341,
	ShftGagL: 0.5,
	ShftTilt: 0,
	NacCMxn:  0, NacCMyn: 0, NacCMzn: 0.6,
	NcIMUxn: 0, NcIMUyn: 0, NcIMUzn: 0,
	Twr2Shft:  0.6,
	TowerHt:   24.4,
	TowerBsHt: 0,
	PtfmCMxt:  0, PtfmCMyt: 0, PtfmCMzt: -0,
	PtfmRefzt: -0,
	TipMass1:  5.9, TipMass2: 5.9, TipMass3: 5.9,
	HubMass:   247.3,
	HubIner:   9,
	GenIner:   10,
	NacMass:   1747,
	NacYIner:  976.3,
	YawBrMass: 0,
	PtfmMass:  0,
	PtfmRIner: 0,
	PtfmPIner: 0,
	PtfmYIner: 0,
	BldNodes:  10,
	BldFile1:  "../AOC/AOC_Blade.dat",
	BldFile2:  "../AOC/AOC_Blade.dat",
	BldFile3:  "../AOC/AOC_Blade.dat",
	TeetMod:   0,
	TeetDmpP:  0,
	TeetDmp:   0,
	TeetCDmp:  0,
	TeetSStP:  0,
	TeetHStP:  0,
	TeetSSSp:  0,
	TeetHSSp:  0,
	GBoxEff:   100,
	GBRatio:   28.25,
	DTTorSpr:  600000,
	DTTorDmp:  100000,
	Furling:   false,
	FurlFile:  "unused",
	TwrNodes:  11,
	TwrFile:   "../AOC/AOC_Tower.dat",
	SumPrint:  true,
	OutFile:   1,
	TabDelim:  true,
	OutFmt:    "ES10.3E2",
	TStart:    5,
	DecFact:   10,
	NTwGages:  0,
	TwrGagNd:  []int{0},
	NBlGages:  5,
	BldGagNd:  []int{2, 4, 6, 8, 9},
	OutList: []string{
		"TipDxb3", "TipDyb3", "TipRDxb3", "TipRDyb3", "Spn5ALxb1", "Spn5ALyb1",
		"RotSpeed", "LSSGagV", "HSShftV", "RootFxb3", "RootFyb3", "RootMEdg3",
		"RootMFlp3", "Spn4MLxb1", "Spn4MLyb1", "LSSGagFxs", "LSSGagFys", "LSSGagFzs",
		"LSShftTq", "HSShftTq", "LSShftPwr", "HSShftPwr",
	},
	Defaults: map[string]struct{}{},
}
