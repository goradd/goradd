package db

import (
	"context"
	"github.com/goradd/goradd/pkg/orm/op"
	. "github.com/goradd/goradd/pkg/orm/query"
	"github.com/goradd/maps"
)

// LimitInfo is the information needed to limit the rows being requested.
type LimitInfo struct {
	MaxRowCount int
	Offset      int
}

// QueryBuilder is a helper to implement the QueryBuilderI interface in various builder classes.
// It is designed to be embedded in a database specific implementation.
// It gathers the builder instructions as the query is built. It leaves the implementation
// of the functions that actually query a database -- Load, Delete, Count -- to the containing structure.
type QueryBuilder struct {
	Ctx           context.Context // The context that will be used in all the queries
	Joins         []NodeI
	OrderBys      []NodeI
	ConditionNode NodeI
	IsDistinct    bool
	AliasNodes    *AliasNodesType
	// Adds a COUNT(*) to the select list
	GroupBys   []NodeI
	Selects    []NodeI
	LimitInfo  *LimitInfo
	HavingNode NodeI
	IsSubquery bool
}

type AliasNodesType = maps.SliceMap[string, Aliaser]

// Init initializes the QueryBuilder.
func (b *QueryBuilder) Init(ctx context.Context) {
	b.Ctx = ctx
}

// Context returns the context.
func (b *QueryBuilder) Context() context.Context {
	return b.Ctx
}

// Join will attach the given reference node to the builder.
func (b *QueryBuilder) Join(n NodeI, condition NodeI) {
	// Possible TBD: If we ever want to support joining the same tables multiple
	// times with different conditions, we could use an alias to name each join. We would
	// then need to create an Alias node to specify which join is meant in different clauses.

	if b.Joins != nil {
		if !NodeIsReferenceI(n) {
			panic("you can only join Reference, ReverseReference and ManyManyReference nodes")
		}

		if NodeTableName(RootNode(n)) != NodeTableName(b.Joins[0]) {
			panic("you can only join nodes starting from the same table as the root node. This node must start from " + NodeTableName(b.Joins[0]))
		}
	}

	NodeSetCondition(n, condition)
	b.Joins = append(b.Joins, n)
}

// Alias adds a node that is given a manual alias name. This is usually some kind of operation, but it can be
// any query.Aliaser kind of node.
func (b *QueryBuilder) Alias(name string, n NodeI) {
	if b.AliasNodes == nil {
		b.AliasNodes = new(AliasNodesType)
	}
	a := n.(Aliaser)
	a.SetAlias(name)
	b.AliasNodes.Set(name, a)
}

// Expand expands an array type node so that it will produce individual rows instead of an array of items
func (b *QueryBuilder) Expand(n NodeI) {
	if typ := NodeGetType(n); !(typ == ReverseReferenceNodeType || typ == ManyManyNodeType) {
		panic("you can only expand a node that is a ReverseReference or ManyMany node.")
	} else {
		n.(Expander).Expand()
		b.Join(n, nil)
	}
}

// Condition adds the condition of the Where clause. If a condition already exists, it will be anded to the previous condition.
func (b *QueryBuilder) Condition(c NodeI) {
	if b.ConditionNode == nil {
		b.ConditionNode = c
	} else {
		b.ConditionNode = op.And(b.ConditionNode, c)
	}
}

// OrderBy adds the order by nodes. If these are table type nodes, the primary key of the table will be used.
// These nodes can be modified using Ascending and Descending calls.
func (b *QueryBuilder) OrderBy(nodes ...NodeI) {
	b.OrderBys = append(b.OrderBys, nodes...)
}

// Limit sets the limit parameters of what is returned.
func (b *QueryBuilder) Limit(maxRowCount int, offset int) {
	if b.LimitInfo != nil {
		panic("Query already has a limit")
	}
	b.LimitInfo = &LimitInfo{maxRowCount, offset}
}

// Select specifies what specific nodes are selected. This is an optimization in order to limit the amount
// of data returned by the query. Without this, the query will expand all the join items to return every
// column of each table joined.
func (b *QueryBuilder) Select(nodes ...NodeI) {
	if b.GroupBys != nil {
		panic("You cannot have Select and GroupBy statements in the same query. The GroupBy columns will automatically be selected.")
	}
	for _, n := range nodes {
		if NodeGetType(n) != ColumnNodeType {
			panic("you can only select column nodes")
		}
	}
	b.Selects = append(b.Selects, nodes...)
}

