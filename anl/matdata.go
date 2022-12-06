package anl

import (
	"bytes"
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"regexp"
	"sort"
	"strings"

	"gonum.org/v1/gonum/mat"
)

type MatData struct {
	LinData    []*LinData
	NumStep    int
	NumStates  int
	NumStates2 int
	NumInputs  int
	NumOutputs int
	NumDOF1    int
	NumDOF2    int
	Azimuth    *mat.VecDense
	Omega      *mat.VecDense
	OmegaDot   *mat.VecDense
	WindSpeed  *mat.VecDense
	A, B, C, D []*mat.Dense
	OpX        []*mat.VecDense
	OpXd       []*mat.VecDense
	AvgA       *mat.Dense
	AvgOpX     *mat.VecDense
	AvgOpXd    *mat.VecDense
	Rotation   RotationTriplets
	Modes      []*ModeResults
}

type RotationTriplets struct {
	TripletsStates2 [][]int
	PermuteStates2  []int
	TripletsStates1 [][]int
	PermuteStates1  []int
	TripletsInputs  [][]int
	PermuteInputs   []int
	TripletsOutputs [][]int
	PermuteOutputs  []int
}

func collectMatrixData(linData []*LinData) (*MatData, error) {

	var err error

	// Sort linearization data by azimuth
	sort.Slice(linData, func(i, j int) bool {
		return linData[i].Azimuth < linData[j].Azimuth
	})

	// Get number of steps
	numSteps := len(linData)

	// Use last linearization data for initialization
	initData := linData[numSteps-1]

	// Create MBC structure
	md := &MatData{
		NumStep:    numSteps,
		LinData:    linData,
		NumStates:  initData.NumX,
		NumStates2: initData.NumX2,
		NumDOF1:    initData.NumX - initData.NumX2,
		NumDOF2:    initData.NumX - initData.NumX2/2,
		NumInputs:  initData.NumU,
		NumOutputs: initData.NumY,
		Azimuth:    mat.NewVecDense(numSteps, nil),
		Omega:      mat.NewVecDense(numSteps, nil),
		OmegaDot:   mat.NewVecDense(numSteps, nil),
		WindSpeed:  mat.NewVecDense(numSteps, nil),
		A:          make([]*mat.Dense, numSteps),
	}

	// If number of states is greater than zero
	if md.NumStates > 0 {

		// Create slice of vectors
		md.OpX = make([]*mat.VecDense, numSteps)
		for i := range md.OpX {
			md.OpX[i] = mat.NewVecDense(md.NumStates, nil)
		}

		// Create slice of vectors
		md.OpXd = make([]*mat.VecDense, numSteps)
		for i := range md.OpXd {
			md.OpXd[i] = mat.NewVecDense(md.NumStates, nil)
		}

		md.AvgA = mat.NewDense(md.NumStates, md.NumStates, nil)
		md.AvgOpX = mat.NewVecDense(md.NumStates, nil)
		md.AvgOpXd = mat.NewVecDense(md.NumStates, nil)

		// Find blade triplets
		md.Rotation.TripletsStates2 = findBladeTriplets(initData.X[:md.NumDOF2])
		md.Rotation.TripletsStates1 = findBladeTriplets(initData.X[:md.NumDOF1])

		// Find permutations
		md.Rotation.PermuteStates2, err = tripletsToPermutations(md.NumDOF2, md.Rotation.TripletsStates2)
		if err != nil {
			return nil, err
		}

		md.Rotation.PermuteStates1, err = tripletsToPermutations(md.NumDOF1, md.Rotation.TripletsStates1)
		if err != nil {
			return nil, err
		}
	}

	numBlades := 3

	numFixFrameStates2 := md.NumDOF2 - len(md.Rotation.TripletsStates2)*numBlades
	numFixFrameStates1 := md.NumDOF1 - len(md.Rotation.TripletsStates1)*numBlades
	numFixFrameInputs := md.NumInputs - len(md.Rotation.TripletsInputs)*numBlades
	numFixFrameoutputs := md.NumOutputs - len(md.Rotation.TripletsOutputs)*numBlades

	// Get state permutation slice
	permuteStates := md.Rotation.PermuteStates2
	for _, v := range md.Rotation.PermuteStates2 {
		permuteStates = append(permuteStates, v+md.NumDOF2)
	}
	for _, v := range md.Rotation.PermuteStates1 {
		permuteStates = append(permuteStates, v+md.NumStates2)
	}
	P := mat.NewDense(len(permuteStates), len(permuteStates), nil)
	P.Permutation(len(permuteStates), permuteStates)

	// Loop through linearization data
	for i, ld := range linData {

		// Rotor speed in radians/sec and rotor speed squared
		omega := ld.RotorSpeed
		omega2 := omega * omega
		omegaDot := 0.0

		md.Omega.SetVec(i, omega)
		md.OmegaDot.SetVec(i, omegaDot)
		md.Azimuth.SetVec(i, ld.Azimuth*180/math.Pi)
		md.WindSpeed.SetVec(i, ld.WindSpeed)

		cosCol := make([]float64, numBlades)
		sinCol := make([]float64, numBlades)
		cosColNeg := make([]float64, numBlades)
		sinColNeg := make([]float64, numBlades)

		for i := 0; i < numBlades; i++ {
			az := ld.Azimuth + 2*math.Pi*float64(i)/float64(numBlades)
			sinCol[i], cosCol[i] = math.Sincos(az)
			sinColNeg[i] = -sinCol[i]
			cosColNeg[i] = -cosCol[i]
		}

		ones := NewOnesVec(numBlades)

		// Eq. 9, t_tilde
		tt := mat.NewDense(3, numBlades, nil)
		tt.SetCol(0, ones)
		tt.SetCol(1, cosCol)
		tt.SetCol(2, sinCol)

		// t_tilde inverse
		c1, c2, c3 := cosCol[0], cosCol[1], cosCol[2]
		s1, s2, s3 := sinCol[0], sinCol[1], sinCol[2]
		ttv := mat.NewDense(3, 3, []float64{
			c2*s3 - s2*c3, c3*s1 - s3*c1, c1*s2 - s1*c2,
			s2 - s3, s3 - s1, s1 - s2,
			c3 - c2, c1 - c3, c2 - c1,
		})
		ttv.Scale(1/(1.5*math.Sqrt(3)), ttv)

		// Eq. 16 a, t_tilde_2
		tt2 := mat.NewDense(3, numBlades, nil)
		tt2.SetCol(1, sinColNeg)
		tt2.SetCol(2, cosCol)

		// Eq. 16 b, t_tilde_3
		tt3 := mat.NewDense(3, numBlades, nil)
		tt3.SetCol(1, cosColNeg)
		tt3.SetCol(2, sinColNeg)

		// Eq. 11 for second-order states only
		T1 := blockDiag(eye(numFixFrameStates2),
			Repeat(tt, len(md.Rotation.TripletsStates2))...)

		// Inverse of T1
		T1v := blockDiag(eye(numFixFrameStates2),
			Repeat(ttv, len(md.Rotation.TripletsStates2))...)

		// Eq. 14  for second-order states only
		T2 := blockDiag(mat.NewDense(numFixFrameStates2, numFixFrameStates2, nil),
			Repeat(tt2, len(md.Rotation.TripletsStates2))...)
		T2_omega := &mat.Dense{}
		T2_omega.Scale(omega, T2)
		T2_2omega := &mat.Dense{}
		T2_2omega.Scale(2*omega, T2)

		// Eq. 11 for first-order states (eq. 8 in MBC3 Update document)
		T1q := blockDiag(eye(numFixFrameStates1),
			Repeat(tt, len(md.Rotation.TripletsStates1))...)

		// Inverse of T1q
		T1qv := blockDiag(eye(numFixFrameStates1),
			Repeat(ttv, len(md.Rotation.TripletsStates1))...)

		// Eq. 14 for first-order states (eq.  9 in MBC3 Update document)
		T2q_omega := &mat.Dense{}
		if numFixFrameStates1 > 0 {
			T2q := blockDiag(mat.NewDense(numFixFrameStates1, numFixFrameStates1, nil),
				Repeat(tt2, len(md.Rotation.TripletsStates1))...)
			T2q_omega.Scale(omega, T2q)
		}

		// Eq. 15
		T3 := blockDiag(mat.NewDense(numFixFrameStates2, numFixFrameStates2, nil),
			Repeat(tt3, len(md.Rotation.TripletsStates2))...)
		T3_omega2 := &mat.Dense{}
		T3_omega2.Scale(omega2, T3)

		T1c := &mat.Dense{}
		if numFixFrameInputs > 0 {
			T1c = blockDiag(eye(numFixFrameInputs),
				Repeat(tt, len(md.Rotation.TripletsInputs))...)
		}

		// Inverse of T1q
		T1ov := &mat.Dense{}
		if numFixFrameoutputs > 0 {
			T1ov = blockDiag(eye(numFixFrameoutputs),
				Repeat(ttv, len(md.Rotation.TripletsOutputs))...)
		}

		_ = T1c
		_ = T1ov

		// Copy A matrix from linearization data
		A := &mat.Dense{}
		A.CloneFrom(ld.A)
		A.Mul(P, A)
		A.Mul(A, P)

		L := blockDiag(T1, T1, T1q)
		L.Slice(0, md.NumDOF2, md.NumDOF2, md.NumStates2).(*mat.Dense).Scale(omega, T2)

		R := blockDiag(T2_omega, T2_2omega, T2q_omega)
		tmp1, tmp2 := &mat.Dense{}, &mat.Dense{}
		tmp1.Scale(omega2, T3)
		tmp2.Scale(omegaDot, T2)
		R.Slice(0, md.NumDOF2, md.NumDOF2, md.NumStates2).(*mat.Dense).Add(tmp1, tmp2)

		AL := &mat.Dense{}
		AL.Mul(A, L)
		ALmR := &mat.Dense{}
		ALmR.Sub(AL, R)

		ANR := &mat.Dense{}
		ANR.Mul(blockDiag(T1v, T1v, T1qv), ALmR)
		ANR.Mul(P, ANR)
		ANR.Mul(ANR, P)

		md.A[i] = ANR

		if false {
			toCSV(tt, "mat-tt.csv")
			toCSV(ttv, "mat-ttv.csv")
			toCSV(tt2, "mat-tt2.csv")
			toCSV(tt3, "mat-tt3.csv")
			toCSV(T1, "mat-T1.csv")
			toCSV(T1v, "mat-T1v.csv")
			toCSV(T2, "mat-T2.csv")
			toCSV(T1q, "mat-T1q.csv")
			toCSV(T1qv, "mat-T1qv.csv")
			toCSV(T1q, "mat-T1q.csv")
			toCSV(T2q_omega, "mat-T2q.csv")
			toCSV(T3, "mat-T3.csv")
			toCSV(T1c, "mat-T1c.csv")
			toCSV(T1ov, "mat-T1ov.csv")
			toCSV(L, "mat-L.csv")
			toCSV(R, "mat-R.csv")
			toCSV(A, "mat-A.csv")
			toCSV(ANR, "mat-ANR.csv")
		}

		// for j, opx := range ld.X {
		// 	md.OpX[i].SetVec(j, opx.OperPoint)
		// }

		// for j, opxd := range ld.Xd {
		// 	md.OpXd[i].SetVec(j, opxd.OperPoint)
		// }

		// Reorder A, opx, opxd
	}

	// Average the A matrix
	for _, A := range md.A {
		md.AvgA.Add(md.AvgA, A)
	}
	md.AvgA.Scale(1/float64(len(md.A)), md.AvgA)

	// Average X operating points
	for _, op := range md.OpX {
		md.AvgOpX.AddVec(md.AvgOpX, op)
	}
	md.AvgOpX.ScaleVec(1/float64(len(md.OpX)), md.AvgOpX)

	// Average Xd operating points
	for _, op := range md.OpXd {
		md.AvgOpXd.AddVec(md.AvgOpXd, op)
	}
	md.AvgOpXd.ScaleVec(1/float64(len(md.OpXd)), md.AvgOpXd)

	// Eigenvalue/eigenvector analysis
	eig := mat.Eigen{}
	if ok := eig.Factorize(md.AvgA, mat.EigenRight); !ok {
		return nil, fmt.Errorf("error computing eigenvalues")
	}
	eigvecs := &mat.CDense{}
	eig.VectorsTo(eigvecs)

	// Eigenvector columns to keep based on degrees of freedom
	vecRows := []int{}
	for i := 0; i < md.NumDOF2; i++ {
		vecRows = append(vecRows, i)
	}
	for i := 2*md.NumDOF2 + 1; i < md.NumStates; i++ {
		vecRows = append(vecRows, i)
	}

	// Collect mode results
	for i, ev := range eig.Values(nil) {
		if imag(ev) > 0 {

			evAbs := cmplx.Abs(ev)

			// Create mode
			mode := &ModeResults{
				EigenValue:     ev,
				NaturalFreqRaw: evAbs,
				NaturalFreqHz:  evAbs / (2 * math.Pi),
				DampedFreqRaw:  imag(ev),
				DampedFreqHz:   imag(ev) / (2 * math.Pi),
				DampingRatio:   -real(ev) / evAbs,
				EigenVector:    make([]complex128, len(vecRows)),
				Magnitudes:     make([]float64, len(vecRows)),
				Phases:         make([]float64, len(vecRows)),
				Shape:          make([]float64, len(vecRows)),
			}

			// Extract relevant eigenvector values
			for j, r := range vecRows {
				v := eigvecs.At(r, i)
				mode.EigenVector[j] = v
				mode.Magnitudes[j] = cmplx.Abs(v)
				mode.Phases[j] = cmplx.Phase(v) * 180 / math.Pi
			}

			// Normalize magnitudes to get mode shape
			maxMag := mode.Magnitudes[0]
			for _, m := range mode.Magnitudes {
				if math.Abs(m) > math.Abs(maxMag) {
					maxMag = m
				}
			}
			for j, m := range mode.Magnitudes {
				mode.Shape[j] = m / maxMag
			}

			// Add mode to slice of modes
			md.Modes = append(md.Modes, mode)
		}
	}

	return md, nil
}

