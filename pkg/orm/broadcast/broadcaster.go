package broadcast

import (
	"context"
	"github.com/goradd/goradd/pkg/watcher"
)

// Broadcaster is the injected broadcaster that the generated forms use to notify the application
// that the database has changed. The application will start with the given default below, but
// you can change it if needed.
var Broadcaster BroadcasterI

type BroadcasterI interface {
	Insert(ctx context.Context, dbId string, table string, pk string)
	Update(ctx context.Context, dbId string, table string, pk string, fieldnames ...string)
	Delete(ctx context.Context, dbId string, table string, pk string)
	BulkChange(ctx context.Context, dbId string, table string)
}

// DefaultBroadcaster broadcasts database changes to the client using the watcher's pub/sub
// mechanism
type DefaultBroadcaster struct {
}

func (b DefaultBroadcaster) Insert(ctx context.Context, dbId string, table string, pk string) {
	watcher.BroadcastInsert(ctx, dbId, table, pk)
}

func (b DefaultBroadcaster) Update(ctx context.Context, dbId string, table string, pk string, fieldnames ...string) {
	watcher.BroadcastUpdate(ctx, dbId, table, pk, fieldnames)
}

func (b DefaultBroadcaster) Delete(ctx context.Context, dbId string, table string, pk string) {
	watcher.BroadcastDelete(ctx, dbId, table, pk)
}

func (b DefaultBroadcaster) BulkChange(ctx context.Context, dbId string, table string) {
	watcher.BroadcastBulkChange(ctx, dbId, table)
}

func Insert(ctx context.Context, dbId string, table string, pk string) {
	if Broadcaster != nil {
		Broadcaster.Insert(ctx, dbId, table, pk)
	}
}

func Update(ctx context.Context, dbId string, table string, pk string, fieldnames ...string) {
	if Broadcaster != nil {
		Broadcaster.Update(ctx, dbId, table, pk, fieldnames...)
	}
}

func Delete(ctx context.Context, dbId string, table string, pk string) {
	if Broadcaster != nil {
		Broadcaster.Delete(ctx, dbId, table, pk)
	}
}

func BulkChange(ctx context.Context, dbId string, table string) {
	if Broadcaster != nil {
		Broadcaster.BulkChange(ctx, dbId, table)
	}
}





