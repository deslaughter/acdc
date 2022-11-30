package input

import (
	"fmt"
	"io"
	"reflect"
)

var AD15AirfoilInfoSchema = NewSchema("AD15AirfoilInfo", []SchemaEntry{
	// {Heading: "AeroDyn 15 AirfoilInfo Input File"},
	{Keyword: "InterpOrd", Type: Int, Desc: `Interpolation order to use for quasi-steady table lookup {1=linear; 3=cubic spline; "default"} [default=3]`, CanBeDefault: true},
	{Keyword: "NonDimArea", Type: Float, Desc: `The non-dimensional area of the airfoil (set to 1.0 if unsure or unneeded)`, Unit: "area/chord^2"},
	{Keyword: "NumCoords", Type: Int, Desc: `The number of coordinates in the airfoil shape file.  Set to zero if coordinates not included.`},
	{Keyword: "BL_file", Type: String, Desc: `The file name including the boundary layer characteristics of the profile. Ignored if the aeroacoustic module is not called.`},
	{Keyword: "NumTabs", Type: Int, Desc: `Number of airfoil tables in this file.  Each table must have lines for Re and Ctrl.`},
	// {Heading: "Data for Table 1"},
	{Keyword: "Re", Type: Float, Desc: `Reynolds numbers in millions`},
	{Keyword: "Ctrl", Type: Int, Desc: `Control setting (must be 0 for current AirfoilInfo)`},
	{Keyword: "InclUAdata", Type: Bool, Desc: `Is unsteady aerodynamics data included in this table? If TRUE, then include 30 UA coeffs below this line`},
	{Keyword: "alpha0", Type: Float, Desc: `0-lift angle of attack, depends on airfoil.`, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "alpha1", Type: Float, Desc: `Angle of attack at f=0.7, (approximately the stall angle) for AOA>alpha0.`, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "alpha2", Type: Float, Desc: `Angle of attack at f=0.7, (approximately the stall angle) for AOA<alpha0.`, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "eta_e", Type: Float, Desc: `Recovery factor in the range [0.85 - 0.95] used only for UAMOD=1, it is set to 1 in the code when flookup=True.`, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "C_nalpha", Type: Float, Desc: `Slope of the 2D normal force coefficient curve in the linear region of the polar.`, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "T_f0", Type: Float, Desc: `Intial value of the time constant associated with Df in the expression of Df and f'. Default value = 3.`, CanBeDefault: true, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "T_V0", Type: Float, Desc: `Intial value of the time constant associated with the vortex lift decay process; it is used in the expression of Cvn. It depends on Re,M, and airfoil class. Default value= 6.`, CanBeDefault: true, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "T_p", Type: Float, Desc: `Boundary-layer,leading edge pressure gradient time constant in the expression of Dp. It should be tuned based on airfoil experimental data. Default =1.7.`, CanBeDefault: true, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "T_VL", Type: Float, Desc: `Intial value of the time constant associated with the vortex advection process; it represents the non-dimensional time in semi-chords, needed for a vortex to travel from LE to trailing edge (TE); it is used in the expression of Cvn. It depends on Re, M (weakly), and airfoil. Value's range = [6; 13]; default value= 11.`, CanBeDefault: true, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "b1", Type: Float, Desc: `Constant in the expression of phi_alpha^c and phi_q^c;  from experimental results, it was set to 0.14. This value is relatively insensitive for thin airfoils, but may be different for turbine airfoils.`, CanBeDefault: true, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "b2", Type: Float, Desc: `Constant in the expression of phi_alpha^c and phi_q^c;  from experimental results, it was set to 0.53. This value is relatively insensitive for thin airfoils, but may be different for turbine airfoils.`, CanBeDefault: true, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "b5", Type: Float, Desc: `Constant in the expression of K'''_q,Cm_q^nc, and k_m,q; from  experimental results, it was set to 5.`, CanBeDefault: true, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "A1", Type: Float, Desc: `Constant in the expression of phi_alpha^c and phi_q^c;  from experimental results, it was set to 0.3. This value is relatively insensitive for thin airfoils, but may be different for turbine airfoils.`, CanBeDefault: true, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "A2", Type: Float, Desc: `Constant in the expression of phi_alpha^c and phi_q^c;  from experimental results, it was set to 0.7. This value is relatively insensitive for thin airfoils, but may be different for turbine airfoils.`, CanBeDefault: true, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "A5", Type: Float, Desc: `Constant in the expression of K'''_q,Cm_q^nc, and k_m,q; from  experimental results, it was set to 1.`, CanBeDefault: true, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "S1", Type: Float, Desc: `Constant in the f curve bestfit for alpha0<=AOA<=alpha1;by definition it depends on the airfoil. Ignored if UAMod<>1.`, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "S2", Type: Float, Desc: `Constant in the f curve bestfit for         AOA>alpha1;by definition it depends on the airfoil. Ignored if UAMod<>1.`, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "S3", Type: Float, Desc: `Constant in the f curve bestfit for alpha2<=AOA<alpha0;by definition it depends on the airfoil. Ignored if UAMod<>1.`, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "S4", Type: Float, Desc: `Constant in the f curve bestfit for         AOA<alpha2;by definition it depends on the airfoil. Ignored if UAMod<>1.`, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "Cn1", Type: Float, Desc: `Critical value of C0n at leading edge separation. It should be extracted from airfoil data at a given Mach and Reynolds number. It can be calculated from the static value of Cn at either the break in the pitching moment or the loss of chord force at the onset of stall. It is close to the condition of maximum lift of the airfoil at low Mach numbers.`, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "Cn2", Type: Float, Desc: `As Cn1 for negative AOAs.`, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "St_sh", Type: Float, Desc: `Strouhal's shedding frequency constant; default =0.19.`, CanBeDefault: true, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "Cd0", Type: Float, Desc: `2D drag coefficient value at 0-lift.`, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "Cm0", Type: Float, Desc: `2D pitching moment coeffcient about 1/4-chord location, at 0-lift, positive if nose up.`, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "k0", Type: Float, Desc: `Constant in the \hat(x)_cp curve best-fit; = (\hat(x)_AC-0.25). Ignored if UAMod<>1.`, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "k1", Type: Float, Desc: `Constant in the \hat(x)_cp curve best-fit. Ignored if UAMod<>1.`, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "k2", Type: Float, Desc: `Constant in the \hat(x)_cp curve best-fit. Ignored if UAMod<>1.`, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "k3", Type: Float, Desc: `Constant in the \hat(x)_cp curve best-fit. Ignored if UAMod<>1.`, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "k1_hat", Type: Float, Desc: `Constant in the expression of Cc due to leading edge vortex effects. Ignored if UAMod<>1.`, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "x_cp_bar", Type: Float, Desc: `Constant in the expression of \hat(x)_cp^v. Default value =0.2. Ignored if UAMod<>1.`, CanBeDefault: true, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "UACutout", Type: Float, Desc: `Angle of attack above which unsteady aerodynamics are disabled (deg). [Specifying the string "Default" sets UACutout to 45 degrees]`, CanBeDefault: true, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "filtCutOff", Type: Float, Desc: `Cut-off frequency (-3 dB corner frequency) for low-pass filtering the AoA input to UA, as well as the 1st and 2nd derivatives (Hz) [default = 20]`, CanBeDefault: true, SkipIf: []Condition{{"InclUAdata", "==", false}}},
	{Keyword: "NumAlf", Type: Int, Desc: `Number of data lines in the following CoeffData table`},
	{Keyword: "CoeffData", Dims: 1, Parse: parseAD15AirfoilCoeffs, Format: formatAD15AirfoilCoeffs,
		Table: &Table{
			Columns: []TableColumn{
				{Keyword: "Alpha", Type: Float, Unit: "deg"},
				{Keyword: "Cl", Type: Float, Unit: "-"},
				{Keyword: "Cd", Type: Float, Unit: "-"},
				{Keyword: "Cm", Type: Float, Unit: "-"},
			},
		},
	},
})

