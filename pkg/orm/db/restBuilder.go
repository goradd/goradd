package db

import (
	"context"
	"encoding/json"
	. "github.com/goradd/goradd/pkg/orm/query"
)

type QueryOperation int

const (
	QueryOperationUnknown = iota
	QueryOperationLoad
	QueryOperationDelete
	QueryOperationCount
)

type restBuilder struct {
	db *Rest

	/* The variables below are populated while defining the query */
	QueryBuilder

	/* The variables below are populated during the sql build process */

	op QueryOperation
}

// NewrestBuilder creates a new restBuilder object.
func NewRestBuilder(db *Rest) *restBuilder {
	b := &restBuilder{
		db: db,
	}
	b.QueryBuilder.Init(b)
	return b
}

func (b *restBuilder) Load(ctx context.Context) (result []map[string]interface{}) {
	b.op = QueryOperationLoad

	exp := ExportQuery(&b.QueryBuilder)
	r := RestBuilderExport{*exp, b.op}
	js, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	// send query to rest point as a GET
	_ = js
	return nil
}

func (b *restBuilder) Delete(ctx context.Context) {
	b.op = QueryOperationDelete
}

func (b *restBuilder) Count(ctx context.Context, distinct bool, nodes ...NodeI) uint {
	b.op = QueryOperationCount
	return 0
}

type RestBuilderExport struct {
	QueryExport
	Op QueryOperation
}
