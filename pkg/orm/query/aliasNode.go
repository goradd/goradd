package query

import (
	"bytes"
	"encoding/gob"
	"log"
	"strings"
)

type AliasNodeI interface {
	NodeI
	Aliaser
}

// An AliasNode is a reference to a prior aliased operation later in a query. An alias is a name given
// to a computed value.
type AliasNode struct {
	nodeAlias
}

// Alias returns an AliasNode type, which allows you to refer to a prior created named alias operation.
func Alias(goName string) *AliasNode {
	return &AliasNode{
		nodeAlias{
			alias: goName,
		},
	}
}

func (n *AliasNode) nodeType() NodeType {
	return AliasNodeType
}

func (n *AliasNode) tableName() string {
	return ""
}

func (n *AliasNode) databaseKey() string {
	return ""
}


// Equals returns true if the given node points to the same alias value as receiver.
func (n *AliasNode) Equals(n2 NodeI) bool {
	if a, ok := n2.(*AliasNode); ok {
		return a.alias == n.alias
	}
	return false
}

func (n *AliasNode) log(level int) {
	tabs := strings.Repeat("\t", level)
	log.Print(tabs + "Alias: " + n.alias)
}

// Return the name as a capitalized object name
/* I don't think an alias should have a go name
func (n *AliasNode) goName() string {
	return n.alias
}
*/

func (n *AliasNode) GobEncode() (data []byte, err error) {
	var buf bytes.Buffer
	e := gob.NewEncoder(&buf)

	if err = e.Encode(&n.alias); err != nil {
		panic(err)
	}
	data = buf.Bytes()
	return
}


func (n *AliasNode) GobDecode(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err = dec.Decode(&n.alias); err != nil {
		panic(err)
	}
	return
}

func init() {
	gob.Register(&AliasNode{})
}
