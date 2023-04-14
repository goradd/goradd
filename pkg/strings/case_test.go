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

func ExampleCamelToKebab() {
	a := CamelToKebab("AbcDef")
	fmt.Println(a)
	b := CamelToKebab("AbcDEFghi")
	fmt.Println(b)
	//Output: abc-def
	//abc-de-fghi
}

func ExampleCamelToSnake() {
	a := CamelToSnake("AbcDef")
	fmt.Println(a)
	b := CamelToSnake("AbcDEFghi")
	fmt.Println(b)
	//Output: abc_def
	//abc_de_fghi
}

func TestCamelToKebab(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"a", "a", "a"},
		{"A", "A", "a"},
		{"ab", "ab", "ab"},
		{"AB", "AB", "ab"},
		{"Ab", "Ab", "ab"},
		{"aB", "aB", "a-b"},
		{"Abc", "Abc", "abc"},
		{"AbC", "AbC", "ab-c"},
		{"ABc", "ABc", "a-bc"},
		{"a1b", "a1b", "ab"},
		{"ABC", "ABC", "abc"},
		{"ABCd", "ABCd", "ab-cd"},
		{"AbCdE", "ABCdE", "ab-cd-e"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CamelToKebab(tt.s); got != tt.want {
				t.Errorf("SnakeToKebab() = %v, want %v", got, tt.want)
			}
		})
	}
}
