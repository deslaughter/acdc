package gui_test

import (
	"acdc/gui"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	if err := gui.Run(os.DirFS(".")); err != nil {
		t.Fatal(err)
	}
}
