package util

import (
	"testing"
	"fmt"
	"math/rand"
)

func TestRandomHtmlValueString(t *testing.T) {
	rand.Seed(1)	// reset random seed

	s := RandomHtmlValueString(40);
	fmt.Printf(s)

	if len(s) != 40 {
		t.Error("Wrong size")
	}
}