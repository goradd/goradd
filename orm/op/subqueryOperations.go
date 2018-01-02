package op

import . 	"github.com/spekary/goradd/orm/query"

func Subquery(b QueryBuilderI) *SubqueryNode {
	return NewSubqueryNode(b)
}
