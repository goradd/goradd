package javascript_test

import (
	"encoding/json"
	"testing"

	"github.com/goradd/gengen/pkg/maps"
	. "github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/goradd/pkg/time"
	"github.com/stretchr/testify/assert"
)

func TestToJavaScript(t *testing.T) {
	m1 := maps.NewSliceMap()
	m1.Set("a", `Hi "`)
	m1.Set("b", NoQuoteKey{JsCode{"There"}})
	m1.Set("c", 4)

	m2 := maps.NewStringSliceMap()
	m2.Set("a", `Hi "`)
	m2.Set("b", "There")
	m2.Set("c", "4")

	m3 := maps.NewSliceMap()

	tests := []struct {
		name string
		arg  interface{}
		want string
	}{
		{"JavaScripter", JsCode{"Test"}, "Test"},
		{"Undefined", Undefined{}, "undefined"},
		{"String", `Hal"s /super/ \fine`, `"Hal\"s /super/ \\fine"`},
		{"String Slice", []string{`a / ' b`, `C & "D"`}, `["a / ' b","C \u0026 \"D\""]`},
		{"Interface Slice", []interface{}{"Hi", JsCode{"There"}}, `["Hi",There]`},
		{"Interface String Map", map[string]interface{}{"a": `Hi "`, "b": NoQuoteKey{JsCode{"There"}}, "c": 4}, `{"a":"Hi \"",b:There,"c":4}`},
		{"Interface Int Map", map[int]interface{}{1: `Hi "`, 2: JsCode{"There"}, 3: 4}, `{1:"Hi \"",2:There,3:4}`},
		{"MapI", m1, `{"a":"Hi \"",b:There,"c":4}`},
		{"Empty map", map[string]interface{}{}, `{}`},
		{"Empty int map", map[int]interface{}{}, `{}`},
		{"Empty MapI", m3, `{}`},
		{"StringMapI", m2, `{"a":"Hi \"","b":"There","c":"4"}`},
		{"Int", 1, `1`},
		{"Null", nil, `null`},
		{"NewFunctionCall Arguments", Arguments([]interface{}{3, "me"}), `3,"me"`},
		{"Time", time.NewDateTime(2022, 3, 2, 6, 5, 7, 0), `new Date(1646201107000)`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToJavaScript(tt.arg); got != tt.want {
				t.Errorf("ToJavaScript() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumberInt(t *testing.T) {
	tests := []struct {
		name string
		arg  interface{}
		want int
	}{
		{"json.Number", json.Number("1"), 1},
		{"String", "2", 2},
		{"Bad value", 2, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NumberInt(tt.arg); got != tt.want {
				t.Errorf("NumberInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumberFloat(t *testing.T) {
	tests := []struct {
		name string
		arg  interface{}
		want float64
	}{
		{"json.Number", json.Number("1"), 1},
		{"json.Number2", json.Number("3.33"), 3.33},
		{"String", "2.1", 2.1},
		{"Bad value", 2, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NumberFloat(tt.arg); got != tt.want {
				t.Errorf("NumberFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumberString(t *testing.T) {
	tests := []struct {
		name string
		arg  interface{}
		want string
	}{
		{"json.Number", json.Number("1"), "1"},
		{"json.Number2", json.Number("3.33"), "3.33"},
		{"String", "2.1", "2.1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NumberString(tt.arg); got != tt.want {
				t.Errorf("NumberString() = %v, want %v", got, tt.want)
			}
		})
	}

	assert.Panics(t, func() {
		NumberString(4)
	})
}

func TestUndefined_MarshalJSON(t *testing.T) {
	c := Undefined{}
	a := []interface{}{c}
	b, err := json.Marshal(a)
	assert.NoError(t, err)
	_ = b
}
