package db

import (
	"strings"
	"log"
)

// An operation is a general purpose structure that specs an operation on a node or group of nodes
// The operation could be arithmetic, boolean, or a function.
type SubqueryNode struct {
	Node
	b QueryBuilderI
}

func NewSubqueryNode(b QueryBuilderI) *SubqueryNode {
	n := &SubqueryNode {
		b: b,
	}
	return n
}


func (n *SubqueryNode) nodeType() NodeType {
	return SUBQUERY_NODE
}


func (n *SubqueryNode) Equals(n2 NodeI) bool {
	if cn,ok := n2.(*SubqueryNode); ok {
		return cn.b == n.b
	}
	return false
}

func (n *SubqueryNode) containedNodes() (nodes []NodeI) {
	return n.b.nodes()
}

func (n *SubqueryNode) tableName() string {
	return ""
}


func (n *SubqueryNode) log(level int) {
	tabs := strings.Repeat("\t", level)
	log.Print(tabs + "Subquery: ")
}



