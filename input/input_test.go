package input_test

import (
	"acdc/input"
	"encoding/json"
	"os"
	"testing"
)

func TestReadInputFiles(t *testing.T) {

	model, err := input.ReadFiles("testdata/Model_GE-1.5SLE/GE-1.5SLE.fst")
	if err != nil {
		t.Fatal(err)
	}

	// if act, exp := inp.Fast.EDFile, "AOC_WSt_ElastoDyn.dat"; act != exp {
	// 	t.Fatalf("Fast.EDFile = %v, expected %v", act, exp)
	// }

	b, err := json.MarshalIndent(model, "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile("testdata/model.json", b, 0777)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWriteFiles(t *testing.T) {

	inp, err := input.ReadFiles("testdata/Model_GE-1.5SLE/GE-1.5SLE.fst")
	if err != nil {
		t.Fatal(err)
	}

	if err := inp.WriteFiles("testdata/write_test/ge_1p5.fst"); err != nil {
		t.Fatal(err)
	}
}
