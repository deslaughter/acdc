package input

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

var TurbSimSchema = NewSchema("TurbSim", []SchemaEntry{
	{Heading: "TurbSim Input File"},
	{Keyword: "Title", Type: String, Format: formatTitle, Parse: parseTitle},
	{Heading: "Runtime Options"},
	{Keyword: "Echo", Type: Bool, Desc: `Echo input data to <RootName>.ech (flag)`},
	{Keyword: "RandSeed1", Type: Int, Desc: `First random seed  (-2147483648 to 2147483647)`},
	{Keyword: "RandSeed2", Type: String, Desc: `Second random seed (-2147483648 to 2147483647) for intrinsic pRNG, or an alternative pRNG: "RanLux" or "RNSNLW"`},
	{Keyword: "WrBHHTP", Type: Bool, Desc: `Output hub-height turbulence parameters in binary form?  (Generates RootName.bin)`},
	{Keyword: "WrFHHTP", Type: Bool, Desc: `Output hub-height turbulence parameters in formatted form?  (Generates RootName.dat)`},
	{Keyword: "WrADHH", Type: Bool, Desc: `Output hub-height time-series data in AeroDyn form?  (Generates RootName.hh)`},
	{Keyword: "WrADFF", Type: Bool, Desc: `Output full-field time-series data in TurbSim/AeroDyn form? (Generates RootName.bts)`},
	{Keyword: "WrBLFF", Type: Bool, Desc: `Output full-field time-series data in BLADED/AeroDyn form?  (Generates RootName.wnd)`},
	{Keyword: "WrADTWR", Type: Bool, Desc: `Output tower time-series data? (Generates RootName.twr)`},
	{Keyword: "WrFMTFF", Type: Bool, Desc: `Output full-field time-series data in formatted (readable) form?  (Generates RootName.u, RootName.v, RootName.w)`},
	{Keyword: "WrACT", Type: Bool, Desc: `Output coherent turbulence time steps in AeroDyn form? (Generates RootName.cts)`},
	{Keyword: "Clockwise", Type: Bool, Desc: `Clockwise rotation looking downwind? (used only for full-field binary files - not necessary for AeroDyn)`},
	{Keyword: "ScaleIEC", Type: Int, Desc: `Scale IEC turbulence models to exact target standard deviation? [0=no additional scaling; 1=use hub scale uniformly; 2=use individual scales]`},
	{Heading: "-"},
	{Heading: "Turbine/Model Specifications"},
	{Keyword: "NumGrid_Z", Type: Int, Desc: `Vertical grid-point matrix dimension`},
	{Keyword: "NumGrid_Y", Type: Int, Desc: `Horizontal grid-point matrix dimension`},
	{Keyword: "TimeStep", Type: Float, Desc: `Time step [seconds]`},
	{Keyword: "AnalysisTime", Type: Float, Desc: `Length of analysis time series [seconds] (program will add time if necessary: AnalysisTime = MAX(AnalysisTime, UsableTime+GridWidth/MeanHHWS) )`},
	{Keyword: "UsableTime", Type: Float, Format: formatUsableTime, Parse: parseUsableTime, Desc: `Usable length of output time series [seconds] (program will add GridWidth/MeanHHWS seconds unless UsableTime is "ALL")`},
	{Keyword: "HubHt", Type: Float, Desc: `Hub height [m] (should be > 0.5*GridHeight)`},
	{Keyword: "GridHeight", Type: Float, Desc: `Grid height [m]`},
	{Keyword: "GridWidth", Type: Float, Desc: `Grid width [m] (should be >= 2*(RotorRadius+ShaftLength))`},
	{Keyword: "VFlowAng", Type: Float, Desc: `Vertical mean flow (uptilt) angle [degrees]`},
	{Keyword: "HFlowAng", Type: Float, Desc: `Horizontal mean flow (skew) angle [degrees]`},
	{Heading: "-"},
	{Heading: "Meteorological Boundary Conditions"},
	{Keyword: "TurbModel", Type: String, Desc: `Turbulence model ("IECKAI","IECVKM","GP_LLJ","NWTCUP","SMOOTH","WF_UPW","WF_07D","WF_14D","TIDAL","API","USRINP","TIMESR", or "NONE")`},
	{Keyword: "UserFile", Type: String, Desc: `Name of the file that contains inputs for user-defined spectra or time series inputs (used only for "USRINP" and "TIMESR" models)`},
	{Keyword: "IECstandard", Type: String, Desc: `Number of IEC 61400-x standard (x=1,2, or 3 with optional 61400-1 edition number (i.e. "1-Ed2") )`},
	{Keyword: "IECturbc", Type: String, Desc: `IEC turbulence characteristic ("A", "B", "C" or the turbulence intensity in percent) ("KHTEST" option with NWTCUP model, not used for other models)`},
	{Keyword: "IEC_WindType", Type: String, Desc: `IEC turbulence type ("NTM"=normal, "xETM"=extreme turbulence, "xEWM1"=extreme 1-year wind, "xEWM50"=extreme 50-year wind, where x=wind turbine class 1, 2, or 3)`},
	{Keyword: "ETMc", Type: Float, CanBeDefault: true, Desc: `IEC Extreme Turbulence Model "c" parameter [m/s]`},
	{Keyword: "WindProfileType", Type: String, CanBeDefault: true, Desc: `Velocity profile type ("LOG";"PL"=power law;"JET";"H2L"=Log law for TIDAL model;"API";"USR";"TS";"IEC"=PL on rotor disk, LOG elsewhere; or "default")`},
	{Keyword: "ProfileFile", Type: String, Desc: `Name of the file that contains input profiles for WindProfileType="USR" and/or TurbModel="USRVKM" [-]`},
	{Keyword: "RefHt", Type: Float, Desc: `Height of the reference velocity (URef) [m]`},
	{Keyword: "URef", Type: Float, CanBeDefault: true, Desc: `Mean (total) velocity at the reference height [m/s] (or "default" for JET velocity profile) [must be 1-hr mean for API model; otherwise is the mean over AnalysisTime seconds]`},
	{Keyword: "ZJetMax", Type: Float, CanBeDefault: true, Desc: `Jet height [m] (used only for JET velocity profile, valid 70-490 m)`},
	{Keyword: "PLExp", Type: Float, CanBeDefault: true, Desc: `Power law exponent [-] (or "default")`},
	{Keyword: "Z0", Type: Float, CanBeDefault: true, Desc: `Surface roughness length [m] (or "default")`},
	{Heading: "-"},
	{Heading: "Non-IEC Meteorological Boundary Conditions"},
	{Keyword: "Latitude", Type: Float, CanBeDefault: true, Desc: `Site latitude [degrees] (or "default")`},
	{Keyword: "RICH_NO", Type: Float, Desc: `Gradient Richardson number [-]`},
	{Keyword: "UStar", Type: Float, CanBeDefault: true, Desc: `Friction or shear velocity [m/s] (or "default")`},
	{Keyword: "ZI", Type: Float, CanBeDefault: true, Desc: `Mixing layer depth [m] (or "default")`},
	{Keyword: "PC_UW", Type: Float, CanBeDefault: true, Desc: `Hub mean uw Reynolds stress [m^2/s^2] (or "default" or "none")`},
	{Keyword: "PC_UV", Type: Float, CanBeDefault: true, Desc: `Hub mean uv Reynolds stress [m^2/s^2] (or "default" or "none")`},
	{Keyword: "PC_VW", Type: Float, CanBeDefault: true, Desc: `Hub mean vw Reynolds stress [m^2/s^2] (or "default" or "none")`},
	{Heading: "-"},
	{Heading: "Spatial Coherence Parameters"},
	{Keyword: "SCMod1", Type: String, CanBeDefault: true, Desc: `u-component coherence model ("GENERAL", "IEC", "API", "NONE", or "default")`},
	{Keyword: "SCMod2", Type: String, CanBeDefault: true, Desc: `v-component coherence model ("GENERAL", "IEC", "NONE", or "default")`},
	{Keyword: "SCMod3", Type: String, CanBeDefault: true, Desc: `w-component coherence model ("GENERAL", "IEC", "NONE", or "default")`},
	{Keyword: "InCDec1", Type: Float, CanBeDefault: true, Desc: `u-component coherence parameters for general or IEC models [-, m^-1] (e.g. "10.0  0.3e-3" in quotes) (or "default")`},
	{Keyword: "InCDec2", Type: Float, CanBeDefault: true, Desc: `v-component coherence parameters for general or IEC models [-, m^-1] (e.g. "10.0  0.3e-3" in quotes) (or "default")`},
	{Keyword: "InCDec3", Type: Float, CanBeDefault: true, Desc: `w-component coherence parameters for general or IEC models [-, m^-1] (e.g. "10.0  0.3e-3" in quotes) (or "default")`},
	{Keyword: "CohExp", Type: Float, CanBeDefault: true, Desc: `Coherence exponent for general model [-] (or "default")`},
	{Heading: "-"},
	{Heading: "Coherent Turbulence Scaling Parameters"},
	{Keyword: "CTEventPath", Type: String, Desc: `Name of the path where event data files are located`},
	{Keyword: "CTEventFile", Type: String, Desc: `Type of event files ("LES", "DNS", or "RANDOM")`},
	{Keyword: "Randomize", Type: Bool, Desc: `Randomize the disturbance scale and locations? (true/false)`},
	{Keyword: "DistScl", Type: Float, Desc: `Disturbance scale [-] (ratio of event dataset height to rotor disk). (Ignored when Randomize = true.)`},
	{Keyword: "CTLy", Type: Float, Desc: `Fractional location of tower centerline from right [-] (looking downwind) to left side of the dataset. (Ignored when Randomize = true.)`},
	{Keyword: "CTLz", Type: Float, Desc: `Fractional location of hub height from the bottom of the dataset. [-] (Ignored when Randomize = true.)`},
	{Keyword: "CTStartTime", Type: Float, Desc: `Minimum start time for coherent structures in RootName.cts [seconds]`},
})

func formatUsableTime(w io.Writer, s, field any, entry SchemaEntry) error {
	if value := field.(float64); value == -1 {
		fmt.Fprintf(w, "%12s    %-14s - %s\n", "ALL", entry.Keyword, entry.Desc)
	} else {
		fmt.Fprintf(w, "%12v    %-14s - %s\n", value, entry.Keyword, entry.Desc)
	}
	return nil
}

func parseUsableTime(s any, v reflect.Value, lines []string) ([]string, error) {
	line, lines, err := findKeywordLine("UsableTime", lines)
	if err != nil {
		return nil, err
	}
	value := strings.Fields(line)[0]
	if strings.ToUpper(value) == "ALL" {
		v.SetFloat(-1.0)
	} else {
		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing float '%s'", value)
		}
		v.SetFloat(num)
	}
	return lines, nil
}
