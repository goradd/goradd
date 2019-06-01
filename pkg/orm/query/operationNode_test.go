package query

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOperationNodeInterfaces(t *testing.T) {
	n := NewOperationNode(OpAdd, 4, 5)

	assert.Implements(t, (*OperationNodeI)(nil), n)
}
