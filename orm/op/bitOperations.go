package op

import . 	"github.com/spekary/goradd/orm/query"




func BitAnd(arg1, arg2 interface{}) *OperationNode {
	return NewOperationNode(OpBitAnd, arg1, arg2)
}

func BitOr(arg1, arg2 interface{}) *OperationNode {
	return NewOperationNode(OpBitOr, arg1, arg2)
}

func BitXor(arg1, arg2 interface{}) *OperationNode {
	return NewOperationNode(OpBitXor, arg1, arg2)
}

func BitShiftLeft(arg1, arg2 interface{}) *OperationNode {
	return NewOperationNode(OpShiftLeft, arg1, arg2)
}

func BitShiftRight(arg1, arg2 interface{}) *OperationNode {
	return NewOperationNode(OpShiftRight, arg1, arg2)
}

func BitInvert(n interface{}) *OperationNode {
	return NewOperationNode(OpBitInvert, n)
}
