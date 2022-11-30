package input

import (
	"fmt"
	"io"
	"reflect"
)

var ElastoDynBladeSchema = NewSchema("ElastoDynBlade", []SchemaEntry{
	{Heading: "ElastoDyn Blade Input File"},
	{Keyword: "Title", Type: String, Format: formatTitle, Parse: parseTitle},
	{Heading: "Blade Parameters"},
	{Keyword: "NBlInpSt", Type: Int, Desc: `Number of blade input stations`, Unit: "-"},
	{Keyword: "BldFlDmp(1)", Type: Float, Desc: `Blade flap mode #1 structural damping in percent of critical`, Unit: "%"},
	{Keyword: "BldFlDmp(2)", Type: Float, Desc: `Blade flap mode #2 structural damping in percent of critical`, Unit: "%"},
	{Keyword: "BldEdDmp(1)", Type: Float, Desc: `Blade edge mode #1 structural damping in percent of critical`, Unit: "%"},
	{Heading: "Blade Adjustment Factors"},
	{Keyword: "FlStTunr(1)", Type: Float, Desc: `Blade flapwise modal stiffness tuner, 1st mode`, Unit: "-"},
	{Keyword: "FlStTunr(2)", Type: Float, Desc: `Blade flapwise modal stiffness tuner, 2nd mode`, Unit: "-"},
	{Keyword: "AdjBlMs", Type: Float, Desc: `Factor to adjust blade mass density`, Unit: "-"},
	{Keyword: "AdjFlSt", Type: Float, Desc: `Factor to adjust blade flap stiffness`, Unit: "-"},
	{Keyword: "AdjEdSt", Type: Float, Desc: `Factor to adjust blade edge stiffness`, Unit: "-"},
	{Heading: "Distributed Blade Properties"},
	{Keyword: "BlInpSt", Dims: 1, Parse: parseBlInpSt, Format: formatBlInpSt,
		Table: &Table{
			Columns: []TableColumn{
				{Keyword: "BlFract", Type: Float, Unit: "-"},
				{Keyword: "PitchAxis", Type: Float, Unit: "-"},
				{Keyword: "StrcTwst", Type: Float, Unit: "deg"},
				{Keyword: "BMassDen", Type: Float, Unit: "kg/m"},
				{Keyword: "FlpStff", Type: Float, Unit: "Nm^2"},
				{Keyword: "EdgStff", Type: Float, Unit: "Nm^2"},
			},
		},
	},
	{Heading: "Blade Mode Shapes"},
	{Keyword: "BldFl1Sh(2)", Type: Float, Desc: `Flap mode 1, coeff of x^2`},
	{Keyword: "BldFl1Sh(3)", Type: Float, Desc: `           , coeff of x^3`},
	{Keyword: "BldFl1Sh(4)", Type: Float, Desc: `           , coeff of x^4`},
	{Keyword: "BldFl1Sh(5)", Type: Float, Desc: `           , coeff of x^5`},
	{Keyword: "BldFl1Sh(6)", Type: Float, Desc: `           , coeff of x^6`},
	{Keyword: "BldFl2Sh(2)", Type: Float, Desc: `Flap mode 2, coeff of x^2`},
	{Keyword: "BldFl2Sh(3)", Type: Float, Desc: `           , coeff of x^3`},
	{Keyword: "BldFl2Sh(4)", Type: Float, Desc: `           , coeff of x^4`},
	{Keyword: "BldFl2Sh(5)", Type: Float, Desc: `           , coeff of x^5`},
	{Keyword: "BldFl2Sh(6)", Type: Float, Desc: `           , coeff of x^6`},
	{Keyword: "BldEdgSh(2)", Type: Float, Desc: `Edge mode 1, coeff of x^2`},
	{Keyword: "BldEdgSh(3)", Type: Float, Desc: `           , coeff of x^3`},
	{Keyword: "BldEdgSh(4)", Type: Float, Desc: `           , coeff of x^4`},
	{Keyword: "BldEdgSh(5)", Type: Float, Desc: `           , coeff of x^5`},
	{Keyword: "BldEdgSh(6)", Type: Float, Desc: `           , coeff of x^6`},
})

func formatBlInpSt(w io.Writer, s, field any, entry SchemaEntry) error {
	fmt.Fprintf(w, "%14s %14s %14s %14s %14s %14s\n", "BlFract", "PitchAxis", "StrcTwst", "BMassDen", "FlpStff", "EdgStff")
	fmt.Fprintf(w, "%14s %14s %14s %14s %14s %14s\n", "(-)", "(-)", "(deg)", "(kg/m)", "(Nm^2)", "(Nm^2)")
	edb := s.(*ElastoDynBlade)
	for _, BlInpSt := range edb.BlInpSt {
		fmt.Fprintf(w, "%14g %14g %14g %14g %14g %14g\n",
			BlInpSt.BlFract, BlInpSt.PitchAxis, BlInpSt.StrcTwst,
			BlInpSt.BMassDen, BlInpSt.FlpStff, BlInpSt.EdgStff)
	}
	return nil
}

func parseBlInpSt(s any, v reflect.Value, lines []string) ([]string, error) {
	edb := s.(*ElastoDynBlade)
	edb.BlInpSt = make([]ElastoDynBladeBlInpSt, edb.NBlInpSt)
	if len(lines) < edb.NBlInpSt+3 {
		return nil, fmt.Errorf("insufficient lines to parse distributed Blade properties")
	}
	lines = lines[3:]
	for i := range edb.BlInpSt {
		n, err := fmt.Sscan(lines[0], &edb.BlInpSt[i].BlFract,
			&edb.BlInpSt[i].PitchAxis, &edb.BlInpSt[i].StrcTwst,
			&edb.BlInpSt[i].BMassDen, &edb.BlInpSt[i].FlpStff,
			&edb.BlInpSt[i].EdgStff)
		if n < 6 {
			return nil, fmt.Errorf("error parsing distributed Blade properties: %s", err)
		}
		lines = lines[1:]
	}
	return lines, nil
}