type ModeResults struct {
	EigenValue     complex128
	NaturalFreqRaw float64
	NaturalFreqHz  float64
	DampedFreqRaw  float64
	DampedFreqHz   float64
	DampingRatio   float64
	EigenVector    []complex128
	Magnitudes     []float64
	Phases         []float64
	Shape          []float64
}

func tripletsToPermutations(ndof int, triplets [][]int) ([]int, error) {

	tripletDOFs := map[int]struct{}{}
	tripletsPerms := make([]int, 0, len(triplets)*3)
	for _, triplet := range triplets {
		if len(triplet) != 3 {
			return nil, fmt.Errorf("number of values in triplet must be 3: %v", triplet)
		}
		for _, rc := range triplet {
			tripletDOFs[rc] = struct{}{}
			tripletsPerms = append(tripletsPerms, rc-1)
		}
	}

	// Create slice of permutations
	permutations := make([]int, 0, ndof)

	// Add permutations for DOFs not in triplets
	for i := 1; i <= ndof; i++ {
		if _, ok := tripletDOFs[i]; !ok {
			permutations = append(permutations, i-1)
		}
	}

	// Add permutations for triplet DOFs
	permutations = append(permutations, tripletsPerms...)

	return permutations, nil
}

// Regular expressions to find blades in operating point descriptions
var bladeRe = []*regexp.Regexp{
	regexp.MustCompile(`(?i)blade\s+\d`),
	regexp.MustCompile(`(?i)blade root \d`),
	regexp.MustCompile(`(?i)PitchBearing\d`),
	regexp.MustCompile(`(?i)BD_\d`),
	regexp.MustCompile(`(?i)BD\d`),
}

