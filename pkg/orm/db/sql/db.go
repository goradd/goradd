// Package sql contains helper functions that connect a standard Go database/sql object
// to the GoRADD system.
//
// Most of the functionality in this package is used by database implementations. GoRADD users would
// not normally directly call functions in this package.
package sql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/goradd/goradd/pkg/goradd"
	"github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/orm/db"
	. "github.com/goradd/goradd/pkg/orm/query"
	"strings"
	"time"
)

// The DbI interface describes the interface that a sql database needs to implement so that
// it will work with the Builder object. If you know a DatabaseI object is
// standard Go database/sql database, you can
// cast it to this type to gain access to the low level SQL driver and send it raw SQL commands.
type DbI interface {
	// Exec executes a query that does not expect to return values
	Exec(ctx context.Context, sql string, args ...interface{}) (r sql.Result, err error)
	// Query executes a query that returns values
	Query(ctx context.Context, sql string, args ...interface{}) (r *sql.Rows, err error)
	// GenerateSelectSql will generate the select sql from the builder. This sql can be specific to the database used.
	GenerateSelectSql(QueryBuilderI) (sql string, args []interface{})
	// GenerateDeleteSql will generate delete sql from the given builder.
	GenerateDeleteSql(QueryBuilderI) (sql string, args []interface{})
}

// ProfileEntry contains the data collected during sql profiling
type ProfileEntry struct {
	DbKey     string
	BeginTime time.Time
	EndTime   time.Time
	Typ       string
	Sql       string
}

// sqlContext is what is stored in the current context to keep track of queries. You must save a copy of this in the
// current context with the sqlContext key before calling database functions in order to use transactions or
// database profiling, or anything else the context is required for. The framework does this for you, but you will need
// to do this yourself if using the orm without the framework.
type sqlContext struct {
	tx       *sql.Tx
	txCount  int // Keeps track of when to close a transaction
	profiles []ProfileEntry
}

// DbHelper is a mixin for SQL database drivers. It implements common code needed by all SQL database drivers.
type DbHelper struct {
	dbKey     string  // key of the database as used in the global database map
	db        *sql.DB // Internal copy of a Go database/sql object
	profiling bool
}

// NewSqlDb creates a default DbHelper mixin.
func NewSqlDb(dbKey string, db *sql.DB) DbHelper {
	s := DbHelper{
		dbKey: dbKey,
		db:    db,
	}
	return s
}

// Begin starts a transaction. You should immediately defer a Rollback using the returned transaction id.
// If you Commit before the Rollback happens, no Rollback will occur. The Begin-Commit-Rollback pattern is nestable.
func (s *DbHelper) Begin(ctx context.Context) (txid db.TransactionID) {
	c := s.getContext(ctx)
	if c == nil {
		panic("Can't use transactions without pre-loading a context")
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
	return db.TransactionID(c.txCount)
}

// Commit commits the transaction, and if an error occurs, will panic with the error.
func (s *DbHelper) Commit(ctx context.Context, txid db.TransactionID) {
	c := s.getContext(ctx)
	if c == nil {
		panic("Can't use transactions without pre-loading a context")
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
func (s *DbHelper) Rollback(ctx context.Context, txid db.TransactionID) {
	c := s.getContext(ctx)
	if c == nil {
		panic("Can't use transactions without pre-loading a context")
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
func (s *DbHelper) Exec(ctx context.Context, sql string, args ...interface{}) (r sql.Result, err error) {
	c := s.getContext(ctx)
	log.FrameworkDebug("Exec: ", sql, args)

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
func (s *DbHelper) Prepare(ctx context.Context, sql string) (r *sql.Stmt, err error) {
	var c *sqlContext
	i := ctx.Value(goradd.sqlContext)
	if i != nil {
		c = i.(*sqlContext)
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
func (s *DbHelper) Query(ctx context.Context, sql string, args ...interface{}) (r *sql.Rows, err error) {
	c := s.getContext(ctx)
	log.FrameworkDebug("Query: ", sql, args)

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

// PutBlankContext puts a blank context into the context chain to track transactions and other
// special database situations.
func (s *DbHelper) PutBlankContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, s.contextKey(), &sqlContext{})
}

func (s *DbHelper) contextKey() goradd.ContextKey {
	return goradd.ContextKey("goradd.sql-" + s.DbKey())
}

func (s *DbHelper) getContext(ctx context.Context) *sqlContext {
	i := ctx.Value(s.contextKey())
	if i != nil {
		if c, ok := i.(*sqlContext); ok {
			return c
		}
	}
	return nil
}

// DbKey returns the key of the database in the global database store.
func (s *DbHelper) DbKey() string {
	return s.dbKey
}

// SqlDb returns the underlying database/sql database object.
func (s *DbHelper) SqlDb() *sql.DB {
	return s.db
}

// StartProfiling will start the database profiling process.
func (s *DbHelper) StartProfiling() {
	s.profiling = true
}

// GetProfiles returns currently collected profile information
// TODO: Move profiles to a session variable so we can access ajax queries too
func (s *DbHelper) GetProfiles(ctx context.Context) []ProfileEntry {
	c := s.getContext(ctx)
	if c == nil {
		panic("Profiling requires a preloaded context.")
	}

	p := c.profiles
	c.profiles = nil
	return p
}
