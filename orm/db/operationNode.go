package db

import (
	"log"
	"strings"
)

type Operator string

const (
	// Standard logical operators
	OpEqual Operator = "="
	OpNotEqual Operator = "<>"
	OpAnd = "AND"
	OpOr = "OR"
	OpXor = "XOR"
	OpGreater = ">"
	OpGreaterEqual = ">="
	OpLess = "<"
	OpLessEqual = "<="

	// Unary logical
	OpNot = "NOT"

	// Math operators
	OpAdd = "+"
	OpSubtract = "-"
	OpMultiply = "*"
	OpDivide = "/"
	OpModulo = "%"

	// Unary math
	OpNegate = " -"

	// Bit operators
	OpBitAnd = "&"
	OpBitOr = "|"
	OpBitXor = "^"
	OpShiftLeft = "<<"
	OpShiftRight = ">>"

	// Unary bit
	OpBitInvert = "~"

	// Function operator
	// The function name is followed by the operators in parenthesis
	OpFunc = "func"

	// SQL functions that act like operators in that the operator is put in between the operands
	OpLike = "LIKE"
	OpIn = "IN"
	OpNotIn = "NOT IN"

	// Special NULL tests
	OpNull = "NULL"
	OpNotNull = "NOT NULL"
)

func (o Operator) String() string {
	return string(o)
}

// An operation is a general purpose structure that specs an operation on a node or group of nodes
// The operation could be arithmetic, boolean, or a function.
type OperationNode struct {
	Node
	op Operator
	operands []NodeI
	functionName string // for function operations specific to the db driver
	isAggregate bool // requires that an aggregation clause be present in the query
	sortDescending bool
	distinct bool // some aggregate queries, particularly count, allow this inside the function
}

func NewOperationNode (op Operator, operands... interface{}) *OperationNode {
	n := &OperationNode {
		op:       op,
	}
	n.assignOperands(operands...)
	return n
}

func NewFunctionNode (functionName string,operands... interface{}) *OperationNode {
	n := &OperationNode {
		op:       OpFunc,
		functionName: functionName,
	}
	n.assignOperands(operands...)
	return n
}

// NewCountNode creates a Count function node. If no operands are given, it will use * as the parameter to the function
// which means it will count nulls. To NOT count nulls, at least one column name needs to be specified.
func NewCountNode(operands... NodeI) *OperationNode {
	n := &OperationNode {
		op:       OpFunc,
		functionName: "COUNT",
	}
	for _,op := range operands {
		n.operands = append(n.operands, op)
	}

	return n
}

// process the list of operands at run time, making sure all static values are escaped
func (n *OperationNode) assignOperands(operands... interface{}) {
	var op interface{}

	for _,op = range operands {
		if ni,ok := op.(NodeI); ok {
			n.operands = append(n.operands, ni)
		} else {
			n.operands = append(n.operands, NewValueNode(op))
		}
	}
}

func (n *OperationNode) Ascending() *OperationNode {
	n.sortDescending = false
	return n
}

func (n *OperationNode) Descending() *OperationNode {
	n.sortDescending = true
	return n
}

func (n *OperationNode) Distinct() *OperationNode {
	n.distinct = true
	return n
}

func (n *OperationNode) sortDesc() bool {
	return n.sortDescending
}

func (n *OperationNode) nodeType() NodeType {
	return OPERATION_NODE
}


func (n *OperationNode) Equals(n2 NodeI) bool {
	if cn,ok := n2.(*OperationNode); ok {
		if cn.op != n.op {
			return false
		}
		if cn.functionName != n.functionName {
			return false
		}
		if cn.isAggregate != n.isAggregate {
			return false
		}
		if cn.sortDescending != n.sortDescending {
			return false
		}
		if cn.operands == nil && n.operands == nil {
			return true // neither side has operands, so no need to check further
		}
		if len(cn.operands) != len(n.operands) {
			return false
		}

		for i,o := range n.operands {
			if !o.Equals(cn.operands[i]) {
				return false
			}
		}
		return true
	}
	return false
}

func (n *OperationNode) containedNodes() (nodes []NodeI) {
	for _,op := range n.operands {
		if nc,ok := op.(nodeContainer); ok {
			nodes = append(nodes, nc.containedNodes()...)
		} else {
			nodes = append(nodes,op)
		}
	}
	return
}

func (n *OperationNode) tableName() string {
	return ""
}



func (n *OperationNode) log(level int) {
	tabs := strings.Repeat("\t", level)
	log.Print(tabs + "Op: " + n.op.String())
}