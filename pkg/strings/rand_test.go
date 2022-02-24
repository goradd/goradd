package strings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordString(t *testing.T) {
	assert.Panics(t, func() {
		PasswordString(0)
	})
	assert.Panics(t, func() {
		PasswordString(3)
	})

	assert.Equal(t, 10, len(PasswordString(10)))
}

func TestRandomString(t *testing.T) {
	assert.Equal(t, "", RandomString("a", 0))
	assert.Equal(t, "a", RandomString("a", 1))
	assert.Equal(t, 1, len(RandomString("abc", 1)))
	assert.Equal(t, 10, len(RandomString("abc", 10)))
}
