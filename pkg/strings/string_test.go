package strings

import (
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

func TestStartsWithString(t *testing.T) {
	if !StartsWith(".45", ".45") {
		t.Fail()
	}
	if !StartsWith("a", "a") {
		t.Fail()
	}
	if !StartsWith("abc", "a") {
		t.Fail()
	}
	if StartsWith("234f asd fa", "a") {
		t.Fail()
	}
	if StartsWith("asdfsaf sdabc", "abc") {
		t.Fail()
	}
	if StartsWith("bc", "abc") {
		t.Fail()
	}
	if StartsWith("ab", "abc") {
		t.Fail()
	}
	if StartsWith("", "abc") {
		t.Fail()
	}
}
