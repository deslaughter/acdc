package anl

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/mat"
)

type LinData struct {
	FilePath       string
	SimTime        float64
	RotorSpeed     float64
	Azimuth        float64
	WindSpeed      float64
	NumX           int
	NumX2          int
	NumXd          int
	NumZ           int
	NumU           int
	NumY           int
	X, Xd, Z, U, Y []OperPointData
	A, B, C, D     *mat.Dense
}

type OperPointData struct {
	RC         int
	OperPoint  float64
	IsRotating bool
	DerivOrder int
	Desc       string
}

func ReadLinData(filePath string) (*LinData, error) {

	linData := &LinData{FilePath: filePath}

	linFile, err := os.Open(linData.FilePath)
	if err != nil {
		return linData, err
	}
	defer linFile.Close()

	scanner := bufio.NewScanner(linFile)

	//--------------------------------------------------------------------------
	// Header
	//--------------------------------------------------------------------------

	hasJacobians := false

	for scanner.Scan() {

		// Get line without leading/trailing whitespace, skip if empty
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		// Split line into fields
		fields := strings.Fields(line)

		if strings.HasPrefix(line, "Simulation time") {
			if linData.SimTime, err = strconv.ParseFloat(fields[2], 64); err != nil {
				return nil, fmt.Errorf("error parsing Simulation time: %w", err)
			}
		} else if strings.HasPrefix(line, "Rotor Speed") {
			if linData.RotorSpeed, err = strconv.ParseFloat(fields[2], 64); err != nil {
				return nil, fmt.Errorf("error parsing Rotor Speed: %w", err)
			}
		} else if strings.HasPrefix(line, "Azimuth") {
			if linData.Azimuth, err = strconv.ParseFloat(fields[1], 64); err != nil {
				return nil, fmt.Errorf("error parsing Azimuth: %w", err)
			}
		} else if strings.HasPrefix(line, "Wind Speed") {
			if linData.WindSpeed, err = strconv.ParseFloat(fields[2], 64); err != nil {
				return nil, fmt.Errorf("error parsing Wind Speed: %w", err)
			}
		} else if strings.HasPrefix(line, "Number of continuous states") {
			if linData.NumX, err = strconv.Atoi(fields[4]); err != nil {
				return nil, fmt.Errorf("error parsing Number of continuous states: %w", err)
			}
		} else if strings.HasPrefix(line, "Number of discrete states") {
			if linData.NumXd, err = strconv.Atoi(fields[4]); err != nil {
				return nil, fmt.Errorf("error parsing Number of discrete states: %w", err)
			}
		} else if strings.HasPrefix(line, "Number of constraint states") {
			if linData.NumZ, err = strconv.Atoi(fields[4]); err != nil {
				return nil, fmt.Errorf("error parsing Number of constraint states: %w", err)
			}
		} else if strings.HasPrefix(line, "Number of inputs") {
			if linData.NumU, err = strconv.Atoi(fields[3]); err != nil {
				return nil, fmt.Errorf("error parsing Number of inputs: %w", err)
			}
		} else if strings.HasPrefix(line, "Number of outputs") {
			if linData.NumY, err = strconv.Atoi(fields[3]); err != nil {
				return nil, fmt.Errorf("error parsing Number of outputs: %w", err)
			}
		} else if strings.HasPrefix(line, "Jacobians included") {
			if fields[5] == "Yes" {
				hasJacobians = true
			}
			break
		}
	}

	_ = hasJacobians

	//--------------------------------------------------------------------------
	// Operating points
	//--------------------------------------------------------------------------

	var currentOP *[]OperPointData
	hasDeriv := false
	defaultDeriv := 0

	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		fields := strings.Fields(line)

		if strings.HasPrefix(line, "Row/Column") {
			hasDeriv = strings.Contains(line, "Derivative Order")
		} else if strings.HasPrefix(line, "Order of continuous states") {
			currentOP = &linData.X
			defaultDeriv = 2
		} else if strings.HasPrefix(line, "Order of continuous state derivatives") {
			currentOP = &linData.Xd
			defaultDeriv = 2
		} else if strings.HasPrefix(line, "Linearized state matrices") {
			break
		} else {

			// Get first column as integer, skip line if not valid
			rc, err := strconv.Atoi(fields[0])
			if err != nil {
				continue
			}

			op, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return nil, err
			}

			isRotating, err := strconv.ParseBool(fields[2])
			if err != nil {
				return nil, err
			}

			derivOrder := defaultDeriv
			if hasDeriv {
				derivOrder, err = strconv.Atoi(fields[3])
				if err != nil {
					return nil, err
				}
			}

			*currentOP = append(*currentOP, OperPointData{
				RC:         rc,
				OperPoint:  op,
				IsRotating: isRotating,
				DerivOrder: derivOrder,
				Desc:       strings.Join(fields[4:], " "),
			})
		}
	}

	// Sum number of second order continuous states
	for _, op := range linData.X {
		if op.DerivOrder == 2 {
			linData.NumX2++
		}
	}

	//--------------------------------------------------------------------------
	// State Matrices
	//--------------------------------------------------------------------------

	var matrix *mat.Dense
	iRow := 0
	for scanner.Scan() {

		// Get line with whitespace removed
		line := strings.TrimSpace(scanner.Text())

		// If nothing on line, continue
		if len(line) == 0 {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) == 4 && fields[2] == "x" {
			rows, _ := strconv.Atoi(fields[1])
			cols, _ := strconv.Atoi(fields[3])
			matrix = mat.NewDense(rows, cols, nil)
			iRow = 0
		}

		// Switch based on beginning of line
		switch line[0] {
		case 'A':
			linData.A = matrix
		case 'B':
			linData.B = matrix
		case 'C':
			linData.C = matrix
		case 'D':
			linData.D = matrix
		default:
			row := make([]float64, len(fields))
			for i, s := range fields {
				if row[i], err = strconv.ParseFloat(s, 64); err != nil {
					return nil, err
				}
			}
			matrix.SetRow(iRow, row)
			iRow++
		}
	}

	// Get error from scanner
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return linData, nil
}
