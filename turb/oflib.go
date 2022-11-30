//go:build oflib

package turb

// #cgo CFLAGS: -I/Users/dslaught/projects/openfast-mac/install/include
// #cgo LDFLAGS: -L/Users/dslaught/projects/openfast-mac/install/lib -lopenfastlib
//
// #include <stdlib.h>
// #include "FAST_Library.h"
import "C"
import (
	"bytes"
	"fmt"
	"math"
	"path/filepath"
	"sync"
	"unsafe"
)

type FASTLib struct {
	numTurbines int
}

type Turbines []*Turbine

func NewTurbines(numTurbines int) (Turbines, error) {

	// Convert numTurbines to C int
	numTurbinesC := C.int(numTurbines)

	// Error handling variables
	errStat := C.int(0)
	errMsg := make([]byte, C.INTERFACE_STRING_LENGTH)

	// Call FAST routine to allocate turbines
	C.FAST_AllocateTurbines(&numTurbinesC, &errStat,
		(*C.char)(unsafe.Pointer(&errMsg[0])))

	// If error status is not none, return message
	if errStat != C.ErrID_None {
		return nil, fmt.Errorf("%s", bytes.TrimSpace(errMsg[:len(errMsg)-1]))
	}

	// Create slice to represent each turbine
	turbines := make(Turbines, numTurbines)
	for i := range turbines {
		turbines[i] = &Turbine{
			Num:           C.int(i),
			numInputs:     C.NumFixedInputs,
			abortErrLevel: C.ErrID_Fatal,
		}
	}

	return turbines, nil
}

// Delete deallocates memory associated with the turbines
func (ts *Turbines) Delete() error {

	// Error handling variables
	errStat := C.int(0)
	errMsg := make([]byte, C.INTERFACE_STRING_LENGTH)

	// Deallocate memory
	C.FAST_DeallocateTurbines(&errStat, (*C.char)(unsafe.Pointer(&errMsg[0])))
	if errStat != C.ErrID_None {
		return fmt.Errorf("%s", bytes.TrimSpace(errMsg[:len(errMsg)-1]))
	}

	// Nullify slice of turbines so they can't be reused
	*ts = nil

	return nil
}

type Turbine struct {
	Num                      C.int
	ChannelNames             []string
	InputPath                string
	errStat                  C.int
	errMsg                   [C.INTERFACE_STRING_LENGTH]byte
	abortErrLevel            C.int
	TimeStep, TimeStepOutput C.double
	TimeMax                  C.double
	numInputs                C.int
	inputData                []C.double
	numOutputs               C.int
	outputData               [][]C.double
	NumOutputSteps           int
	NumUpdateSteps           int
	stop                     sync.Once
}

type TurbineOpts struct {
	MinOutput bool
}

func (t *Turbine) Initialize(inputPath string, opts *TurbineOpts) error {

	var err error

	// Get input path as absolute path
	t.InputPath, err = filepath.Abs(inputPath)
	if err != nil {
		return fmt.Errorf("error converting '%s' to absolute path: %s", inputPath, err)
	}

	// Convert absolute path to bytes
	absInputPathBytes := []byte(t.InputPath)

	// Channel names byte array
	var channelNames [C.CHANNEL_LENGTH * C.MAXIMUM_OUTPUTS]byte

	// Initialize simulation
	C.FAST_Sizes(&t.Num,
		(*C.char)(unsafe.Pointer(&absInputPathBytes[0])),
		&t.abortErrLevel, &t.numOutputs, &t.TimeStep, &t.TimeStepOutput,
		&t.TimeMax, &t.errStat, (*C.char)(unsafe.Pointer(&t.errMsg[0])),
		(*C.char)(unsafe.Pointer(&channelNames[0])), nil, nil)

	// If error, return message
	if t.errStat != C.ErrID_None {
		return fmt.Errorf("%s", bytes.TrimSpace(t.errMsg[:len(t.errMsg)-1]))
	}

	// Extract channel names from byte slice
	t.ChannelNames = make([]string, t.numOutputs)
	for i := range t.ChannelNames {
		channelName := channelNames[i*C.CHANNEL_LENGTH : (i+1)*C.CHANNEL_LENGTH]
		t.ChannelNames[i] = string(bytes.TrimSpace(channelName))
	}

	// Calculate number of update steps
	t.NumUpdateSteps = int(math.Ceil(float64(t.TimeMax/t.TimeStep)) + 1)

	// If options are not nill and min output flag is true, set number of
	// output steps to one
	if opts != nil && opts.MinOutput {
		t.TimeStepOutput = t.TimeMax
	}

	// Calculate number of output steps
	t.NumOutputSteps = int(math.Ceil(float64(t.TimeMax/t.TimeStepOutput)) + 1)

	// Allocate storage for input
	t.inputData = make([]C.double, t.numInputs)

	// Allocate storage for output
	t.outputData = make([][]C.double, t.NumOutputSteps)
	for i := range t.outputData {
		t.outputData[i] = make([]C.double, t.numOutputs)
	}

	return nil
}

func (t *Turbine) Run() error {

	// Flag to end simulation early
	EndSimulationEarly := C.bool(false)

	// Calculate output frequency relative to step time
	output_frequency := int(math.Round(float64(t.TimeStepOutput / t.TimeStep)))

	// Start simulation
	C.FAST_Start(&t.Num, &t.numInputs, &t.numOutputs, &t.inputData[0],
		&t.outputData[0][0], &t.errStat, (*C.char)(unsafe.Pointer(&t.errMsg[0])))
	if t.errStat != C.ErrID_None {
		return fmt.Errorf("%s", bytes.TrimSpace(t.errMsg[:len(t.errMsg)-1]))
	}

	// Loop through update steps
	for iStep, iOut := 0, 0; iStep < t.NumUpdateSteps; iStep++ {

		// Update solution, return if error
		C.FAST_Update(&t.Num, &t.numInputs, &t.numOutputs,
			&t.inputData[0], &(t.outputData[iOut])[0],
			&EndSimulationEarly, &t.errStat,
			(*C.char)(unsafe.Pointer(&t.errMsg[0])))
		if t.errStat != C.ErrID_None {
			return fmt.Errorf("%s", bytes.TrimSpace(t.errMsg[:len(t.errMsg)-1]))
		}

		// Update output data index
		if iStep%output_frequency == 0 {
			iOut++
		}

		// If simulation requests early end, break
		if EndSimulationEarly {
			break
		}
	}

	return nil
}

// Stop ends the turbine simulation.
func (t *Turbine) Stop() {
	t.stop.Do(func() {
		stopProgram := C.bool(false)
		C.FAST_End(&t.Num, &stopProgram)
	})
}
