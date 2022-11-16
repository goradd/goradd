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
	"github.com/kenshaw/snaker"
	"strconv"
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

// GenerateSelectSql generates SQL for a SELECT clause.
// It returns the clause plus the arguments that substitute for values.
func (m *DB) GenerateSelectSql(qb QueryBuilderI) (sql string, args []interface{}) {
	b := qb.(*sql2.Builder)

	var s string
	var a []interface{}

	if b.IsDistinct {
		sql = "SELECT DISTINCT\n"
	} else {
		sql = "SELECT\n"
	}

	s, a = m.generateColumnListWithAliases(b)
	sql += s
	args = append(args, a...)

	s, a = m.generateFromSql(b)
	sql += s
	args = append(args, a...)

	s, a = m.generateWhereSql(b)
	sql += s
	args = append(args, a...)

	s, a = m.generateGroupBySql(b)
	sql += s
	args = append(args, a...)

	s, a = m.generateHaving(b)
	sql += s
	args = append(args, a...)

	s, a = m.generateOrderBySql(b)
	sql += s
	args = append(args, a...)

	sql += m.generateLimitSql(b)

	return
}

// GenerateDeleteSql generates SQL for a DELETE clause.
// It returns the generated SQL plus the arguments for value substitutions.
func (m *DB) GenerateDeleteSql(qb QueryBuilderI) (sql string, args []interface{}) {
	b := qb.(*sql2.Builder)

	var s string
	var a []interface{}

	j := b.RootJoinTreeItem

	sql = "DELETE " + iq(j.Alias) + " "

	s, a = m.generateFromSql(b)
	sql += s
	args = append(args, a...)

	s, a = m.generateWhereSql(b)
	sql += s
	args = append(args, a...)

	s, a = m.generateOrderBySql(b)
	sql += s
	args = append(args, a...)

	sql += m.generateLimitSql(b)

	return
}

// iq surrounds the given value with sql identifier quotes.
// It will split it if it contains a schema.
func iq(v string) string {
	parts := strings.Split(v, ".")
	if len(parts) == 2 {
		return `"` + parts[0] + `"."` + parts[1] + `"`
	}
	return `"` + v + `"`
}

func (m *DB) generateColumnListWithAliases(b *sql2.Builder) (sql string, args []interface{}) {
	b.ColumnAliases.Range(func(key string, j *sql2.JoinTreeItem) bool {
		sql += m.generateColumnNodeSql(j.Parent.Alias, j.Node) + " AS " + key + ",\n"
		return true
	})

	if b.AliasNodes != nil {
		b.AliasNodes.Range(func(key string, v Aliaser) bool {
			node := v.(NodeI)
			aliaser := v.(Aliaser)
			s, a := m.generateNodeSql(b, node, false)
			sql += s
			alias := aliaser.GetAlias()
			if alias != "" {
				// This happens in a subquery
				sql += " AS " + iq(alias)
			}
			sql += ",\n"
			args = append(args, a...)
			return true
		})
	}

	sql = strings.TrimSuffix(sql, ",\n")
	sql += "\n"
	return
}

func (m *DB) generateFromSql(b *sql2.Builder) (sql string, args []interface{}) {
	var s string
	var a []interface{}

	sql = "FROM\n"

	j := b.RootJoinTreeItem
	sql += iq(NodeTableName(j.Node)) + " AS " + j.Alias + "\n"

	for _, cj := range j.ChildReferences {
		s, a = m.generateJoinSql(b, cj)
		sql += s
		args = append(args, a...)
	}
	return
}

