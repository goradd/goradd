// Code generated by GoRADD. DO NOT EDIT.

package node

import (
	"bytes"
	"encoding/gob"

	"github.com/goradd/goradd/pkg/orm/query"
)

type personWithLockNode struct {
	query.ReferenceNodeI
}

func PersonWithLock() *personWithLockNode {
	n := personWithLockNode{
		query.NewTableNode("goradd", "public.person_with_lock", "PersonWithLock"),
	}
	query.SetParentNode(&n, nil)
	return &n
}

func (n *personWithLockNode) SelectNodes_() (nodes []*query.ColumnNode) {
	nodes = append(nodes, n.ID())
	nodes = append(nodes, n.FirstName())
	nodes = append(nodes, n.LastName())
	nodes = append(nodes, n.SysTimestamp())
	return nodes
}
func (n *personWithLockNode) PrimaryKeyNode() *query.ColumnNode {
	return n.ID()
}
func (n *personWithLockNode) EmbeddedNode_() query.NodeI {
	return n.ReferenceNodeI
}
func (n *personWithLockNode) Copy_() query.NodeI {
	return &personWithLockNode{query.CopyNode(n.ReferenceNodeI)}
}

// ID represents the id column in the database.
func (n *personWithLockNode) ID() *query.ColumnNode {
	cn := query.NewColumnNode(
		"goradd",
		"public.person_with_lock",
		"id",
		"ID",
		query.ColTypeString,
		true,
	)
	query.SetParentNode(cn, n)
	return cn
}

// FirstName represents the first_name column in the database.
func (n *personWithLockNode) FirstName() *query.ColumnNode {
	cn := query.NewColumnNode(
		"goradd",
		"public.person_with_lock",
		"first_name",
		"FirstName",
		query.ColTypeString,
		false,
	)
	query.SetParentNode(cn, n)
	return cn
}

// LastName represents the last_name column in the database.
func (n *personWithLockNode) LastName() *query.ColumnNode {
	cn := query.NewColumnNode(
		"goradd",
		"public.person_with_lock",
		"last_name",
		"LastName",
		query.ColTypeString,
		false,
	)
	query.SetParentNode(cn, n)
	return cn
}

// SysTimestamp represents the sys_timestamp column in the database.
func (n *personWithLockNode) SysTimestamp() *query.ColumnNode {
	cn := query.NewColumnNode(
		"goradd",
		"public.person_with_lock",
		"sys_timestamp",
		"SysTimestamp",
		query.ColTypeTime,
		false,
	)
	query.SetParentNode(cn, n)
	return cn
}

type personWithLockNodeEncoded struct {
	RefNode query.ReferenceNodeI
}

func (n *personWithLockNode) GobEncode() (data []byte, err error) {
	var buf bytes.Buffer
	e := gob.NewEncoder(&buf)

	s := personWithLockNodeEncoded{
		RefNode: n.ReferenceNodeI,
	}

	if err = e.Encode(s); err != nil {
		panic(err)
	}
	data = buf.Bytes()
	return
}

func (n *personWithLockNode) GobDecode(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	var s personWithLockNodeEncoded
	if err = dec.Decode(&s); err != nil {
		panic(err)
	}
	n.ReferenceNodeI = s.RefNode
	query.SetParentNode(n, query.ParentNode(n)) // Reinforce types
	return
}

func init() {
	gob.RegisterName("personWithLockNode2", &personWithLockNode{})
}
