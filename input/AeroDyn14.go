package input

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

var AeroDyn14Schema = NewSchema("AeroDyn14", []SchemaEntry{
	{Heading: "AeroDyn 14 Input File"},
	{Keyword: "Title", Type: String, Parse: parseTitle, Format: formatTitle},
	{Keyword: "StallMod", Type: String, Desc: `Dynamic stall included`, Unit: "unquoted string", Options: []Option{{"BEDDOES", "BEDDOES"}, {"STEADY", "STEADY"}}},
	{Keyword: "UseCm", Type: String, Desc: `Use aerodynamic pitching moment model?`, Unit: "unquoted string", Options: []Option{{"USE_CM", "USE_CM"}, {"NO_CM", "NO_CM"}}},
	{Keyword: "InfModel", Type: String, Desc: `Inflow model [DYNIN or EQUIL]`, Unit: "unquoted string"},
	{Keyword: "IndModel", Type: String, Desc: `Induction-factor model [NONE or WAKE or SWIRL]`, Unit: "unquoted string"},
	{Keyword: "AToler", Type: Float, Desc: `Induction-factor tolerance (convergence criteria) (-)`},
	{Keyword: "TLModel", Type: String, Desc: `Tip-loss model (EQUIL only) [PRANDtl, GTECH, or NONE]`, Unit: "unquoted string"},
	{Keyword: "HLModel", Type: String, Desc: `Hub-loss model (EQUIL only) [PRANdtl or NONE]`, Unit: "unquoted string"},
	{Keyword: "TwrShad", Type: Float, Desc: `Tower-shadow velocity deficit (-)`},
	{Keyword: "ShadHWid", Type: Float, Desc: `Tower-shadow half width (m)`},
	{Keyword: "T_Shad_Refpt", Type: Float, Desc: `Tower-shadow reference point (m)`},
	{Keyword: "AirDens", Type: Float, Desc: `Air density (kg/m^3)`},
	{Keyword: "KinVisc", Type: Float, Desc: `Kinematic air viscosity [CURRENTLY IGNORED] (m^2/sec)`},
	{Keyword: "DTAero", Type: Float, Desc: `Time interval for aerodynamic calculations (sec)`},
	{Keyword: "NumFoil", Type: Int, Desc: `Number of airfoil files (-)`},
	{Keyword: "FoilNm", Type: String, Dims: 1, Parse: parseFoilNm, Format: formatFoilNm, Desc: `Names of the airfoil files [NumFoil lines] (quoted strings)`},
	{Keyword: "BldNodes", Type: Int, Desc: `Number of blade nodes used for analysis (-)`},
	{Keyword: "BlNd", Dims: 1, Parse: parseAeroDyn14BlNd, Format: formatAeroDyn14BlNd,
		Table: &Table{
			Columns: []TableColumn{
				{Keyword: "RNodes", Type: Float},
				{Keyword: "AeroTwst", Type: Float},
				{Keyword: "DRNodes", Type: Float},
				{Keyword: "Chord", Type: Float},
				{Keyword: "NFoil", Type: Int},
				{Keyword: "PrnElm", Type: String},
			},
		},
	},
})

func formatFoilNm(w io.Writer, s, field any, entry SchemaEntry) error {
	vals := field.([]string)
	fmt.Fprintf(w, "%-12s    %-14s - %s\n",
		`"`+vals[0]+`"`, entry.Keyword, entry.Desc)
	for _, v := range vals[1:] {
		fmt.Fprintf(w, "%-12s\n", `"`+v+`"`)
	}
	return nil
}

func parseFoilNm(s any, v reflect.Value, lines []string) ([]string, error) {
	names := make([]string, s.(*AeroDyn14).NumFoil)
	for {
		if len(lines) == 0 {
			return lines, fmt.Errorf("label FoilNm not found")
		}
		if strings.Contains(strings.ToLower(lines[0]), "foilnm") {
			break
		}
		lines = lines[1:]
	}
	if len(lines) < len(names) {
		return lines, fmt.Errorf("insufficient lines for FoilNm (%d)", len(lines))
	}
	line := ""
	for i := range names {
		line, lines = strings.Trim(lines[0], ` "`), lines[1:]
		if strings.Contains(line, `"`) {
			names[i], _, _ = strings.Cut(line, `"`)
		} else {
			names[i] = strings.Fields(line)[0]
		}
	}
	v.Set(reflect.ValueOf(names))
	return lines, nil
}

func formatAeroDyn14BlNd(w io.Writer, s, field any, entry SchemaEntry) error {
	fmt.Fprintf(w, "%12s %12s %12s %12s %8s %10s\n",
		"RNodes", "AeroTwst", "DRNodes", "Chord", "NFoil", "PrnElm")
	ad := s.(*AeroDyn14)
	for _, BlNd := range ad.BlNd {
		fmt.Fprintf(w, "%12g %12g %12g %12g %8d %10s\n", BlNd.RNodes,
			BlNd.AeroTwst, BlNd.DRNodes, BlNd.Chord, BlNd.NFoil, BlNd.PrnElm)
	}
	return nil
}

func parseAeroDyn14BlNd(s any, v reflect.Value, lines []string) ([]string, error) {
	ad := s.(*AeroDyn14)
	ad.BlNd = make([]AeroDyn14BlNd, ad.BldNodes)
	line := ""
	lines = lines[1:]
	for i := range ad.BlNd {
		line, lines = lines[0], lines[1:]
		num, err := fmt.Sscan(line, &ad.BlNd[i].RNodes, &ad.BlNd[i].AeroTwst,
			&ad.BlNd[i].DRNodes, &ad.BlNd[i].Chord, &ad.BlNd[i].NFoil,
			&ad.BlNd[i].PrnElm)
		if num < 6 {
			return nil, err
		}
	}
	return lines, nil
}
