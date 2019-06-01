package op

import . "github.com/goradd/goradd/pkg/orm/query"

func StartsWith(arg1 interface{}, arg2 string) *OperationNode {
	return NewOperationNode(OpStartsWith, arg1, arg2)
}

func EndsWith(arg1 interface{}, arg2 string) *OperationNode {
	return NewOperationNode(OpEndsWith, arg1, arg2)
}

func Contains(arg1 interface{}, arg2 string) *OperationNode {
	return NewOperationNode(OpContains, arg1, arg2)
}
