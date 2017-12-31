package driver

import (
	"database/sql"
	"context"
)

type SqlDb struct {
	dbKey   string  // key of the database as used in the global database map
	db      *sql.DB // Internal copy of golang database
	tx      *sql.Tx
	txCount int // Keeps track of when to close a transaction

	// codegen options
	typeTableSuffix        string // Primarily for sql tables
	associationTableSuffix string // Primarily for sql tables
	idSuffix               string // suffix to strip off the ends of names of foreign keys when converting them to internal names

	// These codegen options may be moved higher up in hierarchy some day
	goStructPrefix         string // Helps differentiate objects when different databases have the same name.
	associatedObjectPrefix string // Helps differentiate between objects and local values
}


func NewSqlDb(dbKey string) SqlDb {
	s := SqlDb{
		dbKey:dbKey,
		typeTableSuffix:"_type",
		associationTableSuffix:"_assn",
		idSuffix: "_id",
	}
	return s
}


func (s *SqlDb) Begin() {
	s.txCount++

	if s.txCount == 1 {
		var err error

		s.tx, err = s.db.Begin()
		if err != nil {
			panic(err.Error())
		}
	}
}

func (s *SqlDb) Commit() {
	s.txCount--
	if s.txCount < 0 {
		panic("Called Commit without a matching Begin")
	}
	if s.txCount == 0 {
		err := s.tx.Commit()
		if err != nil {
			panic(err.Error())
		}
		s.tx = nil
	}
}

func (s *SqlDb) Rollback() {
	if s.tx != nil {
		err := s.tx.Rollback()
		if err != nil {
			panic(err.Error())
		}
		s.tx = nil
		s.txCount = 0
	}
}

func (s *SqlDb) Exec(ctx context.Context, sql string, args ...interface{}) (r sql.Result, err error) {
	if s.tx != nil {
		r, err = s.tx.ExecContext(ctx, sql, args...)
	} else {
		r, err = s.db.ExecContext(ctx, sql, args...)
	}
	return
}

func (s *SqlDb) Prepare(sql string) (r *sql.Stmt, err error) {
	if s.tx != nil {
		r, err = s.tx.Prepare(sql)
	} else {
		r, err = s.db.Prepare(sql)
	}
	return
}

func (s *SqlDb) Query(ctx context.Context, sql string, args ...interface{}) (r *sql.Rows, err error) {
	if s.tx != nil {
		r, err = s.tx.QueryContext(ctx, sql, args...)
	} else {
		r, err = s.db.QueryContext(ctx, sql, args...)
	}
	return
}


func (s *SqlDb) DbKey() string {
	return s.dbKey
}

func (s *SqlDb) SetTypeTableSuffix(suffix string) {
	s.typeTableSuffix = suffix
}

func (s *SqlDb) SetAssociationTableSuffix(suffix string) {
	s.associationTableSuffix = suffix
}

func (s *SqlDb) TypeTableSuffix() string {
	return s.typeTableSuffix
}

func (s *SqlDb) AssociationTableSuffix() string {
	return s.associationTableSuffix
}

func (s *SqlDb) SetGoStructPrefix(prefix string) {
	s.goStructPrefix = prefix
}

func (s *SqlDb) SetAssociatedObjectPrefix(prefix string) {
	s.associatedObjectPrefix = prefix
}

func (s *SqlDb) GoStructPrefix() string {
	return s.goStructPrefix
}

func (s *SqlDb) AssociatedObjectPrefix() string {
	return s.associatedObjectPrefix
}

func (s *SqlDb) IdSuffix() string {
	return s.idSuffix
}
