package input_test

import (
	"acdc/input"
	"os"
	"testing"
)

func TestElastoDynBladeFormat(t *testing.T) {

	text, err := ElastoDynBladeExp.Format()
	if err != nil {
		t.Fatal(err)
	}

	act := input.NewElastoDynBlade()
	if err := act.Parse(text); err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, ElastoDynBladeExp); err != nil {
		t.Fatal(err)
	}

	os.WriteFile("testdata/test_ElastoDynBlade.dat", text, 0777)
}

func TestElastoDynBladeParse(t *testing.T) {

	act, err := input.ReadElastoDynBlade("testdata/AOC_Blade.dat")
	if err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, ElastoDynBladeExp); err != nil {
		t.Fatal(err)
	}
}

var ElastoDynBladeExp = &input.ElastoDynBlade{
	Title:     "AOC 15/50 blade file.  GJStiff -> EdgEAof are mostly lies.",
	NBlInpSt:  11,
	BldFlDmp1: 4,
	BldFlDmp2: 4,
	BldEdDmp1: 4,
	FlStTunr1: 1,
	FlStTunr2: 1,
	AdjBlMs:   1,
	AdjFlSt:   1,
	AdjEdSt:   1,
	BlInpSt: []input.ElastoDynBladeBlInpSt{
		{0, 0.25, 7.69, 49.75, 8244850, 7643160},
		{0.1, 0.25, 5.35, 33.57, 5729700, 10364760},
		{0.2, 0.25, 4.66, 21.68, 2690000, 10999800},
		{0.3, 0.25, 4.31, 22.89, 1923350, 12632760},
		{0.4, 0.25, 3.91, 18.99, 1398800, 11158560},
		{0.5, 0.25, 3.25, 16.8, 887700, 5896800},
		{0.6, 0.25, 2.56, 15.59, 498995, 3832920},
		{0.7, 0.25, 1.86, 12.37, 250170, 2336040},
		{0.8, 0.25, 1.16, 11.24, 82314, 1174824},
		{0.9, 0.25, 0.73, 10.11, 55818, 907200},
		{1, 0.25, 0.35, 8.98, 29456, 641844},
	},
	BldFl1Sh2: 0.2506,
	BldFl1Sh3: 1.215,
	BldFl1Sh4: -2.0261,
	BldFl1Sh5: 2.7203,
	BldFl1Sh6: -1.1598,
	BldFl2Sh2: -2.3421,
	BldFl2Sh3: 5.0047,
	BldFl2Sh4: -25.9119,
	BldFl2Sh5: 40.8648,
	BldFl2Sh6: -16.6154,
	BldEdgSh2: 1.8381,
	BldEdgSh3: -2.0103,
	BldEdgSh4: 0.9662,
	BldEdgSh5: 0.9933,
	BldEdgSh6: -0.7874,
	Defaults:  map[string]struct{}{},
}
