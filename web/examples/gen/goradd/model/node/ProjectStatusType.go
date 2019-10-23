package node

// Code generated by goradd. DO NOT EDIT.

import (
    "encoding/gob"
	"github.com/goradd/goradd/pkg/orm/query"
)

type projectStatusTypeNode struct {
	query.ReferenceNodeI
}


func (n *projectStatusTypeNode) SelectNodes_() (nodes []*query.ColumnNode) {
	nodes = append(nodes, n.ID())
	nodes = append(nodes, n.Name())
	nodes = append(nodes, n.Description())
	nodes = append(nodes, n.Guidelines())
	nodes = append(nodes, n.IsActive())
	return nodes
}

func (n *projectStatusTypeNode) PrimaryKeyNode_() (*query.ColumnNode) {
	return n.ID()
}

func (n *projectStatusTypeNode) EmbeddedNode_() query.NodeI {
	return n.ReferenceNodeI
}

func (n *projectStatusTypeNode) Copy_() query.NodeI {
	return &projectStatusTypeNode{query.CopyNode(n.ReferenceNodeI)}
}

func (n *projectStatusTypeNode) ID() *query.ColumnNode {

	cn := query.NewColumnNode (
		"goradd",
		"project_status_type",
		"id",
		"ID",
		query.ColTypeUnsigned,
		true,
	)
	query.SetParentNode(cn, n)
	return cn
}
func (n *projectStatusTypeNode) Name() *query.ColumnNode {

	cn := query.NewColumnNode (
		"goradd",
		"project_status_type",
		"name",
		"Name",
		query.ColTypeString,
		false,
	)
	query.SetParentNode(cn, n)
	return cn
}
func (n *projectStatusTypeNode) Description() *query.ColumnNode {

	cn := query.NewColumnNode (
		"goradd",
		"project_status_type",
		"description",
		"Description",
		query.ColTypeString,
		false,
	)
	query.SetParentNode(cn, n)
	return cn
}
func (n *projectStatusTypeNode) Guidelines() *query.ColumnNode {

	cn := query.NewColumnNode (
		"goradd",
		"project_status_type",
		"guidelines",
		"Guidelines",
		query.ColTypeString,
		false,
	)
	query.SetParentNode(cn, n)
	return cn
}
func (n *projectStatusTypeNode) IsActive() *query.ColumnNode {

	cn := query.NewColumnNode (
		"goradd",
		"project_status_type",
		"is_active",
		"IsActive",
		query.ColTypeBool,
		false,
	)
	query.SetParentNode(cn, n)
	return cn
}
func init() {
   gob.RegisterName("projectStatusTypeNode2", &projectStatusTypeNode{})
}