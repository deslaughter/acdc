package input

import (
	"fmt"
	"io"
	"reflect"
)

var BeamDynBladeSchema = NewSchema("BeamDynBlade", []SchemaEntry{
	{Heading: "BeamDyn Blade Input File"},
	{Keyword: "Title", Type: String, Format: formatTitle, Parse: parseTitle},
	{Heading: "Blade Parameters"},
	{Keyword: "station_total", Type: Int, Desc: `Number of blade input stations`},
	{Keyword: "damp_type", Type: Int, Desc: `Damping type: 0: no damping; 1: damped`},
	{Heading: "Damping Coefficient"},
	{Keyword: "mu", Type: Float, Dims: 1, Parse: parseBeamDynBladeMu, Format: formatBeamDynBladeMu},
	{Heading: "Distributed Properties"},
	{Keyword: "DistProps", Dims: 1,
		Table: &Table{
			Columns: []TableColumn{
				{Keyword: "station_eta", Type: Float},
				{Keyword: "stiffness_matrix", Type: Float, Dims: 2},
				{Keyword: "mass_matrix", Type: Float, Dims: 2},
			},
		},
		Parse: parseBeamDynBladeDistProps, Format: formatBeamDynBladeDistProps,
	},
})

func parseBeamDynBladeMu(s any, v reflect.Value, lines []string) ([]string, error) {
	bd := s.(*BeamDynBlade)
	line := ""
	line, lines = lines[3], lines[4:]
	bd.Mu = make([]float64, 6)
	num, err := fmt.Sscan(line, &bd.Mu[0], &bd.Mu[1], &bd.Mu[2], &bd.Mu[3], &bd.Mu[4], &bd.Mu[5])
	if num < 6 || err != nil {
		return nil, err
	}
	return lines, nil
}

func formatBeamDynBladeMu(w io.Writer, s, field any, entry SchemaEntry) error {
	bd := s.(*BeamDynBlade)
	fmt.Fprintf(w, "%9s %9s %9s %9s %9s %9s\n",
		"mu1", "mu2", "mu3", "mu4", "mu5", "mu6")
	fmt.Fprintf(w, "%9s %9s %9s %9s %9s %9s\n", "(-)", "(-)", "(-)", "(-)", "(-)", "(-)")
	fmt.Fprintf(w, "%9g %9g %9g %9g %9g %9g\n",
		bd.Mu[0], bd.Mu[1], bd.Mu[2], bd.Mu[3], bd.Mu[4], bd.Mu[5])
	return nil
}

func formatBeamDynBladeDistProps(w io.Writer, s, field any, entry SchemaEntry) error {
	bd := s.(*BeamDynBlade)
	for _, dp := range bd.DistProps {
		fmt.Fprintf(w, "%15.7g\n", dp.Station_eta)
		for _, row := range dp.Stiffness_matrix {
			for _, v := range row {
				fmt.Fprintf(w, "%15.7g", v)
			}
			fmt.Fprint(w, "\n")
		}
		fmt.Fprint(w, "\n")

		for _, row := range dp.Mass_matrix {
			for _, v := range row {
				fmt.Fprintf(w, "%15.7g", v)
			}
			fmt.Fprint(w, "\n")
		}
		fmt.Fprint(w, "\n")
	}
	return nil
}

func parseBeamDynBladeDistProps(s any, v reflect.Value, lines []string) ([]string, error) {
	bd := s.(*BeamDynBlade)
	bd.DistProps = make([]BeamDynBladeDistProps, bd.Station_total)
	line := ""
	lines = lines[1:]

	for i := range bd.DistProps {

		line, lines = lines[0], lines[1:]
		num, err := fmt.Sscan(line, &bd.DistProps[i].Station_eta)
		if num < 1 || err != nil {
			return nil, err
		}

		bd.DistProps[i].Stiffness_matrix = make([][]float64, 6)
		for j := range bd.DistProps[i].Stiffness_matrix {
			row := make([]float64, 6)
			line, lines = lines[0], lines[1:]
			num, err := fmt.Sscan(line, &row[0], &row[1], &row[2], &row[3], &row[4], &row[5])
			if num < 6 || err != nil {
				return nil, err
			}
			bd.DistProps[i].Stiffness_matrix[j] = row
		}
		lines = lines[1:]

		bd.DistProps[i].Mass_matrix = make([][]float64, 6)
		for j := range bd.DistProps[i].Mass_matrix {
			row := make([]float64, 6)
			line, lines = lines[0], lines[1:]
			num, err := fmt.Sscan(line, &row[0], &row[1], &row[2], &row[3], &row[4], &row[5])
			if num < 6 || err != nil {
				return nil, err
			}
			bd.DistProps[i].Mass_matrix[j] = row
		}
		lines = lines[1:]
	}
	return lines, nil
}
