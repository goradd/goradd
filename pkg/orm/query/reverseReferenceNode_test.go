package query

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReverseReferenceNodeInterfaces(t *testing.T) {
	n := NewReverseReferenceNode("db", "table", "dbKey", "dbCol", "goName", "table2", "col2", false)

	n.SetAlias("alias")

	assert.Equal(t, "alias", n.GetAlias())
}
