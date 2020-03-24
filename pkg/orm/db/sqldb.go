package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/goradd/goradd/pkg/goradd"
	"github.com/goradd/goradd/pkg/log"
	. "github.com/goradd/goradd/pkg/orm/query"
	"strings"
	"time"
)

// The SqlDbI interface describes the interface that a sql database needs to implement so that
// it will work with the sqlBuilder object.
type SqlDbI interface {
	// Exec executes a query that does not expect to return values
	Exec(ctx context.Context, sql string, args ...interface{}) (r sql.Result, err error)
	// Exec executes a query that returns values
	Query(ctx context.Context, sql string, args ...interface{}) (r *sql.Rows, err error)
	// TypeTableSuffix returns the suffix used in a table name to indicate that the table is a type table. By default this is "_type".
	TypeTableSuffix() string
	// AssociationTableSuffix returns the suffix used in a table name to indicate that the table is an association table. By default this is "_assn".
	AssociationTableSuffix() string

	// generateSelectSql will generate the select sql from the builder. This sql can be specific to the database used.
	generateSelectSql(QueryBuilderI) (sql string, args []interface{})
	// generateDeleteSql will generate delete sql from the given builder.
	generateDeleteSql(QueryBuilderI) (sql string, args []interface{})
}

// ProfileEntry contains the data collected during sql profiling
type ProfileEntry struct {
	DbKey     string
	BeginTime time.Time
	EndTime   time.Time
	Typ       string
	Sql       string
}

// SqlContext is what is stored in the current context to keep track of queries. You must save a copy of this in the
// current context with the SqlContext key before calling database functions in order to use transactions or
// database profiling, or anything else the context is required for. The framework does this for you, but you will need
// to do this yourself if using the orm without the framework.
type SqlContext struct {
	tx      *sql.Tx
	txCount int // Keeps track of when to close a transaction

	profiles []ProfileEntry
}

// SqlDb is a mixin for SQL database drivers. It implements common code needed by all SQL database drivers.
type SqlDb struct {
	dbKey string  // key of the database as used in the global database map
	db    *sql.DB // Internal copy of golang database

	// codegen options
	typeTableSuffix        string // Primarily for sql tables
	associationTableSuffix string // Primarily for sql tables
	idSuffix               string // suffix to strip off the ends of names of foreign keys when converting them to internal names

	// These codegen options may be moved higher up in hierarchy some day
	associatedObjectPrefix string // Helps differentiate between objects and local values

	profiling bool
}

// NewSqlDb creates a default SqlDb mixin.
func NewSqlDb(dbKey string) SqlDb {
	s := SqlDb{
		dbKey:                  dbKey,
		typeTableSuffix:        "_type",
		associationTableSuffix: "_assn",
		idSuffix:               "_id",
	}
	return s
}

// Begin starts a transaction. You should immediately defer a Rollback using the returned transaction id.
// If you Commit before the Rollback happens, no Rollback will occur. The Begin-Commit-Rollback pattern is nestable.
func (s *SqlDb) Begin(ctx context.Context) (txid TransactionID) {
	var c *SqlContext

	i := ctx.Value(goradd.SqlContext)
	if i == nil {
		panic("Can't use transactions without pre-loading a context")
	} else {
		c = i.(*SqlContext)
	}
	c.txCount++

	if c.txCount == 1 {
		var err error

		c.tx, err = s.db.Begin()
		if err != nil {
			_ = c.tx.Rollback()
			c.txCount-- // transaction did not begin
			panic(err.Error())
		}
	}
	return TransactionID(c.txCount)
}

// Commit commits the transaction, and if an error occurs, will panic with the error.
func (s *SqlDb) Commit(ctx context.Context, txid TransactionID) {
	var c *SqlContext
	i := ctx.Value(goradd.SqlContext)
	if i == nil {
		panic("Can't use transactions without pre-loading a context")
	} else {
		c = i.(*SqlContext)
	}

	if c.txCount != int(txid) {
		panic("Missing Rollback after previous Begin")
	}

	if c.txCount == 0 {
		panic("Called Commit without a matching Begin")
	}
	if c.txCount == 1 {
		err := c.tx.Commit()
		if err != nil {
			panic(err.Error())
		}
		c.tx = nil
	}
	c.txCount--
}

// Rollback will rollback the transaction if the transaction is still pointing to the given txid. This gives the effect
// that if you call Rollback on a transaction that has already been committed, no Rollback will happen. This makes it easier
// to implement a transaction management scheme, because you simply always defer a Rollback after a Begin. Pass the txid
// that you got from the Begin to the Rollback. To trigger a Rollback, simply panic.
func (s *SqlDb) Rollback(ctx context.Context, txid TransactionID) {
	var c *SqlContext
	i := ctx.Value(goradd.SqlContext)
	if i == nil {
		panic("Can't use transactions without pre-loading a context")
	} else {
		c = i.(*SqlContext)
	}

	if c.txCount == int(txid) {
		err := c.tx.Rollback()
		c.txCount = 0
		c.tx = nil
		if err != nil {
			panic(err.Error())
		}
	}
}

