package node

// Code generated by goradd. DO NOT EDIT.

import (
	"bytes"
	"encoding/gob"

	"github.com/goradd/goradd/pkg/orm/query"
)

type employeeInfoNode struct {
	query.ReferenceNodeI
}

func EmployeeInfo() *employeeInfoNode {
	n := employeeInfoNode{
		query.NewTableNode("goradd", "employee_info", "EmployeeInfo"),
	}
	query.SetParentNode(&n, nil)
	return &n
}

func (n *employeeInfoNode) SelectNodes_() (nodes []*query.ColumnNode) {
	nodes = append(nodes, n.ID())
	nodes = append(nodes, n.PersonID())
	nodes = append(nodes, n.EmployeeNumber())
	return nodes
}
func (n *employeeInfoNode) PrimaryKeyNode() *query.ColumnNode {
	return n.ID()
}
func (n *employeeInfoNode) EmbeddedNode_() query.NodeI {
	return n.ReferenceNodeI
}
func (n *employeeInfoNode) Copy_() query.NodeI {
	return &employeeInfoNode{query.CopyNode(n.ReferenceNodeI)}
}

// ID represents the id column in the database.
func (n *employeeInfoNode) ID() *query.ColumnNode {
	cn := query.NewColumnNode(
		"goradd",
		"employee_info",
		"id",
		"ID",
		query.ColTypeString,
		true,
	)
	query.SetParentNode(cn, n)
	return cn
}

// PersonID represents the person_id column in the database.
func (n *employeeInfoNode) PersonID() *query.ColumnNode {
	cn := query.NewColumnNode(
		"goradd",
		"employee_info",
		"person_id",
		"PersonID",
		query.ColTypeString,
		false,
	)
	query.SetParentNode(cn, n)
	return cn
}

// Person represents the link to the Person object.
func (n *employeeInfoNode) Person() *personNode {
	cn := &personNode{
		query.NewReferenceNode(
			"goradd",
			"employee_info",
			"person_id",
			"PersonID",
			"Person",
			"person",
			"id",
			false,
			query.ColTypeString,
		),
	}
	query.SetParentNode(cn, n)
	return cn
}

// EmployeeNumber represents the employee_number column in the database.
func (n *employeeInfoNode) EmployeeNumber() *query.ColumnNode {
	cn := query.NewColumnNode(
		"goradd",
		"employee_info",
		"employee_number",
		"EmployeeNumber",
		query.ColTypeInteger,
		false,
	)
	query.SetParentNode(cn, n)
	return cn
}

type employeeInfoNodeEncoded struct {
	RefNode query.ReferenceNodeI
}

func (n *employeeInfoNode) GobEncode() (data []byte, err error) {
	var buf bytes.Buffer
	e := gob.NewEncoder(&buf)

	s := employeeInfoNodeEncoded{
		RefNode: n.ReferenceNodeI,
	}

	if err = e.Encode(s); err != nil {
		panic(err)
	}
	data = buf.Bytes()
	return
}

func (n *employeeInfoNode) GobDecode(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	var s employeeInfoNodeEncoded
	if err = dec.Decode(&s); err != nil {
		panic(err)
	}
	n.ReferenceNodeI = s.RefNode
	query.SetParentNode(n, query.ParentNode(n)) // Reinforce types
	return
}

func init() {
	gob.RegisterName("employeeInfoNode2", &employeeInfoNode{})
}
