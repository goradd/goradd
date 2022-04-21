package mysql

import (
	sqldb "database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/goradd/goradd/pkg/orm/db"
	sql2 "github.com/goradd/goradd/pkg/orm/db/sql"
	"github.com/goradd/goradd/pkg/reflect"
	"github.com/goradd/goradd/pkg/stringmap"
	"strings"
	"time"

	//"goradd/orm/query"
	"context"
	. "github.com/goradd/goradd/pkg/orm/query"
	"github.com/kenshaw/snaker"
	"strconv"
)

// DB is the goradd driver for mysql databases. It works through the excellent go-sql-driver driver,
// to supply functionality above go's built in driver. To use it, call NewMysqlDB, but afterwards,
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
// master and slave to keep datetimes in sync. Also be aware that if you are using a scaling service that is global,
// it too may change the local timezone of the server, which may be different from the timezone of the database.
// Add to this the possibility that your users may be accessing the servers from different timezones than either the
// database or server, and you get quite a tangle.
//
// Add to that the TIMESTAMP has a max year of 2038, so TIMESTAMP itself is going to have to change soon.
//
// So, as a general rule, use DATETIME types to represent a date combined with a time, like an appointment in
// a calendar or a recurring event that happens is entered in the current timezone is and that is editable. If you
// change timezones, the time will change too.
// Use TIMESTAMP or DATETIME types to store data that records when an event happened in world time. Use separate DATE and TIME
// values to record a date and time that should always be thought of in the perspective of the viewer, and
// that if the viewer changes timezones, the time will not change. 9 am in one timezone is 9 am in the other(An alarm
// for example.)
//
// Also, set the Loc configuration parameter to be the same as the server's timezone. By default, its UTC.
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
type DB struct {
	sql2.DbHelper
	goraddDatabase *db.Database
	databaseName   string
}

// NewMysqlDB returns a new DB database object that you can add to the datastore.
func NewMysqlDB(dbKey string, params string, config *mysql.Config) *DB {
	if params == "" && config == nil {
		panic("must specify how to connect to the database")
	}
	if params == "" {
		params = config.FormatDSN()
	} else {
		var err error
		config, err = mysql.ParseDSN(params)
		if err != nil {
			panic("could not parse the connection string")
		}
	}

	db3, err := sqldb.Open("mysql", params)
	if err != nil {
		panic("Could not open database: " + err.Error())
	}
	err = db3.Ping()
	if err != nil {
		panic("Could not ping database: " + err.Error())
	}

	m := DB{
		DbHelper: sql2.NewSqlDb(dbKey, db3),
	}
	m.databaseName = config.DBName // save off the database name for later use
	m.loadDescription()
	return &m
}

// OverrideConfigSettings will use a map read in from a json file to modify
// the given config settings
func OverrideConfigSettings(config *mysql.Config, jsonContent map[string]interface{}) {
	for k, v := range jsonContent {
		switch k {
		case "dbname":
			config.DBName = v.(string)
		case "user":
			config.User = v.(string)
		case "password":
			config.Passwd = v.(string)
		case "net":
			config.Net = v.(string) // Typically, tcp or unix (for unix sockets).
		case "address":
			config.Addr = v.(string) // Note: if you set address, you MUST set net also.
		case "params":
			config.Params = stringmap.ToStringStringMap(v.(map[string]interface{}))
		case "collation":
			config.Collation = v.(string)
		case "maxAllowedPacket":
			config.MaxAllowedPacket = int(v.(float64))
		case "serverPubKey":
			config.ServerPubKey = v.(string)
		case "tlsConfig":
			config.TLSConfig = v.(string)
		case "timeout":
			config.Timeout = time.Duration(int(v.(float64))) * time.Second
		case "readTimeout":
			config.ReadTimeout = time.Duration(int(v.(float64))) * time.Second
		case "writeTimeout":
			config.WriteTimeout = time.Duration(int(v.(float64))) * time.Second
		case "allowAllFiles":
			config.AllowAllFiles = v.(bool)
		case "allowCleartextPasswords":
			config.AllowCleartextPasswords = v.(bool)
		case "allowNativePasswords":
			config.AllowNativePasswords = v.(bool)
		case "allowOldPasswords":
			config.AllowOldPasswords = v.(bool)
		}
	}

	// The other config options effect how queries work, and so should be set before
	// calling this function, as they will change how the GO code for these queries will
	// need to be written.
}

