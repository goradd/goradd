package op

import (
	. "github.com/goradd/goradd/pkg/orm/query"
	"reflect"
)

func Equal(arg1 interface{}, arg2 interface{}) *OperationNode {
	return NewOperationNode(OpEqual, arg1, arg2)
}

func NotEqual(arg1 interface{}, arg2 interface{}) *OperationNode {
	return NewOperationNode(OpNotEqual, arg1, arg2)
}

func GreaterThan(arg1 interface{}, arg2 interface{}) *OperationNode {
	return NewOperationNode(OpGreater, arg1, arg2)
}

func GreaterOrEqual(arg1 interface{}, arg2 interface{}) *OperationNode {
	return NewOperationNode(OpGreaterEqual, arg1, arg2)
}

func LessThan(arg1 interface{}, arg2 interface{}) *OperationNode {
	return NewOperationNode(OpLess, arg1, arg2)
}

func LessOrEqual(arg1 interface{}, arg2 interface{}) *OperationNode {
	return NewOperationNode(OpLessEqual, arg1, arg2)
}

func And(args ...interface{}) *OperationNode {
	return NewOperationNode(OpAnd, args...)
}

func Or(args ...interface{}) *OperationNode {
	return NewOperationNode(OpOr, args...)
}

func Xor(arg1, arg2 interface{}) *OperationNode {
	return NewOperationNode(OpXor, arg1, arg2)
}

func Not(n interface{}) *OperationNode {
	return NewOperationNode(OpNot, n)
}

// All is a placeholder for when you need to return something that represents selecting everything
func All() *OperationNode {
	return NewOperationNode(OpAll)
}

func None() *OperationNode {
	return NewOperationNode(OpNone)
}

func Like(n interface{}, pattern string) *OperationNode {
	return NewOperationNode(OpLike, n, NewValueNode(pattern))
}

// In tests to see if the given node is in the "what" list
func In(n NodeI, what ...interface{}) *OperationNode {
	k := reflect.TypeOf(what[0]).Kind()
	if k == reflect.Array || k == reflect.Slice {
		// first item is array-like. Since that doesn't make sense in this context, we will assume
		// we are trying to send in an array of something rather than list them out.
		return NewOperationNode(OpIn, n, what[0]) // by passing an array as what here, it will cause the values to be output as a list
	}
	return NewOperationNode(OpIn, n, what) // by passing an array as what here, it will cause the values to be output as a list
}

func NotIn(n NodeI, what ...interface{}) *OperationNode {
	k := reflect.TypeOf(what[0]).Kind()
	if k == reflect.Array || k == reflect.Slice {
		// first item is array-like. Since that doesn't make sense in this context, we will assume
		// we are trying to send in an array of something rather than list them out.
		return NewOperationNode(OpNotIn, n, what[0]) // by passing an array as what here, it will cause the values to be output as a list
	}
	return NewOperationNode(OpNotIn, n, what) // by passing an array as what here, it will cause the values to be output as a list
}

func IsNull(n interface{}) *OperationNode {
	return NewOperationNode(OpNull, n)
}

func IsNotNull(n interface{}) *OperationNode {
	return NewOperationNode(OpNotNull, n)
}
