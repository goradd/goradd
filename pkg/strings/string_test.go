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

func TestIndent(t *testing.T) {
	if Indent("a\nb\nc") != "\ta\n\tb\n\tc" {
		t.Fail()
	}
	if Indent("\na\nb\nc") != "\t\n\ta\n\tb\n\tc" {
		t.Fail()
	}
	if Indent("a\nb\nc\n") != "\ta\n\tb\n\tc\n" {
		t.Fail()
	}
}

func TestKebabToCamel(t *testing.T) {
	if KebabToCamel("ab-cd-ef") != "AbCdEf" {
		t.Fail()
	}
}

