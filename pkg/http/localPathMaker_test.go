package http

import (
	"testing"

	"github.com/goradd/goradd/pkg/config"
)

func TestMakeLocalPath(t *testing.T) {
	tests := []struct {
		name      string
		proxyPath string
		path      string
		want      string
	}{
		{"empty", "", "/test", "/test"},
		{"wity proxy", "/proxy", "/test", "/proxy/test"},
		{"dir", "/proxy", "/test/", "/proxy/test/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.ProxyPath = tt.proxyPath
			if got := MakeLocalPath(tt.path); got != tt.want {
				t.Errorf("MakeLocalPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
