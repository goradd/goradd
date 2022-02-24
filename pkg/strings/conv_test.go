package strings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAto(t *testing.T) {
	assert.Equal(t, uint(23), AtoUint("23"))
	assert.Equal(t, uint64(23), AtoUint64("23"))
	assert.Equal(t, int64(23), AtoInt64("23"))
}
