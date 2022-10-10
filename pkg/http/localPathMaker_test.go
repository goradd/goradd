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
		{"empty local", "", "here/test", "here/test"},
		{"with proxy", "/proxy", "/test", "/proxy/test"},
		{"dir", "/proxy", "/test/", "/proxy/test/"},
		{"with proxy local", "/proxy", "here/test", "here/test"},
		{"dir local", "/proxy", "here/test/", "here/test/"},
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
