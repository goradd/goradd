package query

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestColumnNodeInterfaces(t *testing.T) {
	n := NewColumnNode("db", "table", "dbName", "goName", ColTypeString, true)
	n.SetAlias("alias")

	//var i ColumnNodeI = n

	assert.Implements(t, (*ColumnNodeI)(nil), n)
	assert.Equal(t, "alias", n.GetAlias())
}
