package query

import (
	"bytes"
	"encoding/gob"
	"github.com/goradd/goradd/pkg/log"
	"strings"
)

type OperationNodeI interface {
	nodeContainer
	Aliaser
}

// Operator is used internally by the framework to specify an operation to be performed by the database.
// Not all databases can perform all the operations. It will be up to the database driver to sort this out.
type Operator string

const (
	// Standard logical operators

	OpEqual        Operator = "="
	OpNotEqual     Operator = "<>"
	OpAnd                   = "AND"
	OpOr                    = "OR"
	OpXor                   = "XOR"
	OpGreater               = ">"
	OpGreaterEqual          = ">="
	OpLess                  = "<"
	OpLessEqual             = "<="

	// Unary logical
	OpNot = "NOT"

	OpAll  = "1=1"
	OpNone = "1=0"

	// Math operators
	OpAdd      = "+"
	OpSubtract = "-"
	OpMultiply = "*"
	OpDivide   = "/"
	OpModulo   = "%"

	// Unary math
	OpNegate = " -"

	// Bit operators
	OpBitAnd     = "&"
	OpBitOr      = "|"
	OpBitXor     = "^"
	OpShiftLeft  = "<<"
	OpShiftRight = ">>"

	// Unary bit
	OpBitInvert = "~"

	// Function operator
	// The function name is followed by the operators in parenthesis
	OpFunc = "func"

	// SQL functions that act like operators in that the operator is put in between the operands
	OpLike  = "LIKE" // This is very SQL specific and may not be supported in NoSql
	OpIn    = "IN"
	OpNotIn = "NOT IN"

	// Special NULL tests
	OpNull    = "NULL"
	OpNotNull = "NOT NULL"

	// Our own custom operators for universal support
	OpStartsWith     = "StartsWith"
	OpEndsWith       = "EndsWith"
	OpContains       = "Contains"
	OpDateAddSeconds = "AddSeconds" // Adds the given number of seconds to a datetime
)

// String returns a string representation of the Operator type. For convenience, this also corresponds to the SQL
// representation of an operator
func (o Operator) String() string {
	return string(o)
}

// An OperationNode is a general purpose structure that specifies an operation on a node or group of nodes.
// The operation could be arithmetic, boolean, or a function.
type OperationNode struct {
	nodeAlias
	op           Operator
	operands     []NodeI
	functionName string // for function operations specific to the db driver
	distinct     bool   // some aggregate queries, particularly count, allow this inside the function
}

// NewOperationNode returns a new operation.
func NewOperationNode(op Operator, operands ...interface{}) *OperationNode {
	n := &OperationNode{
		op: op,
	}
	n.assignOperands(operands...)
	return n
}

// NewFunctionNode returns an operation node that executes a database function.
func NewFunctionNode(functionName string, operands ...interface{}) *OperationNode {
	n := &OperationNode{
		op:           OpFunc,
		functionName: functionName,
	}
	n.assignOperands(operands...)
	return n
}

// NewCountNode creates a Count function node. If no operands are given, it will use * as the parameter to the function
// which means it will count nulls. To NOT count nulls, at least one table name needs to be specified.
func NewCountNode(operands ...NodeI) *OperationNode {
	n := &OperationNode{
		op:           OpFunc,
		functionName: "COUNT",
	}
	for _, op := range operands {
		n.operands = append(n.operands, op)
	}

	return n
}

func (n *OperationNode) nodeType() NodeType {
	return OperationNodeType
}

// assignOperands processes the list of operands at run time, making sure all static values are escaped
func (n *OperationNode) assignOperands(operands ...interface{}) {
	var op interface{}

	if operands != nil {
		for _, op = range operands {
			if ni, ok := op.(NodeI); ok {
				n.operands = append(n.operands, ni)
			} else {
				n.operands = append(n.operands, NewValueNode(op))
			}
		}
	}
}

/*
func (n *OperationNode) Ascending() *OperationNode {
	n.sortDescending = false
	return n
}

func (n *OperationNode) Descending() *OperationNode {
	n.sortDescending = true
	return n
}
*/

// Distinct sets the operation to return distinct results
func (n *OperationNode) Distinct() *OperationNode {
	n.distinct = true
	return n
}

/*
func (n *OperationNode) sortDesc() bool {
	return n.sortDescending
}
*/

// Equals is used internally by the framework to tell if two nodes are equal
func (n *OperationNode) Equals(n2 NodeI) bool {
	if cn, ok := n2.(*OperationNode); ok {
		if cn.op != n.op {
			return false
		}
		if cn.functionName != n.functionName {
			return false
		}
		/*
			if cn.isAggregate != n.isAggregate {
				return false
			}
		*/
		/*
			if cn.sortDescending != n.sortDescending {
				return false
			}*/
		if cn.operands == nil && n.operands == nil {
			return true // neither side has operands, so no need to check further
		}
		if len(cn.operands) != len(n.operands) {
			return false
		}

		for i, o := range n.operands {
			if !o.Equals(cn.operands[i]) {
				return false
			}
		}
		return true
	}
	return false
}

func (n *OperationNode) containedNodes() (nodes []NodeI) {
	for _, op := range n.operands {
		if nc, ok := op.(nodeContainer); ok {
			nodes = append(nodes, nc.containedNodes()...)
		} else {
			nodes = append(nodes, op)
		}
	}
	return
}

func (n *OperationNode) tableName() string {
	return ""
}

func (n *OperationNode) databaseKey() string {
	return ""
}

func (n *OperationNode) log(level int) {
	tabs := strings.Repeat("\t", level)
	log.FrameworkDebug(tabs + "Op: " + n.op.String())
}

func (n *OperationNode) GobEncode() (data []byte, err error) {
	var buf bytes.Buffer
	e := gob.NewEncoder(&buf)

	if err = e.Encode(n.alias); err != nil {
		panic(err)
	}
	if err = e.Encode(n.op); err != nil {
		panic(err)
	}
	if err = e.Encode(n.operands); err != nil {
		panic(err)
	}
	if err = e.Encode(n.functionName); err != nil {
		panic(err)
	}
	if err = e.Encode(n.distinct); err != nil {
		panic(err)
	}
	data = buf.Bytes()
	return
}

func (n *OperationNode) GobDecode(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err = dec.Decode(&n.alias); err != nil {
		panic(err)
	}
	if err = dec.Decode(&n.op); err != nil {
		panic(err)
	}
	if err = dec.Decode(&n.operands); err != nil {
		panic(err)
	}
	if err = dec.Decode(&n.functionName); err != nil {
		panic(err)
	}
	if err = dec.Decode(&n.distinct); err != nil {
		panic(err)
	}
	return
}

func init() {
	gob.Register(&OperationNode{})
}

// OperationNodeOperator is used internally by the framework to get the operator.
func OperationNodeOperator(n *OperationNode) Operator {
	return n.op
}

// OperationNodeOperands is used internally by the framework to get the operands.
func OperationNodeOperands(n *OperationNode) []NodeI {
	return n.operands
}

// OperationNodeFunction is used internally by the framework to get the function.
func OperationNodeFunction(n *OperationNode) string {
	return n.functionName
}

// OperationNodeDistinct is used internally by the framework to get the distinct value.
func OperationNodeDistinct(n *OperationNode) bool {
	return n.distinct
}

/*
func OperationIsAggregate(n *OperationNode) bool {
	return n.isAggregate
}
*/
/*
func OperationSortDescending(n *OperationNode) bool {
	return n.sortDescending
}
*/