func parseAD15AirfoilCoeffs(s any, v reflect.Value, lines []string) ([]string, error) {
	ai := s.(*AD15AirfoilInfo)
	ai.CoeffData = make([]AD15AirfoilInfoCoeffData, ai.NumAlf)
	line := ""
	lines = lines[2:]
	for i := range ai.CoeffData {
		line, lines = lines[0], lines[1:]
		num, err := fmt.Sscan(line, &ai.CoeffData[i].Alpha, &ai.CoeffData[i].Cl, &ai.CoeffData[i].Cd, &ai.CoeffData[i].Cm)
		if num < 4 || err != nil {
			return nil, err
		}
	}
	return lines, nil
}

func formatAD15AirfoilCoeffs(w io.Writer, s, field any, entry SchemaEntry) error {
	fmt.Fprintf(w, "! %14s %14s %14s %14s\n", "Alpha", "Cl", "Cd", "Cm")
	fmt.Fprintf(w, "! %14s %14s %14s %14s\n", "(deg)", "(-)", "(-)", "(-)")
	ai := s.(*AD15AirfoilInfo)
	for _, cd := range ai.CoeffData {
		fmt.Fprintf(w, "  %14g %14g %14g %14g\n",
			cd.Alpha, cd.Cl, cd.Cd, cd.Cm)
	}
	return nil

}
