package op

import . "github.com/goradd/goradd/pkg/orm/query"

func Subquery(b QueryBuilderI) *SubqueryNode {
	return NewSubqueryNode(b)
}