func findBladeTriplets(opd []OperPointData) [][]int {

	// Create local copy of operating point data containing only rotating frame points
	opl := []OperPointData{}
	for _, op := range opd {
		if op.IsRotating {
			opl = append(opl, op)
		}
	}

	// Sort operating points by description
	sort.SliceStable(opl, func(i, j int) bool {
		return opl[i].Desc < opl[j].Desc
	})

	// Find triplets based on descriptions
	prefix := ""
	triplets := [][]int{}
	for _, op := range opl {
		if prefix != "" && strings.HasPrefix(op.Desc, prefix) {
			triplets[0] = append(triplets[0], op.RC)
			continue
		}
		desc := op.Desc
		if i := strings.Index(desc, "("); i != -1 {
			if j := strings.LastIndex(desc, ")"); j != -1 {
				desc = desc[:i] + desc[j+1:]
			}
		}
		for _, re := range bladeRe {
			if loc := re.FindStringIndex(desc); loc != nil {
				prefix = desc[:loc[0]]
				triplets = append([][]int{{op.RC}}, triplets...)
				break
			}
		}
	}

	// Sort triplets by first number
	sort.Slice(triplets, func(i, j int) bool {
		return triplets[i][0] < triplets[j][0]
	})

	return triplets
}

func NewOnesVec(n int) []float64 {
	s := make([]float64, n)
	for i := range s {
		s[i] = 1
	}
	return s
}

