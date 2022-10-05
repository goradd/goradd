package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedirect(t *testing.T) {
	js := Redirect("http://a/b/c").RenderScript(RenderParams{})
	assert.Equal(t, `goradd.redirect("http://a/b/c");`, js)
}
