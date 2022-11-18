package pgsql

import (
	"context"
	sqldb "database/sql"
	"fmt"
	"github.com/goradd/goradd/pkg/orm/db"
	sql2 "github.com/goradd/goradd/pkg/orm/db/sql"
	. "github.com/goradd/goradd/pkg/orm/query"
	"github.com/goradd/goradd/pkg/reflect"
	"github.com/goradd/goradd/pkg/stringmap"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"strings"
)

// DB is the goradd driver for postgresql databases.
type DB struct {
	sql2.DbHelper
	model   *db.Model
	schemas []string
}

// NewDB returns a new Postgresql DB database object based on the pgx driver
// that you can add to the datastore.
// If connectionString is set, it will be used to create the configuration. Otherwise,
// use a config setting. Using a configSetting can potentially give you access to the
// underlying pgx database for advanced operations.
//
// The postgres driver specifies that you must use ParseConfig
// to create the initial configuration, although that can be sent a blank string to
// gather initial values from environment variables. You can then change items in
// the configuration structure. For example:
//
//	config,_ := pgx.ParseConfig(connectionString)
//	config.Password = "mysecret"
//	db := pgsql.NewDB(key, "", config)
func NewDB(dbKey string,
	connectionString string,
	config *pgx.ConnConfig) *DB {
	if connectionString == "" && config == nil {
		panic("must specify how to connect to the database")
	}

	if connectionString == "" {
		connectionString = stdlib.RegisterConnConfig(config)
	}

	db3, err := sqldb.Open("pgx", connectionString)
	if err != nil {
		panic("Could not open database: " + err.Error())
	}
	err = db3.Ping()
	if err != nil {
		panic("Could not ping database " + dbKey + ":" + err.Error())
	}

	m := DB{
		DbHelper: sql2.NewSqlDb(dbKey, db3),
	}
	return &m
}

// OverrideConfigSettings will use a map read in from a json file to modify
// the given config settings
func OverrideConfigSettings(config *pgx.ConnConfig, jsonContent map[string]interface{}) {
	for k, v := range jsonContent {
		switch k {
		case "database":
			config.Database = v.(string)
		case "user":
			config.User = v.(string)
		case "password":
			config.Password = v.(string)
		case "host":
			config.Host = v.(string) // Typically, tcp or unix (for unix sockets).
		case "port":
			config.Port = uint16(v.(float64))
		case "runtimeParams":
			config.RuntimeParams = stringmap.ToStringStringMap(v.(map[string]interface{}))
		case "kerberosServerName":
			config.KerberosSrvName = v.(string)
		case "kerberosSPN":
			config.KerberosSpn = v.(string)
		}
	}
}

// NewBuilder returns a new query builder to build a query that will be processed by the database.
func (m *DB) NewBuilder(ctx context.Context) QueryBuilderI {
	return sql2.NewSqlBuilder(ctx, m)
}

// Model returns the database description object
func (m *DB) Model() *db.Model {
	return m.model
}

func iq(v string) string {
	parts := strings.Split(v, ".")

	// if the identifier has a schema, quote the parts separately
	if len(parts) == 2 {
		return `"` + parts[0] + `"."` + parts[1] + `"`
	}
	return `"` + v + `"`
}

// QuoteIdentifier surrounds the given identifier with quote characters
// appropriate for Postgres
func (m *DB) QuoteIdentifier(v string) string {
	return iq(v)
}

// FormatArgument formats the given argument number for embedding in a SQL statement.
func (m *DB) FormatArgument(n int) string {
	return fmt.Sprintf(`$%d`, n)
}

// OperationSql provides Postgres specific SQL for certain operators.
func (m *DB) OperationSql(op Operator, operandStrings []string) (sql string) {
	switch op {
	case OpDateAddSeconds:
		// Modifying a datetime in the query
		// Only works on date, datetime and timestamps. Not times.
		s := operandStrings[0]
		s2 := operandStrings[1]
		return fmt.Sprintf(`(%s + make_interval(seconds => %s))`, s, s2)
	}
	return
}

