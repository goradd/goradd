package strings

import (
	"fmt"
	"testing"
)

func ExampleKebabToCamel() {
	a := KebabToCamel("abc-def")
	fmt.Println(a)
	//Output: AbcDef
}

func TestKebabToCamel(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"-a", "-a", "A"},
		{"b-a", "b-a", "BA"},
		{"123-abc", "123-abc", "123Abc"},
		{"this-that", "this-that", "ThisThat"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KebabToCamel(tt.s); got != tt.want {
				t.Errorf("KebabToCamel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleSnakeToKebab() {
	a := SnakeToKebab("abc_def")
	fmt.Println(a)
	//Output: abc-def
}

func TestSnakeToKebab(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"a_b", "a_b", "a-b"},
		{"-b", "-b", "-b"},
		{"_b", "_b", "-b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SnakeToKebab(tt.s); got != tt.want {
				t.Errorf("SnakeToKebab() = %v, want %v", got, tt.want)
			}
		})
	}
}
