package query

import (
	"fmt"
	"github.com/spekary/goradd/datetime"
	"log"
	"reflect"
	"strings"
	"time"
)

type ValueNode struct {
	Node

	value interface{}
}

// Shortcut for converting a constant value to a node
func Value(i interface{}) NodeI {
	return NewValueNode(i)
}

func NewValueNode(i interface{}) NodeI {
	n := &ValueNode{
		value: i,
	}
	switch v := i.(type) {
	// do nothings
	case string:
	case int:
	case uint:
	case uint64:
	case int64:
	case float64:
	case time.Time:

		// casts
	case []byte:
		n.value = string(v[:])
	case uint8:
		n.value = uint(v)
	case uint16:
		n.value = uint(v)
	case uint32:
		n.value = uint(v)
	case int8:
		n.value = int(v)
	case int16:
		n.value = int(v)
	case int32:
		n.value = int(v)
	case float32:
		n.value = float64(v)
	case datetime.DateTime:
		n.value = v.Time
	default:
		// Arrays of items
		if reflect.TypeOf(v).Kind() == reflect.Slice || reflect.TypeOf(v).Kind() == reflect.Array {
			ary := []NodeI{}
			s := reflect.ValueOf(v)
			for i := 0; i < s.Len(); i++ {
				// TODO: Handle NodeI's here too? Prevent more than one level deep?
				ary = append(ary, NewValueNode(s.Index(i).Interface()))
			}
			n.value = ary
		} else {
			panic("Can't use this type as a value node.")
		}
	}
	return n
}

func (n *ValueNode) nodeType() NodeType {
	return VALUE_NODE
}

func (n *ValueNode) Equals(n2 NodeI) bool {
	if cn, ok := n2.(*ValueNode); ok {
		return cn.value == n.value
	}
	return false
}

func (n *ValueNode) tableName() string {
	return ""
}

func (n *ValueNode) log(level int) {
	tabs := strings.Repeat("\t", level)
	var alias string
	if n.alias != "" {
		alias = " as " + n.alias
	}
	log.Print(tabs + "Val: " + fmt.Sprint(n.value) + alias)
}

func ValueNodeGetValue(n *ValueNode) interface{} {
	return n.value
}
