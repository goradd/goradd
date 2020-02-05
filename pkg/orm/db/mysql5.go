package db

import (
	sqldb "database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/goradd/goradd/pkg/reflect"
	"strings"
	//"goradd/orm/query"
	"context"
	. "github.com/goradd/goradd/pkg/orm/query"
	"github.com/knq/snaker"
	"strconv"
)

// Mysql5 is the goradd driver for mysql databases. It works through the excellent go-sql-driver driver,
// to supply functionality above go's built in driver. To use it, call NewMysql5, but afterwards,
// work through the DB parent interface so that the underlying database can be swapped out later if needed.
//
// Timezones
// Timezones are always tricky. Mysql has some interesting quirks:
//  - Datetime types are internally stored in the timezone of the server, and then returned based on the
// timezone of the client.
//  - Timestamp types are internally stored in UTC and returned in the timezone of the client.
// The benefit of this is that you can move your database
// to a server in another timezone, and the times will automatically change to the correct timezone.
//  - The mysql-go-driver has the ability to set a default timezone in the Loc configuration parameter
// It appears to convert all times to this timezone before sending them
// to the database, and then when receiving times, it will set this as the timezone of the date.
//
// These issues are further compounded by the fact that MYSQL can initialize date and time values to what it
// believes is the current date and time in its server's timezone, but will not save the timezone itself.
// If the database gets replicated around the world, you must explicitly set the timezone of each database
// master and slave to keep datetime's in sync. Also be aware that if you are using a scaling service that is global,
// it too may change the local timezone of the server, which may be different than the timezone of the database.
// Add to this the possibility that your users may be accessing the servers from different timezones than either the
// database or server, and you get quite a tangle.
//
// So, as a general rule, use DATETIME types to represent a date combined with a time, like an appointment in
// a calendar or a recurring event that happens is entered in the current timezone is and that is editable. If you
// change timezones, the time will change too.
// Use TIMESTAMP types to store data that records when an event happened in world time. Use separate DATE and TIME
// values to record a date and time that should always be thought of in the perspective of the viewer, and
// that if the viewer changes timezones, the time will not change. 9 am in one timezone is 9 am in the other(i.e. An alarm
// for example.)
//
// Also, set the Loc configuration parameter to be the same as the server's timezone. By default its UTC.
// That will make it so all dates and times are in the same timezone as those automatically generated by MYSQL.
// It is best to set this and your database to UTC, as this will make your database portable to other timezones.
//
// Set the ParseTime configuration parameter to TRUE so that the driver will parse the times into the correct
// timezone, navigating the GO server and database server timezones. Otherwise, we
// can only assume that the database is in UTC time, since we will not get any timezone info from the server.
//
// The driver will return times in the timezone of the mysql server. This will mean that you can save data in local time,
// but you will need to convert to local time in some situations. Be particularly careful of DATE and TIME types, since
// these have no timezone information, and will always be in server time; converting to local time may have unintended
// effects.
//
// You need to be aware that when you view the data in the SQL, it will appear in whatever
// timezone the MYSQL server is set to.
type Mysql5 struct {
	SqlDb
	goraddDatabase *Database
	config      *mysql.Config
}

// New Mysql5 returns a new Mysql5 database object that you can add to the datastore.
func NewMysql5(dbKey string, params string, config *mysql.Config) *Mysql5 {
	var err error

	m := Mysql5{
		SqlDb: NewSqlDb(dbKey),
	}

	if params == "" && config == nil {
		panic("Must specify how to connect to the database.")
	}
	if params == "" {
		params = config.FormatDSN()
		m.config = config
	} else {
		m.config, err = mysql.ParseDSN(params)
		if err != nil {
			panic("Could not parse the connection string.")
		}
	}
	m.db, err = sqldb.Open("mysql", params)
	if err != nil {
		panic("Could not open database: " + err.Error())
	}
	err = m.db.Ping()
	if err != nil {
		panic("Could not ping database: " + err.Error())
	}
	m.loadDescription()
	return &m
}

