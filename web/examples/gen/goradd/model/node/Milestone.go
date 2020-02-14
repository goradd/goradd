package node

// Code generated by goradd. DO NOT EDIT.

import (
	"bytes"
	"encoding/gob"

	"github.com/goradd/goradd/pkg/orm/query"
)

type milestoneNode struct {
	query.ReferenceNodeI
}

func Milestone() *milestoneNode {
	n := milestoneNode{
		query.NewTableNode("goradd", "milestone", "Milestone"),
	}
	query.SetParentNode(&n, nil)
	return &n
}

func (n *milestoneNode) SelectNodes_() (nodes []*query.ColumnNode) {
	nodes = append(nodes, n.ID())
	nodes = append(nodes, n.ProjectID())
	nodes = append(nodes, n.Name())
	return nodes
}
func (n *milestoneNode) PrimaryKeyNode() *query.ColumnNode {
	return n.ID()
}
func (n *milestoneNode) EmbeddedNode_() query.NodeI {
	return n.ReferenceNodeI
}
func (n *milestoneNode) Copy_() query.NodeI {
	return &milestoneNode{query.CopyNode(n.ReferenceNodeI)}
}

// ID represents the id column in the database.
func (n *milestoneNode) ID() *query.ColumnNode {
	cn := query.NewColumnNode(
		"goradd",
		"milestone",
		"id",
		"ID",
		query.ColTypeString,
		true,
	)
	query.SetParentNode(cn, n)
	return cn
}

// ProjectID represents the project_id column in the database.
func (n *milestoneNode) ProjectID() *query.ColumnNode {
	cn := query.NewColumnNode(
		"goradd",
		"milestone",
		"project_id",
		"ProjectID",
		query.ColTypeString,
		false,
	)
	query.SetParentNode(cn, n)
	return cn
}

// Project represents the link to the Project object.
func (n *milestoneNode) Project() *projectNode {
	cn := &projectNode{
		query.NewReferenceNode(
			"goradd",
			"milestone",
			"project_id",
			"ProjectID",
			"Project",
			"project",
			"id",
			false,
		),
	}
	query.SetParentNode(cn, n)
	return cn
}

// Name represents the name column in the database.
func (n *milestoneNode) Name() *query.ColumnNode {
	cn := query.NewColumnNode(
		"goradd",
		"milestone",
		"name",
		"Name",
		query.ColTypeString,
		false,
	)
	query.SetParentNode(cn, n)
	return cn
}

type milestoneNodeEncoded struct {
	RefNode query.ReferenceNodeI
}

func (n *milestoneNode) GobEncode() (data []byte, err error) {
	var buf bytes.Buffer
	e := gob.NewEncoder(&buf)

	s := milestoneNodeEncoded{
		RefNode: n.ReferenceNodeI,
	}

	if err = e.Encode(s); err != nil {
		panic(err)
	}
	data = buf.Bytes()
	return
}

func (n *milestoneNode) GobDecode(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	var s milestoneNodeEncoded
	if err = dec.Decode(&s); err != nil {
		panic(err)
	}
	n.ReferenceNodeI = s.RefNode
	query.SetParentNode(n, query.ParentNode(n)) // Reinforce types
	return
}

func init() {
	gob.RegisterName("milestoneNode2", &milestoneNode{})
}
