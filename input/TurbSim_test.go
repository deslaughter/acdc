package input_test

import (
	"os"
	"testing"

	"github.com/deslaughter/acdc/input"
)

func TestTurbSimFormat(t *testing.T) {

	text, err := TurbSimExp.Format()
	if err != nil {
		t.Fatal(err)
	}

	act := input.NewTurbSim()
	if err := act.Parse(text); err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, TurbSimExp); err != nil {
		t.Fatal(err)
	}

	os.WriteFile("testdata/test_TurbSim.dat", text, 0777)
}

func TestTurbSimParse(t *testing.T) {

	bs, err := os.ReadFile("testdata/turbsim.in")
	if err != nil {
		t.Fatal(err)
	}

	act := input.NewTurbSim()
	if err := act.Parse(bs); err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, TurbSimExp); err != nil {
		t.Fatal(err)
	}
}

var TurbSimExp = &input.TurbSim{
	Title:           " Turbsim input file",
	Echo:            false,
	RandSeed1:       1501552846,
	RandSeed2:       "RANLUX",
	WrBHHTP:         false,
	WrFHHTP:         false,
	WrADHH:          false,
	WrADFF:          true,
	WrBLFF:          false,
	WrADTWR:         false,
	WrFMTFF:         false,
	WrACT:           false,
	Clockwise:       false,
	ScaleIEC:        0,
	NumGrid_Z:       25,
	NumGrid_Y:       25,
	TimeStep:        0.05,
	AnalysisTime:    150,
	UsableTime:      -1,
	HubHt:           79.405,
	GridHeight:      158.809,
	GridWidth:       158.809,
	VFlowAng:        0,
	HFlowAng:        0,
	TurbModel:       "IECKAI",
	UserFile:        "unused",
	IECstandard:     "1-ED3",
	IECturbc:        "A",
	IEC_WindType:    "NTM",
	ETMc:            0,
	WindProfileType: "PL",
	ProfileFile:     "unused",
	RefHt:           79.405,
	URef:            11,
	ZJetMax:         0,
	PLExp:           0.2,
	Z0:              0,
	Latitude:        0,
	RICH_NO:         0.05,
	UStar:           0,
	ZI:              0,
	PC_UW:           0,
	PC_UV:           0,
	PC_VW:           0,
	SCMod1:          "",
	SCMod2:          "",
	SCMod3:          "",
	InCDec1:         0,
	InCDec2:         0,
	InCDec3:         0,
	CohExp:          0,
	CTEventPath:     "unused",
	CTEventFile:     "RANDOM",
	Randomize:       true,
	DistScl:         1,
	CTLy:            0.5,
	CTLz:            0.5,
	CTStartTime:     30,
	Defaults: map[string]struct{}{
		"ETMc":     {},
		"UStar":    {},
		"ZI":       {},
		"PC_UW":    {},
		"InCDec3":  {},
		"ZJetMax":  {},
		"SCMod1":   {},
		"SCMod3":   {},
		"CohExp":   {},
		"InCDec2":  {},
		"Z0":       {},
		"Latitude": {},
		"PC_UV":    {},
		"PC_VW":    {},
		"SCMod2":   {},
		"InCDec1":  {},
	},
}
