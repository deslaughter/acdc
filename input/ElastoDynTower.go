package input

import (
	"fmt"
	"io"
	"reflect"
)

var ElastoDynTowerSchema = NewSchema("ElastoDynTower", []SchemaEntry{
	{Heading: "ElastoDyn Tower Input File"},
	{Keyword: "Title", Type: String, Format: formatTitle, Parse: parseTitle},
	{Heading: "Tower Parameters"},
	{Keyword: "NTwInpSt", Type: Int, Desc: `Number of input stations to specify tower geometry`},
	{Keyword: "TwrFADmp(1)", Type: Float, Desc: `Tower 1st fore-aft mode structural damping ratio`, Unit: "%"},
	{Keyword: "TwrFADmp(2)", Type: Float, Desc: `Tower 2nd fore-aft mode structural damping ratio`, Unit: "%"},
	{Keyword: "TwrSSDmp(1)", Type: Float, Desc: `Tower 1st side-to-side mode structural damping ratio`, Unit: "%"},
	{Keyword: "TwrSSDmp(2)", Type: Float, Desc: `Tower 2nd side-to-side mode structural damping ratio`, Unit: "%"},
	{Heading: "Tower Adjustment Factors"},
	{Keyword: "FAStTunr(1)", Type: Float, Desc: `Tower fore-aft modal stiffness tuner, 1st mode`, Unit: "-"},
	{Keyword: "FAStTunr(2)", Type: Float, Desc: `Tower fore-aft modal stiffness tuner, 2nd mode`, Unit: "-"},
	{Keyword: "SSStTunr(1)", Type: Float, Desc: `Tower side-to-side stiffness tuner, 1st mode`, Unit: "-"},
	{Keyword: "SSStTunr(2)", Type: Float, Desc: `Tower side-to-side stiffness tuner, 2nd mode`, Unit: "-"},
	{Keyword: "AdjTwMa", Type: Float, Desc: `Factor to adjust tower mass density`, Unit: "-"},
	{Keyword: "AdjFASt", Type: Float, Desc: `Factor to adjust tower fore-aft stiffness`, Unit: "-"},
	{Keyword: "AdjSSSt", Type: Float, Desc: `Factor to adjust tower side-to-side stiffness`, Unit: "-"},
	{Heading: "Distributed Tower Properties"},
	{Keyword: "TwInpSt", Dims: 1, Format: formatElastoDynTowerTwInpSt, Parse: parseElastoDynTowerTwInpSt,
		Table: &Table{
			Columns: []TableColumn{
				{Keyword: "HtFract", Type: Float, Desc: `Fractional height of the flexible portion of tower for a given input station`, Unit: "-"},
				{Keyword: "TMassDen", Type: Float, Desc: `Tower mass density for a given input station`, Unit: "kg/m"},
				{Keyword: "TwFAStif", Type: Float, Desc: `Tower fore-aft stiffness for a given input station`, Unit: "Nm^2"},
				{Keyword: "TwSSStif", Type: Float, Desc: `Tower side-to-side stiffness for a given input station`, Unit: "Nm^2"},
			},
		},
	},
	{Heading: "Tower Fore-Aft Mode Shapes"},
	{Keyword: "TwFAM1Sh(2)", Type: Float, Desc: `Mode 1, coefficient of x^2 term`},
	{Keyword: "TwFAM1Sh(3)", Type: Float, Desc: `      , coefficient of x^3 term`},
	{Keyword: "TwFAM1Sh(4)", Type: Float, Desc: `      , coefficient of x^4 term`},
	{Keyword: "TwFAM1Sh(5)", Type: Float, Desc: `      , coefficient of x^5 term`},
	{Keyword: "TwFAM1Sh(6)", Type: Float, Desc: `      , coefficient of x^6 term`},
	{Keyword: "TwFAM2Sh(2)", Type: Float, Desc: `Mode 2, coefficient of x^2 term`},
	{Keyword: "TwFAM2Sh(3)", Type: Float, Desc: `      , coefficient of x^3 term`},
	{Keyword: "TwFAM2Sh(4)", Type: Float, Desc: `      , coefficient of x^4 term`},
	{Keyword: "TwFAM2Sh(5)", Type: Float, Desc: `      , coefficient of x^5 term`},
	{Keyword: "TwFAM2Sh(6)", Type: Float, Desc: `      , coefficient of x^6 term`},
	{Heading: "Tower Side-To-Side Mode Shapes"},
	{Keyword: "TwSSM1Sh(2)", Type: Float, Desc: `Mode 1, coefficient of x^2 term`},
	{Keyword: "TwSSM1Sh(3)", Type: Float, Desc: `      , coefficient of x^3 term`},
	{Keyword: "TwSSM1Sh(4)", Type: Float, Desc: `      , coefficient of x^4 term`},
	{Keyword: "TwSSM1Sh(5)", Type: Float, Desc: `      , coefficient of x^5 term`},
	{Keyword: "TwSSM1Sh(6)", Type: Float, Desc: `      , coefficient of x^6 term`},
	{Keyword: "TwSSM2Sh(2)", Type: Float, Desc: `Mode 2, coefficient of x^2 term`},
	{Keyword: "TwSSM2Sh(3)", Type: Float, Desc: `      , coefficient of x^3 term`},
	{Keyword: "TwSSM2Sh(4)", Type: Float, Desc: `      , coefficient of x^4 term`},
	{Keyword: "TwSSM2Sh(5)", Type: Float, Desc: `      , coefficient of x^5 term`},
	{Keyword: "TwSSM2Sh(6)", Type: Float, Desc: `      , coefficient of x^6 term`},
})

func formatElastoDynTowerTwInpSt(w io.Writer, s, field any, entry SchemaEntry) error {
	fmt.Fprintf(w, "%14s %14s %14s %14s\n", "HtFract", "TMassDen", "TwFAStif", "TwSSStif")
	fmt.Fprintf(w, "%14s %14s %14s %14s\n", "(-)", "(kg/m)", "(Nm^2)", "(Nm^2)")
	edt := s.(*ElastoDynTower)
	for _, TwInpSt := range edt.TwInpSt {
		fmt.Fprintf(w, "%14g %14g %14g %14g\n",
			TwInpSt.HtFract, TwInpSt.TMassDen, TwInpSt.TwFAStif, TwInpSt.TwSSStif)
	}
	return nil
}

func parseElastoDynTowerTwInpSt(s any, v reflect.Value, lines []string) ([]string, error) {
	edt := s.(*ElastoDynTower)
	edt.TwInpSt = make([]ElastoDynTowerTwInpSt, edt.NTwInpSt)
	if len(lines) < edt.NTwInpSt+3 {
		return nil, fmt.Errorf("insufficient lines to parse TwInpSt")
	}
	lines = lines[3:]
	for i := range edt.TwInpSt {
		n, err := fmt.Sscan(lines[0], &edt.TwInpSt[i].HtFract,
			&edt.TwInpSt[i].TMassDen, &edt.TwInpSt[i].TwFAStif,
			&edt.TwInpSt[i].TwSSStif)
		if n < 4 {
			return nil, fmt.Errorf("error parsing TwInpSt: %s", err)
		}
		lines = lines[1:]
	}
	return lines, nil
}
