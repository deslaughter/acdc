package input

import (
	"fmt"
	"io"
	"reflect"
)

var ServoDynSchema = NewSchema("ServoDyn", []SchemaEntry{
	{Heading: "ServoDyn Input File"},
	{Keyword: "Title", Type: String, Format: formatTitle, Parse: parseTitle},
	{Heading: "Simulation Control"},
	{Keyword: "Echo", Type: Bool, Desc: `Echo input data to <RootName>.ech (flag)`},
	{Keyword: "DT", Type: Float, CanBeDefault: true, Desc: `Communication interval for controllers (s) (or "default")`},
	{Heading: "Pitch Control"},
	{Keyword: "PCMode", Type: Int, Desc: `Pitch control mode {0: none, 3: user-defined from routine PitchCntrl, 4: user-defined from Simulink/Labview, 5: user-defined from Bladed-style DLL} (switch)`},
	{Keyword: "TPCOn", Type: Float, Desc: `Time to enable active pitch control (s) [unused when PCMode=0]`},
	{Keyword: "TPitManS(1)", Type: Float, Desc: `Time to start override pitch maneuver for blade 1 and end standard pitch control (s)`},
	{Keyword: "TPitManS(2)", Type: Float, Desc: `Time to start override pitch maneuver for blade 2 and end standard pitch control (s)`},
	{Keyword: "TPitManS(3)", Type: Float, Desc: `Time to start override pitch maneuver for blade 3 and end standard pitch control (s) [unused for 2 blades]`},
	{Keyword: "PitManRat(1)", Type: Float, Desc: `Pitch rate at which override pitch maneuver heads toward final pitch angle for blade 1 (deg/s)`},
	{Keyword: "PitManRat(2)", Type: Float, Desc: `Pitch rate at which override pitch maneuver heads toward final pitch angle for blade 2 (deg/s)`},
	{Keyword: "PitManRat(3)", Type: Float, Desc: `Pitch rate at which override pitch maneuver heads toward final pitch angle for blade 3 (deg/s) [unused for 2 blades]`},
	{Keyword: "BlPitchF(1)", Type: Float, Desc: `Blade 1 final pitch for pitch maneuvers (degrees)`},
	{Keyword: "BlPitchF(2)", Type: Float, Desc: `Blade 2 final pitch for pitch maneuvers (degrees)`},
	{Keyword: "BlPitchF(3)", Type: Float, Desc: `Blade 3 final pitch for pitch maneuvers (degrees) [unused for 2 blades]`},
	{Heading: "Generator And Torque Control"},
	{Keyword: "VSContrl", Type: Int, Desc: `Variable-speed control mode {0: none, 1: simple VS, 3: user-defined from routine UserVSCont, 4: user-defined from Simulink/Labview, 5: user-defined from Bladed-style DLL} (switch)`},
	{Keyword: "GenModel", Type: Int, Desc: `Generator model {1: simple, 2: Thevenin, 3: user-defined from routine UserGen} (switch) [used only when VSContrl=0]`},
	{Keyword: "GenEff", Type: Float, Desc: `Generator efficiency [ignored by the Thevenin and user-defined generator models] (%)`},
	{Keyword: "GenTiStr", Type: Bool, Desc: `Method to start the generator {T: timed using TimGenOn, F: generator speed using SpdGenOn} (flag)`},
	{Keyword: "GenTiStp", Type: Bool, Desc: `Method to stop the generator {T: timed using TimGenOf, F: when generator power = 0} (flag)`},
	{Keyword: "SpdGenOn", Type: Float, Desc: `Generator speed to turn on the generator for a startup (HSS speed) (rpm) [used only when GenTiStr=False]`},
	{Keyword: "TimGenOn", Type: Float, Desc: `Time to turn on the generator for a startup (s) [used only when GenTiStr=True]`},
	{Keyword: "TimGenOf", Type: Float, Desc: `Time to turn off the generator (s) [used only when GenTiStp=True]`},
	{Heading: "Simple Variable-Speed Torque Control"},
	{Keyword: "VS_RtGnSp", Type: Float, Desc: `Rated generator speed for simple variable-speed generator control (HSS side) (rpm) [used only when VSContrl=1]`},
	{Keyword: "VS_RtTq", Type: Float, Desc: `Rated generator torque/constant generator torque in Region 3 for simple variable-speed generator control (HSS side) (N-m) [used only when VSContrl=1]`},
	{Keyword: "VS_Rgn2K", Type: Float, Desc: `Generator torque constant in Region 2 for simple variable-speed generator control (HSS side) (N-m/rpm^2) [used only when VSContrl=1]`},
	{Keyword: "VS_SlPc", Type: Float, Desc: `Rated generator slip percentage in Region 2 1/2 for simple variable-speed generator control (%) [used only when VSContrl=1]`},
	{Heading: "Simple Induction Generator"},
	{Keyword: "SIG_SlPc", Type: Float, Desc: `Rated generator slip percentage (%) [used only when VSContrl=0 and GenModel=1]`},
	{Keyword: "SIG_SySp", Type: Float, Desc: `Synchronous (zero-torque) generator speed (rpm) [used only when VSContrl=0 and GenModel=1]`},
	{Keyword: "SIG_RtTq", Type: Float, Desc: `Rated torque (N-m) [used only when VSContrl=0 and GenModel=1]`},
	{Keyword: "SIG_PORt", Type: Float, Desc: `Pull-out ratio (Tpullout/Trated) (-) [used only when VSContrl=0 and GenModel=1]`},
	{Heading: "Thevenin-Equivalent Induction Generator"},
	{Keyword: "TEC_Freq", Type: Float, Desc: `Line frequency [50 or 60] (Hz) [used only when VSContrl=0 and GenModel=2]`},
	{Keyword: "TEC_NPol", Type: Int, Desc: `Number of poles [even integer > 0] (-) [used only when VSContrl=0 and GenModel=2]`},
	{Keyword: "TEC_SRes", Type: Float, Desc: `Stator resistance (ohms) [used only when VSContrl=0 and GenModel=2]`},
	{Keyword: "TEC_RRes", Type: Float, Desc: `Rotor resistance (ohms) [used only when VSContrl=0 and GenModel=2]`},
	{Keyword: "TEC_VLL", Type: Float, Desc: `Line-to-line RMS voltage (volts) [used only when VSContrl=0 and GenModel=2]`},
	{Keyword: "TEC_SLR", Type: Float, Desc: `Stator leakage reactance (ohms) [used only when VSContrl=0 and GenModel=2]`},
	{Keyword: "TEC_RLR", Type: Float, Desc: `Rotor leakage reactance (ohms) [used only when VSContrl=0 and GenModel=2]`},
	{Keyword: "TEC_MR", Type: Float, Desc: `Magnetizing reactance (ohms) [used only when VSContrl=0 and GenModel=2]`},
	{Heading: "High-Speed Shaft Brake"},
	{Keyword: "HSSBrMode", Type: Int, Desc: `HSS brake model {0: none, 1: simple, 3: user-defined from routine UserHSSBr, 4: user-defined from Simulink/Labview, 5: user-defined from Bladed-style DLL} (switch)`},
	{Keyword: "THSSBrDp", Type: Float, Desc: `Time to initiate deployment of the HSS brake (s)`},
	{Keyword: "HSSBrDT", Type: Float, Desc: `Time for HSS-brake to reach full deployment once initiated (sec) [used only when HSSBrMode=1]`},
	{Keyword: "HSSBrTqF", Type: Float, Desc: `Fully deployed HSS-brake torque (N-m)`},
	{Heading: "Nacelle-Yaw Control"},
	{Keyword: "YCMode", Type: Int, Desc: `Yaw control mode {0: none, 3: user-defined from routine UserYawCont, 4: user-defined from Simulink/Labview, 5: user-defined from Bladed-style DLL} (switch)`},
	{Keyword: "TYCOn", Type: Float, Desc: `Time to enable active yaw control (s) [unused when YCMode=0]`},
	{Keyword: "YawNeut", Type: Float, Desc: `Neutral yaw position--yaw spring force is zero at this yaw (degrees)`},
	{Keyword: "YawSpr", Type: Float, Desc: `Nacelle-yaw spring constant (N-m/rad)`},
	{Keyword: "YawDamp", Type: Float, Desc: `Nacelle-yaw damping constant (N-m/(rad/s))`},
	{Keyword: "TYawManS", Type: Float, Desc: `Time to start override yaw maneuver and end standard yaw control (s)`},
	{Keyword: "YawManRat", Type: Float, Desc: `Yaw maneuver rate (in absolute value) (deg/s)`},
	{Keyword: "NacYawF", Type: Float, Desc: `Final yaw angle for override yaw maneuvers (degrees)`},
	{Heading: "Aerodynamic Flow Control"},
	{Keyword: "AfCmode", Type: Int, Desc: `Airfoil control mode {0: none, 1: cosine wave cycle, 4: user-defined from Simulink/Labview, 5: user-defined from Bladed-style DLL} (switch)`},
	{Keyword: "AfC_Mean", Type: Float, Desc: `Mean level for cosine cycling or steady value (-) [used only with AfCmode==1]`},
	{Keyword: "AfC_Amp", Type: Float, Desc: `Amplitude for for cosine cycling of flap signal (-) [used only with AfCmode==1]`},
	{Keyword: "AfC_Phase", Type: Float, Desc: `Phase relative to the blade azimuth (0 is vertical) for for cosine cycling of flap signal (deg) [used only with AfCmode==1]`},
	{Heading: "Structural Control"},
	{Keyword: "NumBStC", Type: Int, Desc: `Number of blade structural controllers (integer)`},
	{Keyword: "BStCfiles", Type: String, Desc: `Name of the files for blade structural controllers (quoted strings) [unused when NumBStC==0]`},
	{Keyword: "NumNStC", Type: Int, Desc: `Number of nacelle structural controllers (integer)`},
	{Keyword: "NStCfiles", Type: String, Desc: `Name of the files for nacelle structural controllers (quoted strings) [unused when NumNStC==0]`},
	{Keyword: "NumTStC", Type: Int, Desc: `Number of tower structural controllers (integer)`},
	{Keyword: "TStCfiles", Type: String, Desc: `Name of the files for tower structural controllers (quoted strings) [unused when NumTStC==0]`},
	{Keyword: "NumSStC", Type: Int, Desc: `Number of substructure structural controllers (integer)`},
	{Keyword: "SStCfiles", Type: String, Desc: `Name of the files for substructure structural controllers (quoted strings) [unused when NumSStC==0]`},
	{Heading: "Cable Control"},
	{Keyword: "CCmode", Type: Int, Desc: `Cable control mode {0: none, 4: user-defined from Simulink/Labview, 5: user-defined from Bladed-style DLL} (switch)`},
	{Heading: "BLADED Interface"},
	{Keyword: "DLL_FileName", Type: String, Desc: `Name/location of the dynamic library {.dll [Windows] or .so [Linux]} in the Bladed-DLL format (-) [used only with Bladed Interface]`},
	{Keyword: "DLL_InFile", Type: String, Desc: `Name of input file sent to the DLL (-) [used only with Bladed Interface]`},
	{Keyword: "DLL_ProcName", Type: String, Desc: `Name of procedure in DLL to be called (-) [case sensitive; used only with DLL Interface]`},
	{Keyword: "DLL_DT", Type: Float, CanBeDefault: true, Desc: `Communication interval for dynamic library (s) (or "default") [used only with Bladed Interface]`},
	{Keyword: "DLL_Ramp", Type: Bool, Desc: `Whether a linear ramp should be used between DLL_DT time steps [introduces time shift when true] (flag) [used only with Bladed Interface]`},
	{Keyword: "BPCutoff", Type: Float, Desc: `Cutoff frequency for low-pass filter on blade pitch from DLL (Hz) [used only with Bladed Interface]`},
	{Keyword: "NacYaw_North", Type: Float, Desc: `Reference yaw angle of the nacelle when the upwind end points due North (deg) [used only with Bladed Interface]`},
	{Keyword: "Ptch_Cntrl", Type: Int, Desc: `Record 28: Use individual pitch control {0: collective pitch; 1: individual pitch control} (switch) [used only with Bladed Interface]`},
	{Keyword: "Ptch_SetPnt", Type: Float, Desc: `Record  5: Below-rated pitch angle set-point (deg) [used only with Bladed Interface]`},
	{Keyword: "Ptch_Min", Type: Float, Desc: `Record  6: Minimum pitch angle (deg) [used only with Bladed Interface]`},
	{Keyword: "Ptch_Max", Type: Float, Desc: `Record  7: Maximum pitch angle (deg) [used only with Bladed Interface]`},
	{Keyword: "PtchRate_Min", Type: Float, Desc: `Record  8: Minimum pitch rate (most negative value allowed) (deg/s) [used only with Bladed Interface]`},
	{Keyword: "PtchRate_Max", Type: Float, Desc: `Record  9: Maximum pitch rate  (deg/s) [used only with Bladed Interface]`},
	{Keyword: "Gain_OM", Type: Float, Desc: `Record 16: Optimal mode gain (Nm/(rad/s)^2) [used only with Bladed Interface]`},
	{Keyword: "GenSpd_MinOM", Type: Float, Desc: `Record 17: Minimum generator speed (rpm) [used only with Bladed Interface]`},
	{Keyword: "GenSpd_MaxOM", Type: Float, Desc: `Record 18: Optimal mode maximum speed (rpm) [used only with Bladed Interface]`},
	{Keyword: "GenSpd_Dem", Type: Float, Desc: `Record 19: Demanded generator speed above rated (rpm) [used only with Bladed Interface]`},
	{Keyword: "GenTrq_Dem", Type: Float, Desc: `Record 22: Demanded generator torque above rated (Nm) [used only with Bladed Interface]`},
	{Keyword: "GenPwr_Dem", Type: Float, Desc: `Record 13: Demanded power (W) [used only with Bladed Interface]`},
	{Heading: "BLADED Interface Torque-Speed Look-Up Table"},
	{Keyword: "DLL_NumTrq", Type: Int, Desc: `Record 26: No. of points in torque-speed look-up table {0 = none and use the optimal mode parameters; nonzero = ignore the optimal mode PARAMETERs by setting Record 16 to 0.0} (-) [used only with Bladed Interface]`},
	{Keyword: "GenSpdTrq", Dims: 1,
		Table: &Table{
			Columns: []TableColumn{
				{Keyword: "GenSpd_TLU", Type: Float, Dims: 1, Desc: `GenSpd_TLU`},
				{Keyword: "GenTrq_TLU", Type: Float, Dims: 1, Desc: `GenTrq_TLU`},
			},
		},
		Parse: parseGenSpdTrq, Format: formatGenSpdTrq,
	},
	{Heading: "Output"},
	{Keyword: "SumPrint", Type: Bool, Desc: `Print summary data to <RootName>.sum (flag) (currently unused)`},
	{Keyword: "OutFile", Type: Int, Desc: `Switch to determine where output will be placed: {1: in module output file only; 2: in glue code output file only; 3: both} (currently unused)`},
	{Keyword: "TabDelim", Type: Bool, Desc: `Use tab delimiters in text tabular output file? (flag) (currently unused)`},
	{Keyword: "OutFmt", Type: String, Desc: `Format used for text tabular output (except time).  Resulting field should be 10 characters. (quoted string) (currently unused)`},
	{Keyword: "TStart", Type: Float, Desc: `Time to begin tabular output (s) (currently unused)`},
	{Keyword: "OutList", Type: String, Dims: 1, Format: formatOutList, Parse: parseOutList, Desc: `The next line(s) contains a list of output parameters.  See OutListParameters.xlsx for a listing of available output channels, (-)`},
})

func formatGenSpdTrq(w io.Writer, s, field any, entry SchemaEntry) error {
	fmt.Fprintf(w, "%14s %14s\n", "GenSpd_TLU", "GenTrq_TLU")
	fmt.Fprintf(w, "%14s %14s\n", "(rpm)", "(Nm)")
	sd := s.(*ServoDyn)
	for _, gst := range sd.GenSpdTrq {
		fmt.Fprintf(w, "%14g %14g\n", gst.GenSpd_TLU, gst.GenTrq_TLU)
	}
	return nil
}

func parseGenSpdTrq(s any, v reflect.Value, lines []string) ([]string, error) {
	sd := s.(*ServoDyn)
	sd.GenSpdTrq = make([]ServoDynGenSpdTrq, sd.DLL_NumTrq)
	line := ""
	lines = lines[2:]
	for i := range sd.GenSpdTrq {
		line, lines = lines[0], lines[1:]
		num, err := fmt.Sscan(line, &sd.GenSpdTrq[i].GenSpd_TLU, &sd.GenSpdTrq[i].GenTrq_TLU)
		if num < 2 {
			return nil, err
		}
	}
	return lines, nil
}
