package db

import (
	sqldb "database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"strings"
	//"goradd/orm/query"
	"github.com/knq/snaker"
	//"github.com/spekary/goradd/util"
	"context"
	. "github.com/spekary/goradd/orm/query"
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
//  - Timestamp types are internally stored in UTC. The benefit of this is that you can move your database
// to a server in another timezone, and the times will automatically change to the correct timezone.
//
// So, as a general rule, use Datetime types to represent a date combined with a time, like an appointment in
// a calendar or a recurring event that happens in whatever the current timezone is. Use Timestamp types to
// store data that records when an event happened.
//
// A result of this whole thing is that when you save a datetime.DateTime type, or time.Time type into
// the database, and then read it back, you might get the time back in a different timezone than you saved it,
// since timezones are not recorded by mysql, but it will still represent the same
// moment in time. In other words, if you save 5:00 UTC+2, you might get back 3:00 UTC. We try to adjust
// times so that they are always in local server time so that the most common case will not require a change
// in timezone.
//
// Setting up your server to save in UTC time is a good thing. UTC takes out the ambiguity of daylight
// savings time. Your server is likely that way by default, but its good to check.
type Mysql5 struct {
	SqlDb
	description *DatabaseDescription
	config      *mysql.Config
}

// New Mysql5 creates a new Mysql5 object and returns its matching interface
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

func (s *Mysql5) NewBuilder() QueryBuilderI {
	return NewSqlBuilder(s)
}

func (m *Mysql5) Describe() *DatabaseDescription {
	return m.description
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

	n := b.rootNode

	sql = "DELETE " + n.GetAlias() + " "

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
	b.columnAliases.Range(func(key string, v interface{}) bool {
		node := v.(*ColumnNode)
		sql += m.generateColumnNodeSql(node, false) + " AS `" + node.GetAlias() + "`,\n"
		return true
	})

	b.aliasNodes.Range(func(key string, v interface{}) bool {
		node := v.(NodeI)
		s, a := m.generateNodeSql(node, false)
		sql += s + " AS `" + node.GetAlias() + "`,\n"
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

	n := b.rootNode
	sql += "`" + NodeTableName(n) + "` AS `" + n.GetAlias() + "`\n"

	var childNodes []NodeI
	var cn NodeI
	if childNodes = ChildNodes(n); childNodes != nil {
		for _, cn = range childNodes {
			s, a = m.generateJoinSql(b, cn)
			sql += s
			args = append(args, a...)
		}
	}
	return
}

func (m *Mysql5) generateJoinSql(b *sqlBuilder, n NodeI) (sql string, args []interface{}) {
	var tn TableNodeI
	var ok bool

	if tn, ok = n.(TableNodeI); !ok {
		return
	}

	switch node := tn.EmbeddedNode_().(type) {
	case *ReferenceNode:
		sql = "LEFT JOIN "
		sql += "`" + ReferenceNodeRefTable(node) + "` AS `" +
			node.GetAlias() + "` ON `" + ParentNode(node).GetAlias() + "`.`" +
			ReferenceNodeDbColumnName(node) + "` = `" + node.GetAlias() + "`.`" + ReferenceNodeRefColumn(node) + "`"
		if condition := NodeCondition(node); condition != nil {
			s, a := m.generateNodeSql(condition, false)
			sql += " AND " + s
			args = append(args, a...)
		}
	case *ReverseReferenceNode:
		if b.limitInfo != nil && node.IsArray() {
			panic("We do not currently support limited queries with an array join.")
		}

		sql = "LEFT JOIN "
		sql += "`" + ReverseReferenceNodeRefTable(node) + "` AS `" +
			node.GetAlias() + "` ON `" + ParentNode(node).GetAlias() + "`.`" +
			ReverseReferenceNodeDbColumnName(node) + "` = `" + node.GetAlias() + "`.`" + ReverseReferenceNodeRefColumn(node) + "`"
		if condition := NodeCondition(node); condition != nil {
			s, a := m.generateNodeSql(condition, false)
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
			pk = snaker.CamelToSnake(m.Describe().TypeTableDescription(ManyManyNodeRefTable(node)).PkField)
		} else {
			pk = m.Describe().TableDescription(ManyManyNodeRefTable(node)).PrimaryKeyColumn.DbName
		}

		sql += "`" + ManyManyNodeDbTable(node) + "` AS `" + node.GetAlias() + "a` ON `" +
			ParentNode(node).GetAlias() + "`.`" +
			ColumnNodeDbName(ParentNode(node).(TableNodeI).PrimaryKeyNode_()) +
			"` = `" + node.GetAlias() + "a`.`" + ManyManyNodeDbColumn(node) + "`\n"
		sql += "LEFT JOIN `" + ManyManyNodeRefTable(node) + "` AS `" + node.GetAlias() + "` ON `" + node.GetAlias() + "a`.`" + ManyManyNodeRefColumn(node) +
			"` = `" + node.GetAlias() + "`.`" + pk + "`"

		if condition := NodeCondition(node); condition != nil {
			s, a := m.generateNodeSql(condition, false)
			sql += " AND " + s
			args = append(args, a...)
		}
	default:
		return
	}
	sql += "\n"
	if childNodes := ChildNodes(n); childNodes != nil {
		for _, cn := range childNodes {
			s, a := m.generateJoinSql(b, cn)
			sql += s
			args = append(args, a...)

		}
	}
	return
}

func (m *Mysql5) generateNodeSql(n NodeI, useAlias bool) (sql string, args []interface{}) {
	switch node := n.(type) {
	case *ValueNode:
		sql = "?"
		args = append(args, ValueNodeGetValue(node))
	case *OperationNode:
		sql, args = m.generateOperationSql(node, useAlias)
	case *ColumnNode:
		sql = m.generateColumnNodeSql(node, useAlias)
	case *AliasNode:
		sql = "`" + node.GetAlias() + "`"
	case *SubqueryNode:
		sql, args = m.generateSubquerySql(node)
	default:
		if tn, ok := n.(TableNodeI); ok {
			sql = m.generateColumnNodeSql(tn.PrimaryKeyNode_(), false)
		} else {
			panic("Can't generate sql from node type.")
		}

	}
	return
}

func (m *Mysql5) generateSubquerySql(node *SubqueryNode) (sql string, args []interface{}) {
	sql, args = m.generateSelectSql(SubqueryBuilder(node).(*sqlBuilder))
	sql = "(" + sql + ")"
	return
}

func (m *Mysql5) generateOperationSql(n *OperationNode, useAlias bool) (sql string, args []interface{}) {
	if useAlias && n.GetAlias() != "" {
		sql = n.GetAlias()
		return
	}
	switch OperationNodeOperator(n) {
	case OpFunc:
		if len(OperationNodeOperands(n)) > 0 {
			for _, o := range OperationNodeOperands(n) {
				s, a := m.generateNodeSql(o, useAlias)
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
		s, a := m.generateNodeSql(OperationNodeOperands(n)[0], useAlias)
		sql = s + " IS " + OperationNodeOperator(n).String()
		args = append(args, a...)
		sql = "(" + sql + ") "

	case OpNot:
		s, a := m.generateNodeSql(OperationNodeOperands(n)[0], useAlias)
		sql = OperationNodeOperator(n).String() + " " + s
		args = append(args, a...)
		sql = "(" + sql + ") "

	case OpIn:
		fallthrough
	case OpNotIn:
		s, a := m.generateNodeSql(OperationNodeOperands(n)[0], useAlias)
		sql = s + " " + OperationNodeOperator(n).String() + " ("
		args = append(args, a...)

		for _, o := range ValueNodeGetValue(OperationNodeOperands(n)[1].(*ValueNode)).([]NodeI) {
			s, a = m.generateNodeSql(o, useAlias)
			sql += s + ","
			args = append(args, a...)
		}
		sql = strings.TrimSuffix(sql, ",") + ") "

	case OpAll:
		fallthrough
	case OpNone:
		sql = "(" + OperationNodeOperator(n).String() + ") "

	default:
		for _, o := range OperationNodeOperands(n) {
			s, a := m.generateNodeSql(o, useAlias)
			sql += s + " " + OperationNodeOperator(n).String() + " "
			args = append(args, a...)
		}

		sql = strings.TrimSuffix(sql, " "+OperationNodeOperator(n).String()+" ")
		sql = "(" + sql + ") "

	}
	return
}

func (m *Mysql5) generateColumnNodeSql(n *ColumnNode, useAlias bool) (sql string) {
	if useAlias {
		sql = "`" + n.GetAlias() + "`"
	} else {
		sql = "`" + ParentNode(n).GetAlias() + "`.`" + ColumnNodeDbName(n) + "`"
	}
	return
}

func (m *Mysql5) generateNodeListSql(nodes []NodeI, useAlias bool) (sql string, args []interface{}) {
	for _, node := range nodes {
		s, a := m.generateNodeSql(node, useAlias)
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
			s, a := m.generateNodeSql(n, true)
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
			s, a := m.generateNodeSql(n, true)
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
		s, args = m.generateNodeSql(b.condition, false)
		sql += s + "\n"
	}
	return
}

func (m *Mysql5) generateHaving(b *sqlBuilder) (sql string, args []interface{}) {
	if b.having != nil {
		sql = "HAVING "
		var s string
		s, args = m.generateNodeSql(b.having, false)
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

func (m *Mysql5) Update(ctx context.Context, table string, fields map[string]interface{}, pkName string, pkValue string) {
	var sql = "UPDATE " + table + "\n"
	var args = []interface{}{}
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

func (m *Mysql5) Insert(ctx context.Context, table string, fields map[string]interface{}) string {
	var sql = "INSERT " + table + "\n"
	var args = []interface{}{}
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

func (m *Mysql5) Delete(ctx context.Context, table string, pkName string, pkValue interface{}) {
	var sql = "DELETE FROM " + table + "\n"
	var args = []interface{}{}
	sql += "WHERE " + pkName + " = ?"
	args = append(args, pkValue)
	_, e := m.Exec(ctx, sql, args...)
	if e != nil {
		panic(e.Error())
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
