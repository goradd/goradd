package query

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReferenceNodeInterfaces(t *testing.T) {
	n := NewReferenceNode("db", "table", "dbCol", "goCol", "goName", "table2", "col2", false, ColTypeInteger)
	n.SetAlias("alias")

	assert.Implements(t, (*ReferenceNodeI)(nil), n)
	assert.Equal(t, "alias", n.GetAlias())
}
