package strings

import (
	"fmt"
	"math/rand"
	"testing"
)


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

