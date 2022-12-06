package input_test

import (
	"os"
	"testing"

	"github.com/deslaughter/acdc/input"
)

func TestElastoDynTowerFormat(t *testing.T) {

	text, err := ElastoDynTowerExp.Format()
	if err != nil {
		t.Fatal(err)
	}

	act := input.NewElastoDynTower()
	if err := act.Parse(text); err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, ElastoDynTowerExp); err != nil {
		t.Fatal(err)
	}

	os.WriteFile("testdata/test_ElastoDynTower.dat", text, 0777)
}

func TestElastoDynTowerParse(t *testing.T) {

	act, err := input.ReadElastoDynTower("testdata/AOC_Tower.dat")
	if err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, ElastoDynTowerExp); err != nil {
		t.Fatal(err)
	}
}

var ElastoDynTowerExp = &input.ElastoDynTower{
	Title:     "AOC tower data.  This is pure fiction.",
	NTwInpSt:  11,
	TwrFADmp1: 3,
	TwrFADmp2: 3,
	TwrSSDmp1: 3,
	TwrSSDmp2: 3,
	FAStTunr1: 1,
	FAStTunr2: 1,
	SSStTunr1: 1,
	SSStTunr2: 1,
	AdjTwMa:   1,
	AdjFASt:   1,
	AdjSSSt:   1,
	TwInpSt: []input.ElastoDynTowerTwInpSt{
		{0, 151.1, 2.35e+09, 2.35e+09},
		{0.1, 141.7, 2.12e+09, 2.12e+09},
		{0.2, 135.6, 1.88e+09, 1.88e+09},
		{0.3, 132.2, 1.65e+09, 1.65e+09},
		{0.4, 130.7, 1.42e+09, 1.42e+09},
		{0.5, 130.3, 1.19e+09, 1.19e+09},
		{0.6, 130.3, 9.57e+08, 9.57e+08},
		{0.7, 129.8, 7.26e+08, 7.26e+08},
		{0.8, 128.2, 4.94e+08, 4.94e+08},
		{0.9, 124.6, 2.62e+08, 2.62e+08},
		{1, 118.4, 3.00e+07, 3.00e+07},
	},
	TwFAM1Sh2: 1.0495,
	TwFAM1Sh3: 0.0694,
	TwFAM1Sh4: -0.289,
	TwFAM1Sh5: 0.3003,
	TwFAM1Sh6: -0.1301,
	TwFAM2Sh2: -25.1012,
	TwFAM2Sh3: 20.1243,
	TwFAM2Sh4: 0.9012,
	TwFAM2Sh5: 16.6452,
	TwFAM2Sh6: -11.5696,
	TwSSM1Sh2: 1.0495,
	TwSSM1Sh3: 0.0694,
	TwSSM1Sh4: -0.289,
	TwSSM1Sh5: 0.3003,
	TwSSM1Sh6: -0.1301,
	TwSSM2Sh2: -25.1012,
	TwSSM2Sh3: 20.1243,
	TwSSM2Sh4: 0.9012,
	TwSSM2Sh5: 16.6452,
	TwSSM2Sh6: -11.5696,
	Defaults:  map[string]struct{}{},
}
