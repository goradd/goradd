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
	a := &argList{}
	sql = m.generateSelectSql(qb, a)
	args = a.args()
	return
}

func (m *DB) generateSelectSql(qb QueryBuilderI, args argLister) (sql string) {
	b := qb.(*sql2.Builder)

	if b.IsDistinct {
		sql = "SELECT DISTINCT\n"
	} else {
		sql = "SELECT\n"
	}

	sql += m.generateColumnListWithAliases(b, args)

	sql += m.generateFromSql(b, args)

	sql += m.generateWhereSql(b, args)

	sql += m.generateGroupBySql(b, args)

	sql += m.generateHaving(b, args)

	sql += m.generateOrderBySql(b, args)

	sql += m.generateLimitSql(b)

	return
}

// GenerateDeleteSql generates SQL for a DELETE clause.
// It returns the generated SQL plus the arguments for value substitutions.
func (m *DB) GenerateDeleteSql(qb QueryBuilderI) (sql string, args []any) {
	a := &argList{}
	sql = m.generateDeleteSql(qb, a)
	args = a.args()
	return
}

func (m *DB) generateDeleteSql(qb QueryBuilderI, args argLister) (sql string) {
	b := qb.(*sql2.Builder)

	sql = "DELETE "

	sql += m.generateFromSql(b, args)

	sql += m.generateWhereSql(b, args)

	sql += m.generateOrderBySql(b, args)

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

func (m *DB) generateColumnListWithAliases(b *sql2.Builder, args argLister) (sql string) {
	b.ColumnAliases.Range(func(key string, j *sql2.JoinTreeItem) bool {
		sql += m.generateColumnNodeSql(j.Parent.Alias, j.Node) + " AS " + key + ",\n"
		return true
	})

	if b.AliasNodes != nil {
		b.AliasNodes.Range(func(key string, v Aliaser) bool {
			node := v.(NodeI)
			aliaser := v.(Aliaser)
			sql += m.generateNodeSql(b, node, false, args)
			alias := aliaser.GetAlias()
			if alias != "" {
				// This happens in a subquery
				sql += " AS " + iq(alias)
			}
			sql += ",\n"
			return true
		})
	}

	sql = strings.TrimSuffix(sql, ",\n")
	sql += "\n"
	return
}

func (m *DB) generateFromSql(b *sql2.Builder, args argLister) (sql string) {
	sql = "FROM\n"

	j := b.RootJoinTreeItem
	sql += iq(NodeTableName(j.Node)) + " AS " + j.Alias + "\n"

	for _, cj := range j.ChildReferences {
		sql += m.generateJoinSql(b, cj, args)
	}
	return
}

func (m *DB) generateJoinSql(b *sql2.Builder, j *sql2.JoinTreeItem, args argLister) (sql string) {
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
			s := m.generateNodeSql(b, j.JoinCondition, false, args)
			sql += " AND " + s
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
			s := m.generateNodeSql(b, j.JoinCondition, false, args)
			sql += " AND " + s
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
			s := m.generateNodeSql(b, j.JoinCondition, false, args)
			sql += " AND " + s
		}
	default:
		return
	}
	sql += "\n"
	for _, cj := range j.ChildReferences {
		s := m.generateJoinSql(b, cj, args)
		sql += s
	}
	return
}

func (m *DB) generateNodeSql(b *sql2.Builder, n NodeI, useAlias bool, args argLister) (sql string) {
	switch node := n.(type) {
	case *ValueNode:
		return args.addArg(ValueNodeGetValue(node))
	case *OperationNode:
		return m.generateOperationSql(b, node, useAlias, args)
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
		sql = m.generateSubquerySql(node, args)
	case TableNodeI:
		tj := b.GetItemFromNode(node)
		sql = m.generateColumnNodeSql(tj.Alias, node.PrimaryKeyNode())
	default:
		panic("Can't generate sql from node type.")
	}
	return
}

func (m *DB) generateSubquerySql(node *SubqueryNode, args argLister) (sql string) {
	sql = m.generateSelectSql(SubqueryBuilder(node).(*sql2.Builder), args)
	sql = "(" + sql + ")"
	return
}

