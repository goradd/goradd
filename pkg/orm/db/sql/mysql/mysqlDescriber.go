package mysql

import (
	"database/sql"
	"fmt"
	"github.com/goradd/goradd/pkg/orm/db"
	sql2 "github.com/goradd/goradd/pkg/orm/db/sql"
	. "github.com/goradd/goradd/pkg/orm/query"
	"github.com/goradd/goradd/pkg/stringmap"
	strings2 "github.com/goradd/goradd/pkg/strings"
	"log"
	"math"
	"sort"
	"strings"
)

/*
This file contains the code that parses the data structure found in a MySQL database into
our own cross-platform internal database description object.
*/

type mysqlTable struct {
	name                string
	columns             []mysqlColumn
	indexes             []mysqlIndex
	fkMap               map[string]mysqlForeignKey
	comment             string
	options             map[string]interface{}
	supportsForeignKeys bool
}

type mysqlColumn struct {
	name            string
	defaultValue    sql2.SqlReceiver
	isNullable      string
	dataType        string
	dataLen         int
	characterMaxLen sql.NullInt64
	columnType      string
	key             string
	extra           string
	comment         string
	options         map[string]interface{}
}

type mysqlIndex struct {
	name       string
	nonUnique  bool
	tableName  string
	columnName string
}

type mysqlForeignKey struct {
	constraintName       string
	tableName            string
	columnName           string
	referencedSchema     sql.NullString
	referencedTableName  sql.NullString
	referencedColumnName sql.NullString
	updateRule           sql.NullString
	deleteRule           sql.NullString
}

/*
type DB struct {
	dbKey   string
	db      *sql.DB
	config  mysql.Config
	options DbOptions
}



func NewMysql2 (dbKey string , options DbOptions, config *mysql.Config) (*DB, error) {
	source := DB{}
	db, err := sql.Open("mysql", config.FormatDSN())

	source.dbKey = DbKey
	source.db = db
	source.config = *config
	source.options = options

	// Ping?

	return &source,err
}*/

func (m *DB) Analyze(options Options) {
	rawTables := m.getRawTables()
	description := m.descriptionFromRawTables(rawTables, options)
	m.model = db.NewModel(m.DbKey(), options.ForeignKeySuffix, description)
}

func (m *DB) getRawTables() map[string]mysqlTable {
	var tableMap = make(map[string]mysqlTable)

	indexes, err := m.getIndexes()
	if err != nil {
		return nil
	}

	foreignKeys, err := m.getForeignKeys()
	if err != nil {
		return nil
	}

	tables := m.getTables()
	for _, table := range tables {
		// Do some processing on the foreign keys
		for _, fk := range foreignKeys[table.name] {
			if fk.referencedColumnName.Valid && fk.referencedTableName.Valid {
				if _, ok := table.fkMap[fk.columnName]; ok {
					log.Printf("Warning: Column %s:%s multi-table foreign keys are not supported.", table.name, fk.columnName)
					delete(table.fkMap, fk.columnName)
				} else {
					table.fkMap[fk.columnName] = fk
				}
			}
		}

		columns, err2 := m.getColumns(table.name)
		if err2 != nil {
			return nil
		}

		table.indexes = indexes[table.name]
		table.columns = columns
		tableMap[table.name] = table
	}

	return tableMap

}