// Distinct sets the distinct bit, causing the query to not return duplicates.
func (b *QueryBuilder) Distinct() {
	b.IsDistinct = true
}

// GroupBy sets the nodes that are grouped. According to SQL rules, these then are the only nodes that can be
// selected, and they MUST be selected.
func (b *QueryBuilder) GroupBy(nodes ...NodeI) {
	if b.Selects != nil {
		panic("You cannot have Select and GroupBy statements in the same query. The GroupBy columns will automatically be selected.")
	}
	b.GroupBys = append(b.GroupBys, nodes...)
}

// Having adds a HAVING condition, which is a filter that acts on the results of a query.
// In particular its useful for filtering after aggregate functions have done their work.
func (b *QueryBuilder) Having(node NodeI) {
	b.HavingNode = node // should be a condition node?
}

// Subquery adds a subquery node, which is like a mini query builder that should result in a single value.
func (b *QueryBuilder) Subquery() *SubqueryNode {
	n := NewSubqueryNode(b)
	b.IsSubquery = true
	return n
}

// Load is a stub that helps the QueryBuilder implement the query.QueryBuilderI interface so it can be included in sub-queries.
func (b *QueryBuilder) Load() []map[string]interface{} {
	return nil
}

// Delete is a stub that helps the QueryBuilder implement the query.QueryBuilderI interface so it can be included in sub-queries.
func (b *QueryBuilder) Delete() {
}

// Count is a stub that helps the QueryBuilder implement the query.QueryBuilderI interface so it can be included in sub-queries.
func (b *QueryBuilder) Count(_ bool, _ ...NodeI) uint {
	return 0
}

// LoadCursor is a stub that helps the QueryBuilder implement the query.QueryBuilderI interface so it can be included in sub-queries.
func (b *QueryBuilder) LoadCursor() CursorI {
	return nil
}

// Nodes is used by the build process to return the nodes referred to in the query.
// Some nodes will be container nodes, and so will have nodes
// inside them, but every node is either referred to, or contained in the returned nodes.
// Only query builders normally need to call this.
func Nodes(b QueryBuilder) []NodeI {
	var nodes []NodeI
	for _, n := range b.Joins {
		nodes = append(nodes, n)
		if c := NodeCondition(n); c != nil {
			nodes = append(nodes, c)
		}
	}
	nodes = append(nodes, b.OrderBys...)

	if b.ConditionNode != nil {
		nodes = append(nodes, b.ConditionNode)
	}

	for _, n := range b.GroupBys {
		if NodeIsTableNodeI(n) {
			n = NodePrimaryKey(n) // Allow table nodes, but then actually have them be the pk in this context
		}
		nodes = append(nodes, n)
	}

	if b.HavingNode != nil {
		nodes = append(nodes, b.HavingNode)
	}
	nodes = append(nodes, b.Selects...)

	if b.AliasNodes != nil {
		b.AliasNodes.Range(func(key string, value Aliaser) bool {
			nodes = append(nodes, value.(NodeI))
			return true
		})
	}

	return nodes
}

type QueryExport struct {
	Joins      []NodeI
	OrderBys   []NodeI
	Condition  NodeI
	Distinct   bool
	AliasNodes *AliasNodesType
	// Adds a COUNT(*) to the select list
	GroupBys   []NodeI
	Selects    []NodeI
	LimitInfo  *LimitInfo
	Having     NodeI
	IsSubquery bool
}

func ExportQuery(b *QueryBuilder) *QueryExport {
	return &QueryExport{
		Joins:      b.Joins,
		OrderBys:   b.OrderBys,
		Condition:  b.ConditionNode,
		Distinct:   b.IsDistinct,
		AliasNodes: b.AliasNodes,
		GroupBys:   b.GroupBys,
		Selects:    b.Selects,
		LimitInfo:  b.LimitInfo,
		Having:     b.HavingNode,
		IsSubquery: b.IsSubquery,
	}
}
