package input

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Model struct {
	AeroDyn14       *AeroDyn14
	AeroDyn15       *AeroDyn15
	AD15AirfoilInfo []*AD15AirfoilInfo
	AeroDynBlade    []*AeroDynBlade
	BeamDyn         []*BeamDyn
	BeamDynBlade    []*BeamDynBlade
	ElastoDyn       *ElastoDyn
	ElastoDynBlade  []*ElastoDynBlade
	ElastoDynTower  *ElastoDynTower
	FAST            *FAST
	InflowWind      *InflowWind
	ServoDyn        *ServoDyn
}

func ReadFiles(fastFilePath string) (*Model, error) {

	var err error

	// Create inputs structure
	inp := Model{
		FAST: NewFAST(),
	}

	rootDir := filepath.Dir(fastFilePath)

	numBlades := 0

	//--------------------------------------------------------------------------
	// FAST file
	//--------------------------------------------------------------------------

	if inp.FAST, err = ReadFAST(fastFilePath); err != nil {
		return nil, err
	}

	//--------------------------------------------------------------------------
	// ElastoDyn file
	//--------------------------------------------------------------------------

	if inp.FAST.CompElast > 0 {

		// ElastoDyn File
		ElastoDynFilePath := filepath.Join(rootDir, inp.FAST.EDFile)
		if inp.ElastoDyn, err = ReadElastoDyn(ElastoDynFilePath); err != nil {
			return nil, fmt.Errorf("error reading ElastoDyn File: %w", err)
		}
		numBlades = inp.ElastoDyn.NumBl

		// ElastoDyn Blades (not using BeamDyn)
		if inp.FAST.CompElast < 2 {

			ElastoDynBladeFiles := []string{
				inp.ElastoDyn.BldFile1,
				inp.ElastoDyn.BldFile2,
				inp.ElastoDyn.BldFile3,
			}[:numBlades]
			fileMap := map[string]struct{}{}
			for _, file := range ElastoDynBladeFiles[1:] {
				fileMap[file] = struct{}{}
			}
			if len(fileMap) == 1 {
				ElastoDynBladeFiles = ElastoDynBladeFiles[:1]
			}
			inp.ElastoDynBlade = make([]*ElastoDynBlade, len(ElastoDynBladeFiles))
			for i, ElastoDynBladeFile := range ElastoDynBladeFiles {
				ElastoDynBladeFilePath := filepath.Join(filepath.Dir(ElastoDynFilePath), ElastoDynBladeFile)
				inp.ElastoDynBlade[i], err = ReadElastoDynBlade(ElastoDynBladeFilePath)
				if err != nil {
					return nil, fmt.Errorf("error reading ElastoDyn Blade File: %w", err)
				}
			}
		}

		// ElastoDyn Tower
		ElastoDynTowerFilePath := filepath.Join(filepath.Dir(ElastoDynFilePath), inp.ElastoDyn.TwrFile)
		inp.ElastoDynTower, err = ReadElastoDynTower(ElastoDynTowerFilePath)
		if err != nil {
			return nil, err
		}
	}

	//--------------------------------------------------------------------------
	// BeamDynBlade files
	//--------------------------------------------------------------------------

	if inp.FAST.CompElast == 2 {

		BeamDynFiles := []string{
			inp.FAST.BDBldFile1,
			inp.FAST.BDBldFile2,
			inp.FAST.BDBldFile3,
		}[:numBlades]
		fileMap := map[string]struct{}{}
		for _, file := range BeamDynFiles[1:] {
			fileMap[file] = struct{}{}
		}
		if len(fileMap) == 1 {
			BeamDynFiles = BeamDynFiles[:1]
		}
		inp.BeamDyn = make([]*BeamDyn, len(BeamDynFiles))
		inp.BeamDynBlade = make([]*BeamDynBlade, len(BeamDynFiles))
		for i, BeamDynFile := range BeamDynFiles {

			// Read BeamDyn file
			BeamDynFilePath := filepath.Join(rootDir, BeamDynFile)
			inp.BeamDyn[i], err = ReadBeamDyn(BeamDynFilePath)
			if err != nil {
				return nil, fmt.Errorf("error reading BeamDyn File %d: %w", i+1, err)
			}

			// Read BeamDyn Blade file
			BeamDynBladeFilePath := filepath.Join(filepath.Dir(BeamDynFilePath), inp.BeamDyn[i].BldFile)
			inp.BeamDynBlade[i], err = ReadBeamDynBlade(BeamDynBladeFilePath)
			if err != nil {
				return nil, fmt.Errorf("error reading BeamDyn Blade File %d: %w", i+1, err)
			}
		}
	}

	//--------------------------------------------------------------------------
	// InflowWind file
	//--------------------------------------------------------------------------

	if inp.FAST.CompInflow > 0 {
		InflowWindFilePath := filepath.Join(rootDir, inp.FAST.InflowFile)
		if inp.InflowWind, err = ReadInflowWind(InflowWindFilePath); err != nil {
			return nil, fmt.Errorf("error reading InflowWind File: %w", err)
		}
	}

	//--------------------------------------------------------------------------
	// AeroDyn14 file
	//--------------------------------------------------------------------------

	if inp.FAST.CompAero == 1 {
		AeroDynFilePath := filepath.Join(rootDir, inp.FAST.AeroFile)
		if inp.AeroDyn14, err = ReadAeroDyn14(AeroDynFilePath); err != nil {
			return nil, fmt.Errorf("error reading AeroDyn14 File: %w", err)
		}
	}

	//--------------------------------------------------------------------------
	// AeroDyn15 file
	//--------------------------------------------------------------------------

	if inp.FAST.CompAero == 2 {

		AeroDynFilePath := filepath.Join(rootDir, inp.FAST.AeroFile)
		if inp.AeroDyn15, err = ReadAeroDyn15(AeroDynFilePath); err != nil {
			return nil, fmt.Errorf("error reading AeroDyn15 File: %w", err)
		}

		// Read blade files
		inp.AeroDynBlade = make([]*AeroDynBlade, numBlades)
		AeroDynBladeFiles := []string{
			inp.AeroDyn15.ADBlFile1,
			inp.AeroDyn15.ADBlFile2,
			inp.AeroDyn15.ADBlFile3,
		}[:numBlades]
		for i, AeroDynBladeFile := range AeroDynBladeFiles {
			AeroDynBladeFilePath := filepath.Join(filepath.Dir(AeroDynFilePath), AeroDynBladeFile)
			inp.AeroDynBlade[i], err = ReadAeroDynBlade(AeroDynBladeFilePath)
			if err != nil {
				return nil, fmt.Errorf("error reading AeroDyn15 Blade File %d: %w", i+1, err)
			}
		}

		// Read airfoil info files
		inp.AD15AirfoilInfo = make([]*AD15AirfoilInfo, len(inp.AeroDyn15.AFNames))
		for i, AirfoilInfoFile := range inp.AeroDyn15.AFNames {
			AirfoilInfoFilePath := filepath.Join(filepath.Dir(AeroDynFilePath), AirfoilInfoFile)
			inp.AD15AirfoilInfo[i], err = ReadAD15AirfoilInfo(AirfoilInfoFilePath)
			if err != nil {
				return nil, fmt.Errorf("error reading AeroDyn15 AirfoilInfo File %d: %w", i+1, err)
			}
		}
	}

	//--------------------------------------------------------------------------
	// ServoDyn file
	//--------------------------------------------------------------------------

	if inp.FAST.CompServo == 1 {
		ServoDynFilePath := filepath.Join(rootDir, inp.FAST.ServoFile)
		inp.ServoDyn, err = ReadServoDyn(ServoDynFilePath)
		if err != nil {
			return nil, fmt.Errorf("error reading ServoDyn File: %w", err)
		}
	}

	//--------------------------------------------------------------------------
	// HydroDyn file
	//--------------------------------------------------------------------------

	// if inp.Fast.CompHydro == 1 {

	// 	text, err := os.ReadFile(filepath.Join(rootDir, inp.Fast.HydroFile))
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error reading '%s': %w", fastFilePath, err)
	// 	}

	// inp.HydroDyn = NewHydroDyn()
	// 	if err = inp.HydroDyn.Parse(text); err != nil {
	// 		return nil, fmt.Errorf("error parsing '%s': %w", fastFilePath, err)
	// 	}
	// }

	//--------------------------------------------------------------------------
	// SubDyn file
	//--------------------------------------------------------------------------

	// if inp.Fast.CompSub == 1 {

	// 	text, err := os.ReadFile(filepath.Join(rootDir, inp.Fast.SubFile))
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error reading '%s': %w", fastFilePath, err)
	// 	}

	// inp.SubDyn = NewSubDyn()
	// 	if err = inp.SubDyn.Parse(text); err != nil {
	// 		return nil, fmt.Errorf("error parsing '%s': %w", fastFilePath, err)
	// 	}
	// }

	//--------------------------------------------------------------------------
	// Ice file
	//--------------------------------------------------------------------------

	// if inp.Fast.CompMooring == 1 {

	// 	text, err := os.ReadFile(filepath.Join(rootDir, inp.Fast.MooringFile))
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error reading '%s': %w", fastFilePath, err)
	// 	}

	// inp.Mooring = NewMooring()
	// 	if err = inp.Mooring.Parse(text); err != nil {
	// 		return nil, fmt.Errorf("error parsing '%s': %w", fastFilePath, err)
	// 	}
	// }

	//--------------------------------------------------------------------------
	// Ice file
	//--------------------------------------------------------------------------

	// if inp.Fast.CompIce == 1 {

	// 	text, err := os.ReadFile(filepath.Join(rootDir, inp.Fast.IceFile))
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error reading '%s': %w", fastFilePath, err)
	// 	}

	// inp.Ice = NewIce()
	// 	if err = inp.Ice.Parse(text); err != nil {
	// 		return nil, fmt.Errorf("error parsing '%s': %w", fastFilePath, err)
	// 	}
	// }

	return &inp, nil
}

