package input

import (
	"fmt"
	"io"
	"reflect"
)

var AeroDynBladeSchema = NewSchema("AeroDynBlade", []SchemaEntry{
	{Heading: "AeroDyn Blade Input File"},
	{Keyword: "Title", Type: String, Format: formatTitle, Parse: parseTitle},
	{Heading: "Blade Properties"},
	{Keyword: "NumBlNds", Type: Int, Desc: `Number of blade nodes used in the analysis (-)`},
	{Keyword: "BlNd", Dims: 1, Format: formatAeroDynBladeBlNd, Parse: parseAeroDynBladeBlNd,
		Table: &Table{
			Columns: []TableColumn{
				{Keyword: "BlSpn", Type: Float, Desc: `(m)`},
				{Keyword: "BlCrvAC", Type: Float, Desc: `(m)`},
				{Keyword: "BlSwpAC", Type: Float, Desc: `(m)`},
				{Keyword: "BlCrvAng", Type: Float, Desc: `(deg)`},
				{Keyword: "BlTwist", Type: Float, Desc: `(deg)`},
				{Keyword: "BlChord", Type: Float, Desc: `(m)`},
				{Keyword: "BlAFID", Type: Int, Desc: `(-)`},
			},
		},
	},
})

func formatAeroDynBladeBlNd(w io.Writer, s, field any, entry SchemaEntry) error {
	fmt.Fprintf(w, "%12s %12s %12s %12s %12s %12s %10s\n",
		"BlSpn", "BlCrvAC", "BlSwpAC", "BlCrvAng", "BlTwist", "BlChord", "BlAFID")
	fmt.Fprintf(w, "%12s %12s %12s %12s %12s %12s %10s\n",
		"(m)", "(m)", "(m)", "(deg)", "(deg)", "(m)", "(-)")
	adb := s.(*AeroDynBlade)
	for _, BlNd := range adb.BlNd {
		fmt.Fprintf(w, "%12g %12g %12g %12g %12g %12g %10d\n",
			BlNd.BlSpn, BlNd.BlCrvAC, BlNd.BlSwpAC, BlNd.BlCrvAng,
			BlNd.BlTwist, BlNd.BlChord, BlNd.BlAFID)
	}
	return nil
}

func parseAeroDynBladeBlNd(s any, v reflect.Value, lines []string) ([]string, error) {
	adb := s.(*AeroDynBlade)
	adb.BlNd = make([]AeroDynBladeBlNd, adb.NumBlNds)
	line := ""
	lines = lines[2:]
	for i := range adb.BlNd {
		line, lines = lines[0], lines[1:]
		num, err := fmt.Sscan(line, &adb.BlNd[i].BlSpn, &adb.BlNd[i].BlCrvAC,
			&adb.BlNd[i].BlSwpAC, &adb.BlNd[i].BlCrvAng, &adb.BlNd[i].BlTwist,
			&adb.BlNd[i].BlChord, &adb.BlNd[i].BlAFID)
		if num < 7 {
			return nil, err
		}
	}
	return lines, nil
}
