// Code generated by GoRADD. DO NOT EDIT.

package node

import (
	"bytes"
	"encoding/gob"

	"github.com/goradd/goradd/pkg/orm/query"
)

type loginNode struct {
	query.ReferenceNodeI
}

func Login() *loginNode {
	n := loginNode{
		query.NewTableNode("goradd", "login", "Login"),
	}
	query.SetParentNode(&n, nil)
	return &n
}

func (n *loginNode) SelectNodes_() (nodes []*query.ColumnNode) {
	nodes = append(nodes, n.ID())
	nodes = append(nodes, n.PersonID())
	nodes = append(nodes, n.Username())
	nodes = append(nodes, n.Password())
	nodes = append(nodes, n.IsEnabled())
	return nodes
}
func (n *loginNode) PrimaryKeyNode() *query.ColumnNode {
	return n.ID()
}
func (n *loginNode) EmbeddedNode_() query.NodeI {
	return n.ReferenceNodeI
}
func (n *loginNode) Copy_() query.NodeI {
	return &loginNode{query.CopyNode(n.ReferenceNodeI)}
}

// ID represents the id column in the database.
func (n *loginNode) ID() *query.ColumnNode {
	cn := query.NewColumnNode(
		"goradd",
		"login",
		"id",
		"ID",
		query.ColTypeString,
		true,
	)
	query.SetParentNode(cn, n)
	return cn
}

// PersonID represents the person_id column in the database.
func (n *loginNode) PersonID() *query.ColumnNode {
	cn := query.NewColumnNode(
		"goradd",
		"login",
		"person_id",
		"PersonID",
		query.ColTypeString,
		false,
	)
	query.SetParentNode(cn, n)
	return cn
}

// Person represents the link to the Person object.
func (n *loginNode) Person() *personNode {
	cn := &personNode{
		query.NewReferenceNode(
			"goradd",
			"login",
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

// Username represents the username column in the database.
func (n *loginNode) Username() *query.ColumnNode {
	cn := query.NewColumnNode(
		"goradd",
		"login",
		"username",
		"Username",
		query.ColTypeString,
		false,
	)
	query.SetParentNode(cn, n)
	return cn
}

// Password represents the password column in the database.
func (n *loginNode) Password() *query.ColumnNode {
	cn := query.NewColumnNode(
		"goradd",
		"login",
		"password",
		"Password",
		query.ColTypeString,
		false,
	)
	query.SetParentNode(cn, n)
	return cn
}

// IsEnabled represents the is_enabled column in the database.
func (n *loginNode) IsEnabled() *query.ColumnNode {
	cn := query.NewColumnNode(
		"goradd",
		"login",
		"is_enabled",
		"IsEnabled",
		query.ColTypeBool,
		false,
	)
	query.SetParentNode(cn, n)
	return cn
}

type loginNodeEncoded struct {
	RefNode query.ReferenceNodeI
}

func (n *loginNode) GobEncode() (data []byte, err error) {
	var buf bytes.Buffer
	e := gob.NewEncoder(&buf)

	s := loginNodeEncoded{
		RefNode: n.ReferenceNodeI,
	}

	if err = e.Encode(s); err != nil {
		panic(err)
	}
	data = buf.Bytes()
	return
}

func (n *loginNode) GobDecode(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	var s loginNodeEncoded
	if err = dec.Decode(&s); err != nil {
		panic(err)
	}
	n.ReferenceNodeI = s.RefNode
	query.SetParentNode(n, query.ParentNode(n)) // Reinforce types
	return
}

func init() {
	gob.RegisterName("loginNode2", &loginNode{})
}