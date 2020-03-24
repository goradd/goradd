package query

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAliasNodeInterfaces(t *testing.T) {
	n := Alias("test")

	assert.Implements(t, (*AliasNodeI)(nil), n)
	assert.Equal(t, "test", n.GetAlias())
}
