package db

import (
	"context"
)

// The special item to use for named aliases in the result set
const AliasResults = "aliases_"

// The query builder is the primary aid in creating cross-platform, portable queries to the database(s)
type QueryBuilderI interface {
	Join(n NodeI, condition NodeI) QueryBuilderI
	Expand(n NodeI) QueryBuilderI
	Condition(c NodeI) QueryBuilderI
	Having(c NodeI) QueryBuilderI
	OrderBy(nodes... NodeI) QueryBuilderI
	GroupBy(nodes... NodeI) QueryBuilderI
	Limit(maxRowCount int64, offset int64) QueryBuilderI
	Select(nodes... NodeI) QueryBuilderI
	Distinct() QueryBuilderI
	Alias(name string, n NodeI) QueryBuilderI
	Load(ctx context.Context) []map[string]interface{}
	Delete(ctx context.Context)
	Count(ctx context.Context, distinct bool, nodes... NodeI) uint
	Subquery() *SubqueryNode
	nodes() []NodeI
}
