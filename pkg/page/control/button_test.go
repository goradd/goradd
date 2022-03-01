package control

import (
	"github.com/goradd/goradd/pkg/page"
	"testing"
)

func TestNewButton(t *testing.T) {
	form := page.NewMockForm()
	b := NewButton(form, "btnId")
	_ = b
}