// NewBuilder returns a new query builder to build a query that will be processed by the database.
func (m *DB) NewBuilder(ctx context.Context) QueryBuilderI {
	return sql2.NewSqlBuilder(ctx, m)
}

// Describe returns the database description object
func (m *DB) Describe() *db.Database {
	return m.goraddDatabase
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

	sql = "DELETE " + j.Alias + " "

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

func (m *DB) generateColumnListWithAliases(b *sql2.Builder) (sql string, args []interface{}) {
	b.ColumnAliases.Range(func(key string, j *sql2.JoinTreeItem) bool {
		sql += m.generateColumnNodeSql(j.Parent.Alias, j.Node) + " AS `" + key + "`,\n"
		return true
	})

	if b.AliasNodes != nil {
		b.AliasNodes.Range(func(key string, v Aliaser) bool {
			node := v.(NodeI)
			aliaser := v.(Aliaser)
			s, a := m.generateNodeSql(b, node, false)
			sql += s + " AS `" + aliaser.GetAlias() + "`,\n"
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
	sql += "`" + NodeTableName(j.Node) + "` AS `" + j.Alias + "`\n"

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
		sql += "`" + ReferenceNodeRefTable(node) + "` AS `" +
			j.Alias + "` ON `" + j.Parent.Alias + "`.`" +
			ReferenceNodeDbColumnName(node) + "` = `" + j.Alias + "`.`" + ReferenceNodeRefColumn(node) + "`"
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
		sql += "`" + ReverseReferenceNodeRefTable(node) + "` AS `" +
			j.Alias + "` ON `" + j.Parent.Alias + "`.`" +
			ReverseReferenceNodeKeyColumnName(node) + "` = `" + j.Alias + "`.`" + ReverseReferenceNodeRefColumn(node) + "`"
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
			pk = snaker.CamelToSnake(m.Describe().TypeTable(ManyManyNodeRefTable(node)).PkField)
		} else {
			pk = m.Describe().Table(ManyManyNodeRefTable(node)).PrimaryKeyColumn().DbName
		}

		sql += "`" + ManyManyNodeDbTable(node) + "` AS `" + j.Alias + "a` ON `" +
			j.Parent.Alias + "`.`" +
			ColumnNodeDbName(ParentNode(node).(TableNodeI).PrimaryKeyNode()) +
			"` = `" + j.Alias + "a`.`" + ManyManyNodeDbColumn(node) + "`\n"
		sql += "LEFT JOIN `" + ManyManyNodeRefTable(node) + "` AS `" + j.Alias + "` ON `" + j.Alias + "a`.`" + ManyManyNodeRefColumn(node) +
			"` = `" + j.Alias + "`.`" + pk + "`"

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
		sql = "?"
		args = append(args, ValueNodeGetValue(node))
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
		sql = "`" + node.GetAlias() + "`"
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
	return "`" + parentAlias + "`.`" + ColumnNodeDbName(node.(*ColumnNode)) + "`"
}

func (m *DB) generateAlias(alias string) (sql string) {
	return "`" + alias + "`"
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

	sql += "WHERE " + pkName + " = ?"
	args = append(args, pkValue)
	_, e := m.Exec(ctx, sql, args...)
	if e != nil {
		panic(e.Error())
	}
}

// Insert inserts the given data as a new record in the database.
// It returns the record id of the new record.
func (m *DB) Insert(ctx context.Context, table string, fields map[string]interface{}) string {
	var sql = "INSERT " + table + "\n"
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
func (m *DB) Associate(ctx context.Context,
	table string,
	column string,
	pk interface{},
	_ string,
	relatedColumn string,
	relatedPks interface{}) { //relatedPks must be a slice of items

	// TODO: Could optimize by separating out what gets deleted, what gets added, and what stays the same.

	// First delete all previous associations
	var sql = "DELETE FROM " + table + " WHERE " + column + "=?"
	_, e := m.Exec(ctx, sql, pk)
	if e != nil {
		panic(e.Error())
	}
	if relatedPks == nil {
		return
	}

	// Add new associations
	for _, relatedPk := range reflect.InterfaceSlice(relatedPks) {
		sql = "INSERT " + table + " SET " + column + "=?, " + relatedColumn + "=?"
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
		sql += fmt.Sprintf("%s=?, ", k)
		args = append(args, v)
	}

	sql = strings.TrimSuffix(sql, ", ")
	sql += "\n"
	return
}
