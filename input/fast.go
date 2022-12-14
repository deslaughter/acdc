package input

var FASTSchema = NewSchema("FAST", []SchemaEntry{
	{Heading: "OpenFAST Input File"},
	{Keyword: "Title", Type: String, Desc: `File title`, Format: formatTitle, Parse: parseTitle},

	{Heading: "Simulation Control"},
	{Keyword: "Echo", Type: Bool, Desc: `Echo input data to <RootName>.ech`, Unit: "flag"},
	{Keyword: "AbortLevel", Type: String, Desc: `Error level when simulation should abort`, Unit: "string", Options: []Option{{"WARNING", "WARNING"}, {"SEVERE", "SEVERE"}, {"FATAL", "FATAL"}}},
	{Keyword: "TMax", Type: Float, Desc: `Total run time`, Unit: "sec"},
	{Keyword: "DT", Type: Float, Desc: `Recommended module time step`, Unit: "sec"},
	{Keyword: "InterpOrder", Type: Int, Desc: `Interpolation order for input/output time history`, Unit: "-", Options: []Option{{1, "Linear"}, {2, "Quadratic"}}},
	{Keyword: "NumCrctn", Type: Int, Desc: `Number of correction iterations {0=explicit calculation, i.e., no corrections}`, Unit: "-"},
	{Keyword: "DT_UJac", Type: Float, Desc: `Time between calls to get Jacobians`, Unit: "sec"},
	{Keyword: "UJacSclFact", Type: Float, Desc: `Scaling factor used in Jacobians`, Unit: "-"},

	{Heading: "Feature switches and flags"},
	{Keyword: "CompElast", Type: Int, Desc: `Compute structural dynamics`, Unit: "switch", Options: []Option{{1, "ElastoDyn"}, {2, "ElastoDyn + BeamDyn for blades"}}},
	{Keyword: "CompInflow", Type: Int, Desc: `Compute inflow wind velocities`, Unit: "switch", Options: []Option{{0, "Still Air"}, {1, "InflowWind"}, {2, "External (OpenFOAM)"}}},
	{Keyword: "CompAero", Type: Int, Desc: `Compute aerodynamic loads`, Unit: "switch", Options: []Option{{0, "None"}, {1, "AeroDyn v14"}, {2, "AeroDyn v15"}}},
	{Keyword: "CompServo", Type: Int, Desc: `Compute control and electrical-drive dynamics`, Unit: "switch", Options: []Option{{0, "None"}, {1, "ServoDyn"}}},
	{Keyword: "CompHydro", Type: Int, Desc: `Compute hydrodynamic loads`, Unit: "switch", Options: []Option{{0, "None"}, {1, "HydroDyn"}}},
	{Keyword: "CompSub", Type: Int, Desc: `Compute sub-structural dynamics`, Unit: "switch", Options: []Option{{0, "None"}, {1, "SubDyn"}, {2, "External Platform MCKF"}}},
	{Keyword: "CompMooring", Type: Int, Desc: `Compute mooring system`, Unit: "switch", Options: []Option{{0, "None"}, {1, "MAP++"}, {2, "FEAMooring"}, {3, "MoorDyn"}, {4, "OrcaFlex"}}},
	{Keyword: "CompIce", Type: Int, Desc: `Compute ice loads`, Unit: "switch", Options: []Option{{0, "None"}, {1, "IceFloe"}, {2, "IceDyn"}}},
	{Keyword: "MHK", Type: Int, Desc: `MHK turbine type`, Unit: "switch", Options: []Option{{0, "Not an MHK turbine"}, {1, "Fixed MHK turbine"}, {2, "Floating MHK turbine"}}},

	{Heading: "Environmental Conditions"},
	{Keyword: "Gravity", Type: Float, Desc: `Gravitational acceleration`, Unit: "m/s^2"},
	{Keyword: "AirDens", Type: Float, Desc: `Air density`, Unit: "kg/m^3"},
	{Keyword: "WtrDens", Type: Float, Desc: `Water density`, Unit: "kg/m^3"},
	{Keyword: "KinVisc", Type: Float, Desc: `Kinematic viscosity of working fluid`, Unit: "m^2/s"},
	{Keyword: "SpdSound", Type: Float, Desc: `Speed of sound in working fluid`, Unit: "m/s"},
	{Keyword: "Patm", Type: Float, Desc: `Atmospheric pressure [used only for an MHK turbine cavitation check]`, Unit: "Pa"},
	{Keyword: "Pvap", Type: Float, Desc: `Vapour pressure of working fluid [used only for an MHK turbine cavitation check]`, Unit: "Pa"},
	{Keyword: "WtrDpth", Type: Float, Desc: `Water depth`, Unit: "m"},
	{Keyword: "MSL2SWL", Type: Float, Desc: `Offset between still-water level and mean sea level [positive upward]`, Unit: "m"},

	{Heading: "Input Files"},
	{Keyword: "EDFile", Type: String, Desc: `Name of file containing ElastoDyn input parameters`, Unit: "quoted string"},
	{Keyword: "BDBldFile(1)", Type: String, Desc: `Name of file containing BeamDyn input parameters for blade 1`, Unit: "quoted string", Show: []Condition{{"CompElast", "==", 2}}},
	{Keyword: "BDBldFile(2)", Type: String, Desc: `Name of file containing BeamDyn input parameters for blade 2`, Unit: "quoted string", Show: []Condition{{"CompElast", "==", 2}}},
	{Keyword: "BDBldFile(3)", Type: String, Desc: `Name of file containing BeamDyn input parameters for blade 3`, Unit: "quoted string", Show: []Condition{{"CompElast", "==", 2}}},
	{Keyword: "InflowFile", Type: String, Desc: `Name of file containing inflow wind input parameters`, Unit: "quoted string", Show: []Condition{{"CompInflow", "==", 1}}},
	{Keyword: "AeroFile", Type: String, Desc: `Name of file containing aerodynamic input parameters`, Unit: "quoted string", Show: []Condition{{"CompAero", "==", 1}}},
	{Keyword: "ServoFile", Type: String, Desc: `Name of file containing control and electrical-drive input parameters`, Unit: "quoted string", Show: []Condition{{"CompServo", "==", 1}}},
	{Keyword: "HydroFile", Type: String, Desc: `Name of file containing hydrodynamic input parameters`, Unit: "quoted string", Show: []Condition{{"CompHydro", "==", 1}}},
	{Keyword: "SubFile", Type: String, Desc: `Name of file containing sub-structural input parameters`, Unit: "quoted string", Show: []Condition{{"CompSub", "==", 1}}},
	{Keyword: "MooringFile", Type: String, Desc: `Name of file containing mooring system input parameters`, Unit: "quoted string", Show: []Condition{{"CompMooring", ">", 0}}},
	{Keyword: "IceFile", Type: String, Desc: `Name of file containing ice input parameters`, Unit: "quoted string", Show: []Condition{{"CompIce", ">", 0}}},

	{Heading: "Output"},
	{Keyword: "SumPrint", Type: Bool, Desc: `Print summary data to '<RootName>.sum'`, Unit: "flag"},
	{Keyword: "SttsTime", Type: Float, Desc: `Amount of time between screen status messages`, Unit: "sec"},
	{Keyword: "ChkptTime", Type: Float, Desc: `Amount of time between creating checkpoint files for potential restart`, Unit: "sec"},
	{Keyword: "DT_Out", Type: Float, CanBeDefault: true, Desc: `Time step for tabular output (or "default")`, Unit: "sec"},
	{Keyword: "TStart", Type: Float, Desc: `Time to begin tabular output`, Unit: "sec"},
	{Keyword: "OutFileFmt", Type: Int, Desc: `Format for tabular (time-marching) output file`, Unit: "switch", Options: []Option{{0, "uncompressed binary [<RootName>.outb]"}, {1, "text file [<RootName>.out]"}, {2, "binary file [<RootName>.outb]"}, {3, "both 1 and 2"}}},
	{Keyword: "TabDelim", Type: Bool, Desc: `Use tab delimiters in text tabular output file?`, Unit: "flag", Options: []Option{{true, "Tabs"}, {false, "Spaces"}}},
	{Keyword: "OutFmt", Type: String, Desc: `Format used for text tabular output, excluding the time channel.  Resulting field should be 10 characters`, Unit: "quoted string"},

	{Heading: "Linearization"},
	{Keyword: "Linearize", Type: Bool, Desc: `Linearization analysis`},
	{Keyword: "CalcSteady", Type: Bool, Desc: `Calculate a steady-state periodic operating point before linearization?`, Unit: "flag", Show: []Condition{{"Linearize", "==", true}}},
	{Keyword: "TrimCase", Type: Int, Desc: `Controller parameter to be trimmed`, Unit: "-", Options: []Option{{1, "yaw"}, {2, "torque"}, {3, "pitch"}}, Show: []Condition{{"CalcSteady", "==", true}}},
	{Keyword: "TrimTol", Type: Float, Desc: `Tolerance for the rotational speed convergence`, Unit: "-", Show: []Condition{{"CalcSteady", "==", true}}},
	{Keyword: "TrimGain", Type: Float, Desc: `Proportional gain for the rotational speed error (>0) (rad/(rad/s) for yaw or pitch; Nm/(rad/s) for torque)`, Show: []Condition{{"CalcSteady", "==", true}}},
	{Keyword: "Twr_Kdmp", Type: Float, Desc: `Damping factor for the tower`, Unit: "N/(m/s)", Show: []Condition{{"CalcSteady", "==", true}}},
	{Keyword: "Bld_Kdmp", Type: Float, Desc: `Damping factor for the blades`, Unit: "N/(m/s)", Show: []Condition{{"CalcSteady", "==", true}}},
	{Keyword: "NLinTimes", Type: Int, Desc: `Number of times to linearize [>=1]`, Show: []Condition{{"Linearize", "==", true}}},
	{Keyword: "LinTimes", Type: Float, Dims: 1, Desc: `List of times at which to linearize [1 to NLinTimes] [used only when Linearize=True and CalcSteady=False]`, Unit: "sec", Show: []Condition{{"Linearize", "==", true}, {"CalcSteady", "==", false}}},
	{Keyword: "LinInputs", Type: Int, Desc: `Inputs included in linearization`, Unit: "switch", Options: []Option{{0, "none"}, {1, "standard"}, {2, "all module inputs (debug)"}}, Show: []Condition{{"Linearize", "==", true}}},
	{Keyword: "LinOutputs", Type: Int, Desc: `Outputs included in linearization`, Unit: "switch", Options: []Option{{0, "none"}, {1, "from OutList(s)"}, {2, "all module outputs (debug)"}}, Show: []Condition{{"Linearize", "==", true}}},
	{Keyword: "LinOutJac", Type: Bool, Desc: `Include full Jacobians in linearization output (for debug)`, Unit: "flag", Show: []Condition{{"Linearize", "==", true}, {"LinInputs", "==", 2}, {"LinOutputs", "==", 2}}},
	{Keyword: "LinOutMod", Type: Bool, Desc: `Write module-level linearization output files in addition to output for full system?`, Unit: "flag", Show: []Condition{{"Linearize", "==", true}}},

	{Heading: "Visualization"},
	{Keyword: "WrVTK", Type: Int, Desc: `VTK visualization data output`, Unit: "switch", Options: []Option{{0, "none"}, {1, "initialization data only"}, {2, "animation"}, {3, "mode shapes"}}},
	{Keyword: "VTK_type", Type: Int, Desc: `Type of VTK visualization data`, Unit: "switch", Options: []Option{{1, "surfaces"}, {2, "basic meshes (lines/points)"}, {3, "all meshes (debug)"}}, Show: []Condition{{"WrVTK", ">", 0}}},
	{Keyword: "VTK_fields", Type: Bool, Desc: `Write mesh fields to VTK data files?`, Unit: "flag", Show: []Condition{{"WrVTK", ">", 0}}},
	{Keyword: "VTK_fps", Type: Float, Desc: `Frame rate for VTK output (frames per second) {will use closest integer multiple of DT}`, Show: []Condition{{"WrVTK", "in", []int{2, 3}}}},
})
