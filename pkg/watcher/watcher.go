package watcher

import (
	"context"
	"github.com/goradd/goradd/pkg/messageServer"
)

// The injected watcher. See the application initialization process for Watcher creation.
var Watcher WatcherI

type WatcherI interface {
	MakeKey(ctx context.Context, dbKey string, table string, pk string) string
	BroadcastUpdate(ctx context.Context, dbKey string, table string, pk string, fieldKeys []string)
	BroadcastInsert(ctx context.Context, dbKey string, table string, pk string)
	BroadcastDelete(ctx context.Context, dbKey string, table string, pk string)
}

type DefaultWatcher struct {
}

func (*DefaultWatcher) MakeKey(ctx context.Context, dbKey string, table string, pk string) string {
	k := dbKey + "." + table
	if pk != "" {
		k += "." + pk
	}
	return k
}

func (w *DefaultWatcher) BroadcastUpdate(ctx context.Context, dbKey string, table string, pk string, fieldKeys []string)  {
	tableChannel := w.MakeKey(ctx, dbKey, table, "")
	pkChannel := w.MakeKey(ctx, dbKey, table, pk)
	message := make(map[string]interface{})
	message["pk"] = pk
	message["fields"] = fieldKeys
	message["op"] = "upd"
	messageServer.SendMessage (tableChannel, nil)
	messageServer.SendMessage(pkChannel, message)
}

func (w *DefaultWatcher) BroadcastInsert(ctx context.Context, dbKey string, table string, pk string)  {
	tableChannel := w.MakeKey(ctx, dbKey, table, "")
	message := make(map[string]interface{})
	message["pk"] = pk
	message["op"] = "ins"
	messageServer.SendMessage (tableChannel, nil)
}

func (w *DefaultWatcher) BroadcastDelete(ctx context.Context, dbKey string, table string, pk string)  {
	tableChannel := w.MakeKey(ctx, dbKey, table, "")
	message := make(map[string]interface{})
	message["pk"] = pk
	message["op"] = "del"
	messageServer.SendMessage (tableChannel, nil)
}

func BroadcastUpdate(ctx context.Context, dbKey string, table string, pk string, fieldKeys []string)  {
	if Watcher != nil {
		Watcher.BroadcastUpdate(ctx, dbKey, table, pk, fieldKeys)
	}
}

func BroadcastInsert(ctx context.Context, dbKey string, table string, pk string)  {
	if Watcher != nil {
		Watcher.BroadcastInsert(ctx, dbKey, table, pk)
	}
}

func BroadcastDelete(ctx context.Context, dbKey string, table string, pk string)  {
	if Watcher != nil {
		Watcher.BroadcastDelete(ctx, dbKey, table, pk)
	}
}

func MakeKey(ctx context.Context, dbKey string, table string, pk string) string {
	if Watcher == nil {
		return ""
	}
	return Watcher.MakeKey(ctx, dbKey, table, pk)
}


