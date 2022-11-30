package input

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

var AeroDyn15Schema = NewSchema("AeroDyn15", []SchemaEntry{
	{Heading: "AeroDyn 15 Input File"},
	{Keyword: "Title", Type: String, Format: formatTitle, Parse: parseTitle},
	{Heading: "General Options"},
	{Keyword: "Echo", Type: Bool, Desc: `Echo the input to "<rootname>.AD.ech"?`, Unit: "flag"},
	{Keyword: "DTAero", Type: Float, CanBeDefault: true, Desc: `Time interval for aerodynamic calculations {or "default"}`, Unit: "sec"},
	{Keyword: "WakeMod", Type: Int, Desc: `Type of wake/induction model [WakeMod cannot be 2 or 3 when linearizing]`, Unit: "switch", Options: []Option{{0, "none"}, {1, "BEMT"}, {2, "DBEMT"}, {3, "OLAF"}}},
	{Keyword: "AFAeroMod", Type: Int, Desc: `Type of blade airfoil aerodynamics model[AFAeroMod must be 1 when linearizing]`, Unit: "switch", Options: []Option{{1, "steady model"}, {2, "Beddoes-Leishman unsteady model"}}},
	{Keyword: "TwrPotent", Type: Int, Desc: `Type tower influence on wind based on potential flow around the tower`, Unit: "switch", Options: []Option{{0, "none"}, {1, "baseline potential flow"}, {2, "potential flow with Bak correction"}}},
	{Keyword: "TwrShadow", Type: Int, Desc: `Calculate tower influence on wind based on downstream tower shadow`, Unit: "switch", Options: []Option{{0, "none"}, {1, "Powles model"}, {2, "Eames model"}}},
	{Keyword: "TwrAero", Type: Bool, Desc: `Calculate tower aerodynamic loads?`, Unit: "flag"},
	{Keyword: "FrozenWake", Type: Bool, Desc: `Assume frozen wake during linearization? [used only when WakeMod=1 and when linearizing]`, Unit: "flag"},
	{Keyword: "CavitCheck", Type: Bool, Desc: `Perform cavitation check? [AFAeroMod must be 1 when CavitCheck=true]`, Unit: "flag"},
	{Keyword: "CompAA", Type: Bool, Desc: `Flag to compute AeroAcoustics calculation [used only when WakeMod = 1 or 2]`},
	{Keyword: "AA_InputFile", Type: String, Desc: `AeroAcoustics input file [used only when CompAA=true]`},
	{Heading: "Environmental Conditions"},
	{Keyword: "AirDens", Type: Float, CanBeDefault: true, Desc: `Air density`, Unit: "kg/m^3"},
	{Keyword: "KinVisc", Type: Float, CanBeDefault: true, Desc: `Kinematic viscosity of working fluid`, Unit: "m^2/s"},
	{Keyword: "SpdSound", Type: Float, CanBeDefault: true, Desc: `Speed of sound in working fluid`, Unit: "m/s"},
	{Keyword: "Patm", Type: Float, CanBeDefault: true, Desc: `Atmospheric pressure [used only when CavitCheck=True]`, Unit: "Pa", Show: []Condition{{"CavitCheck", "==", true}}},
	{Keyword: "Pvap", Type: Float, CanBeDefault: true, Desc: `Vapour pressure of working fluid [used only when CavitCheck=True]`, Unit: "Pa", Show: []Condition{{"CavitCheck", "==", true}}},
	{Heading: "Blade-Element/Momentum Theory Options [unused when WakeMod=0 or 3]"},
	{Keyword: "SkewMod", Type: Int, Desc: `Type of skewed-wake correction model [unused when WakeMod=0 or 3]`, Unit: "switch", Options: []Option{{1, "uncoupled"}, {2, "Pitt/Peters"}, {3, "coupled"}}},
	{Keyword: "SkewModFactor", Type: Float, CanBeDefault: true, Desc: `Constant used in Pitt/Peters skewed wake model {or "default" is 15/32*pi} (-) [used only when SkewMod=2; unused when WakeMod=0 or 3]`},
	{Keyword: "TipLoss", Type: Bool, Desc: `Use the Prandtl tip-loss model? [unused when WakeMod=0 or 3]`, Unit: "flag"},
	{Keyword: "HubLoss", Type: Bool, Desc: `Use the Prandtl hub-loss model? [unused when WakeMod=0 or 3]`, Unit: "flag"},
	{Keyword: "TanInd", Type: Bool, Desc: `Include tangential induction in BEMT calculations? [unused when WakeMod=0 or 3]`, Unit: "flag"},
	{Keyword: "AIDrag", Type: Bool, Desc: `Include the drag term in the axial-induction calculation? [unused when WakeMod=0 or 3]`, Unit: "flag"},
	{Keyword: "TIDrag", Type: Bool, Desc: `Include the drag term in the tangential-induction calculation? [unused when WakeMod=0,3 or TanInd=FALSE]`, Unit: "flag"},
	{Keyword: "IndToler", Type: Float, CanBeDefault: true, Desc: `Convergence tolerance for BEMT nonlinear solve residual equation {or "default"} (-) [unused when WakeMod=0 or 3]`},
	{Keyword: "MaxIter", Type: Int, Desc: `Maximum number of iteration steps (-) [unused when WakeMod=0]`},
	{Heading: "Dynamic Blade-Element/Momentum Theory Options [used only when WakeMod=2]"},
	{Keyword: "DBEMT_Mod", Type: Int, Desc: `Type of dynamic BEMT (DBEMT) model [used only when WakeMod=2]`, Unit: "switch", Options: []Option{{1, "constant tau1"}, {2, "time-dependent tau1"}}},
	{Keyword: "tau1_const", Type: Float, Desc: `Time constant for DBEMT (s) [used only when WakeMod=2 and DBEMT_Mod=1]`},
	{Heading: "OLAF -- cOnvecting LAgrangian Filaments (Free Vortex Wake) Theory Options [used only when WakeMod=3]"},
	{Keyword: "OLAFInputFileName", Type: String, Desc: `Input file for OLAF [used only when WakeMod=3]`},
	{Heading: "Beddoes-Leishman Unsteady Airfoil Aerodynamics Options [used only when AFAeroMod=2]"},
	{Keyword: "UAMod", Type: Int, Desc: `Unsteady Aero Model Switch (switch) {1=Baseline model (Original), 2=Gonzalez's variant (changes in Cn,Cc,Cm), 3=Minnema/Pierce variant (changes in Cc and Cm)} [used only when AFAeroMod=2]`},
	{Keyword: "FLookup", Type: Bool, Desc: `Flag to indicate whether a lookup for f' will be calculated (TRUE) or whether best-fit exponential equations will be used (FALSE); if FALSE S1-S4 must be provided in airfoil input files (flag) [used only when AFAeroMod=2]`},
	{Heading: "Airfoil Information"},
	{Keyword: "AFTabMod", Type: Int, Desc: `Interpolation method for multiple airfoil tables {1=1D interpolation on AoA (first table only); 2=2D interpolation on AoA and Re; 3=2D interpolation on AoA and UserProp} (-)`},
	{Keyword: "InCol_Alfa", Type: Int, Desc: `The column in the airfoil tables that contains the angle of attack (-)`},
	{Keyword: "InCol_Cl", Type: Int, Desc: `The column in the airfoil tables that contains the lift coefficient (-)`},
	{Keyword: "InCol_Cd", Type: Int, Desc: `The column in the airfoil tables that contains the drag coefficient (-)`},
	{Keyword: "InCol_Cm", Type: Int, Desc: `The column in the airfoil tables that contains the pitching-moment coefficient; use zero if there is no Cm column (-)`},
	{Keyword: "InCol_Cpmin", Type: Int, Desc: `The column in the airfoil tables that contains the Cpmin coefficient; use zero if there is no Cpmin column (-)`},
	{Keyword: "NumAFfiles", Type: Int, Desc: `Number of airfoil files used (-)`},
	{Keyword: "AFNames", Type: String, Dims: 1, Parse: parseAFNames, Format: formatAFNames, Desc: `Airfoil file names (NumAFfiles lines) (quoted strings)`},
	{Heading: "Rotor/Blade Properties"},
	{Keyword: "UseBlCm", Type: Bool, Desc: `Include aerodynamic pitching moment in calculations?  (flag)`},
	{Keyword: "ADBlFile(1)", Type: String, Desc: `Name of file containing distributed aerodynamic properties for Blade #1 (-)`},
	{Keyword: "ADBlFile(2)", Type: String, Desc: `Name of file containing distributed aerodynamic properties for Blade #2 (-) [unused if NumBl < 2]`},
	{Keyword: "ADBlFile(3)", Type: String, Desc: `Name of file containing distributed aerodynamic properties for Blade #3 (-) [unused if NumBl < 3]`},
	{Heading: "Tower Influence and Aerodynamics [used only when TwrPotent/=0, TwrShadow/=0, or TwrAero=True]"},
	{Keyword: "NumTwrNds", Type: Int, Desc: `Number of tower nodes used in the analysis  (-) [used only when TwrPotent/=0, TwrShadow/=0, or TwrAero=True]`},
	{Keyword: "TwrNodes", Dims: 1,
		Table: &Table{
			Columns: []TableColumn{
				{Keyword: "TwrElev", Type: Float, Desc: `(m)`},
				{Keyword: "TwrDiam", Type: Float, Desc: `(m)`},
				{Keyword: "TwrCd", Type: Float, Desc: `(-)`},
				{Keyword: "TwrTI", Type: Float, Desc: `(used only with TwrShadow=2) (-)`},
			},
		},
		Parse: parseTwrNodes, Format: formatTwrNodes,
	},
	{Heading: "Outputs"},
	{Keyword: "SumPrint", Type: Bool, Desc: `Generate a summary file listing input options and interpolated properties to "<rootname>.AD.sum"?  (flag)`},
	{Keyword: "NBlOuts", Type: Int, Desc: `Number of blade node outputs [0 - 9] (-)`},
	{Keyword: "BlOutNd", Dims: 1, Type: Int, Desc: `Blade nodes whose values will be output  (-)`},
	{Keyword: "NTwOuts", Type: Int, Desc: `Number of tower node outputs [0 - 9]  (-)`},
	{Keyword: "TwOutNd", Dims: 1, Type: Int, Desc: `Tower nodes whose values will be output  (-)`},
	{Keyword: "OutList", Type: String, Dims: 1, Parse: parseOutList, Format: formatOutList, Desc: `The next line(s) contains a list of output parameters.  See OutListParameters.xlsx for a listing of available output channels, (-)`},
})

