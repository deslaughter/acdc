package input_test

import (
	"acdc/input"
	"os"
	"testing"
)

func TestBeamDynBladeFormat(t *testing.T) {

	text, err := BeamDynBladeExp.Format()
	if err != nil {
		t.Fatal(err)
	}

	os.WriteFile("testdata/test_BeamDynBlade.dat", text, 0777)

	act := input.NewBeamDynBlade()
	if err := act.Parse(text); err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, BeamDynBladeExp); err != nil {
		t.Fatal(err)
	}
}

func TestBeamDynBladeParse(t *testing.T) {

	act, err := input.ReadBeamDynBlade("testdata/NRELOffshrBsline5MW_BeamDyn_Blade.dat")
	if err != nil {
		t.Fatal(err)
	}

	if err := CompareStructs(act, BeamDynBladeExp); err != nil {
		t.Fatal(err)
	}
}

var BeamDynBladeExp = &input.BeamDynBlade{
	Title:         " Test Format 1",
	Station_total: 49,
	Damp_type:     1,
	Mu:            []float64{0.001, 0.001, 0.001, 0.0014, 0.0022, 0.0022},
	DistProps: []input.BeamDynBladeDistProps{
		{
			Station_eta: 0,
			Stiffness_matrix: [][]float64{
				{9.72948e+08, 0, 0, 0, 0, 0},
				{0, 9.72948e+08, 0, 0, 0, 0},
				{0, 0, 9.72948e+09, 0, 0, 0},
				{0, 0, 0, 1.81136e+10, 0, 0},
				{0, 0, 0, 0, 1.811e+10, 0},
				{0, 0, 0, 0, 0, 5.5644e+09},
			},
			Mass_matrix: [][]float64{
				{678.935, 0, 0, 0, 0, 0},
				{0, 678.935, 0, -0, 0, 0},
				{0, 0, 678.935, 0, 0, 0},
				{0, -0, 0, 973.04, 0, 0},
				{0, 0, 0, 0, 972.86, 0},
				{0, 0, 0, 0, 0, 1945.9},
			},
		},
		{
			Station_eta: 0.00325,
			Stiffness_matrix: [][]float64{
				{9.72948e+08, 0, 0, 0, 0, 0},
				{0, 9.72948e+08, 0, 0, 0, 0},
				{0, 0, 9.72948e+09, 0, 0, 0},
				{0, 0, 0, 1.81136e+10, 0, 0},
				{0, 0, 0, 0, 1.811e+10, 0},
				{0, 0, 0, 0, 0, 5.5644e+09},
			},
			Mass_matrix: [][]float64{
				{678.935, 0, 0, 0, 0, 0},
				{0, 678.935, 0, -0, 0, 0},
				{0, 0, 678.935, 0, 0, 0},
				{0, -0, 0, 973.04, 0, 0},
				{0, 0, 0, 0, 972.86, 0},
				{0, 0, 0, 0, 0, 1945.9},
			},
		},
		{
			Station_eta: 0.01951,
			Stiffness_matrix: [][]float64{
				{1.07895e+09, 0, 0, 0, 0, 0},
				{0, 1.07895e+09, 0, 0, 0, 0},
				{0, 0, 1.07895e+10, 0, 0, 0},
				{0, 0, 0, 1.95586e+10, 0, 0},
				{0, 0, 0, 0, 1.94249e+10, 0},
				{0, 0, 0, 0, 0, 5.43159e+09}},
			Mass_matrix: [][]float64{
				{773.363, 0, 0, 0, 0, 0},
				{0, 773.363, 0, -0, 0, 0},
				{0, 0, 773.363, 0, 0, 0},
				{0, -0, 0, 1066.38, 0, 0},
				{0, 0, 0, 0, 1091.52, 0},
				{0, 0, 0, 0, 0, 2157.9},
			},
		},
		{
			Station_eta: 0.03577,
			Stiffness_matrix: [][]float64{
				{1.006723e+09, 0, 0, 0, 0, 0},
				{0, 1.006723e+09, 0, 0, 0, 0},
				{0, 0, 1.006723e+10, 0, 0, 0},
				{0, 0, 0, 1.94978e+10, 0, 0},
				{0, 0, 0, 0, 1.74559e+10, 0},
				{0, 0, 0, 0, 0, 4.99398e+09}},
			Mass_matrix: [][]float64{
				{740.55, 0, 0, 0, 0, 0},
				{0, 740.55, 0, -0, 0, 0},
				{0, 0, 740.55, 0, 0, 0},
				{0, -0, 0, 1047.36, 0, 0},
				{0, 0, 0, 0, 966.09, 0},
				{0, 0, 0, 0, 0, 2013.45},
			},
		},
		{
			Station_eta: 0.05203,
			Stiffness_matrix: [][]float64{
				{9.86778e+08, 0, 0, 0, 0, 0},
				{0, 9.86778e+08, 0, 0, 0, 0},
				{0, 0, 9.86778e+09, 0, 0, 0},
				{0, 0, 0, 1.97888e+10, 0, 0},
				{0, 0, 0, 0, 1.52874e+10, 0},
				{0, 0, 0, 0, 0, 4.66659e+09}},
			Mass_matrix: [][]float64{
				{740.042, 0, 0, 0, 0, 0},
				{0, 740.042, 0, -0, 0, 0},
				{0, 0, 740.042, 0, 0, 0},
				{0, -0, 0, 1099.75, 0, 0},
				{0, 0, 0, 0, 873.81, 0},
				{0, 0, 0, 0, 0, 1973.56},
			},
		},
		{
			Station_eta: 0.06829,
			Stiffness_matrix: [][]float64{
				{7.60786e+08, 0, 0, 0, 0, 0},
				{0, 7.60786e+08, 0, 0, 0, 0},
				{0, 0, 7.60786e+09, 0, 0, 0},
				{0, 0, 0, 1.48585e+10, 0, 0},
				{0, 0, 0, 0, 1.07824e+10, 0},
				{0, 0, 0, 0, 0, 3.47471e+09}},
			Mass_matrix: [][]float64{
				{592.496, 0, 0, 0, 0, 0},
				{0, 592.496, 0, -0, 0, 0},
				{0, 0, 592.496, 0, 0, 0},
				{0, -0, 0, 873.02, 0, 0},
				{0, 0, 0, 0, 648.55, 0},
				{0, 0, 0, 0, 0, 1521.57},
			},
		},
		{
			Station_eta: 0.08455,
			Stiffness_matrix: [][]float64{
				{5.49126e+08, 0, 0, 0, 0, 0},
				{0, 5.49126e+08, 0, 0, 0, 0},
				{0, 0, 5.49126e+09, 0, 0, 0},
				{0, 0, 0, 1.02206e+10, 0, 0},
				{0, 0, 0, 0, 7.22972e+09, 0},
				{0, 0, 0, 0, 0, 2.32354e+09}},
			Mass_matrix: [][]float64{
				{450.275, 0, 0, 0, 0, 0},
				{0, 450.275, 0, -0, 0, 0},
				{0, 0, 450.275, 0, 0, 0},
				{0, -0, 0, 641.49, 0, 0},
				{0, 0, 0, 0, 456.76, 0},
				{0, 0, 0, 0, 0, 1098.25},
			},
		},
		{
			Station_eta: 0.10081,
			Stiffness_matrix: [][]float64{
				{4.9713e+08, 0, 0, 0, 0, 0},
				{0, 4.9713e+08, 0, 0, 0, 0},
				{0, 0, 4.9713e+09, 0, 0, 0},
				{0, 0, 0, 9.1447e+09, 0, 0},
				{0, 0, 0, 0, 6.30954e+09, 0},
				{0, 0, 0, 0, 0, 1.90787e+09}},
			Mass_matrix: [][]float64{
				{424.054, 0, 0, 0, 0, 0},
				{0, 424.054, 0, -0, 0, 0},
				{0, 0, 424.054, 0, 0, 0},
				{0, -0, 0, 593.73, 0, 0},
				{0, 0, 0, 0, 400.53, 0},
				{0, 0, 0, 0, 0, 994.26},
			},
		},
		{
			Station_eta: 0.11707,
			Stiffness_matrix: [][]float64{
				{4.49395e+08, 0, 0, 0, 0, 0},
				{0, 4.49395e+08, 0, 0, 0, 0},
				{0, 0, 4.49395e+09, 0, 0, 0},
				{0, 0, 0, 8.06316e+09, 0, 0},
				{0, 0, 0, 0, 5.52836e+09, 0},
				{0, 0, 0, 0, 0, 1.57036e+09}},
			Mass_matrix: [][]float64{
				{400.638, 0, 0, 0, 0, 0},
				{0, 400.638, 0, -0, 0, 0},
				{0, 0, 400.638, 0, 0, 0},
				{0, -0, 0, 547.18, 0, 0},
				{0, 0, 0, 0, 351.61, 0},
				{0, 0, 0, 0, 0, 898.79},
			},
		},
		{
			Station_eta: 0.13335,
			Stiffness_matrix: [][]float64{
				{4.0348e+08, 0, 0, 0, 0, 0},
				{0, 4.0348e+08, 0, 0, 0, 0},
				{0, 0, 4.0348e+09, 0, 0, 0},
				{0, 0, 0, 6.88444e+09, 0, 0},
				{0, 0, 0, 0, 4.98006e+09, 0},
				{0, 0, 0, 0, 0, 1.15826e+09}},
			Mass_matrix: [][]float64{
				{382.062, 0, 0, 0, 0, 0},
				{0, 382.062, 0, -0, 0, 0},
				{0, 0, 382.062, 0, 0, 0},
				{0, -0, 0, 490.84, 0, 0},
				{0, 0, 0, 0, 316.12, 0},
				{0, 0, 0, 0, 0, 806.96},
			},
		},
		{
			Station_eta: 0.14959,
			Stiffness_matrix: [][]float64{
				{4.03729e+08, 0, 0, 0, 0, 0},
				{0, 4.03729e+08, 0, 0, 0, 0},
				{0, 0, 4.03729e+09, 0, 0, 0},
				{0, 0, 0, 7.00918e+09, 0, 0},
				{0, 0, 0, 0, 4.93684e+09, 0},
				{0, 0, 0, 0, 0, 1.00212e+09}},
			Mass_matrix: [][]float64{
				{399.655, 0, 0, 0, 0, 0},
				{0, 399.655, 0, -0, 0, 0},
				{0, 0, 399.655, 0, 0, 0},
				{0, -0, 0, 503.86, 0, 0},
				{0, 0, 0, 0, 303.6, 0},
				{0, 0, 0, 0, 0, 807.46},
			},
		},
		{
			Station_eta: 0.16585,
			Stiffness_matrix: [][]float64{
				{4.16972e+08, 0, 0, 0, 0, 0},
				{0, 4.16972e+08, 0, 0, 0, 0},
				{0, 0, 4.16972e+09, 0, 0, 0},
				{0, 0, 0, 7.16768e+09, 0, 0},
				{0, 0, 0, 0, 4.69166e+09, 0},
				{0, 0, 0, 0, 0, 8.559e+08}},
			Mass_matrix: [][]float64{
				{426.321, 0, 0, 0, 0, 0},
				{0, 426.321, 0, -0, 0, 0},
				{0, 0, 426.321, 0, 0, 0},
				{0, -0, 0, 544.7, 0, 0},
				{0, 0, 0, 0, 289.24, 0},
				{0, 0, 0, 0, 0, 833.94},
			},
		},
		{Station_eta: 0.18211,
			Stiffness_matrix: [][]float64{
				{4.08235e+08, 0, 0, 0, 0, 0},
				{0, 4.08235e+08, 0, 0, 0, 0},
				{0, 0, 4.08235e+09, 0, 0, 0},
				{0, 0, 0, 7.27166e+09, 0, 0},
				{0, 0, 0, 0, 3.94946e+09, 0},
				{0, 0, 0, 0, 0, 6.7227e+08}},
			Mass_matrix: [][]float64{
				{416.82, 0, 0, 0, 0, 0},
				{0, 416.82, 0, -0, 0, 0},
				{0, 0, 416.82, 0, 0, 0},
				{0, -0, 0, 569.9, 0, 0},
				{0, 0, 0, 0, 246.57, 0},
				{0, 0, 0, 0, 0, 816.47},
			},
		},
		{Station_eta: 0.19837,
			Stiffness_matrix: [][]float64{
				{4.08597e+08, 0, 0, 0, 0, 0},
				{0, 4.08597e+08, 0, 0, 0, 0},
				{0, 0, 4.08597e+09, 0, 0, 0},
				{0, 0, 0, 7.0817e+09, 0, 0},
				{0, 0, 0, 0, 3.38652e+09, 0},
				{0, 0, 0, 0, 0, 5.4749e+08}},
			Mass_matrix: [][]float64{
				{406.186, 0, 0, 0, 0, 0},
				{0, 406.186, 0, -0, 0, 0},
				{0, 0, 406.186, 0, 0, 0},
				{0, -0, 0, 601.28, 0, 0},
				{0, 0, 0, 0, 215.91, 0},
				{0, 0, 0, 0, 0, 817.19},
			},
		},
		{Station_eta: 0.21465,
			Stiffness_matrix: [][]float64{
				{3.66834e+08, 0, 0, 0, 0, 0},
				{0, 3.66834e+08, 0, 0, 0, 0},
				{0, 0, 3.66834e+09, 0, 0, 0},
				{0, 0, 0, 6.24453e+09, 0, 0},
				{0, 0, 0, 0, 2.93374e+09, 0},
				{0, 0, 0, 0, 0, 4.4884e+08}},
			Mass_matrix: [][]float64{
				{381.42, 0, 0, 0, 0, 0},
				{0, 381.42, 0, -0, 0, 0},
				{0, 0, 381.42, 0, 0, 0},
				{0, -0, 0, 546.56, 0, 0},
				{0, 0, 0, 0, 187.11, 0},
				{0, 0, 0, 0, 0, 733.67},
			},
		},
		{Station_eta: 0.23089,
			Stiffness_matrix: [][]float64{
				{3.14776e+08, 0, 0, 0, 0, 0},
				{0, 3.14776e+08, 0, 0, 0, 0},
				{0, 0, 3.14776e+09, 0, 0, 0},
				{0, 0, 0, 5.04896e+09, 0, 0},
				{0, 0, 0, 0, 2.56896e+09, 0},
				{0, 0, 0, 0, 0, 3.3592e+08}},
			Mass_matrix: [][]float64{
				{352.822, 0, 0, 0, 0, 0},
				{0, 352.822, 0, -0, 0, 0},
				{0, 0, 352.822, 0, 0, 0},
				{0, -0, 0, 468.71, 0, 0},
				{0, 0, 0, 0, 160.84, 0},
				{0, 0, 0, 0, 0, 629.55},
			},
		},
		{Station_eta: 0.24715,
			Stiffness_matrix: [][]float64{
				{3.01158e+08, 0, 0, 0, 0, 0},
				{0, 3.01158e+08, 0, 0, 0, 0},
				{0, 0, 3.01158e+09, 0, 0, 0},
				{0, 0, 0, 4.94849e+09, 0, 0},
				{0, 0, 0, 0, 2.38865e+09, 0},
				{0, 0, 0, 0, 0, 3.1135e+08}},
			Mass_matrix: [][]float64{
				{349.477, 0, 0, 0, 0, 0},
				{0, 349.477, 0, -0, 0, 0},
				{0, 0, 349.477, 0, 0, 0},
				{0, -0, 0, 453.76, 0, 0},
				{0, 0, 0, 0, 148.56, 0},
				{0, 0, 0, 0, 0, 602.32},
			},
		},
		{Station_eta: 0.26341,
			Stiffness_matrix: [][]float64{
				{2.88262e+08, 0, 0, 0, 0, 0},
				{0, 2.88262e+08, 0, 0, 0, 0},
				{0, 0, 2.88262e+09, 0, 0, 0},
				{0, 0, 0, 4.80802e+09, 0, 0},
				{0, 0, 0, 0, 2.27199e+09, 0},
				{0, 0, 0, 0, 0, 2.9194e+08}},
			Mass_matrix: [][]float64{
				{346.538, 0, 0, 0, 0, 0},
				{0, 346.538, 0, -0, 0, 0},
				{0, 0, 346.538, 0, 0, 0},
				{0, -0, 0, 436.22, 0, 0},
				{0, 0, 0, 0, 140.3, 0},
				{0, 0, 0, 0, 0, 576.52},
			},
		},
		{Station_eta: 0.29595,
			Stiffness_matrix: [][]float64{
				{2.61397e+08, 0, 0, 0, 0, 0},
				{0, 2.61397e+08, 0, 0, 0, 0},
				{0, 0, 2.61397e+09, 0, 0, 0},
				{0, 0, 0, 4.5014e+09, 0, 0},
				{0, 0, 0, 0, 2.05005e+09, 0},
				{0, 0, 0, 0, 0, 2.61e+08}},
			Mass_matrix: [][]float64{
				{339.333, 0, 0, 0, 0, 0},
				{0, 339.333, 0, -0, 0, 0},
				{0, 0, 339.333, 0, 0, 0},
				{0, -0, 0, 398.18, 0, 0},
				{0, 0, 0, 0, 124.61, 0},
				{0, 0, 0, 0, 0, 522.79},
			},
		},
		{Station_eta: 0.32846,
			Stiffness_matrix: [][]float64{
				{2.35748e+08, 0, 0, 0, 0, 0},
				{0, 2.35748e+08, 0, 0, 0, 0},
				{0, 0, 2.35748e+09, 0, 0, 0},
				{0, 0, 0, 4.24407e+09, 0, 0},
				{0, 0, 0, 0, 1.82825e+09, 0},
				{0, 0, 0, 0, 0, 2.2882e+08}},
			Mass_matrix: [][]float64{
				{330.004, 0, 0, 0, 0, 0},
				{0, 330.004, 0, -0, 0, 0},
				{0, 0, 330.004, 0, 0, 0},
				{0, -0, 0, 362.08, 0, 0},
				{0, 0, 0, 0, 109.42, 0},
				{0, 0, 0, 0, 0, 471.5},
			},
		},
		{Station_eta: 0.36098,
			Stiffness_matrix: [][]float64{
				{2.14686e+08, 0, 0, 0, 0, 0},
				{0, 2.14686e+08, 0, 0, 0, 0},
				{0, 0, 2.14686e+09, 0, 0, 0},
				{0, 0, 0, 3.99528e+09, 0, 0},
				{0, 0, 0, 0, 1.58871e+09, 0},
				{0, 0, 0, 0, 0, 2.0075e+08}},
			Mass_matrix: [][]float64{
				{321.99, 0, 0, 0, 0, 0},
				{0, 321.99, 0, -0, 0, 0},
				{0, 0, 321.99, 0, 0, 0},
				{0, -0, 0, 335.01, 0, 0},
				{0, 0, 0, 0, 94.36, 0},
				{0, 0, 0, 0, 0, 429.37},
			},
		},
		{Station_eta: 0.3935,
			Stiffness_matrix: [][]float64{
				{1.94409e+08, 0, 0, 0, 0, 0},
				{0, 1.94409e+08, 0, 0, 0, 0},
				{0, 0, 1.94409e+09, 0, 0, 0},
				{0, 0, 0, 3.75076e+09, 0, 0},
				{0, 0, 0, 0, 1.36193e+09, 0},
				{0, 0, 0, 0, 0, 1.7438e+08}},
			Mass_matrix: [][]float64{
				{313.82, 0, 0, 0, 0, 0},
				{0, 313.82, 0, -0, 0, 0},
				{0, 0, 313.82, 0, 0, 0},
				{0, -0, 0, 308.57, 0, 0},
				{0, 0, 0, 0, 80.24, 0},
				{0, 0, 0, 0, 0, 388.81},
			},
		},
		{Station_eta: 0.42602,
			Stiffness_matrix: [][]float64{
				{1.6327e+08, 0, 0, 0, 0, 0},
				{0, 1.6327e+08, 0, 0, 0, 0},
				{0, 0, 1.6327e+09, 0, 0, 0},
				{0, 0, 0, 3.44714e+09, 0, 0},
				{0, 0, 0, 0, 1.10238e+09, 0},
				{0, 0, 0, 0, 0, 1.4447e+08}},
			Mass_matrix: [][]float64{
				{294.734, 0, 0, 0, 0, 0},
				{0, 294.734, 0, -0, 0, 0},
				{0, 0, 294.734, 0, 0, 0},
				{0, -0, 0, 263.87, 0, 0},
				{0, 0, 0, 0, 62.67, 0},
				{0, 0, 0, 0, 0, 326.54},
			},
		},
		{Station_eta: 0.45855,
			Stiffness_matrix: [][]float64{
				{1.4324e+08, 0, 0, 0, 0, 0},
				{0, 1.4324e+08, 0, 0, 0, 0},
				{0, 0, 1.4324e+09, 0, 0, 0},
				{0, 0, 0, 3.13907e+09, 0, 0},
				{0, 0, 0, 0, 8.758e+08, 0},
				{0, 0, 0, 0, 0, 1.1998e+08}},
			Mass_matrix: [][]float64{
				{287.12, 0, 0, 0, 0, 0},
				{0, 287.12, 0, -0, 0, 0},
				{0, 0, 287.12, 0, 0, 0},
				{0, -0, 0, 237.06, 0, 0},
				{0, 0, 0, 0, 49.42, 0},
				{0, 0, 0, 0, 0, 286.48},
			},
		},
		{Station_eta: 0.49106,
			Stiffness_matrix: [][]float64{
				{1.16876e+08, 0, 0, 0, 0, 0},
				{0, 1.16876e+08, 0, 0, 0, 0},
				{0, 0, 1.16876e+09, 0, 0, 0},
				{0, 0, 0, 2.73424e+09, 0, 0},
				{0, 0, 0, 0, 6.813e+08, 0},
				{0, 0, 0, 0, 0, 8.119e+07}},
			Mass_matrix: [][]float64{
				{263.343, 0, 0, 0, 0, 0},
				{0, 263.343, 0, -0, 0, 0},
				{0, 0, 263.343, 0, 0, 0},
				{0, -0, 0, 196.41, 0, 0},
				{0, 0, 0, 0, 37.34, 0},
				{0, 0, 0, 0, 0, 233.75},
			},
		},
		{Station_eta: 0.52358,
			Stiffness_matrix: [][]float64{
				{1.04743e+08, 0, 0, 0, 0, 0},
				{0, 1.04743e+08, 0, 0, 0, 0},
				{0, 0, 1.04743e+09, 0, 0, 0},
				{0, 0, 0, 2.55487e+09, 0, 0},
				{0, 0, 0, 0, 5.3472e+08, 0},
				{0, 0, 0, 0, 0, 6.909e+07}},
			Mass_matrix: [][]float64{
				{253.207, 0, 0, 0, 0, 0},
				{0, 253.207, 0, -0, 0, 0},
				{0, 0, 253.207, 0, 0, 0},
				{0, -0, 0, 180.34, 0, 0},
				{0, 0, 0, 0, 29.14, 0},
				{0, 0, 0, 0, 0, 209.48},
			},
		},
		{Station_eta: 0.5561,
			Stiffness_matrix: [][]float64{
				{9.2295e+07, 0, 0, 0, 0, 0},
				{0, 9.2295e+07, 0, 0, 0, 0},
				{0, 0, 9.2295e+08, 0, 0, 0},
				{0, 0, 0, 2.33403e+09, 0, 0},
				{0, 0, 0, 0, 4.089e+08, 0},
				{0, 0, 0, 0, 0, 5.745e+07}},
			Mass_matrix: [][]float64{
				{241.666, 0, 0, 0, 0, 0},
				{0, 241.666, 0, -0, 0, 0},
				{0, 0, 241.666, 0, 0, 0},
				{0, -0, 0, 162.43, 0, 0},
				{0, 0, 0, 0, 22.16, 0},
				{0, 0, 0, 0, 0, 184.59},
			},
		},
		{Station_eta: 0.58862,
			Stiffness_matrix: [][]float64{
				{7.6082e+07, 0, 0, 0, 0, 0},
				{0, 7.6082e+07, 0, 0, 0, 0},
				{0, 0, 7.6082e+08, 0, 0, 0},
				{0, 0, 0, 1.82873e+09, 0, 0},
				{0, 0, 0, 0, 3.1454e+08, 0},
				{0, 0, 0, 0, 0, 4.592e+07}},
			Mass_matrix: [][]float64{
				{220.638, 0, 0, 0, 0, 0},
				{0, 220.638, 0, -0, 0, 0},
				{0, 0, 220.638, 0, 0, 0},
				{0, -0, 0, 134.83, 0, 0},
				{0, 0, 0, 0, 17.33, 0},
				{0, 0, 0, 0, 0, 152.16},
			},
		},
		{Station_eta: 0.62115,
			Stiffness_matrix: [][]float64{
				{6.4803e+07, 0, 0, 0, 0, 0},
				{0, 6.4803e+07, 0, 0, 0, 0},
				{0, 0, 6.4803e+08, 0, 0, 0},
				{0, 0, 0, 1.5841e+09, 0, 0},
				{0, 0, 0, 0, 2.3863e+08, 0},
				{0, 0, 0, 0, 0, 3.598e+07}},
			Mass_matrix: [][]float64{
				{200.293, 0, 0, 0, 0, 0},
				{0, 200.293, 0, -0, 0, 0},
				{0, 0, 200.293, 0, 0, 0},
				{0, -0, 0, 116.3, 0, 0},
				{0, 0, 0, 0, 13.3, 0},
				{0, 0, 0, 0, 0, 129.6},
			},
		},
		{Station_eta: 0.65366,
			Stiffness_matrix: [][]float64{
				{5.397e+07, 0, 0, 0, 0, 0},
				{0, 5.397e+07, 0, 0, 0, 0},
				{0, 0, 5.397e+08, 0, 0, 0},
				{0, 0, 0, 1.32336e+09, 0, 0},
				{0, 0, 0, 0, 1.7588e+08, 0},
				{0, 0, 0, 0, 0, 2.744e+07}},
			Mass_matrix: [][]float64{
				{179.404, 0, 0, 0, 0, 0},
				{0, 179.404, 0, -0, 0, 0},
				{0, 0, 179.404, 0, 0, 0},
				{0, -0, 0, 97.98, 0, 0},
				{0, 0, 0, 0, 9.96, 0},
				{0, 0, 0, 0, 0, 107.94},
			},
		},
		{Station_eta: 0.68618,
			Stiffness_matrix: [][]float64{
				{5.3115e+07, 0, 0, 0, 0, 0},
				{0, 5.3115e+07, 0, 0, 0, 0},
				{0, 0, 5.3115e+08, 0, 0, 0},
				{0, 0, 0, 1.18368e+09, 0, 0},
				{0, 0, 0, 0, 1.2601e+08, 0},
				{0, 0, 0, 0, 0, 2.09e+07}},
			Mass_matrix: [][]float64{
				{165.094, 0, 0, 0, 0, 0},
				{0, 165.094, 0, -0, 0, 0},
				{0, 0, 165.094, 0, 0, 0},
				{0, -0, 0, 98.93, 0, 0},
				{0, 0, 0, 0, 7.3, 0},
				{0, 0, 0, 0, 0, 106.23},
			},
		},
		{Station_eta: 0.7187,
			Stiffness_matrix: [][]float64{
				{4.6001e+07, 0, 0, 0, 0, 0},
				{0, 4.6001e+07, 0, 0, 0, 0},
				{0, 0, 4.6001e+08, 0, 0, 0},
				{0, 0, 0, 1.02016e+09, 0, 0},
				{0, 0, 0, 0, 1.0726e+08, 0},
				{0, 0, 0, 0, 0, 1.854e+07}},
			Mass_matrix: [][]float64{
				{154.411, 0, 0, 0, 0, 0},
				{0, 154.411, 0, -0, 0, 0},
				{0, 0, 154.411, 0, 0, 0},
				{0, -0, 0, 85.78, 0, 0},
				{0, 0, 0, 0, 6.22, 0},
				{0, 0, 0, 0, 0, 92},
			},
		},
		{Station_eta: 0.75122,
			Stiffness_matrix: [][]float64{
				{3.7575e+07, 0, 0, 0, 0, 0},
				{0, 3.7575e+07, 0, 0, 0, 0},
				{0, 0, 3.7575e+08, 0, 0, 0},
				{0, 0, 0, 7.9781e+08, 0, 0},
				{0, 0, 0, 0, 9.088e+07, 0},
				{0, 0, 0, 0, 0, 1.628e+07}},
			Mass_matrix: [][]float64{
				{138.935, 0, 0, 0, 0, 0},
				{0, 138.935, 0, -0, 0, 0},
				{0, 0, 138.935, 0, 0, 0},
				{0, -0, 0, 69.96, 0, 0},
				{0, 0, 0, 0, 5.19, 0},
				{0, 0, 0, 0, 0, 75.15},
			},
		},
		{Station_eta: 0.78376,
			Stiffness_matrix: [][]float64{
				{3.2889e+07, 0, 0, 0, 0, 0},
				{0, 3.2889e+07, 0, 0, 0, 0},
				{0, 0, 3.2889e+08, 0, 0, 0},
				{0, 0, 0, 7.0961e+08, 0, 0},
				{0, 0, 0, 0, 7.631e+07, 0},
				{0, 0, 0, 0, 0, 1.453e+07}},
			Mass_matrix: [][]float64{
				{129.555, 0, 0, 0, 0, 0},
				{0, 129.555, 0, -0, 0, 0},
				{0, 0, 129.555, 0, 0, 0},
				{0, -0, 0, 61.41, 0, 0},
				{0, 0, 0, 0, 4.36, 0},
				{0, 0, 0, 0, 0, 65.77},
			},
		},
		{Station_eta: 0.81626,
			Stiffness_matrix: [][]float64{
				{2.4404e+07, 0, 0, 0, 0, 0},
				{0, 2.4404e+07, 0, 0, 0, 0},
				{0, 0, 2.4404e+08, 0, 0, 0},
				{0, 0, 0, 5.1819e+08, 0, 0},
				{0, 0, 0, 0, 6.105e+07, 0},
				{0, 0, 0, 0, 0, 9.07e+06}},
			Mass_matrix: [][]float64{
				{107.264, 0, 0, 0, 0, 0},
				{0, 107.264, 0, -0, 0, 0},
				{0, 0, 107.264, 0, 0, 0},
				{0, -0, 0, 45.44, 0, 0},
				{0, 0, 0, 0, 3.36, 0},
				{0, 0, 0, 0, 0, 48.8},
			},
		},
		{Station_eta: 0.84878,
			Stiffness_matrix: [][]float64{
				{2.116e+07, 0, 0, 0, 0, 0},
				{0, 2.116e+07, 0, 0, 0, 0},
				{0, 0, 2.116e+08, 0, 0, 0},
				{0, 0, 0, 4.5487e+08, 0, 0},
				{0, 0, 0, 0, 4.948e+07, 0},
				{0, 0, 0, 0, 0, 8.06e+06}},
			Mass_matrix: [][]float64{
				{98.776, 0, 0, 0, 0, 0},
				{0, 98.776, 0, -0, 0, 0},
				{0, 0, 98.776, 0, 0, 0},
				{0, -0, 0, 39.57, 0, 0},
				{0, 0, 0, 0, 2.75, 0},
				{0, 0, 0, 0, 0, 42.32},
			},
		},
		{Station_eta: 0.8813,
			Stiffness_matrix: [][]float64{
				{1.8152e+07, 0, 0, 0, 0, 0},
				{0, 1.8152e+07, 0, 0, 0, 0},
				{0, 0, 1.8152e+08, 0, 0, 0},
				{0, 0, 0, 3.9512e+08, 0, 0},
				{0, 0, 0, 0, 3.936e+07, 0},
				{0, 0, 0, 0, 0, 7.08e+06}},
			Mass_matrix: [][]float64{
				{90.248, 0, 0, 0, 0, 0},
				{0, 90.248, 0, -0, 0, 0},
				{0, 0, 90.248, 0, 0, 0},
				{0, -0, 0, 34.09, 0, 0},
				{0, 0, 0, 0, 2.21, 0},
				{0, 0, 0, 0, 0, 36.3},
			},
		},
		{
			Station_eta: 0.89756,
			Stiffness_matrix: [][]float64{
				{1.6025e+07, 0, 0, 0, 0, 0},
				{0, 1.6025e+07, 0, 0, 0, 0},
				{0, 0, 1.6025e+08, 0, 0, 0},
				{0, 0, 0, 3.5372e+08, 0, 0},
				{0, 0, 0, 0, 3.467e+07, 0},
				{0, 0, 0, 0, 0, 6.09e+06}},
			Mass_matrix: [][]float64{
				{83.001, 0, 0, 0, 0, 0},
				{0, 83.001, 0, -0, 0, 0},
				{0, 0, 83.001, 0, 0, 0},
				{0, -0, 0, 30.12, 0, 0},
				{0, 0, 0, 0, 1.93, 0},
				{0, 0, 0, 0, 0, 32.05},
			},
		},
		{
			Station_eta: 0.91382,
			Stiffness_matrix: [][]float64{
				{1.0923e+07, 0, 0, 0, 0, 0},
				{0, 1.0923e+07, 0, 0, 0, 0},
				{0, 0, 1.0923e+08, 0, 0, 0},
				{0, 0, 0, 3.0473e+08, 0, 0},
				{0, 0, 0, 0, 3.041e+07, 0},
				{0, 0, 0, 0, 0, 5.75e+06}},
			Mass_matrix: [][]float64{
				{72.906, 0, 0, 0, 0, 0},
				{0, 72.906, 0, -0, 0, 0},
				{0, 0, 72.906, 0, 0, 0},
				{0, -0, 0, 20.15, 0, 0},
				{0, 0, 0, 0, 1.69, 0},
				{0, 0, 0, 0, 0, 21.84},
			},
		},
		{
			Station_eta: 0.93008,
			Stiffness_matrix: [][]float64{
				{1.0008e+07, 0, 0, 0, 0, 0},
				{0, 1.0008e+07, 0, 0, 0, 0},
				{0, 0, 1.0008e+08, 0, 0, 0},
				{0, 0, 0, 2.8142e+08, 0, 0},
				{0, 0, 0, 0, 2.652e+07, 0},
				{0, 0, 0, 0, 0, 5.33e+06}},
			Mass_matrix: [][]float64{
				{68.772, 0, 0, 0, 0, 0},
				{0, 68.772, 0, -0, 0, 0},
				{0, 0, 68.772, 0, 0, 0},
				{0, -0, 0, 18.53, 0, 0},
				{0, 0, 0, 0, 1.49, 0},
				{0, 0, 0, 0, 0, 20.02},
			},
		},
		{
			Station_eta: 0.93821,
			Stiffness_matrix: [][]float64{
				{9.224e+06, 0, 0, 0, 0, 0},
				{0, 9.224e+06, 0, 0, 0, 0},
				{0, 0, 9.224e+07, 0, 0, 0},
				{0, 0, 0, 2.6171e+08, 0, 0},
				{0, 0, 0, 0, 2.384e+07, 0},
				{0, 0, 0, 0, 0, 4.94e+06}},
			Mass_matrix: [][]float64{
				{66.264, 0, 0, 0, 0, 0},
				{0, 66.264, 0, -0, 0, 0},
				{0, 0, 66.264, 0, 0, 0},
				{0, -0, 0, 17.11, 0, 0},
				{0, 0, 0, 0, 1.34, 0},
				{0, 0, 0, 0, 0, 18.45},
			},
		},
		{
			Station_eta: 0.94636,
			Stiffness_matrix: [][]float64{
				{6.323e+06, 0, 0, 0, 0, 0},
				{0, 6.323e+06, 0, 0, 0, 0},
				{0, 0, 6.323e+07, 0, 0, 0},
				{0, 0, 0, 1.5881e+08, 0, 0},
				{0, 0, 0, 0, 1.963e+07, 0},
				{0, 0, 0, 0, 0, 4.24e+06}},
			Mass_matrix: [][]float64{
				{59.34, 0, 0, 0, 0, 0},
				{0, 59.34, 0, -0, 0, 0},
				{0, 0, 59.34, 0, 0, 0},
				{0, -0, 0, 11.55, 0, 0},
				{0, 0, 0, 0, 1.1, 0},
				{0, 0, 0, 0, 0, 12.65},
			},
		},
		{
			Station_eta: 0.95447,
			Stiffness_matrix: [][]float64{
				{5.332e+06, 0, 0, 0, 0, 0},
				{0, 5.332e+06, 0, 0, 0, 0},
				{0, 0, 5.332e+07, 0, 0, 0},
				{0, 0, 0, 1.3788e+08, 0, 0},
				{0, 0, 0, 0, 1.6e+07, 0},
				{0, 0, 0, 0, 0, 3.66e+06}},
			Mass_matrix: [][]float64{
				{55.914, 0, 0, 0, 0, 0},
				{0, 55.914, 0, -0, 0, 0},
				{0, 0, 55.914, 0, 0, 0},
				{0, -0, 0, 9.77, 0, 0},
				{0, 0, 0, 0, 0.89, 0},
				{0, 0, 0, 0, 0, 10.66},
			},
		},
		{
			Station_eta: 0.9626,
			Stiffness_matrix: [][]float64{
				{4.453e+06, 0, 0, 0, 0, 0},
				{0, 4.453e+06, 0, 0, 0, 0},
				{0, 0, 4.453e+07, 0, 0, 0},
				{0, 0, 0, 1.1879e+08, 0, 0},
				{0, 0, 0, 0, 1.283e+07, 0},
				{0, 0, 0, 0, 0, 3.13e+06}},
			Mass_matrix: [][]float64{
				{52.484, 0, 0, 0, 0, 0},
				{0, 52.484, 0, -0, 0, 0},
				{0, 0, 52.484, 0, 0, 0},
				{0, -0, 0, 8.19, 0, 0},
				{0, 0, 0, 0, 0.71, 0},
				{0, 0, 0, 0, 0, 8.9},
			},
		},
		{
			Station_eta: 0.97073,
			Stiffness_matrix: [][]float64{
				{3.69e+06, 0, 0, 0, 0, 0},
				{0, 3.69e+06, 0, 0, 0, 0},
				{0, 0, 3.69e+07, 0, 0, 0},
				{0, 0, 0, 1.0163e+08, 0, 0},
				{0, 0, 0, 0, 1.008e+07, 0},
				{0, 0, 0, 0, 0, 2.64e+06}},
			Mass_matrix: [][]float64{
				{49.114, 0, 0, 0, 0, 0},
				{0, 49.114, 0, -0, 0, 0},
				{0, 0, 49.114, 0, 0, 0},
				{0, -0, 0, 6.82, 0, 0},
				{0, 0, 0, 0, 0.56, 0},
				{0, 0, 0, 0, 0, 7.38},
			},
		},
		{
			Station_eta: 0.97886,
			Stiffness_matrix: [][]float64{
				{2.992e+06, 0, 0, 0, 0, 0},
				{0, 2.992e+06, 0, 0, 0, 0},
				{0, 0, 2.992e+07, 0, 0, 0},
				{0, 0, 0, 8.507e+07, 0, 0},
				{0, 0, 0, 0, 7.55e+06, 0},
				{0, 0, 0, 0, 0, 2.17e+06}},
			Mass_matrix: [][]float64{
				{45.818, 0, 0, 0, 0, 0},
				{0, 45.818, 0, -0, 0, 0},
				{0, 0, 45.818, 0, 0, 0},
				{0, -0, 0, 5.57, 0, 0},
				{0, 0, 0, 0, 0.42, 0},
				{0, 0, 0, 0, 0, 5.99},
			},
		},
		{
			Station_eta: 0.98699,
			Stiffness_matrix: [][]float64{
				{2.131e+06, 0, 0, 0, 0, 0},
				{0, 2.131e+06, 0, 0, 0, 0},
				{0, 0, 2.131e+07, 0, 0, 0},
				{0, 0, 0, 6.426e+07, 0, 0},
				{0, 0, 0, 0, 4.6e+06, 0},
				{0, 0, 0, 0, 0, 1.58e+06}},
			Mass_matrix: [][]float64{
				{41.669, 0, 0, 0, 0, 0},
				{0, 41.669, 0, -0, 0, 0},
				{0, 0, 41.669, 0, 0, 0},
				{0, -0, 0, 4.01, 0, 0},
				{0, 0, 0, 0, 0.25, 0},
				{0, 0, 0, 0, 0, 4.26},
			},
		},
		{
			Station_eta: 0.99512,
			Stiffness_matrix: [][]float64{
				{485000, 0, 0, 0, 0, 0},
				{0, 485000, 0, 0, 0, 0},
				{0, 0, 4.85e+06, 0, 0, 0},
				{0, 0, 0, 6.61e+06, 0, 0},
				{0, 0, 0, 0, 250000, 0},
				{0, 0, 0, 0, 0, 250000}},
			Mass_matrix: [][]float64{
				{11.453, 0, 0, 0, 0, 0},
				{0, 11.453, 0, -0, 0, 0},
				{0, 0, 11.453, 0, 0, 0},
				{0, -0, 0, 0.94, 0, 0},
				{0, 0, 0, 0, 0.04, 0},
				{0, 0, 0, 0, 0, 0.98},
			},
		},
		{
			Station_eta: 1,
			Stiffness_matrix: [][]float64{
				{353000, 0, 0, 0, 0, 0},
				{0, 353000, 0, 0, 0, 0},
				{0, 0, 3.53e+06, 0, 0, 0},
				{0, 0, 0, 5.01e+06, 0, 0},
				{0, 0, 0, 0, 170000, 0},
				{0, 0, 0, 0, 0, 190000}},
			Mass_matrix: [][]float64{
				{10.319, 0, 0, 0, 0, 0},
				{0, 10.319, 0, -0, 0, 0},
				{0, 0, 10.319, 0, 0, 0},
				{0, -0, 0, 0.68, 0, 0},
				{0, 0, 0, 0, 0.02, 0},
				{0, 0, 0, 0, 0, 0.7},
			},
		},
	},
	Defaults: map[string]struct{}{},
}