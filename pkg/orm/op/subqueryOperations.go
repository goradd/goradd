package op

import . "github.com/spekary/goradd/pkg/orm/query"

func Subquery(b QueryBuilderI) *SubqueryNode {
	return NewSubqueryNode(b)
}