// Gets information for a table
func (m *DB) getTables() []mysqlTable {
	var tableName, tableComment, tableEngine string
	var tables []mysqlTable

	// Use the MySQL5 Information Schema to get a list of all the tables in this database
	// (excluding views, etc.)
	dbName := m.databaseName

	rows, err := m.SqlDb().Query(fmt.Sprintf(`
	SELECT
	table_name,
	table_comment,
	engine
	FROM
	information_schema.tables
	WHERE
	table_type <> 'VIEW' AND
	table_schema = '%s';
	`, dbName))

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&tableName, &tableComment, &tableEngine)
		var supportsForeignKeys bool
		if tableEngine == "InnoDB" {
			supportsForeignKeys = true
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Println(tableName)
		table := mysqlTable{
			name:                tableName,
			comment:             tableComment,
			columns:             []mysqlColumn{},
			fkMap:               make(map[string]mysqlForeignKey),
			indexes:             []mysqlIndex{},
			supportsForeignKeys: supportsForeignKeys,
		}
		if table.options, table.comment, err = sql2.ExtractOptions(table.comment); err != nil {
			log.Print("Error in comment options for table " + table.name + " - " + err.Error())
		}

		tables = append(tables, table)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return tables
}

func (m *DB) getColumns(table string) (columns []mysqlColumn, err error) {
	dbName := m.databaseName

	rows, err := m.SqlDb().Query(fmt.Sprintf(`
	SELECT
	column_name,
	column_default,
	is_nullable,
	data_type,
	character_maximum_length,
	column_type,
	column_key,
	extra,
	column_comment
	FROM
	information_schema.columns
	WHERE
	table_name = '%s' AND
	table_schema = '%s'
	ORDER BY
	ordinal_position;
	`, table, dbName))

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var col mysqlColumn

	for rows.Next() {
		col = mysqlColumn{}
		err = rows.Scan(&col.name, &col.defaultValue.R, &col.isNullable, &col.dataType, &col.characterMaxLen, &col.columnType, &col.key, &col.extra, &col.comment)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		if col.options, col.comment, err = sql2.ExtractOptions(col.comment); err != nil {
			log.Print("Error in table comment options for table " + table + ":" + col.name + " - " + err.Error())
		}
		columns = append(columns, col)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return columns, err
}

func (m *DB) getIndexes() (indexes map[string][]mysqlIndex, err error) {

	dbName := m.databaseName
	indexes = make(map[string][]mysqlIndex)

	rows, err := m.SqlDb().Query(fmt.Sprintf(`
	SELECT
	index_name,
	non_unique,
	table_name,
	column_name
	FROM
	information_schema.statistics
	WHERE
	table_schema = '%s';
	`, dbName))

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var index mysqlIndex

	for rows.Next() {
		index = mysqlIndex{}
		err = rows.Scan(&index.name, &index.nonUnique, &index.tableName, &index.columnName)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		tableIndexes := indexes[index.tableName]
		tableIndexes = append(tableIndexes, index)
		indexes[index.tableName] = tableIndexes
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return indexes, err
}

// getForeignKeys gets information on the foreign keys.
//
// Note that querying the information_schema database is SLOW, so we want to do it as few times as possible.
func (m *DB) getForeignKeys() (foreignKeys map[string][]mysqlForeignKey, err error) {
	dbName := m.databaseName
	fkMap := make(map[string]mysqlForeignKey)

	rows, err := m.SqlDb().Query(fmt.Sprintf(`
	SELECT
	constraint_name,
	table_name,
	column_name,
	referenced_table_name,
	referenced_column_name
	FROM
	information_schema.key_column_usage
	WHERE
	constraint_schema = '%s'
	ORDER BY
	ordinal_position;
	`, dbName))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		fk := mysqlForeignKey{}
		err = rows.Scan(&fk.constraintName, &fk.tableName, &fk.columnName, &fk.referencedTableName, &fk.referencedColumnName)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		if fk.referencedColumnName.Valid {
			fkMap[fk.constraintName] = fk
		}
	}

	rows.Close()

	rows, err = m.SqlDb().Query(fmt.Sprintf(`
	SELECT
	constraint_name,
	update_rule,
	delete_rule
	FROM
	information_schema.referential_constraints
	WHERE
	constraint_schema = '%s';
	`, dbName))
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var constraintName string
		var updateRule, deleteRule sql.NullString
		err = rows.Scan(&constraintName, &updateRule, &deleteRule)
		if err != nil {
			log.Fatal(err)
		}
		if fk, ok := fkMap[constraintName]; ok {
			fk.updateRule = updateRule
			fk.deleteRule = deleteRule
			fkMap[constraintName] = fk
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	foreignKeys = make(map[string][]mysqlForeignKey)
	stringmap.Range(fkMap, func(_ string, val interface{}) bool {
		fk := val.(mysqlForeignKey)
		tableKeys := foreignKeys[fk.tableName]
		tableKeys = append(tableKeys, fk)
		foreignKeys[fk.tableName] = tableKeys
		return true
	})
	return foreignKeys, err
}

// Convert the database native type to a more generic sql type, and a go table type.
func (m *DB) processTypeInfo(tableName string, column mysqlColumn, cd *db.ColumnDescription) {
	dataLen := sql2.GetDataDefLength(column.columnType)

	isUnsigned := strings.Contains(column.columnType, "unsigned")
	cd.NativeType = column.dataType

	switch column.dataType {
	case "time":
		cd.GoType = ColTypeTime.GoType()
		cd.SubType = "time"
	case "timestamp":
		cd.GoType = ColTypeTime.GoType()
		cd.SubType = "timestamp"
	case "datetime":
		cd.GoType = ColTypeTime.GoType()
	case "date":
		cd.GoType = ColTypeTime.GoType()
		cd.SubType = "date"
	case "tinyint":
		if dataLen == 1 {
			cd.GoType = ColTypeBool.GoType()
		} else {
			if isUnsigned {
				cd.GoType = ColTypeUnsigned.GoType()
				cd.MinValue = uint64(0)
				cd.MaxValue = uint64(255)
				cd.MaxCharLength = 3
			} else {
				cd.GoType = ColTypeInteger.GoType()
				cd.MinValue = int64(-128)
				cd.MaxValue = int64(127)
				cd.MaxCharLength = 4 // allow for a negative sign
			}
		}

	case "int":
		if isUnsigned {
			cd.GoType = ColTypeUnsigned.GoType()
			cd.MinValue = uint64(0)
			cd.MaxValue = uint64(4294967295)
			cd.MaxCharLength = 10
		} else {
			cd.GoType = ColTypeInteger.GoType()
			cd.MinValue = int64(-2147483648)
			cd.MaxValue = int64(2147483647)
			cd.MaxCharLength = 11
		}

	case "smallint":
		if isUnsigned {
			cd.GoType = ColTypeUnsigned.GoType()
			cd.MinValue = uint64(0)
			cd.MaxValue = uint64(65535)
			cd.MaxCharLength = 5
		} else {
			cd.GoType = ColTypeInteger.GoType()
			cd.MinValue = int64(-32768)
			cd.MaxValue = int64(32767)
			cd.MaxCharLength = 6
		}

	case "mediumint":
		if isUnsigned {
			cd.GoType = ColTypeUnsigned.GoType()
			cd.MinValue = uint64(0)
			cd.MaxValue = uint64(16777215)
			cd.MaxCharLength = 8
		} else {
			cd.GoType = ColTypeInteger.GoType()
			cd.MinValue = int64(-8388608)
			cd.MaxValue = int64(8388607)
			cd.MaxCharLength = 8
		}

	case "bigint": // We need to be explicit about this in go, since int will be whatever the OS native int size is, but go will support int64 always.
		// Also, since Json can only be decoded into float64s, we are limited in our ability to represent large min and max numbers in the json to about 2^53
		if isUnsigned {
			cd.GoType = ColTypeUnsigned64.GoType()
			cd.MinValue = uint64(0)
			cd.MaxValue = uint64(math.MaxUint64)
			cd.MaxCharLength = 20
		} else {
			cd.GoType = ColTypeInteger64.GoType()
			cd.MinValue = int64(math.MinInt64)
			cd.MaxValue = int64(math.MaxInt64)
			cd.MaxCharLength = 20
		}

	case "float":
		cd.GoType = ColTypeFloat32.GoType()
		cd.MinValue = -math.MaxFloat32 // float64 type
		cd.MaxValue = math.MaxFloat32
	case "double":
		cd.GoType = ColTypeFloat64.GoType()
		cd.MinValue = -math.MaxFloat64
		cd.MaxValue = math.MaxFloat64
	case "varchar":
		cd.GoType = ColTypeString.GoType()
		cd.MaxCharLength = uint64(dataLen)

	case "char":
		cd.GoType = ColTypeString.GoType()
		cd.MaxCharLength = uint64(dataLen)

	case "blob":
		cd.GoType = ColTypeBytes.GoType()
		cd.MaxCharLength = 65535
	case "tinyblob":
		cd.GoType = ColTypeBytes.GoType()
		cd.MaxCharLength = 255
	case "mediumblob":
		cd.GoType = ColTypeBytes.GoType()
		cd.MaxCharLength = 16777215
	case "longblob":
		cd.GoType = ColTypeBytes.GoType()
		cd.MaxCharLength = math.MaxUint32

	case "text":
		cd.GoType = ColTypeString.GoType()
		cd.MaxCharLength = 65535
	case "tinytext":
		cd.GoType = ColTypeString.GoType()
		cd.MaxCharLength = 255
	case "mediumtext":
		cd.GoType = ColTypeString.GoType()
		cd.MaxCharLength = 16777215
	case "longtext":
		cd.GoType = ColTypeString.GoType()
		cd.MaxCharLength = math.MaxUint32

	case "decimal": // No native equivalent in Go. See the "Big" go package for support. You will need to shephard numbers into and out of string format to move data to the database
		cd.GoType = ColTypeString.GoType()
		cd.MaxCharLength = uint64(dataLen) + 3

	case "year":
		cd.GoType = ColTypeInteger.GoType()

	case "set":
		log.Print("Note: Using association tables is preferred to using DB SET columns in table " + tableName + ":" + column.name + ".")
		cd.GoType = ColTypeString.GoType()
		cd.MaxCharLength = uint64(column.characterMaxLen.Int64)

	case "enum":
		log.Print("Note: Using type tables is preferred to using DB ENUM columns in table " + tableName + ":" + column.name + ".")
		cd.GoType = ColTypeString.GoType()
		cd.MaxCharLength = uint64(column.characterMaxLen.Int64)

	default:
		cd.GoType = ColTypeString.GoType()
	}

	cd.DefaultValue = column.defaultValue.UnpackDefaultValue(ColTypeFromGoTypeString(cd.GoType))
}

func (m *DB) descriptionFromRawTables(rawTables map[string]mysqlTable, options Options) db.DatabaseDescription {

	dd := db.DatabaseDescription{}

	keys := stringmap.SortedKeys(rawTables)
	for _, tableName := range keys {
		table := rawTables[tableName]
		if table.options["skip"] != nil {
			continue
		}

		if strings2.EndsWith(tableName, options.TypeTableSuffix) {
			t := m.getTypeTableDescription(table)
			dd.Tables = append(dd.Tables, t)
		} else if strings2.EndsWith(tableName, options.AssociationTableSuffix) {
			if mm, ok := m.getManyManyDescription(table, options.TypeTableSuffix); ok {
				dd.MM = append(dd.MM, mm)
			}
		} else {
			t := m.getTableDescription(table)
			dd.Tables = append(dd.Tables, t)
		}
	}
	return dd
}

func (m *DB) getTableDescription(t mysqlTable) db.TableDescription {
	var columnDescriptions []db.ColumnDescription

	// Build the indexes
	pkColumns := make(map[string]bool)
	indexes := make(map[string]*db.IndexDescription)
	uniqueColumns := make(map[string]bool)

	// Fill pkColumns map with the column names of all the pk columns
	// Also file the indexes map with a list of columns for each index
	for _, idx := range t.indexes {
		if idx.name == "PRIMARY" {
			pkColumns[idx.columnName] = true
		} else if i, ok2 := indexes[idx.name]; ok2 {
			i.ColumnNames = append(i.ColumnNames, idx.columnName)
			sort.Strings(i.ColumnNames) // make sure this list stays in a predictable order each time
		} else {
			i = &db.IndexDescription{IsUnique: !idx.nonUnique, ColumnNames: []string{idx.columnName}}
			indexes[idx.name] = i
		}
	}

	// File the uniqueColumns map with all the columns that have a single unique index,
	// including any PK columns. Single indexes are used to determine 1 to 1 relationships.
	for _, i := range indexes {
		if len(i.ColumnNames) == 1 && i.IsUnique {
			uniqueColumns[i.ColumnNames[0]] = true
		}
	}
	if len(pkColumns) == 1 {
		for k, _ := range pkColumns {
			uniqueColumns[k] = true
		}
	}

	var pkCount int
	for _, col := range t.columns {
		cd := m.getColumnDescription(t, col, pkColumns[col.name], uniqueColumns[col.name])

		if cd.IsPk {
			// private keys go first
			// the following code does an insert after whatever previous pks have been found.
			// It is important to do these in order.
			columnDescriptions = append(columnDescriptions, db.ColumnDescription{})
			copy(columnDescriptions[pkCount+1:], columnDescriptions[pkCount:])
			columnDescriptions[pkCount] = cd
			pkCount++
		} else {
			columnDescriptions = append(columnDescriptions, cd)
		}
	}

	td := db.TableDescription{
		Name:                t.name,
		Columns:             columnDescriptions,
		SupportsForeignKeys: t.supportsForeignKeys,
	}

	td.Comment = t.comment
	td.Options = t.options

	// Create the indexes array in index name order so its predictable
	stringmap.Range(indexes, func(key string, val interface{}) bool {
		td.Indexes = append(td.Indexes, *(val.(*db.IndexDescription)))
		return true
	})
	return td
}

func (m *DB) getTypeTableDescription(t mysqlTable) db.TableDescription {
	td := m.getTableDescription(t)

	var columnNames []string
	var columnTypes []GoColumnType

	for i, c := range td.Columns {
		columnNames = append(columnNames, c.Name)
		colType := ColTypeFromGoTypeString(c.GoType)
		if i == 0 {
			colType = ColTypeInteger // Force first value to be treated like an integer
		}
		columnTypes = append(columnTypes, colType)
	}

	result, err := m.SqlDb().Query(`
	SELECT ` +
		"`" + strings.Join(columnNames, "`,`") + "`" +
		`
	FROM ` +
		"`" + td.Name + "`" +
		` ORDER BY ` + "`" + columnNames[0] + "`")

	if err != nil {
		log.Fatal(err)
	}

	values := sql2.SqlReceiveRows(result, columnTypes, columnNames, nil)
	td.TypeData = values
	return td
}

func (m *DB) getColumnDescription(table mysqlTable, column mysqlColumn, isPk bool, isUnique bool) db.ColumnDescription {
	cd := db.ColumnDescription{
		Name: column.name,
	}

	m.processTypeInfo(table.name, column, &cd)

	cd.IsId = strings.Contains(column.extra, "auto_increment")
	cd.IsPk = isPk
	cd.IsNullable = column.isNullable == "YES"
	cd.IsUnique = isUnique

	// indicates that the database is handling update on modify
	// In MySQL this is detectable. In other databases, if you can set this up, but its hard to detect, you can create a comment property to spec this
	if strings.Contains(column.extra, "CURRENT_TIMESTAMP") {
		cd.SubType = "timestamp"
	}

	cd.Comment = column.comment
	cd.Options = column.options

	if fk, ok2 := table.fkMap[cd.Name]; ok2 {
		cd.ForeignKey = &db.ForeignKeyDescription{
			ReferencedTable:  fk.referencedTableName.String,
			ReferencedColumn: fk.referencedColumnName.String,
			UpdateAction:     fkRuleToAction(fk.updateRule),
			DeleteAction:     fkRuleToAction(fk.deleteRule),
		}
	}

	return cd
}

func (m *DB) getManyManyDescription(t mysqlTable, typeTableSuffix string) (mm db.ManyManyDescription, ok bool) {
	td := m.getTableDescription(t)
	if len(td.Columns) != 2 {
		log.Print("Error: table " + td.Name + " must have only 2 primary key columns.")
		return
	}
	var typeIndex = -1
	for i, cd := range td.Columns {
		if !cd.IsPk {
			log.Print("Error: column " + td.Name + ":" + cd.Name + " must be a primary key.")
			return
		}

		if cd.ForeignKey == nil {
			log.Print("Error: column " + td.Name + ":" + cd.Name + " must be a foreign key.")
			return
		}

		if cd.ForeignKey.DeleteAction != db.FKActionCascade {
			log.Print("Warning: column " + td.Name + ":" + cd.Name + " has a DELETE action that is not CASCADE. You will need to manually delete the relationship before the associated object is deleted.")
		}

		if cd.IsNullable {
			log.Print("Error: column " + td.Name + ":" + cd.Name + " cannot be nullable.")
			return
		}

		if strings2.EndsWith(cd.ForeignKey.ReferencedTable, typeTableSuffix) {
			if typeIndex != -1 {
				log.Print("Error: column " + td.Name + ":" + " cannot have two foreign keys to type tables.")
				return
			}
			typeIndex = i
		}
	}

	idx1 := 0
	idx2 := 1
	if typeIndex == 0 {
		idx1 = 1
		idx2 = 0
	}
	options, _, _ := sql2.ExtractOptions(t.columns[idx1].comment)
	mm.Table1 = td.Columns[idx1].ForeignKey.ReferencedTable
	mm.Column1 = td.Columns[idx1].Name
	if opt := options["goName"]; opt != nil {
		if mm.GoName1, ok = opt.(string); !ok {
			log.Print("Error in table comment for table " + t.name + ":" + t.columns[idx1].name + ": goName is not a string")
			return
		}
	}
	if opt := options["goPlural"]; opt != nil {
		if mm.GoPlural1, ok = opt.(string); !ok {
			log.Print("Error in table comment for table " + t.name + ":" + t.columns[idx1].name + ": goPlural is not a string")
			return
		}
	}

	options, _, _ = sql2.ExtractOptions(t.columns[idx2].comment)
	mm.Table2 = td.Columns[idx2].ForeignKey.ReferencedTable
	mm.Column2 = td.Columns[idx2].Name
	if opt := options["goName"]; opt != nil {
		if mm.GoName2, ok = opt.(string); !ok {
			log.Print("Error in table comment for table " + t.name + ":" + t.columns[idx2].name + ": goName is not a string")
			return
		}
	}
	if opt := options["goPlural"]; opt != nil {
		if mm.GoPlural2, ok = opt.(string); !ok {
			log.Print("Error in table comment for table " + t.name + ":" + t.columns[idx2].name + ": goPlural is not a string")
			return
		}
	}

	mm.AssnTableName = t.name
	mm.SupportsForeignKeys = t.supportsForeignKeys
	ok = true
	return
}

func fkRuleToAction(rule sql.NullString) db.FKAction {

	if !rule.Valid {
		return db.FKActionNone // This means we will emulate foreign key actions
	}
	switch strings.ToUpper(rule.String) {
	case "NO ACTION":
		fallthrough
	case "RESTRICT":
		return db.FKActionRestrict
	case "CASCADE":
		return db.FKActionCascade
	case "SET DEFAULT":
		return db.FKActionSetDefault
	case "SET NULL":
		return db.FKActionSetNull

	}
	return db.FKActionNone
}