func (m *DB) generateJoinSql(b *sql2.Builder, j *sql2.JoinTreeItem) (sql string, args []interface{}) {
	var tn TableNodeI
	var ok bool

	if tn, ok = j.Node.(TableNodeI); !ok {
		return
	}

	switch node := tn.EmbeddedNode_().(type) {
	case *ReferenceNode:
		sql = "LEFT JOIN "
		sql += iq(ReferenceNodeRefTable(node)) + " AS " +
			iq(j.Alias) + " ON " + iq(j.Parent.Alias) + "." +
			iq(ReferenceNodeDbColumnName(node)) + " = " + iq(j.Alias) + "." + iq(ReferenceNodeRefColumn(node))
		if j.JoinCondition != nil {
			s, a := m.generateNodeSql(b, j.JoinCondition, false)
			sql += " AND " + s
			args = append(args, a...)
		}
	case *ReverseReferenceNode:
		if b.LimitInfo != nil && ReverseReferenceNodeIsArray(node) {
			panic("We do not currently support limited queries with an array join.")
		}

		sql = "LEFT JOIN "
		sql += iq(ReverseReferenceNodeRefTable(node)) + " AS " +
			iq(j.Alias) + " ON " + iq(j.Parent.Alias) + "." +
			iq(ReverseReferenceNodeKeyColumnName(node)) + " = " + iq(j.Alias) + "." + iq(ReverseReferenceNodeRefColumn(node))
		if j.JoinCondition != nil {
			s, a := m.generateNodeSql(b, j.JoinCondition, false)
			sql += " AND " + s
			args = append(args, a...)
		}
	case *ManyManyNode:
		if b.LimitInfo != nil {
			panic("We do not currently support limited queries with an array join.")
		}

		sql = "LEFT JOIN "

		var pk string
		if ManyManyNodeIsTypeTable(node) {
			pk = snaker.CamelToSnake(m.Model().TypeTable(ManyManyNodeRefTable(node)).PkField)
		} else {
			pk = m.Model().Table(ManyManyNodeRefTable(node)).PrimaryKeyColumn().DbName
		}

		sql += iq(ManyManyNodeDbTable(node)) + " AS " + iq(j.Alias+"a") + " ON " +
			iq(j.Parent.Alias) + "." +
			iq(ColumnNodeDbName(ParentNode(node).(TableNodeI).PrimaryKeyNode())) +
			" = " + iq(j.Alias+"a") + "." + iq(ManyManyNodeDbColumn(node)) + "\n"
		sql += "LEFT JOIN " + iq(ManyManyNodeRefTable(node)) + " AS " + iq(j.Alias) + " ON " + iq(j.Alias+"a") + "." + iq(ManyManyNodeRefColumn(node)) +
			" = " + iq(j.Alias) + "." + iq(pk)

		if j.JoinCondition != nil {
			s, a := m.generateNodeSql(b, j.JoinCondition, false)
			sql += " AND " + s
			args = append(args, a...)
		}
	default:
		return
	}
	sql += "\n"
	for _, cj := range j.ChildReferences {
		s, a := m.generateJoinSql(b, cj)
		sql += s
		args = append(args, a...)

	}
	return
}

func (m *DB) generateNodeSql(b *sql2.Builder, n NodeI, useAlias bool) (sql string, args []interface{}) {
	switch node := n.(type) {
	case *ValueNode:
		args = append(args, ValueNodeGetValue(node))
		sql = fmt.Sprintf("$%d", len(args))
	case *OperationNode:
		sql, args = m.generateOperationSql(b, node, useAlias)
	case *ColumnNode:
		item := b.GetItemFromNode(node)
		if useAlias {
			sql = m.generateAlias(item.Alias)
		} else {
			sql = m.generateColumnNodeSql(item.Parent.Alias, node)
		}
	case *AliasNode:
		sql = iq(node.GetAlias())
	case *SubqueryNode:
		sql, args = m.generateSubquerySql(node)
	case TableNodeI:
		tj := b.GetItemFromNode(node)
		sql = m.generateColumnNodeSql(tj.Alias, node.PrimaryKeyNode())
	default:
		panic("Can't generate sql from node type.")
	}
	return
}

func (m *DB) generateSubquerySql(node *SubqueryNode) (sql string, args []interface{}) {
	sql, args = m.GenerateSelectSql(SubqueryBuilder(node).(*sql2.Builder))
	sql = "(" + sql + ")"
	return
}

