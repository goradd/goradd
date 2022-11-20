package dbtest

import (
	"context"
	"github.com/goradd/goradd/pkg/orm/db"
)

func getContext() context.Context {
	ctx := context.Background()
	for _, d := range db.GetDatabases() {
		ctx = d.PutBlankContext(ctx)
	}
	return ctx
}
