package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfirm(t *testing.T) {
	js := Confirm("a").RenderScript(RenderParams{})
	assert.Equal(t, `if (!window.confirm("a")) return false;`, js)
}
