package db

import (
	"log"
	"strings"
)

// An AliasNode allows reference to a prior aliased operation later in a query.

type AliasNodeI interface {
	NodeI
}

type aliasNode struct {
	Node
}


// AliasNode returns an aliasNode type, which allows you to refer to a prior created named alias operation.
func AliasNode(goName string) *aliasNode {
	return &aliasNode {
		Node: Node {
			alias: goName,
		},
	}
}

func (n *aliasNode) nodeType() NodeType {
	return ALIAS_NODE
}

func (n *aliasNode) tableName() string {
	return ""
}

func (n *aliasNode) Equals(n2 NodeI) bool {
	if n2.nodeType() == ALIAS_NODE {
		return n.GetAlias() == n2.GetAlias()
	}

	return false
}

func (n *aliasNode) log(level int) {
	tabs := strings.Repeat("\t", level)
	log.Print(tabs + "Alias: " + n.GetAlias())
}


// Return the name as a captialized object name
func (n *aliasNode) objectName() string {
	return n.GetAlias()
}