package dbtest

import (
	"context"
	"github.com/goradd/goradd/pkg/orm/db"
)

func getContext() context.Context {
	return db.PutContext(nil)
}

// tearDown will delete all data records possibly inserted
// by a f
func tearDown() {

}
