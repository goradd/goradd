package db

import (
	"context"
	. "github.com/goradd/goradd/pkg/orm/query"
	"github.com/goradd/goradd/pkg/strings"
	"github.com/goradd/maps"
)

type DatabaseMap = maps.SliceMap[string, DatabaseI]

// The dataStore is the central database collection used in code generation and the orm.
var datastore struct {
	databases *DatabaseMap
}

//type LoaderFunc func(QueryBuilderI, map[string]interface{})

type TransactionID int

// DatabaseI is the interface that describes the behaviors required for a database implementation.
type DatabaseI interface {
	// Describe returns a Database object, which is a description of the tables and fields in
	// a database and their relationships. SQL databases can, for the most part, generate this description
	// based on their structure. NoSQL databases would need to get this description some other way, like through
	// a json file.
	Describe() *Database

	// NewBuilder returns a newly created query builder
	NewBuilder(ctx context.Context) QueryBuilderI

	// Update will put the given values into a record that already exists in the database. The "fields" value
	// should include only fields that have changed.
	Update(ctx context.Context, table string, fields map[string]interface{}, pkName string, pkValue interface{})
	// Insert will insert a new record into the database with the given values, and return the new record's primary key value.
	// The fields value should include all the required values in the database.
	Insert(ctx context.Context, table string, fields map[string]interface{}) string
	// Delete will delete the given record from the database
	Delete(ctx context.Context, table string, pkName string, pkValue interface{})
	// Associate sets a many-many relationship to the given values.
	// The values are taken from the ORM, and are treated differently depending on whether this is a SQL or NoSQL database.
	Associate(ctx context.Context,
		table string,
		column string,
		pk interface{},
		relatedTable string,
		relatedColumn string,
		relatedPks interface{})

	// Begin will begin a transaction in the database and return the transaction id
	Begin(ctx context.Context) TransactionID
	// Commit will commit the given transaction
	Commit(ctx context.Context, txid TransactionID)
	// Rollback will roll back the given transaction PROVIDED it has not been committed. If it has been
	// committed, it will do nothing. Rollback can therefore be used in a defer statement as a safeguard in case
	// a transaction fails.
	Rollback(ctx context.Context, txid TransactionID)
	// PutBlankContext is called early in the processing of a response to insert an empty context that the database can use if needed.
	PutBlankContext(ctx context.Context) context.Context
}

// AddDatabase adds a database to the global database store. Only call this during app startup.
func AddDatabase(d DatabaseI, key string) {
	if !strings.HasOnlyLetters(key) {
		panic("data keys can only have letters in them. They are used in titles of variables. Please change " + key)
	}
	if datastore.databases == nil {
		datastore.databases = new(DatabaseMap)
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
