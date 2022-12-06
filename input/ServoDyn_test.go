package input_test

import (
	"os"
	"testing"

	"github.com/deslaughter/acdc/input"
)

func TestServoDynFormat(t *testing.T) {

	text, err := ServoDynExp.Format()
	if err != nil {
		t.Fatal(err)
	}

	act := input.NewServoDyn()
	if err := act.Parse(text); err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, ServoDynExp); err != nil {
		t.Fatal(err)
	}

	os.WriteFile("testdata/test_ServoDyn.dat", text, 0777)
}

func TestServoDynParse(t *testing.T) {

	act, err := input.ReadServoDyn("testdata/AOC_WSt_ServoDyn.dat")
	if err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, ServoDynExp); err != nil {
		t.Fatal(err)
	}
}

var ServoDynExp = &input.ServoDyn{
	Title:        "FAST certification Test #06",
	Echo:         false,
	DT:           0.005,
	PCMode:       0,
	TPCOn:        9999.9,
	TPitManS1:    9999.9,
	TPitManS2:    9999.9,
	TPitManS3:    9999.9,
	PitManRat1:   2,
	PitManRat2:   2,
	PitManRat3:   2,
	BlPitchF1:    1.54,
	BlPitchF2:    1.54,
	BlPitchF3:    1.54,
	VSContrl:     0,
	GenModel:     2,
	GenEff:       89.4,
	GenTiStr:     true,
	GenTiStp:     true,
	SpdGenOn:     9999.9,
	TimGenOn:     6,
	TimGenOf:     25,
	VS_RtGnSp:    9999.9,
	VS_RtTq:      9999.9,
	VS_Rgn2K:     9999.9,
	VS_SlPc:      9999.9,
	SIG_SlPc:     2.222,
	SIG_SySp:     1800,
	SIG_RtTq:     314.3,
	SIG_PORt:     1.75,
	TEC_Freq:     60,
	TEC_NPol:     4,
	TEC_SRes:     0.0492,
	TEC_RRes:     0.000534,
	TEC_VLL:      480,
	TEC_SLR:      0.0001,
	TEC_RLR:      0.0001,
	TEC_MR:       0.00449,
	HSSBrMode:    0,
	THSSBrDp:     9999.9,
	HSSBrDT:      9999.9,
	HSSBrTqF:     9999.9,
	YCMode:       0,
	TYCOn:        9999.9,
	YawNeut:      0,
	YawSpr:       0,
	YawDamp:      0,
	TYawManS:     9999.9,
	YawManRat:    2,
	NacYawF:      0,
	AfCmode:      0,
	AfC_Mean:     0,
	AfC_Amp:      0,
	AfC_Phase:    0,
	NumBStC:      0,
	BStCfiles:    "unused",
	NumNStC:      0,
	NStCfiles:    "unused",
	NumTStC:      0,
	TStCfiles:    "unused",
	NumSStC:      0,
	SStCfiles:    "unused",
	CCmode:       0,
	DLL_FileName: "unused",
	DLL_InFile:   "DISCON.IN",
	DLL_ProcName: "DISCON",
	DLL_DT:       0,
	DLL_Ramp:     false,
	BPCutoff:     9999.9,
	NacYaw_North: 0,
	Ptch_Cntrl:   0,
	Ptch_SetPnt:  0,
	Ptch_Min:     0,
	Ptch_Max:     0,
	PtchRate_Min: 0,
	PtchRate_Max: 0,
	Gain_OM:      0,
	GenSpd_MinOM: 0,
	GenSpd_MaxOM: 0,
	GenSpd_Dem:   0,
	GenTrq_Dem:   0,
	GenPwr_Dem:   0,
	DLL_NumTrq:   0,
	GenSpdTrq:    []input.ServoDynGenSpdTrq{},
	SumPrint:     true,
	OutFile:      1,
	TabDelim:     true,
	OutFmt:       "ES10.3E2",
	TStart:       5,
	OutList:      []string{"GenTq", "GenPwr"},
	Defaults:     map[string]struct{}{"DLL_DT": {}},
}
