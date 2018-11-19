package html

import (
	"fmt"
	"testing"
)

func TestRandomHtmlValueString(t *testing.T) {
	//rand.Seed(1) // reset random seed

	s := RandomString(40)
	fmt.Printf(s + " ")

	if len(s) != 40 {
		t.Error("Wrong size")
	}
}
