package input

import (
	"fmt"
	"io"
	"reflect"
)

var BeamDynSchema = NewSchema("BeamDyn", []SchemaEntry{
	{Heading: "Beamdyn Input File"},
	{Keyword: "Title", Type: String, Format: formatTitle, Parse: parseTitle},
	{Heading: "Simulation Control"},
	{Keyword: "Echo", Type: Bool, Desc: `Echo input data to "<RootName>.ech"? (flag)`},
	{Keyword: "QuasiStaticInit", Type: Bool, Desc: `Use quasi-static pre-conditioning with centripetal accelerations in initialization? (flag) [dynamic solve only]`},
	{Keyword: "rhoinf", Type: Float, Desc: `Numerical damping parameter for generalized-alpha integrator`},
	{Keyword: "quadrature", Type: Int, Desc: `Quadrature method: 1=Gaussian; 2=Trapezoidal (switch)`},
	{Keyword: "refine", Type: Int, CanBeDefault: true, Desc: `Refinement factor for trapezoidal quadrature (-) [DEFAULT = 1; used only when quadrature=2]`},
	{Keyword: "n_fact", Type: Int, CanBeDefault: true, Desc: `Factorization frequency for the Jacobian in N-R iteration(-) [DEFAULT = 5]`},
	{Keyword: "DTBeam", Type: Float, CanBeDefault: true, Desc: `Time step size (s)`},
	{Keyword: "load_retries", Type: Int, CanBeDefault: true, Desc: `Number of factored load retries before quitting the simulation [DEFAULT = 20]`},
	{Keyword: "NRMax", Type: Int, CanBeDefault: true, Desc: `Max number of iterations in Newton-Raphson algorithm (-) [DEFAULT = 10]`},
	{Keyword: "stop_tol", Type: Float, CanBeDefault: true, Desc: `Tolerance for stopping criterion (-) [DEFAULT = 1E-5]`},
	{Keyword: "tngt_stf_fd", Type: Bool, CanBeDefault: true, Desc: `Use finite differenced tangent stiffness matrix? (flag)`},
	{Keyword: "tngt_stf_comp", Type: Bool, CanBeDefault: true, Desc: `Compare analytical finite differenced tangent stiffness matrix? (flag)`},
	{Keyword: "tngt_stf_pert", Type: Float, CanBeDefault: true, Desc: `Perturbation size for finite differencing (-) [DEFAULT = 1E-6]`},
	{Keyword: "tngt_stf_difftol", Type: Float, CanBeDefault: true, Desc: `Maximum allowable relative difference between analytical and fd tangent stiffness (-); [DEFAULT = 0.1]`},
	{Keyword: "RotStates", Type: Bool, Desc: `Orient states in the rotating frame during linearization? (flag) [used only when linearizing] `},
	{Heading: "Geometry Parameter"},
	{Keyword: "member_total", Type: Int, Desc: `Total number of members (-)`},
	{Keyword: "kp_total", Type: Int, Desc: `Total number of key points (-) [must be at least 3]`},
	{Keyword: "Members", Dims: 2,
		Table: &Table{
			Columns: []TableColumn{
				{Keyword: "kp_xr", Type: Float, Unit: "m", Desc: ``},
				{Keyword: "kp_yr", Type: Float, Unit: "m", Desc: ``},
				{Keyword: "kp_zr", Type: Float, Unit: "m", Desc: ``},
				{Keyword: "initial_twist", Type: Float, Unit: "deg", Desc: ``},
			},
		},
		Parse: parseBeamDynMembers, Format: formatBeamDynMembers,
	},
	{Heading: "Mesh Parameter"},
	{Keyword: "order_elem", Type: Int, Desc: `Order of interpolation (basis) function (-)`},
	{Heading: "Material Parameter"},
	{Keyword: "BldFile", Type: String, Desc: `Name of file containing properties for blade (quoted string)`},
	{Heading: "Pitch Actuator Parameters"},
	{Keyword: "UsePitchAct", Type: Bool, Desc: `Whether a pitch actuator should be used (flag)`},
	{Keyword: "PitchJ", Type: Float, Desc: `Pitch actuator inertia (kg-m^2) [used only when UsePitchAct is true]`},
	{Keyword: "PitchK", Type: Float, Desc: `Pitch actuator stiffness (kg-m^2/s^2) [used only when UsePitchAct is true]`},
	{Keyword: "PitchC", Type: Float, Desc: `Pitch actuator damping (kg-m^2/s) [used only when UsePitchAct is true]`},
	{Heading: "Outputs"},
	{Keyword: "SumPrint", Type: Bool, Desc: `Print summary data to "<RootName>.sum" (flag)`},
	{Keyword: "OutFmt", Type: String, Desc: `Format used for text tabular output, excluding the time channel.`},
	{Keyword: "NNodeOuts", Type: Int, Desc: `Number of nodes to output to file [0 - 9] (-)`},
	{Keyword: "OutNd", Dims: 1, Type: Int, Desc: `Nodes whose values will be output  (-)`},
	{Keyword: "OutList", Type: String, Dims: 1, Format: formatOutList, Parse: parseOutList, Desc: `The next line(s) contains a list of output parameters.  See OutListParameters.xlsx for a listing of available output channels, (-)`},
})

func formatBeamDynMembers(w io.Writer, s, field any, entry SchemaEntry) error {
	bd := s.(*BeamDyn)
	for i, m := range bd.Members {
		fmt.Fprintf(w, "%8d %8d %12s - Member number; Number of key points in this member\n", i+1, len(m), "")
		fmt.Fprintf(w, "%13s %13s %13s %13s\n",
			"kp_xr", "kp_yr", "kp_zr", "initial_twist")
		fmt.Fprintf(w, "%13s %13s %13s %13s\n", "(m)", "(m)", "(m)", "(deg)")
		for _, kp := range m {
			fmt.Fprintf(w, "%13g %13g %13g %13g\n",
				kp.Kp_xr, kp.Kp_yr, kp.Kp_zr, kp.Initial_twist)
		}
	}
	return nil
}

func parseBeamDynMembers(s any, v reflect.Value, lines []string) ([]string, error) {
	bd := s.(*BeamDyn)
	bd.Members = make([][]BeamDynMembers, bd.Member_total)
	line := ""
	memberNum := 0
	numKeyPoints := 0
	for i := range bd.Members {
		line, lines = lines[0], lines[1:]
		num, err := fmt.Sscan(line, &memberNum, &numKeyPoints)
		if num < 2 || err != nil {
			return nil, err
		}
		lines = lines[2:]
		keyPoints := make([]BeamDynMembers, numKeyPoints)
		for j := range keyPoints {
			line, lines = lines[0], lines[1:]
			num, err := fmt.Sscan(line, &keyPoints[j].Kp_xr, &keyPoints[j].Kp_yr,
				&keyPoints[j].Kp_zr, &keyPoints[j].Initial_twist)
			if num < 4 || err != nil {
				return nil, err
			}
		}
		bd.Members[i] = keyPoints
	}
	return lines, nil
}
