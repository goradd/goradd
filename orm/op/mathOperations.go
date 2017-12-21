package op

import . 	"github.com/spekary/goradd/orm/db"



func Add(args... interface{}) *OperationNode {
	return NewOperationNode(OpAdd, args...)
}

func Subtract(args... interface{}) *OperationNode {
	return NewOperationNode(OpSubtract, args...)
}

func Multiply(args... interface{}) *OperationNode {
	return NewOperationNode(OpMultiply, args...)
}

func Divide(args... interface{}) *OperationNode {
	return NewOperationNode(OpDivide, args...)
}

func Mod(args... interface{}) *OperationNode {
	return NewOperationNode(OpModulo, args...)
}

func Negative(n interface{}) *OperationNode {
	return NewOperationNode(OpNegate, n)
}