// Update sets specific fields of a record that already exists in the database to the given data.
func (m *DB) Update(ctx context.Context, table string, fields map[string]interface{}, pkName string, pkValue interface{}) {
	var sql = "UPDATE " + table + "\n"
	args := &argList{}
	s := m.makeSetSql(fields, args)
	sql += s

	sql += "WHERE " + iq(pkName) + fmt.Sprintf(" = %s", args.addArg(pkValue))
	_, e := m.Exec(ctx, sql, args.args()...)
	if e != nil {
		panic(e.Error())
	}
}

// Insert inserts the given data as a new record in the database.
// It returns the record id of the new record.
func (m *DB) Insert(ctx context.Context, table string, fields map[string]interface{}) string {
	var sql = "INSERT INTO " + iq(table)
	args := &argList{}
	sql += " " + m.makeInsertSql(fields, args)
	sql += " RETURNING "
	sql += m.Model().Table(table).PrimaryKeyColumn().DbName
	if rows, err := m.Query(ctx, sql, args.args()...); err != nil {
		panic(err.Error())
	} else {
		var id string
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&id)
		}
		if err != nil {
			panic(err.Error())
			return ""
		} else {
			return id
		}
	}
}

// Delete deletes the indicated record from the database.
func (m *DB) Delete(ctx context.Context, table string, pkName string, pkValue interface{}) {
	var sql = "DELETE FROM " + iq(table) + "\n"
	var args []interface{}
	sql += "WHERE " + iq(pkName) + " = $1"
	args = append(args, pkValue)
	_, e := m.Exec(ctx, sql, args...)
	if e != nil {
		panic(e.Error())
	}
}

// Associate sets up the many-many association pointing from the given table and column to another table and column.
// table is the name of the association table.
// column is the name of the column in the association table that contains the pk for the record we are associating.
// pk is the value of the primary key.
// relatedTable is the table the association is pointing to.
// relatedColumn is the column in the association table that points to the relatedTable's pk.
// relatedPks are the new primary keys in the relatedTable we are associating.
func (m *DB) Associate(ctx context.Context,
	table string,
	column string,
	pk interface{},
	_ string,
	relatedColumn string,
	relatedPks interface{}) { //relatedPks must be a slice of items

	// TODO: Could optimize by separating out what gets deleted, what gets added, and what stays the same.

	// TODO: Make this part of a transaction
	// First delete all previous associations
	var sql = "DELETE FROM " + iq(table) + " WHERE " + iq(column) + "=$1"
	_, e := m.Exec(ctx, sql, pk)
	if e != nil {
		panic(e.Error())
	}
	if relatedPks == nil {
		return
	}

	// Add new associations
	for _, relatedPk := range reflect.InterfaceSlice(relatedPks) {
		sql = "INSERT INTO " + iq(table) + "(" + iq(column) + "," + iq(relatedColumn) + ") VALUES ($1, $2)"
		_, e = m.Exec(ctx, sql, pk, relatedPk)
		if e != nil {
			panic(e.Error())
		}
	}
}

func (m *DB) makeSetSql(fields map[string]interface{}, args argLister) (sql string) {
	if len(fields) == 0 {
		panic("No fields to set")
	}
	sql = "SET "
	for k, v := range fields {
		sql += fmt.Sprintf("%s=%s, ", iq(k), args.addArg(v))
	}

	sql = strings.TrimSuffix(sql, ", ")
	sql += "\n"
	return
}

func (m *DB) makeInsertSql(fields map[string]interface{}, args argLister) (sql string) {
	if len(fields) == 0 {
		panic("No fields to set")
	}

	var keys []string
	var values []string

	for k, v := range fields {
		keys = append(keys, iq(k))
		values = append(values, args.addArg(v))
	}

	sql = "(" + strings.Join(keys, ",") + ") VALUES ("
	sql += strings.Join(values, ",") + ")\n"
	return
}
