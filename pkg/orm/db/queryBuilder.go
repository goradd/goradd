package db

import (
	"context"
	"github.com/goradd/gengen/pkg/maps"
	"github.com/goradd/goradd/pkg/orm/op"
	. "github.com/goradd/goradd/pkg/orm/query"
)

type LimitInfo struct {
	maxRowCount int
	offset      int
}

// QueryBuilder is a helper to implement the QueryBuilderI interface above in various builder classes.
// It primarily gathers up the builder instructions as the query is built. It leaves the implementation
// of the functions that actually query a database -- Load, Delete, Count -- to the containing structure.
type QueryBuilder struct {
	self       QueryBuilderI // the subclass object
	joins      []NodeI
	orderBys   []NodeI
	condition  NodeI
	distinct   bool
	aliasNodes *maps.SliceMap
	// Adds a COUNT(*) to the select list
	groupBys   []NodeI
	selects    []NodeI
	limitInfo  *LimitInfo
	having     NodeI
	isSubquery bool
}

// Init initializes the QueryBuilderBase by saving a copy of the subclass to return
// for each of the calls for chaining
func (b *QueryBuilder) Init(self QueryBuilderI) {
	b.self = self
}

// Join will attach the given reference node to the builder
func (b *QueryBuilder) Join(n NodeI, condition NodeI) QueryBuilderI {
	if b.joins != nil {
		if !NodeIsReferenceI(n) {
			panic("you can only join Reference, ReverseReference and ManyManyReference nodes")
		}

		if NodeTableName(RootNode(n)) != NodeTableName(b.joins[0]) {
			panic("you can only join nodes starting from the same table as the root node. This node must start from " + NodeTableName(b.joins[0]))
		}
	}

	NodeSetCondition(n, condition)
	b.joins = append(b.joins, n)

	return b.self
}

// Add a node that is given a manual alias name. This is usually some kind of operation, but it can be
// any Aliaser kind of node.
func (b *QueryBuilder) Alias(name string, n NodeI) QueryBuilderI {
	if b.aliasNodes == nil {
		b.aliasNodes = maps.NewSliceMap()
	}
	a := n.(Aliaser)
	a.SetAlias(name)
	b.aliasNodes.Set(name, a)
	return b.self
}

// Expands an array type node so that it will produce individual rows instead of an array of items
func (b *QueryBuilder) Expand(n NodeI) QueryBuilderI {
	if typ := NodeGetType(n); !(typ == ReverseReferenceNodeType || typ == ManyManyNodeType) {
		panic("you can only expand a node that is a ReverseReference or ManyMany node.")
	} else {
		n.(Expander).Expand()
		b.Join(n, nil)
	}

	return b.self
}

// Condition adds the condition of the Where clause. If a condition already exists, it will be anded to the previous condition.
func (b *QueryBuilder) Condition(c NodeI) QueryBuilderI {
	if b.condition == nil {
		b.condition = c
	} else {
		b.condition = op.And(b.condition, c)
	}
	return b.self
}

// OrderBy adds the order by nodes. If these are table type nodes, the primary key of the table will be used.
// These nodes can be modified using Ascending and Descending calls.
func (b *QueryBuilder) OrderBy(nodes ...NodeI) QueryBuilderI {
	b.orderBys = append(b.orderBys, nodes...)
	return b.self
}

// Limit sets the limit parameters of what is returned.
func (b *QueryBuilder) Limit(maxRowCount int, offset int) QueryBuilderI {
	if b.limitInfo != nil {
		panic("Query already has a limit")
	}
	b.limitInfo = &LimitInfo{maxRowCount, offset}

	return b.self
}

// Select specifies what specific nodes are selected. This is an optimization in order to limit the amount
// of data returned by the query. Without this, the query will expand all the join items to return every
// column of each table joined.
func (b *QueryBuilder) Select(nodes ...NodeI) QueryBuilderI {
	if b.groupBys != nil {
		panic("You cannot have Select and GroupBy statements in the same query. The GroupBy columns will automatically be selected.")
	}
	for _, n := range nodes {
		if NodeGetType(n) != ColumnNodeType {
			panic("you can only select column nodes")
		}
	}
	b.selects = append(b.selects, nodes...)
	return b.self
}

// Distinct sets the distinct bit, causing the query to not return duplicates.
func (b *QueryBuilder) Distinct() QueryBuilderI {
	b.distinct = true
	return b.self
}

// GroupBy sets the nodes that are grouped. According to SQL rules, these then are the only nodes that can be
// selected, and they MUST be selected.
func (b *QueryBuilder) GroupBy(nodes ...NodeI) QueryBuilderI {
	if b.selects != nil {
		panic("You cannot have Select and GroupBy statements in the same query. The GroupBy columns will automatically be selected.")
	}
	b.groupBys = append(b.groupBys, nodes...)
	return b.self
}

// Having adds a HAVING condition, which is a filter that acts on the results of a query.
// In particular its useful for filtering after aggregate functions have done their work.
func (b *QueryBuilder) Having(node NodeI) QueryBuilderI {
	b.having = node // should be a condition node?
	return b.self
}

// Subquery adds a subquery node, which is like a mini query builder that should result in a single value.
func (b *QueryBuilder) Subquery() *SubqueryNode {
	n := NewSubqueryNode(b.self)
	b.isSubquery = true
	return n
}

func (b *QueryBuilder) Load(ctx context.Context) []map[string]interface{} {
	return nil
}
func (b *QueryBuilder) Delete(ctx context.Context) {

}
func (b *QueryBuilder) Count(ctx context.Context, distinct bool, nodes ...NodeI) uint {
	return 0
}

type QueryExport struct {
	Joins      []NodeI
	OrderBys   []NodeI
	Condition  NodeI
	Distinct   bool
	AliasNodes *maps.SliceMap
	// Adds a COUNT(*) to the select list
	GroupBys   []NodeI
	Selects    []NodeI
	LimitInfo  *LimitInfo
	Having     NodeI
	IsSubquery bool
}

func ExportQuery(b *QueryBuilder) *QueryExport {
	return &QueryExport{
		Joins:      b.joins,
		OrderBys:   b.orderBys,
		Condition:  b.condition,
		Distinct:   b.distinct,
		AliasNodes: b.aliasNodes,
		GroupBys:   b.groupBys,
		Selects:    b.selects,
		LimitInfo:  b.limitInfo,
		Having:     b.having,
		IsSubquery: b.isSubquery,
	}
}
