package query

import "context"

// AliasResults is the special item to use for named aliases in the result set
const AliasResults = "aliases_"

// QueryBuilderI is the primary aid in creating cross-platform, portable queries to the database(s)
// The code-generated ORM classes call these functions to build a query. The query will eventually get
// sent to the database for processing, and then unpacked into one of the ORM generated objects. You generally
// will not call these functions directly, but rather will call the matching functions in each of the codegenerated
// ORM classes located in your project directory.
//
// If you are creating a database driver, you will implement these functions in the query builder that you
// provide.
type QueryBuilderI interface {
	Join(n NodeI, condition NodeI)
	Expand(n NodeI)
	Condition(c NodeI)
	Having(c NodeI)
	OrderBy(nodes ...NodeI)
	GroupBy(nodes ...NodeI)
	Limit(maxRowCount int, offset int)
	Select(nodes ...NodeI)
	Distinct()
	Alias(name string, n NodeI)
	// Load terminates the builder, queries the database, and returns the results as an array of interfaces similar in structure to a json structure
	Load() []map[string]interface{}
	Delete()
	Count(distinct bool, nodes ...NodeI) uint
	Subquery() *SubqueryNode
	Context() context.Context
	LoadCursor() CursorI
}

type CursorI interface {
	Next() map[string]interface{}
	Close() error
}
