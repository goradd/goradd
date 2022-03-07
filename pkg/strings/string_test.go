package strings

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestHasOnlyLetters(t *testing.T) {
	if HasOnlyLetters("a-b") {
		t.Fail()
	}
	if !HasOnlyLetters("abc") {
		t.Fail()
	}
	if HasOnlyLetters("123") {
		t.Fail()
	}
}

func TestLcFirst(t *testing.T) {
	assert.Equal(t, "", LcFirst(""))
	assert.Equal(t, "abcDef", LcFirst("AbcDef"))
}

func TestStartsWith(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name      string
		s         string
		beginning string
		want      bool
	}{
		{"same with dot", ".45", ".45", true},
		{"short", "a", "a", true},
		{"short2", "abc", "a", true},
		{"short3", "234f asd fa", "a", false},
		{"mid", "234f abc fa", "abc", false},
		{"smaller", "ab", "abc", false},
		{"smaller2", "abc", "ab", true},
		{"none", "", "abc", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StartsWith(tt.s, tt.beginning); got != tt.want {
				t.Errorf("StartsWith() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndsWith(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		ending string
		want   bool
	}{
		{"same", ".45", ".45", true},
		{"a", "a", "a", true},
		{"long", "234f asd fa", "a", true},
		{"long2", "asdfsaf sdabc", "abc", true},
		{"too short", "bc", "abc", false},
		{"empty", "", "abc", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EndsWith(tt.s, tt.ending); got != tt.want {
				t.Errorf("EndsWith() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleTitle() {
	a := Title("do_i_seeYou")
	fmt.Println(a)
	//Output: Do I See You
}

func TestTitle(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"i", "i", "I"},
		{"iJ", "iJ", "I J"},
		{"i_j", "i_j", "I J"},
		{"iJK", "iJK", "I J K"},
		{"i_J_k", "iJK", "I J K"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Title(tt.s); got != tt.want {
				t.Errorf("Title() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleJoinContent() {
	a := JoinContent("+", "this", "", "that")
	fmt.Println(a)
	//Output: this+that
}

func TestJoinContent(t *testing.T) {
	type args struct {
		sep   string
		items []string
	}
	tests := []struct {
		name  string
		sep   string
		items []string
		want  string
	}{
		{"empty", "", []string{""}, ""},
		{"1", "+", []string{"this"}, "this"},
		{"2", "+", []string{"this", "that"}, "this+that"},
		{"empty sep", "", []string{"this", "that"}, "thisthat"},
		{"empty item", "+", []string{"this", "", "that"}, "this+that"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JoinContent(tt.sep, tt.items...); got != tt.want {
				t.Errorf("JoinContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIf(t *testing.T) {
	assert.Equal(t, "a", If(true, "a", "b"))
	assert.Equal(t, "b", If(false, "a", "b"))
}

func TestContainsAnyStrings(t *testing.T) {
	tests := []struct {
		name     string
		haystack string
		needles  []string
		want     bool
	}{
		{"empty", "", []string{}, false},
		{"empty2", "", []string{"a", "b"}, false},
		{"a", "a", []string{"a", "b"}, true},
		{"b", "b", []string{"a", "b"}, true},
		{"abc", "abc", []string{"h", "bc"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsAnyStrings(tt.haystack, tt.needles...); got != tt.want {
				t.Errorf("ContainsAnyStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasCharType(t *testing.T) {
	type args struct {
		s          string
		wantUpper  bool
		wantLower  bool
		wantDigit  bool
		wantPunc   bool
		wantSymbol bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"lower", args{"a", false, true, false, false, false}, true},
		{"lowerFail", args{"A", false, true, false, false, false}, false},
		{"upper", args{"A", true, false, false, false, false}, true},
		{"upperFail", args{"a", true, false, true, false, false}, false},
		{"digit", args{"1", false, false, true, false, false}, true},
		{"digitFail", args{"A", false, false, true, false, false}, false},
		{"punc", args{",", false, false, false, true, false}, true},
		{"puncFail", args{"a", false, false, false, true, false}, false},
		{"symbol", args{"$", false, false, false, false, true}, true},
		{"symbolFail", args{",", false, false, false, false, true}, false},
		{"mult1", args{"aA", true, true, false, false, false}, true},
		{"mult1Fail", args{"a1", true, true, false, false, false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasCharType(tt.args.s, tt.args.wantUpper, tt.args.wantLower, tt.args.wantDigit, tt.args.wantPunc, tt.args.wantSymbol); got != tt.want {
				t.Errorf("HasCharType() = %v, want %v", got, tt.want)
			}
		})
	}
}