func eye(n int) *mat.Dense {
	if n == 0 {
		return &mat.Dense{}
	}
	d := make([]float64, n*n)
	for i := 0; i < n*n; i += n + 1 {
		d[i] = 1
	}
	return mat.NewDense(n, n, d)
}

func blockDiag(base *mat.Dense, other ...*mat.Dense) *mat.Dense {

	mats := append([]*mat.Dense{base}, other...)

	size := 0
	for _, m := range mats {
		_, c := m.Dims()
		size += c
	}

	if size == 0 {
		return &mat.Dense{}
	}

	M := mat.NewDense(size, size, nil)
	c := 0
	for _, m := range mats {
		_, cm := m.Dims()
		if cm > 0 {
			M.Slice(c, c+cm, c, c+cm).(*mat.Dense).Copy(m)
			c += cm
		}
	}

	return M
}

func Repeat[T any](item T, n int) []T {
	s := make([]T, n)
	for i := range s {
		s[i] = item
	}
	return s
}

func toCSV(m mat.Matrix, path string) error {
	buf := &bytes.Buffer{}
	r, c := m.Dims()
	// fmt.Fprintf(buf, "%d,%d\n", r, c)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			if j > 0 {
				buf.WriteString(",")
			}
			fmt.Fprintf(buf, "%.16e", m.At(i, j))
		}
		buf.WriteString("\n")
	}
	return os.WriteFile(path, buf.Bytes(), 0777)
}
