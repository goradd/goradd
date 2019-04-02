package query

import (
	"log"
	"strings"
)

// A SubqueryNode represents a "select" subquery. Subqueries are not always portable to other databases, and are not
// easily checked for syntax errors, since a subquery can return a scalar, vector, or even an entire table.
// You generally do not create a subquery node directly, but rather you use the codegenerated models to start a
// query on a table, and then end the query with "Subquery()" which will turn the query into a usable subquery node
// that you can embed in other queries.
type SubqueryNode struct {
	nodeAlias
	b QueryBuilderI
}

// NewSubqueryNode creates a new subquery
func NewSubqueryNode(b QueryBuilderI) *SubqueryNode {
	n := &SubqueryNode{
		b: b,
	}
	return n
}

func (n *SubqueryNode) nodeType() NodeType {
	return SubqueryNodeType
}

// Equals is used internally by the framework to determine if two nodes are equal
func (n *SubqueryNode) Equals(n2 NodeI) bool {
	if cn, ok := n2.(*SubqueryNode); ok {
		return cn.b == n.b
	}
	return false
}

/*
func (n *SubqueryNode) containedNodes() (nodes []NodeI) {
	nodes = append(nodes, n) // Return the subquery node itself, because we need to do some work on it

	// must expand the returned nodes one more time
	for _,n2 := range n.b.nodes() {	// Refers back to db package, so do this differently
		if cn,_ := n2.(nodeContainer); cn != nil {
			nodes = append(nodes, cn.containedNodes()...)
		} else {
			nodes = append(nodes, n2)
		}
	}
	return nodes
}
*/

func (n *SubqueryNode) tableName() string {
	return ""
}

func (n *SubqueryNode) log(level int) {
	tabs := strings.Repeat("\t", level)
	log.Print(tabs + "Subquery: ")
}


// SubqueryBuilder is used internally by the framework to return the internal query builder of the subquery
func SubqueryBuilder(n *SubqueryNode) QueryBuilderI {
	return n.b
}
