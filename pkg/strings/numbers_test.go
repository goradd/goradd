package strings

import (
	"fmt"
	"testing"
)

func ExampleExtractNumbers() {
	a := ExtractNumbers("a1b2 c3")
	fmt.Println(a)
	//Output: 123
}

func TestExtractNumbers(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		wantOut string
	}{
		{"empty", "", ""},
		{"123", "123", "123"},
		{"abc", "abc", ""},
		{"a1c", "a1c", "1"},
		{"a1c", "a1c", "1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := ExtractNumbers(tt.in); gotOut != tt.wantOut {
				t.Errorf("ExtractNumbers() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
