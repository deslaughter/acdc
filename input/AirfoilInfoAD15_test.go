package input_test

import (
	"os"
	"testing"

	"github.com/deslaughter/acdc/input"
)

func TestAD15AirfoilInfoFormat(t *testing.T) {

	text, err := AD15AirfoilInfoExp.Format()
	if err != nil {
		t.Fatal(err)
	}

	act := input.NewAD15AirfoilInfo()
	if err := act.Parse(text); err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, AD15AirfoilInfoExp); err != nil {
		t.Fatal(err)
	}

	os.WriteFile("testdata/test_AD15AirfoilInfo.dat", text, 0777)
}

func TestAD15AirfoilInfoParse(t *testing.T) {

	text, err := os.ReadFile("testdata/Airfoils/cylinder_cl_AD15.dat")
	if err != nil {
		t.Fatal(err)
	}

	act := input.NewAD15AirfoilInfo()
	if err := act.Parse(text); err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, AD15AirfoilInfoExp); err != nil {
		t.Fatal(err)
	}
}

var AD15AirfoilInfoExp = &input.AD15AirfoilInfo{
	InterpOrd:  0,
	NonDimArea: 1,
	NumCoords:  0,
	BL_file:    "None",
	NumTabs:    1,
	Re:         1,
	Ctrl:       0,
	InclUAdata: true,
	Alpha0:     0,
	Alpha1:     0,
	Alpha2:     0,
	Eta_e:      0,
	C_nalpha:   0,
	T_f0:       3,
	T_V0:       6,
	T_p:        1.7,
	T_VL:       11,
	B1:         0.14,
	B2:         0.53,
	B5:         5,
	A1:         0.3,
	A2:         0.7,
	A5:         1,
	S1:         0,
	S2:         0,
	S3:         0,
	S4:         0,
	Cn1:        0,
	Cn2:        0,
	St_sh:      0.19,
	Cd0:        0.5,
	Cm0:        0,
	K0:         0,
	K1:         0,
	K2:         0,
	K3:         0,
	K1_hat:     0,
	X_cp_bar:   0.2,
	UACutout:   0,
	FiltCutOff: 0,
	NumAlf:     3,
	CoeffData: []input.AD15AirfoilInfoCoeffData{
		{-180, 0.0001, 0.3, -0.0001},
		{0, 0.0001, 0.3, -0.0001},
		{180, 0.0001, 0.3, -0.0001},
	},
	Defaults: map[string]struct{}{
		"InterpOrd":  {},
		"UACutout":   {},
		"filtCutOff": {},
	},
}
