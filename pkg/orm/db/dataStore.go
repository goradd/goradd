package db

import (
	"context"
	. "github.com/goradd/goradd/pkg/orm/query"
	"log"
)

type LoaderFunc func(QueryBuilderI, map[string]interface{})

type TransactionID int

type DatabaseI interface {
	Describe() *DatabaseDescription

	// For codegen
	GoStructPrefix() string
	AssociatedObjectPrefix() string

	// Aid to build queries and deletes
	NewBuilder() QueryBuilderI

	Update(ctx context.Context, table string, fields map[string]interface{}, pkName string, pkValue string)
	Insert(ctx context.Context, table string, fields map[string]interface{}) string
	Delete(ctx context.Context, table string, pkName string, pkValue interface{})

	Begin(ctx context.Context) TransactionID
	Commit(ctx context.Context, txid TransactionID)
	Rollback(ctx context.Context, txid TransactionID)
}

// The dataStore is the central database collection used in codegeneration and the orm.
var datastore struct {
	// TODO: Change these to OrderedMaps. Since they are used for code generation, we want them to be iterated consistently
	databases  map[string]DatabaseI
	tables     map[string]map[string]*TableDescription
	typeTables map[string]map[string]*TypeTableDescription
}

func AddDatabase(d DatabaseI, key string) {
	if datastore.databases == nil {
		datastore.databases = make(map[string]DatabaseI)
	}

	datastore.databases[key] = d
}

func GetDatabase(key string) DatabaseI {
	d := datastore.databases[key]
	return d
}

func GetDatabases() map[string]DatabaseI {
	return datastore.databases
}

func GetTableDescription(key string, goTypeName string) *TableDescription {
	td, ok := datastore.tables[key][goTypeName]
	if !ok {
		return nil
	}
	return td
}

func GetTypeTableDescription(key string, goTypeName string) *TypeTableDescription {
	td, ok := datastore.typeTables[key][goTypeName]
	if !ok {
		return nil
	}
	return td
}

func AnalyzeDatabases() {
	var dd *DatabaseDescription
	datastore.tables = make(map[string]map[string]*TableDescription)
	datastore.typeTables = make(map[string]map[string]*TypeTableDescription)
	for key, database := range datastore.databases {
		dd = database.Describe()
		datastore.tables[key] = make(map[string]*TableDescription)
		datastore.typeTables[key] = make(map[string]*TypeTableDescription)
		for _, td := range dd.Tables {
			if !td.IsAssociation {
				if _, ok := datastore.tables[key][td.GoName]; ok {
					log.Panic("Table " + key + ":" + td.GoName + " already exists.")
				}
				datastore.tables[key][td.GoName] = td
			}
		}
		for _, td := range dd.TypeTables {
			if _, ok := datastore.typeTables[key][td.GoName]; ok {
				log.Panic("TypeTable " + key + ":" + td.GoName + " already exists.")
			}
			datastore.typeTables[key][td.GoName] = td
		}
	}

}
