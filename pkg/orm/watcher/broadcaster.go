package watcher

type Broadcaster interface {
	Insert(dbId string, table string, pk string)
	Update(dbId string, table string, pk string, fieldnames ...string)
	Delete(dbId string, table string, pk string)
	BulkChange(dbId string, table string)
}
