package input_test

import (
	"os"
	"testing"

	"github.com/deslaughter/acdc/input"
)

func TestAeroDynBladeFormat(t *testing.T) {

	text, err := AeroDynBladeExp.Format()
	if err != nil {
		t.Fatal(err)
	}

	os.WriteFile("testdata/test_AeroDynBlade.dat", text, 0777)

	act := input.NewAeroDynBlade()
	if err := act.Parse(text); err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, AeroDynBladeExp); err != nil {
		t.Fatal(err)
	}
}

func TestAeroDynBladeParse(t *testing.T) {

	act, err := input.ReadAeroDynBlade("testdata/AOC_AeroDyn_blade.dat")
	if err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, AeroDynBladeExp); err != nil {
		t.Fatal(err)
	}
}

var AeroDynBladeExp = &input.AeroDynBlade{
	Title:    "AOC blade aerodynamic parameters",
	NumBlNds: 12,
	BlNd: []input.AeroDynBladeBlNd{
		{0, 0, 0, 0, 7.69, 0.494, 1},
		{0.235, 0, 0, 0, 7.69, 0.494, 1},
		{0.844, 0, 0, 0, 5.04, 0.579, 1},
		{1.594, 0, 0, 0, 4.6, 0.68, 1},
		{2.344, 0, 0, 0, 4.26, 0.744, 1},
		{3.094, 0, 0, 0, 3.85, 0.738, 2},
		{3.84, 0, 0, 0, 3.15, 0.677, 2},
		{4.59, 0, 0, 0, 2.45, 0.616, 3},
		{5.34, 0, 0, 0, 1.75, 0.558, 3},
		{6.09, 0, 0, 0, 1.05, 0.497, 5},
		{6.84, 0, 0, 0, 0.35, 0.436, 5},
		{7.21, 0, 0, 0, 0.35, 0.436, 5},
	},
	Defaults: map[string]struct{}{},
}
