package watcher

import (
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/messageServer"
)

// Watcher is the injected watcher. See the application initialization process for Watcher creation.
var Watcher WatcherI

type WatcherI interface {
	MakeKey(ctx context.Context, dbKey string, table string, pk interface{}) string
	BroadcastUpdate(ctx context.Context, dbKey string, table string, pk interface{}, fieldKeys []string)
	BroadcastInsert(ctx context.Context, dbKey string, table string, pk interface{})
	BroadcastDelete(ctx context.Context, dbKey string, table string, pk interface{})
	BroadcastBulkChange(ctx context.Context, dbKey string, table string)
}

type DefaultWatcher struct {
}

func (*DefaultWatcher) MakeKey(ctx context.Context, dbKey string, table string, pk interface{}) string {
	k := dbKey + "." + table
	if pk != "" {
		k += "." + fmt.Sprint(pk)
	}
	return k
}

func (w *DefaultWatcher) BroadcastUpdate(ctx context.Context, dbKey string, table string, pk interface{}, fieldKeys []string) {
	tableChannel := w.MakeKey(ctx, dbKey, table, "")
	pkChannel := w.MakeKey(ctx, dbKey, table, pk)
	message := make(map[string]interface{})
	message["pk"] = pk
	message["fields"] = fieldKeys
	message["op"] = "upd"
	messageServer.Send(tableChannel, "*")
	messageServer.Send(pkChannel, message)
}

func (w *DefaultWatcher) BroadcastInsert(ctx context.Context, dbKey string, table string, pk interface{}) {
	tableChannel := w.MakeKey(ctx, dbKey, table, "")
	message := make(map[string]interface{})
	message["pk"] = pk
	message["op"] = "ins"
	messageServer.Send(tableChannel, "*")
}

func (w *DefaultWatcher) BroadcastDelete(ctx context.Context, dbKey string, table string, pk interface{}) {
	tableChannel := w.MakeKey(ctx, dbKey, table, "")
	message := make(map[string]interface{})
	message["pk"] = pk
	message["op"] = "del"
	messageServer.Send(tableChannel, "*")
}

func (w *DefaultWatcher) BroadcastBulkChange(ctx context.Context, dbKey string, table string) {
	tableChannel := w.MakeKey(ctx, dbKey, table, "")
	message := make(map[string]interface{})
	message["op"] = "chg"
	messageServer.Send(tableChannel, "*")
}

func BroadcastUpdate(ctx context.Context, dbKey string, table string, pk interface{}, fieldKeys []string) {
	if Watcher != nil {
		Watcher.BroadcastUpdate(ctx, dbKey, table, pk, fieldKeys)
	}
}

func BroadcastInsert(ctx context.Context, dbKey string, table string, pk interface{}) {
	if Watcher != nil {
		Watcher.BroadcastInsert(ctx, dbKey, table, pk)
	}
}

func BroadcastDelete(ctx context.Context, dbKey string, table string, pk interface{}) {
	if Watcher != nil {
		Watcher.BroadcastDelete(ctx, dbKey, table, pk)
	}
}

func BroadcastBulkChange(ctx context.Context, dbKey string, table string) {
	if Watcher != nil {
		Watcher.BroadcastBulkChange(ctx, dbKey, table)
	}
}

func MakeKey(ctx context.Context, dbKey string, table string, pk interface{}) string {
	if Watcher == nil {
		return ""
	}
	return Watcher.MakeKey(ctx, dbKey, table, pk)
}