// NewBuilder returns a new query builder to build a query that will be processed by the database.
func (m *Mysql5) NewBuilder() QueryBuilderI {
	return NewSqlBuilder(m)
}

// Describe returns the database description object
func (m *Mysql5) Describe() *Database {
	return m.goraddDatabase
}

func (m *Mysql5) generateSelectSql(qb QueryBuilderI) (sql string, args []interface{}) {
	b := qb.(*sqlBuilder)

	var s string
	var a []interface{}

	if b.distinct {
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

func (m *Mysql5) generateDeleteSql(qb QueryBuilderI) (sql string, args []interface{}) {
	b := qb.(*sqlBuilder)

	var s string
	var a []interface{}

	j := b.rootJoinTreeItem

	sql = "DELETE " + j.alias + " "

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

func (m *Mysql5) generateColumnListWithAliases(b *sqlBuilder) (sql string, args []interface{}) {
	b.columnAliases.Range(func(key string, j *joinTreeItem) bool {
		sql += m.generateColumnNodeSql(j.parent.alias, j.node) + " AS `" + key + "`,\n"
		return true
	})

	b.aliasNodes.Range(func(key string, v interface{}) bool {
		node := v.(NodeI)
		aliaser := v.(Aliaser)
		s, a := m.generateNodeSql(b, node, false)
		sql += s + " AS `" + aliaser.GetAlias() + "`,\n"
		args = append(args, a...)
		return true
	})

	sql = strings.TrimSuffix(sql, ",\n")
	sql += "\n"
	return
}

func (m *Mysql5) generateFromSql(b *sqlBuilder) (sql string, args []interface{}) {
	var s string
	var a []interface{}

	sql = "FROM\n"

	j := b.rootJoinTreeItem
	sql += "`" + NodeTableName(j.node) + "` AS `" + j.alias + "`\n"

	for _, cj := range j.childReferences {
		s, a = m.generateJoinSql(b, cj)
		sql += s
		args = append(args, a...)
	}
	return
}

func (m *Mysql5) generateJoinSql(b *sqlBuilder, j *joinTreeItem) (sql string, args []interface{}) {
	var tn TableNodeI
	var ok bool

	if tn, ok = j.node.(TableNodeI); !ok {
		return
	}

	switch node := tn.EmbeddedNode_().(type) {
	case *ReferenceNode:
		sql = "LEFT JOIN "
		sql += "`" + ReferenceNodeRefTable(node) + "` AS `" +
			j.alias + "` ON `" + j.parent.alias + "`.`" +
			ReferenceNodeDbColumnName(node) + "` = `" + j.alias + "`.`" + ReferenceNodeRefColumn(node) + "`"
		if j.joinCondition != nil {
			s, a := m.generateNodeSql(b, j.joinCondition, false)
			sql += " AND " + s
			args = append(args, a...)
		}
	case *ReverseReferenceNode:
		if b.limitInfo != nil && ReverseReferenceNodeIsArray(node) {
			panic("We do not currently support limited queries with an array join.")
		}

		sql = "LEFT JOIN "
		sql += "`" + ReverseReferenceNodeRefTable(node) + "` AS `" +
			j.alias + "` ON `" + j.parent.alias + "`.`" +
			ReverseReferenceNodeKeyColumnName(node) + "` = `" + j.alias + "`.`" + ReverseReferenceNodeRefColumn(node) + "`"
		if j.joinCondition != nil {
			s, a := m.generateNodeSql(b, j.joinCondition, false)
			sql += " AND " + s
			args = append(args, a...)
		}
	case *ManyManyNode:
		if b.limitInfo != nil {
			panic("We do not currently support limited queries with an array join.")
		}

		sql = "LEFT JOIN "

		var pk string
		if ManyManyNodeIsTypeTable(node) {
			pk = snaker.CamelToSnake(m.Describe().TypeTable(ManyManyNodeRefTable(node)).PkField)
		} else {
			pk = m.Describe().Table(ManyManyNodeRefTable(node)).PrimaryKeyColumn().DbName
		}

		sql += "`" + ManyManyNodeDbTable(node) + "` AS `" + j.alias + "a` ON `" +
			j.parent.alias + "`.`" +
			ColumnNodeDbName(ParentNode(node).(TableNodeI).PrimaryKeyNode()) +
			"` = `" + j.alias + "a`.`" + ManyManyNodeDbColumn(node) + "`\n"
		sql += "LEFT JOIN `" + ManyManyNodeRefTable(node) + "` AS `" + j.alias + "` ON `" + j.alias + "a`.`" + ManyManyNodeRefColumn(node) +
			"` = `" + j.alias + "`.`" + pk + "`"

		if j.joinCondition != nil {
			s, a := m.generateNodeSql(b, j.joinCondition, false)
			sql += " AND " + s
			args = append(args, a...)
		}
	default:
		return
	}
	sql += "\n"
	for _, cj := range j.childReferences {
		s, a := m.generateJoinSql(b, cj)
		sql += s
		args = append(args, a...)

	}
	return
}

func (m *Mysql5) generateNodeSql(b *sqlBuilder, n NodeI, useAlias bool) (sql string, args []interface{}) {
	switch node := n.(type) {
	case *ValueNode:
		sql = "?"
		args = append(args, ValueNodeGetValue(node))
	case *OperationNode:
		sql, args = m.generateOperationSql(b, node, useAlias)
	case *ColumnNode:
		item := b.getItemFromNode(node)
		if useAlias {
			sql = m.generateAlias(item.alias)
		} else {
			sql = m.generateColumnNodeSql(item.parent.alias, node)
		}
	case *AliasNode:
		sql = "`" + node.GetAlias() + "`"
	case *SubqueryNode:
		sql, args = m.generateSubquerySql(node)
	case TableNodeI:
		tj := b.getItemFromNode(node)
		sql = m.generateColumnNodeSql(tj.alias, node.PrimaryKeyNode())
	default:
		panic("Can't generate sql from node type.")
	}
	return
}

func (m *Mysql5) generateSubquerySql(node *SubqueryNode) (sql string, args []interface{}) {
	sql, args = m.generateSelectSql(SubqueryBuilder(node).(*sqlBuilder))
	sql = "(" + sql + ")"
	return
}

func (m *Mysql5) generateOperationSql(b *sqlBuilder, n *OperationNode, useAlias bool) (sql string, args []interface{}) {
	if useAlias && n.GetAlias() != "" {
		sql = n.GetAlias()
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
func (m *Mysql5) generateColumnNodeSql(parentAlias string, node NodeI) (sql string) {
	return "`" + parentAlias + "`.`" + ColumnNodeDbName(node.(*ColumnNode)) + "`"
}

func (m *Mysql5) generateAlias(alias string) (sql string) {
	return "`" + alias + "`"
}

func (m *Mysql5) generateNodeListSql(b *sqlBuilder, nodes []NodeI, useAlias bool) (sql string, args []interface{}) {
	for _, node := range nodes {
		s, a := m.generateNodeSql(b, node, useAlias)
		sql += s + ","
		args = append(args, a...)
	}
	sql = strings.TrimSuffix(sql, ",")
	return
}

func (m *Mysql5) generateOrderBySql(b *sqlBuilder) (sql string, args []interface{}) {
	if b.orderBys != nil && len(b.orderBys) > 0 {
		sql = "ORDER BY "
		for _, n := range b.orderBys {
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

func (m *Mysql5) generateGroupBySql(b *sqlBuilder) (sql string, args []interface{}) {
	if b.groupBys != nil && len(b.groupBys) > 0 {
		sql = "GROUP BY "
		for _, n := range b.groupBys {
			s, a := m.generateNodeSql(b, n, true)
			sql += s + ","
			args = append(args, a...)
		}
		sql = strings.TrimSuffix(sql, ",")
		sql += "\n"
	}
	return
}

func (m *Mysql5) generateWhereSql(b *sqlBuilder) (sql string, args []interface{}) {
	if b.condition != nil {
		sql = "WHERE "
		var s string
		s, args = m.generateNodeSql(b, b.condition, false)
		sql += s + "\n"
	}
	return
}

func (m *Mysql5) generateHaving(b *sqlBuilder) (sql string, args []interface{}) {
	if b.having != nil {
		sql = "HAVING "
		var s string
		s, args = m.generateNodeSql(b, b.having, false)
		sql += s + "\n"
	}
	return
}

func (m *Mysql5) generateLimitSql(b *sqlBuilder) (sql string) {
	if b.limitInfo == nil {
		return ""
	}
	if b.limitInfo.offset > 0 {
		sql = strconv.Itoa(b.limitInfo.offset) + ","
	}

	if b.limitInfo.maxRowCount > -1 {
		sql += strconv.Itoa(b.limitInfo.maxRowCount)
	}

	if sql != "" {
		sql = "LIMIT " + sql + "\n"
	}

	return
}

// Update sets specific fields of a record that already exists in the database to the given data.
func (m *Mysql5) Update(ctx context.Context, table string, fields map[string]interface{}, pkName string, pkValue interface{}) {
	var sql = "UPDATE " + table + "\n"
	var args []interface{}
	s, a := m.makeSetSql(fields)
	sql += s
	args = append(args, a...)

	sql += "WHERE " + pkName + " = ?"
	args = append(args, pkValue)
	_, e := m.Exec(ctx, sql, args...)
	if e != nil {
		panic(e.Error())
	}
}

// Insert inserts the given data as a new record in the database.
// It returns the record id of the new record.
func (m *Mysql5) Insert(ctx context.Context, table string, fields map[string]interface{}) string {
	var sql = "INSERT " + table + "\n"
	var args []interface{}
	s, a := m.makeSetSql(fields)
	sql += s
	args = append(args, a...)

	if r, err := m.Exec(ctx, sql, args...); err != nil {
		panic(err.Error())
	} else {
		if id, err := r.LastInsertId(); err != nil {
			panic(err.Error())
			return ""
		} else {
			return fmt.Sprint(id)
		}
	}
}


// Delete deletes the indicated record from the database.
func (m *Mysql5) Delete(ctx context.Context, table string, pkName string, pkValue interface{}) {
	var sql = "DELETE FROM " + table + "\n"
	var args []interface{}
	sql += "WHERE " + pkName + " = ?"
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
func (m *Mysql5) Associate(ctx context.Context,
	table string,
	column string,
	pk string,
	relatedTable string,
	relatedColumn string,
	relatedPks interface{}) { //relatedPks must be a slice of items

	// TODO: Could optimize by separating out what gets deleted, what gets added, and what stays the same.

	// First delete all previous associations
	var sql = "DELETE FROM " + table +  " WHERE " + column + "=?"
	_, e := m.Exec(ctx, sql, pk)
	if e != nil {
		panic(e.Error())
	}
	if relatedPks == nil {
		return
	}

	// Add new associations
	for _,relatedPk := range reflect.InterfaceSlice(relatedPks) {
		sql = "INSERT " + table +  " SET " + column + "=?, " + relatedColumn + "=?"
		_, e = m.Exec(ctx, sql, pk, relatedPk)
		if e != nil {
			panic(e.Error())
		}
	}
}



func (m *Mysql5) makeSetSql(fields map[string]interface{}) (sql string, args []interface{}) {
	if len(fields) == 0 {
		panic("No fields to set")
	}
	sql = "SET "
	for k, v := range fields {
		sql += fmt.Sprintf("%s=?, ", k)
		args = append(args, v)
	}

	sql = strings.TrimSuffix(sql, ", ")
	sql += "\n"
	return
}
