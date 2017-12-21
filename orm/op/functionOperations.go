package op

import . 	"github.com/spekary/goradd/orm/db"



// Function lets you use any custom function which your database driver supports. Simply tell it the name and give
// it the arguments, and the result of the function will be used in the query.
func Function(funcName string, args... interface{}) *OperationNode {
	return NewFunctionNode(funcName, args...)
}

func Round(n... interface{}) *OperationNode {
	return NewFunctionNode("ROUND", n...)
}

func Abs(n... interface{}) *OperationNode {
	return NewFunctionNode("ABS", n...)
}

func Ceil(n... interface{}) *OperationNode {
	return NewFunctionNode("CEIL", n...)
}

func Floor(n... interface{}) *OperationNode {
	return NewFunctionNode("FLOOR", n...)
}

func Exp(n... interface{}) *OperationNode {
	return NewFunctionNode("EXP", n...)
}

func Ln(n interface{}) *OperationNode {
	return NewFunctionNode("LN", n)
}

func Power(n... interface{}) *OperationNode {
	return NewFunctionNode("POWER", n...)
}

func Sqrt(n interface{}) *OperationNode {
	return NewFunctionNode("SQRT", n)
}
