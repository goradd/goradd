package http

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseValueAndParams(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name       string
		arg        string
		wantValue  string
		wantParams map[string]string
	}{
		{"empty", "", "", nil},
		{"value", "test/test", "test/test", nil},
		{"valueWithSemi", "test/test;", "test/test", nil},
		{"value with param", "test/test;a=b", "test/test", map[string]string{"a": "b"}},
		{"value with 2 params", "test/test; a=b; c=d", "test/test", map[string]string{"a": "b", "c": "d"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotParams := ParseValueAndParams(tt.arg)
			assert.Equalf(t, tt.wantValue, gotValue, "ParseValueAndParams(%v)", tt.arg)
			assert.Equalf(t, tt.wantParams, gotParams, "ParseValueAndParams(%v)", tt.arg)
		})
	}
}

func TestParseAuthorizationHeader(t *testing.T) {

	tests := []struct {
		name       string
		arg        string
		wantScheme string
		wantParams string
	}{
		{"empty", "", "", ""},
		{"one item", "abc", "abc", ""},
		{"one item with whitespace", "abc  ", "abc", ""},
		{"two items", "abc def", "abc", "def"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotScheme, gotParams := ParseAuthorizationHeader(tt.arg)
			assert.Equalf(t, tt.wantScheme, gotScheme, "ParseAuthorizationHeader(%v)", tt.arg)
			assert.Equalf(t, tt.wantParams, gotParams, "ParseAuthorizationHeader(%v)", tt.arg)
		})
	}
}