func (m *DB) generateOperationSql(b *sql2.Builder, n *OperationNode, useAlias bool, args argLister) (sql string) {
	if useAlias && n.GetAlias() != "" {
		sql = iq(n.GetAlias())
		return
	}
	switch OperationNodeOperator(n) {
	case OpFunc:
		if len(OperationNodeOperands(n)) > 0 {
			for _, o := range OperationNodeOperands(n) {
				s := m.generateNodeSql(b, o, useAlias, args)
				sql += s + ","
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
		s := m.generateNodeSql(b, OperationNodeOperands(n)[0], useAlias, args)
		sql = s + " IS " + OperationNodeOperator(n).String()
		sql = "(" + sql + ") "

	case OpNot:
		s := m.generateNodeSql(b, OperationNodeOperands(n)[0], useAlias, args)
		sql = OperationNodeOperator(n).String() + " " + s
		sql = "(" + sql + ") "

	case OpIn:
		fallthrough
	case OpNotIn:
		s := m.generateNodeSql(b, OperationNodeOperands(n)[0], useAlias, args)
		sql = s + " " + OperationNodeOperator(n).String() + " ("

		for _, o := range ValueNodeGetValue(OperationNodeOperands(n)[1].(*ValueNode)).([]NodeI) {
			s = m.generateNodeSql(b, o, useAlias, args)
			sql += s + ","
		}
		sql = strings.TrimSuffix(sql, ",") + ") "

	case OpAll:
		fallthrough
	case OpNone:
		sql = "(" + OperationNodeOperator(n).String() + ") "
	case OpStartsWith:
		// SQL supports this with a LIKE operation
		operands := OperationNodeOperands(n)
		s := m.generateNodeSql(b, operands[0], useAlias, args)
		v := args.addArg(ValueNodeGetValue(operands[1].(*ValueNode)))
		v += "%"
		sql = fmt.Sprintf(`(%s LIKE %s)`, s, v)
	case OpEndsWith:
		// SQL supports this with a LIKE operation
		operands := OperationNodeOperands(n)
		s := m.generateNodeSql(b, operands[0], useAlias, args)
		v := args.addArg(ValueNodeGetValue(operands[1].(*ValueNode)))
		v = "%" + v
		sql = fmt.Sprintf(`(%s LIKE %s)`, s, v)
	case OpContains:
		// SQL supports this with a LIKE operation
		operands := OperationNodeOperands(n)
		s := m.generateNodeSql(b, operands[0], useAlias, args)
		v := args.addArg(ValueNodeGetValue(operands[1].(*ValueNode)))
		v = "%" + v + "%"
		sql = fmt.Sprintf(`(%s LIKE %s)`, s, v)

	case OpDateAddSeconds:
		// Modifying a datetime in the query
		// Only works on date, datetime and timestamps. Not times.
		operands := OperationNodeOperands(n)
		s := m.generateNodeSql(b, operands[0], useAlias, args)
		s2 := m.generateNodeSql(b, operands[1], useAlias, args)
		sql = fmt.Sprintf(`DATE_ADD(%s, INTERVAL (%s) SECOND)`, s, s2)

	case OpXor:
		// PGSQL does not have an XOR operator, so we have to manually implement the code
		operands := OperationNodeOperands(n)
		s := m.generateNodeSql(b, operands[0], useAlias, args)
		s2 := m.generateNodeSql(b, operands[1], useAlias, args)
		sql = fmt.Sprintf(`(((%[1]s) AND NOT (%[2]s)) OR (NOT (%[1]s) AND (%[2]s)))`, s, s2)

	default:
		for _, o := range OperationNodeOperands(n) {
			s := m.generateNodeSql(b, o, useAlias, args)
			sql += s + " " + OperationNodeOperator(n).String() + " "
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

func (m *DB) generateNodeListSql(b *sql2.Builder, nodes []NodeI, useAlias bool, args argLister) (sql string) {
	for _, node := range nodes {
		s := m.generateNodeSql(b, node, useAlias, args)
		sql += s + ","
	}
	sql = strings.TrimSuffix(sql, ",")
	return
}

func (m *DB) generateOrderBySql(b *sql2.Builder, args argLister) (sql string) {
	if b.OrderBys != nil && len(b.OrderBys) > 0 {
		sql = "ORDER BY "
		for _, n := range b.OrderBys {
			s := m.generateNodeSql(b, n, true, args)
			if sorter, ok := n.(NodeSorter); ok {
				if NodeSorterSortDesc(sorter) {
					s += " DESC"
				}
			}
			sql += s + ","
		}
		sql = strings.TrimSuffix(sql, ",")
		sql += "\n"
	}
	return
}

func (m *DB) generateGroupBySql(b *sql2.Builder, args argLister) (sql string) {
	if b.GroupBys != nil && len(b.GroupBys) > 0 {
		sql = "GROUP BY "
		for _, n := range b.GroupBys {
			s := m.generateNodeSql(b, n, true, args)
			sql += s + ","
		}
		sql = strings.TrimSuffix(sql, ",")
		sql += "\n"
	}
	return
}

func (m *DB) generateWhereSql(b *sql2.Builder, args argLister) (sql string) {
	if b.ConditionNode != nil {
		sql = "WHERE "
		var s string
		s = m.generateNodeSql(b, b.ConditionNode, false, args)
		sql += s + "\n"
	}
	return
}

func (m *DB) generateHaving(b *sql2.Builder, args argLister) (sql string) {
	if b.HavingNode != nil {
		sql = "HAVING "
		var s string
		s = m.generateNodeSql(b, b.HavingNode, false, args)
		sql += s + "\n"
	}
	return
}

func (m *DB) generateLimitSql(b *sql2.Builder) (sql string) {
	if b.LimitInfo == nil {
		return ""
	}

	if b.LimitInfo.MaxRowCount > -1 {
		sql += fmt.Sprintf("LIMIT %d ", b.LimitInfo.MaxRowCount)
	}

	if b.LimitInfo.Offset > 0 {
		sql += fmt.Sprintf("OFFSET %d ", b.LimitInfo.Offset)
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
