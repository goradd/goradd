package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRefresh(t *testing.T) {
	js := Refresh("a").RenderScript(RenderParams{})
	assert.Equal(t, `goradd.postAjax({"controlID":"","eventId":0,"actionValues":{"event":eventData}});`, js)
}

func TestSetControlValue(t *testing.T) {
	js := SetControlValue("a", "b", "c").RenderScript(RenderParams{})
	assert.Equal(t, `goradd.setControlValue("a","b","c");`, js)
}