func (m *DB) generateOperationSql(b *sql2.Builder, n *OperationNode, useAlias bool) (sql string, args []interface{}) {
	if useAlias && n.GetAlias() != "" {
		sql = iq(n.GetAlias())
		return
	}
	switch OperationNodeOperator(n) {
	case OpFunc:
		if len(OperationNodeOperands(n)) > 0 {
			for _, o := range OperationNodeOperands(n) {
				s, a := m.generateNodeSql(b, o, useAlias)
				sql += s + ","
				args = append(args, a...)
			}
			sql = sql[:len(sql)-1]
		} else {
			if OperationNodeFunction(n) == "COUNT" {
				sql = "*"
			}
		}

		if OperationNodeDistinct(n) {
			sql = "DISTINCT " + sql
		}
		sql = OperationNodeFunction(n) + "(" + sql + ") "

	case OpNull:
		fallthrough
	case OpNotNull:
		s, a := m.generateNodeSql(b, OperationNodeOperands(n)[0], useAlias)
		sql = s + " IS " + OperationNodeOperator(n).String()
		args = append(args, a...)
		sql = "(" + sql + ") "

	case OpNot:
		s, a := m.generateNodeSql(b, OperationNodeOperands(n)[0], useAlias)
		sql = OperationNodeOperator(n).String() + " " + s
		args = append(args, a...)
		sql = "(" + sql + ") "

	case OpIn:
		fallthrough
	case OpNotIn:
		s, a := m.generateNodeSql(b, OperationNodeOperands(n)[0], useAlias)
		sql = s + " " + OperationNodeOperator(n).String() + " ("
		args = append(args, a...)

		for _, o := range ValueNodeGetValue(OperationNodeOperands(n)[1].(*ValueNode)).([]NodeI) {
			s, a = m.generateNodeSql(b, o, useAlias)
			sql += s + ","
			args = append(args, a...)
		}
		sql = strings.TrimSuffix(sql, ",") + ") "

	case OpAll:
		fallthrough
	case OpNone:
		sql = "(" + OperationNodeOperator(n).String() + ") "
	case OpStartsWith:
		// SQL supports this with a LIKE operation
		operands := OperationNodeOperands(n)
		s, a := m.generateNodeSql(b, operands[0], useAlias)
		v := ValueNodeGetValue(operands[1].(*ValueNode)).(string)
		v += "%"

		args = append(args, a...)
		args = append(args, v)

		sql = fmt.Sprintf(`(%s LIKE ?)`, s)
	case OpEndsWith:
		// SQL supports this with a LIKE operation
		operands := OperationNodeOperands(n)
		s, a := m.generateNodeSql(b, operands[0], useAlias)
		v := ValueNodeGetValue(operands[1].(*ValueNode)).(string)
		v = "%" + v

		args = append(args, a...)
		args = append(args, v)

		sql = fmt.Sprintf(`(%s LIKE ?)`, s)
	case OpContains:
		// SQL supports this with a LIKE operation
		operands := OperationNodeOperands(n)
		s, a := m.generateNodeSql(b, operands[0], useAlias)
		v := ValueNodeGetValue(operands[1].(*ValueNode)).(string)
		v = "%" + v + "%"

		args = append(args, a...)
		args = append(args, v)

		sql = fmt.Sprintf(`(%s LIKE ?)`, s)

	case OpDateAddSeconds:
		// Modifying a datetime in the query
		// Only works on date, datetime and timestamps. Not times.
		operands := OperationNodeOperands(n)
		s, a := m.generateNodeSql(b, operands[0], useAlias)
		s2, a2 := m.generateNodeSql(b, operands[1], useAlias)

		args = append(args, a...)
		args = append(args, a2...)

		sql = fmt.Sprintf(`DATE_ADD(%s, INTERVAL (%s) SECOND)`, s, s2)

	default:
		for _, o := range OperationNodeOperands(n) {
			s, a := m.generateNodeSql(b, o, useAlias)
			sql += s + " " + OperationNodeOperator(n).String() + " "
			args = append(args, a...)
		}

		sql = strings.TrimSuffix(sql, " "+OperationNodeOperator(n).String()+" ")
		sql = "(" + sql + ") "

	}
	return
}

// Generate the column node sql.
func (m *DB) generateColumnNodeSql(parentAlias string, node NodeI) (sql string) {
	return parentAlias + "." + iq(ColumnNodeDbName(node.(*ColumnNode)))
}

func (m *DB) generateAlias(alias string) (sql string) {
	return iq(alias)
}

func (m *DB) generateNodeListSql(b *sql2.Builder, nodes []NodeI, useAlias bool) (sql string, args []interface{}) {
	for _, node := range nodes {
		s, a := m.generateNodeSql(b, node, useAlias)
		sql += s + ","
		args = append(args, a...)
	}
	sql = strings.TrimSuffix(sql, ",")
	return
}