func (inp *Model) WriteFiles(fastFilePath string) error {

	rootDir := filepath.Dir(fastFilePath)
	baseName := strings.TrimSuffix(filepath.Base(fastFilePath), filepath.Ext(fastFilePath))

	// Create root directory
	if err := os.MkdirAll(rootDir, 0777); err != nil {
		return fmt.Errorf("error creating '%s': %w", rootDir, err)
	}

	//--------------------------------------------------------------------------
	// ElastoDyn file
	//--------------------------------------------------------------------------

	if inp.FAST.CompElast > 0 {

		inp.FAST.EDFile = baseName + "_ElastoDyn.dat"
		ElastoDynFilePath := filepath.Join(rootDir, inp.FAST.EDFile)

		// Write ElastoDyn Blades (not using BeamDyn)
		if inp.FAST.CompElast < 2 {

			ElastoDynBladeFiles := []*string{
				&inp.ElastoDyn.BldFile1,
				&inp.ElastoDyn.BldFile2,
				&inp.ElastoDyn.BldFile3,
			}
			for i := range ElastoDynBladeFiles {
				if len(inp.ElastoDynBlade) == 1 {
					*ElastoDynBladeFiles[i] = fmt.Sprintf("%s_ElastoDyn_Blade.dat", baseName)
				} else {
					*ElastoDynBladeFiles[i] = fmt.Sprintf("%s_ElastoDyn_Blade_%02d.dat", baseName, i+1)
				}
			}
			for i, ElastoDynBlade := range inp.ElastoDynBlade {
				ElastoDynBladeFilePath := filepath.Join(rootDir, *ElastoDynBladeFiles[i])
				if err := ElastoDynBlade.Write(ElastoDynBladeFilePath); err != nil {
					return err
				}
			}
		}

		// Write ElastoDynTower file
		inp.ElastoDyn.TwrFile = baseName + "_ElastoDynTower.dat"
		if err := inp.ElastoDynTower.Write(filepath.Join(rootDir, inp.ElastoDyn.TwrFile)); err != nil {
			return err
		}

		// Write ElastoDyn file
		if err := inp.ElastoDyn.Write(ElastoDynFilePath); err != nil {
			return err
		}
	}

	//--------------------------------------------------------------------------
	// BeamDynBlade files
	//--------------------------------------------------------------------------

	if inp.FAST.CompElast == 2 {

		BeamDynFiles := []*string{
			&inp.FAST.BDBldFile1,
			&inp.FAST.BDBldFile2,
			&inp.FAST.BDBldFile3,
		}
		for i := range BeamDynFiles {
			if len(inp.BeamDyn) == 1 {
				*BeamDynFiles[i] = fmt.Sprintf("%s_BeamDyn.dat", baseName)
			} else {
				*BeamDynFiles[i] = fmt.Sprintf("%s_BeamDyn_%02d.dat", baseName, i+1)
			}
		}
		for i, BeamDyn := range inp.BeamDyn {

			BeamDyn.BldFile = strings.ReplaceAll(*BeamDynFiles[i], "_BeamDyn", "_BeamDynBlade")
			BeamDynBladeFilePath := filepath.Join(rootDir, BeamDyn.BldFile)
			if err := inp.BeamDynBlade[i].Write(BeamDynBladeFilePath); err != nil {
				return err
			}

			BeamDynFilePath := filepath.Join(rootDir, *BeamDynFiles[i])
			if err := BeamDyn.Write(BeamDynFilePath); err != nil {
				return err
			}
		}
	}

	//--------------------------------------------------------------------------
	// InflowWind file
	//--------------------------------------------------------------------------

	if inp.FAST.CompInflow > 0 {
		inp.FAST.InflowFile = baseName + "_InflowWind.dat"
		InflowWindFilePath := filepath.Join(rootDir, inp.FAST.InflowFile)
		if err := inp.InflowWind.Write(InflowWindFilePath); err != nil {
			return err
		}
	}

	//--------------------------------------------------------------------------
	// AeroDyn14 file
	//--------------------------------------------------------------------------

	if inp.FAST.CompAero == 1 {
		inp.FAST.AeroFile = baseName + "_AeroDyn.dat"
		AeroDynFilePath := filepath.Join(rootDir, inp.FAST.AeroFile)
		if err := inp.AeroDyn14.Write(AeroDynFilePath); err != nil {
			return err
		}
	}

	//--------------------------------------------------------------------------
	// AeroDyn15 file
	//--------------------------------------------------------------------------

	if inp.FAST.CompAero == 2 {

		// Create airfoil directory
		airfoilDir := filepath.Join(rootDir, "Airfoils")
		if err := os.MkdirAll(airfoilDir, 0777); err != nil {
			return fmt.Errorf("error creating '%s': %w", airfoilDir, err)
		}

		// Create path to AeroDyn file
		inp.FAST.AeroFile = baseName + "_AeroDyn15.dat"
		AeroDynFilePath := filepath.Join(rootDir, inp.FAST.AeroFile)

		// Write blade files
		bladeFiles := []*string{
			&inp.AeroDyn15.ADBlFile1,
			&inp.AeroDyn15.ADBlFile2,
			&inp.AeroDyn15.ADBlFile3,
		}
		for i, AeroDynBlade := range inp.AeroDynBlade {
			*bladeFiles[i] = fmt.Sprintf("%s_AeroDyn15Blade_%02d.dat", baseName, i+1)
			bladeFilePath := filepath.Join(rootDir, *bladeFiles[i])
			if err := AeroDynBlade.Write(bladeFilePath); err != nil {
				return err
			}
		}

		// Write airfoil info files
		for i, AirfoilInfo := range inp.AD15AirfoilInfo {

			// Create path to airfoil file
			airfoilFileName := fmt.Sprintf("%s_Airfoil_%02d", baseName, i+1)
			airfoilFilePath := filepath.Join(airfoilDir, airfoilFileName)

			// Generate relative path from root directory to airfoil file
			relativePath, err := filepath.Rel(rootDir, airfoilFilePath)
			if err != nil {
				return fmt.Errorf("error generating path '%s': %w", airfoilFilePath, err)
			}
			inp.AeroDyn15.AFNames[i] = relativePath

			if err := AirfoilInfo.Write(airfoilFilePath); err != nil {
				return err
			}
		}

		// Write AeroDyn file
		if err := inp.AeroDyn15.Write(AeroDynFilePath); err != nil {
			return err
		}
	}

	//--------------------------------------------------------------------------
	// ServoDyn file
	//--------------------------------------------------------------------------

	if inp.FAST.CompServo == 1 {
		inp.FAST.ServoFile = baseName + "_ServoDyn.dat"
		ServoDynFilePath := filepath.Join(rootDir, inp.FAST.ServoFile)
		if err := inp.ServoDyn.Write(ServoDynFilePath); err != nil {
			return err
		}
	}

	//--------------------------------------------------------------------------
	// HydroDyn file
	//--------------------------------------------------------------------------

	// if inp.Fast.CompHydro == 1 {

	// 	text, err := os.ReadFile(filepath.Join(rootDir, inp.Fast.HydroFile))
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error reading '%s': %w", fastFilePath, err)
	// 	}

	// inp.HydroDyn = NewHydroDyn()
	// 	if err = inp.HydroDyn.Parse(text); err != nil {
	// 		return nil, fmt.Errorf("error parsing '%s': %w", fastFilePath, err)
	// 	}
	// }

	//--------------------------------------------------------------------------
	// SubDyn file
	//--------------------------------------------------------------------------

	// if inp.Fast.CompSub == 1 {

	// 	text, err := os.ReadFile(filepath.Join(rootDir, inp.Fast.SubFile))
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error reading '%s': %w", fastFilePath, err)
	// 	}

	// inp.SubDyn = NewSubDyn()
	// 	if err = inp.SubDyn.Parse(text); err != nil {
	// 		return nil, fmt.Errorf("error parsing '%s': %w", fastFilePath, err)
	// 	}
	// }

	//--------------------------------------------------------------------------
	// Ice file
	//--------------------------------------------------------------------------

	// if inp.Fast.CompMooring == 1 {

	// 	text, err := os.ReadFile(filepath.Join(rootDir, inp.Fast.MooringFile))
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error reading '%s': %w", fastFilePath, err)
	// 	}

	// inp.Mooring = NewMooring()
	// 	if err = inp.Mooring.Parse(text); err != nil {
	// 		return nil, fmt.Errorf("error parsing '%s': %w", fastFilePath, err)
	// 	}
	// }

	//--------------------------------------------------------------------------
	// Ice file
	//--------------------------------------------------------------------------

	// if inp.Fast.CompIce == 1 {

	// 	text, err := os.ReadFile(filepath.Join(rootDir, inp.Fast.IceFile))
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error reading '%s': %w", fastFilePath, err)
	// 	}

	// inp.Ice = NewIce()
	// 	if err = inp.Ice.Parse(text); err != nil {
	// 		return nil, fmt.Errorf("error parsing '%s': %w", fastFilePath, err)
	// 	}
	// }

	//--------------------------------------------------------------------------
	// FAST file
	//--------------------------------------------------------------------------

	if err := inp.FAST.Write(fastFilePath); err != nil {
		return err
	}

	return nil
}