// Exec executes the given SQL code, without returning any result rows.
func (s *SqlDb) Exec(ctx context.Context, sql string, args ...interface{}) (r sql.Result, err error) {
	var c *SqlContext
	i := ctx.Value(goradd.SqlContext)
	if i != nil {
		c = i.(*SqlContext)
	}
	log.FrameworkDebug("Exec: ", sql)

	var beginTime = time.Now()

	if c != nil && c.tx != nil {
		r, err = c.tx.ExecContext(ctx, sql, args...)
	} else {
		r, err = s.db.ExecContext(ctx, sql, args...)
	}

	var endTime = time.Now()

	if c != nil && s.profiling {
		if args != nil {
			for _, arg := range args {
				sql = strings.TrimSpace(sql)
				sql += fmt.Sprintf(",\n%#v", arg)
			}
		}
		c.profiles = append(c.profiles, ProfileEntry{DbKey: s.dbKey, BeginTime: beginTime, EndTime: endTime, Typ: "Exec", Sql: sql})
	}

	return
}

/*
func (s *SqlDb) Prepare(ctx context.Context, sql string) (r *sql.Stmt, err error) {
	var c *SqlContext
	i := ctx.Value(goradd.SqlContext)
	if i != nil {
		c = i.(*SqlContext)
	}

	var beginTime = time.Now()
	if c != nil && c.tx != nil {
		r, err = c.tx.Prepare(sql)
	} else {
		r, err = s.db.Prepare(sql)
	}
	var endTime = time.Now()
	if c != nil && s.profiling {
		c.profiles = append(c.profiles, ProfileEntry{DbKey: s.dbKey, BeginTime: beginTime, EndTime: endTime, Typ: "Prepare", Sql: sql})
	}

	return
}*/

// Query executes the given sql, and returns a row result set.
func (s *SqlDb) Query(ctx context.Context, sql string, args ...interface{}) (r *sql.Rows, err error) {
	var c *SqlContext
	i := ctx.Value(goradd.SqlContext)
	if i != nil {
		c = i.(*SqlContext)
	}
	log.FrameworkDebug("Query: ", sql)

	var beginTime = time.Now()
	if c != nil && c.tx != nil {
		r, err = c.tx.QueryContext(ctx, sql, args...)
	} else {
		r, err = s.db.QueryContext(ctx, sql, args...)
	}
	var endTime = time.Now()
	if c != nil && s.profiling {
		if args != nil {
			for _, arg := range args {
				sql = strings.TrimSpace(sql)
				sql += fmt.Sprintf(",\n%#v", arg)
			}
		}
		c.profiles = append(c.profiles, ProfileEntry{DbKey: s.dbKey, BeginTime: beginTime, EndTime: endTime, Typ: "Query", Sql: sql})
	}

	return
}

// DbKey returns the database key used in the datastore.
func (s *SqlDb) DbKey() string {
	return s.dbKey
}

// SetTypeTableSuffix sets the suffix used to identify type tables.
func (s *SqlDb) SetTypeTableSuffix(suffix string) {
	s.typeTableSuffix = suffix
}

// SetAssociationTableSuffix sets the suffix used to identify association tables.
func (s *SqlDb) SetAssociationTableSuffix(suffix string) {
	s.associationTableSuffix = suffix
}

// TypeTableSuffix returns the suffix used to identify type tables.
func (s *SqlDb) TypeTableSuffix() string {
	return s.typeTableSuffix
}

// AssociationTableSuffix returns the suffix used to identify association tables.
func (s *SqlDb) AssociationTableSuffix() string {
	return s.associationTableSuffix
}

// SetAssociatedObjectPrefix sets the prefix string used in code generation to indicate a variable is a database object.
func (s *SqlDb) SetAssociatedObjectPrefix(prefix string) {
	s.associatedObjectPrefix = prefix
}

// SetAssociatedObjectPrefix returns the prefix string used in code generation to indicate a variable is a database object.
func (s *SqlDb) AssociatedObjectPrefix() string {
	return s.associatedObjectPrefix
}

// IdSuffix is the suffix used to indicate that a field is a foreign ky to another table.
func (s *SqlDb) IdSuffix() string {
	return s.idSuffix
}

// StartProfiling will start the database profiling process.
func (s *SqlDb) StartProfiling() {
	s.profiling = true
}

// IsProfiling returns true if we are currently collection SQL database profiling information.
func IsProfiling(ctx context.Context) bool {
	var c *SqlContext

	if ctx == nil { // testing
		return false
	}
	i := ctx.Value(goradd.SqlContext)
	if i != nil {
		c = i.(*SqlContext)
		return c.profiles != nil
	}
	return false
}

// GetProfiles returns currently collected profile information
// TODO: Move profiles to a session variable so we can access ajax queries too
func GetProfiles(ctx context.Context) []ProfileEntry {
	var c *SqlContext
	i := ctx.Value(goradd.SqlContext)
	if i == nil {
		panic("Profiling requires a preloaded context.")
	} else {
		c = i.(*SqlContext)
	}

	if c != nil {
		p := c.profiles
		c.profiles = nil
		return p
	}
	return nil
}
