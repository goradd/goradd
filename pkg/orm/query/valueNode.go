package query

import (
	"fmt"
	"github.com/goradd/goradd/pkg/datetime"
	"log"
	"reflect"
	"strings"
	"time"
)

// ValueNode represents a value for a built-in type that is to be used in a query.
type ValueNode struct {
	value interface{}
}

// Shortcut for converting a constant value to a node
func Value(i interface{}) NodeI {
	return NewValueNode(i)
}

// NewValueNode returns a new ValueNode that wraps the given value.
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
	case float32:
	case time.Time:

		// casts
	case []byte:
		n.value = string(v[:])
	case datetime.DateTime:
		n.value = v.GoTime()
	case nil:
		panic("You cannot use nil as an operator. If you are testing for a NULL, use the IsNull function.")
	default:
		// Use reflection to do various conversions
		typ := reflect.TypeOf(v)
		k := typ.Kind()
		val := reflect.ValueOf(v)

		switch k {
		case reflect.Int:fallthrough
		case reflect.Int8:fallthrough
		case reflect.Int16:fallthrough
		case reflect.Int32:fallthrough
		case reflect.Int64:
			n.value = int(val.Int())
		case reflect.Uint:fallthrough
		case reflect.Uint8:fallthrough
		case reflect.Uint16:fallthrough
		case reflect.Uint32:fallthrough
		case reflect.Uint64:
			n.value = uint(val.Uint())
		case reflect.Bool:
			n.value = val.Bool()
		case reflect.Float32:
			// converting float32 to float64 might cause problems in the final sql statement, so we leave the type as float32
			n.value = float32(val.Float())
		case reflect.Float64:
			n.value = val.Float()
		case reflect.Slice:fallthrough
		case reflect.Array:
			var ary []NodeI
			for i := 0; i < val.Len(); i++ {
				// TODO: Handle NodeI's here too? Prevent more than one level deep?
				ary = append(ary, NewValueNode(val.Index(i).Interface()))
			}
			n.value = ary
		case reflect.String:
			n.value = val.String()
		default:
			panic("Can't use this type as a value node.")
		}
	}
	return n
}

func (n *ValueNode) Equals(n2 NodeI) bool {
	if cn, ok := n2.(*ValueNode); ok {
		if an2, ok := cn.value.([]NodeI); ok {
			if an1, ok := n.value.([]NodeI); !ok {
				return false
			} else if len(an2) != len(an1) {
				return false
			} else {
				for i,n := range an1 {
					if !n.Equals(an2[i]) {
						return false
					}
				}
			}
			return true
		}
		return cn.value == n.value
	}
	return false
}

func (n *ValueNode) tableName() string {
	return ""
}

func (n *ValueNode) log(level int) {
	tabs := strings.Repeat("\t", level)
	log.Print(tabs + "Val: " + fmt.Sprint(n.value))
}

// ValueNodeGetValue is used internally by the framework to get the node's internal value.
func ValueNodeGetValue(n *ValueNode) interface{} {
	return n.value
}

func (n *ValueNode) nodeType() NodeType {
	return ValueNodeType
}