func (m *DB) generateOrderBySql(b *sql2.Builder) (sql string, args []interface{}) {
	if b.OrderBys != nil && len(b.OrderBys) > 0 {
		sql = "ORDER BY "
		for _, n := range b.OrderBys {
			s, a := m.generateNodeSql(b, n, true)
			if sorter, ok := n.(NodeSorter); ok {
				if NodeSorterSortDesc(sorter) {
					s += " DESC"
				}
			}
			sql += s + ","
			args = append(args, a...)
		}
		sql = strings.TrimSuffix(sql, ",")
		sql += "\n"
	}
	return
}

func (m *DB) generateGroupBySql(b *sql2.Builder) (sql string, args []interface{}) {
	if b.GroupBys != nil && len(b.GroupBys) > 0 {
		sql = "GROUP BY "
		for _, n := range b.GroupBys {
			s, a := m.generateNodeSql(b, n, true)
			sql += s + ","
			args = append(args, a...)
		}
		sql = strings.TrimSuffix(sql, ",")
		sql += "\n"
	}
	return
}

func (m *DB) generateWhereSql(b *sql2.Builder) (sql string, args []interface{}) {
	if b.ConditionNode != nil {
		sql = "WHERE "
		var s string
		s, args = m.generateNodeSql(b, b.ConditionNode, false)
		sql += s + "\n"
	}
	return
}

func (m *DB) generateHaving(b *sql2.Builder) (sql string, args []interface{}) {
	if b.HavingNode != nil {
		sql = "HAVING "
		var s string
		s, args = m.generateNodeSql(b, b.HavingNode, false)
		sql += s + "\n"
	}
	return
}

func (m *DB) generateLimitSql(b *sql2.Builder) (sql string) {
	if b.LimitInfo == nil {
		return ""
	}
	if b.LimitInfo.Offset > 0 {
		sql = strconv.Itoa(b.LimitInfo.Offset) + ","
	}

	if b.LimitInfo.MaxRowCount > -1 {
		sql += strconv.Itoa(b.LimitInfo.MaxRowCount)
	}

	if sql != "" {
		sql = "LIMIT " + sql + "\n"
	}

	return
}

// Update sets specific fields of a record that already exists in the database to the given data.
func (m *DB) Update(ctx context.Context, table string, fields map[string]interface{}, pkName string, pkValue interface{}) {
	var sql = "UPDATE " + table + "\n"
	var args []interface{}
	s, a := m.makeSetSql(fields)
	sql += s
	args = append(args, a...)

	sql += "WHERE " + iq(pkName) + fmt.Sprintf(" = $%d", len(args)+1)
	args = append(args, pkValue)
	_, e := m.Exec(ctx, sql, args...)
	if e != nil {
		panic(e.Error())
	}
}

// Insert inserts the given data as a new record in the database.
// It returns the record id of the new record.
func (m *DB) Insert(ctx context.Context, table string, fields map[string]interface{}) string {
	var sql = "INSERT INTO " + iq(table)
	var args []interface{}
	s, a := m.makeSetSql(fields)
	sql += s
	args = append(args, a...)

	if r, err := m.Exec(ctx, sql, args...); err != nil {
		panic(err.Error())
	} else {
		if id, err2 := r.LastInsertId(); err2 != nil {
			panic(err2.Error())
			return ""
		} else {
			return fmt.Sprint(id)
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

func (m *DB) makeSetSql(fields map[string]interface{}) (sql string, args []interface{}) {
	if len(fields) == 0 {
		panic("No fields to set")
	}
	sql = "SET "
	for k, v := range fields {
		sql += fmt.Sprintf("%s=$%d, ", k, len(args)+1)
		args = append(args, v)
	}

	sql = strings.TrimSuffix(sql, ", ")
	sql += "\n"
	return
}

func (m *DB) makeInsertSql(fields map[string]interface{}) (sql string, args []interface{}) {
	if len(fields) == 0 {
		panic("No fields to set")
	}

	var keys []string
	var values []string

	for k, v := range fields {
		keys = append(keys, k)
		args = append(args, v)
		values = append(values, fmt.Sprintf("$%d", len(args)))
	}

	sql = "(" + strings.Join(keys, ",") + ") VALUES ("
	sql += strings.Join(values, ",") + ")\n"
	return
}
