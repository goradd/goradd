package op

import . 	"github.com/spekary/goradd/orm/db"



func Subquery(b QueryBuilderI) *SubqueryNode {
	return NewSubqueryNode(b)
}