func formatAFNames(w io.Writer, s, field any, entry SchemaEntry) error {
	vals := field.([]string)
	fmt.Fprintf(w, "\"%-12s\"   %-14s - %s\n", vals[0], entry.Keyword, entry.Desc)
	for _, v := range vals[1:] {
		fmt.Fprintf(w, "\"%-12s\"\n", v)
	}
	return nil
}

func formatTwrNodes(w io.Writer, s, field any, entry SchemaEntry) error {
	fmt.Fprintf(w, "%12s %12s %12s %12s\n", "TwrElev", "TwrDiam", "TwrCd", "TwrTI")
	fmt.Fprintf(w, "%12s %12s %12s %12s\n", "(m)", "(m)", "(-)", "(-)")
	ad := s.(*AeroDyn15)
	for _, ta := range ad.TwrNodes {
		fmt.Fprintf(w, "%12g %12g %12g %12g\n", ta.TwrElev, ta.TwrDiam, ta.TwrCd, ta.TwrTI)
	}
	return nil
}

func parseAFNames(s any, v reflect.Value, lines []string) ([]string, error) {
	line, lines, err := findKeywordLine("AFNames", lines)
	if err != nil {
		return nil, err
	}
	names := make([]string, s.(*AeroDyn15).NumAFfiles)
	if len(lines) < len(names) {
		return lines, fmt.Errorf("insufficient lines for AFNames (%d)", len(lines))
	}
	for i := range names {
		line = strings.Trim(line, "\"\t ")
		if j := strings.Index(line, `"`); j > 0 {
			names[i] = line[:j]
		} else {
			names[i] = strings.Fields(line)[0]
		}
		line, lines = lines[0], lines[1:]
	}
	v.Set(reflect.ValueOf(names))
	return lines, nil
}

func parseTwrNodes(s any, v reflect.Value, lines []string) ([]string, error) {
	ad := s.(*AeroDyn15)
	ad.TwrNodes = make([]AeroDyn15TwrNodes, ad.NumTwrNds)
	line := ""
	lines = lines[2:]
	for i := range ad.TwrNodes {
		line, lines = lines[0], lines[1:]
		num, err := fmt.Sscan(line, &ad.TwrNodes[i].TwrElev,
			&ad.TwrNodes[i].TwrDiam, &ad.TwrNodes[i].TwrCd,
			&ad.TwrNodes[i].TwrTI)
		if num < 3 {
			return nil, err
		}
	}
	return lines, nil
}
