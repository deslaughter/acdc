package anl_test

import (
	"testing"

	"github.com/deslaughter/acdc/anl"
)

func TestPerformMBC(t *testing.T) {

	turb := anl.Turbine{
		Name: "turb_01",
		Dir:  "../turb_01",
	}

	_, err := turb.PerformMBC()
	if err != nil {
		t.Fatal(err)
	}

}
