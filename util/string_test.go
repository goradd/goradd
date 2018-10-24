package util

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestRandomHtmlValueString(t *testing.T) {
	rand.Seed(1) // reset random seed

	s := RandomHtmlValueString(40)
	fmt.Printf(s)

	if len(s) != 40 {
		t.Error("Wrong size")
	}
}

func TestEndsWithString(t *testing.T) {
	if !EndsWith(".45", ".45") {
		t.Fail()
	}
	if !EndsWith("a", "a") {
		t.Fail()
	}
	if !EndsWith("234f asd fa", "a") {
		t.Fail()
	}
	if !EndsWith("asdfsaf sdabc", "abc") {
		t.Fail()
	}
	if EndsWith("bc", "abc") {
		t.Fail()
	}
	if EndsWith("", "abc") {
		t.Fail()
	}
}

