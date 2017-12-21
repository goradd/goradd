package db

import (
	"log"
	"context"
)

type LoaderFunc func(QueryBuilderI, map[string]interface{})

// The dataStore is the central database collection used in codegneration and the orm.
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

}

var datastore struct {
	databases map[string]DatabaseI
	tables map[string]*TableDescription
	typeTables map[string]*TypeTableDescription
}


func AddDatabase(d DatabaseI, key string) {
	if datastore.databases == nil {
		datastore.databases = make (map[string]DatabaseI)
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


func GetTableDescription(goTypeName string) *TableDescription {
	td, ok := datastore.tables[goTypeName]
	if !ok {
		return nil
	}
	return td
}

func GetTypeTableDescription(goTypeName string) *TypeTableDescription {
	td, ok := datastore.typeTables[goTypeName]
	if !ok {
		return nil
	}
	return td
}

func AnalyzeDatabases() {
	var dd *DatabaseDescription
	for _, database := range datastore.databases {
		dd = database.Describe()
		datastore.tables = make(map[string]*TableDescription)
		datastore.typeTables = make(map[string]*TypeTableDescription)
		for _,td := range dd.Tables {
			if !td.IsAssociation {	// association tables are private to individual databases
				if _,ok := datastore.tables[td.GoName]; ok {
					log.Panic("Table " + td.GoName + " already exists.")
				}
				datastore.tables[td.GoName] = td
			}
		}
		for _,td := range dd.TypeTables {
			if _,ok := datastore.typeTables[td.GoName]; ok {
				log.Panic("TypeTable " + td.GoName + " already exists.")
			}
			datastore.typeTables[td.GoName] = td
		}
	}

}

