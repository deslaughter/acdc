package gui_test

import (
	"os"
	"testing"

	"github.com/deslaughter/acdc/gui"
)

func TestRun(t *testing.T) {
	if err := gui.Run(os.DirFS(".")); err != nil {
		t.Fatal(err)
	}
}
