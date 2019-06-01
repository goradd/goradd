package query

import (
	"context"
)

// The special item to use for named aliases in the result set
const AliasResults = "aliases_"

// The query builder is the primary aid in creating cross-platform, portable queries to the database(s)
// The code-generated ORM classes call these functions to build a query. The query will eventually get
// sent to the database for processing, and then unpacked into one of the ORM generated objects. You generally
// will not call these functions directly, but rather will call the matching functions in each of the codegenerated
// ORM classes located in your project directory.
//
// If you are creating a database driver, you will implement these functions in the query builder that you
// provide.
type QueryBuilderI interface {
	Join(n NodeI, condition NodeI) QueryBuilderI
	Expand(n NodeI) QueryBuilderI
	Condition(c NodeI) QueryBuilderI
	Having(c NodeI) QueryBuilderI
	OrderBy(nodes ...NodeI) QueryBuilderI
	GroupBy(nodes ...NodeI) QueryBuilderI
	Limit(maxRowCount int, offset int) QueryBuilderI
	Select(nodes ...NodeI) QueryBuilderI
	Distinct() QueryBuilderI
	Alias(name string, n NodeI) QueryBuilderI
	Load(ctx context.Context) []map[string]interface{}
	Delete(ctx context.Context)
	Count(ctx context.Context, distinct bool, nodes ...NodeI) uint
	Subquery() *SubqueryNode
}
