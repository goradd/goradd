package json

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
