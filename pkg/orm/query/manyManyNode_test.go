package query

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestManyManyNodeInterfaces(t *testing.T) {
	n := NewManyManyNode("db", "table", "dbCol", "goName", "table2", "col2", "tablePk", false)
	n.SetAlias("alias")

	assert.Equal(t, "alias", n.GetAlias())
}
