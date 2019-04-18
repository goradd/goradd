package db

import (
	"context"
	. "github.com/goradd/goradd/pkg/orm/query"
	"log"
)

// The dataStore is the central database collection used in codegeneration and the orm.
var datastore struct {
	// TODO: Change these to OrderedMaps. Since they are used for code generation, we want them to be iterated consistently
	databases  *DatabaseISliceMap
	tables     map[string]map[string]*TableDescription
	typeTables map[string]map[string]*TypeTableDescription
}

//type LoaderFunc func(QueryBuilderI, map[string]interface{})

type TransactionID int

// DatabaseI is the interface that describes the behaviors required for a database implementation.
type DatabaseI interface {
	// Describe returns a DatabaseDescription object, which is a complete description of the tables and fields in
	// a database and their potential relationships. SQL databases can, for the most part, generate this description
	// based on their structure. NoSQL databases would need to get this description some other way, like through
	// a json file.
	Describe() *DatabaseDescription

	// AssociatedObjectPrefix is a prefix we add to all variables that point to ORM objects. By default this is an "o".
	AssociatedObjectPrefix() string

	// NewBuilder returns a newly created query builder
	NewBuilder() QueryBuilderI

	// Update will put the given values into a record that already exists in the database. The "fields" value
	// should include only fields that have changed.
	Update(ctx context.Context, table string, fields map[string]interface{}, pkName string, pkValue string)
	// Insert will insert a new record into the database with the given values, and return the new record's primary key value.
	// The fields value should include all the required values in the database.
	Insert(ctx context.Context, table string, fields map[string]interface{}) string
	// Delete will delete the given record from the database
	Delete(ctx context.Context, table string, pkName string, pkValue interface{})

	// Begin will begin a transaction in the database and return the transaction id
	Begin(ctx context.Context) TransactionID
	// Commit will commit the given transaction
	Commit(ctx context.Context, txid TransactionID)
	// Rollback will roll back the given transaction PROVIDED it has not been committed. If it has been
	// committed, it will do nothing. Rollback can therefore be used in a defer statement as a safeguard in case
	// a transaction fails.
	Rollback(ctx context.Context, txid TransactionID)
}

// AddDatabase adds a database to the global database store. Only call this during app startup.
func AddDatabase(d DatabaseI, key string) {
	if datastore.databases == nil {
		datastore.databases = NewDatabaseISliceMap()
	}

	datastore.databases.Set(key, d)
}

// GetDatabase returns the database given the database's key.
func GetDatabase(key string) DatabaseI {
	return datastore.databases.Get(key)
}

// GetDatabases returns all databases in the datastore
func GetDatabases() []DatabaseI {
	return datastore.databases.Values()
}

// GetTableDescription returns a table description given a database key and the struct name corresponding to the table.
// You must call AnalyzeDatabases first to use this.
func GetTableDescription(key string, goTypeName string) *TableDescription {
	td, ok := datastore.tables[key][goTypeName]
	if !ok {
		return nil
	}
	return td
}

// GetTypeTableDescription returns a type table description given a database key and the struct name corresponding to the table.
// You must call AnalyzeDatabases first to use this.
func GetTypeTableDescription(key string, goTypeName string) *TypeTableDescription {
	td, ok := datastore.typeTables[key][goTypeName]
	if !ok {
		return nil
	}
	return td
}

// AnalyzeDatabases will open all the databases and extra the table information found there. Its primary use
// is to extract the schema out of SQL databases so that it can be exported into a database neutral format,
// perhaps a json or yaml representation.
func AnalyzeDatabases() {
	var dd *DatabaseDescription
	datastore.tables = make(map[string]map[string]*TableDescription)
	datastore.typeTables = make(map[string]map[string]*TypeTableDescription)
	datastore.databases.Range(func(key string, database DatabaseI) bool {
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
		return true
	})

}
