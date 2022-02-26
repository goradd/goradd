package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActionGroup_GetCallbackAction(t *testing.T) {
	js := Group(AddClass("a", "b"), RemoveClass("c", "d")).RenderScript(RenderParams{})
	assert.Equal(t, `g$('a').class("+b");g$('c').class("-d");`, js)
}

func TestActionGroup_HasCallbackAction(t *testing.T) {
	assert.False(t, Group(AddClass("a", "b"), RemoveClass("c", "d")).HasCallbackAction())
}

func TestActionGroup_HasServerAction(t *testing.T) {
	assert.False(t, Group(AddClass("a", "b"), RemoveClass("c", "d")).HasServerAction())
}

func TestActionGroup_GetCallbackAction1(t *testing.T) {
	assert.Nil(t, Group().GetCallbackAction())
	assert.Nil(t, Group(AddClass("a", "b"), RemoveClass("c", "d")).GetCallbackAction())

}
