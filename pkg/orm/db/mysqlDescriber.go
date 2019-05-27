package db

import (
	"database/sql"
	"fmt"
	"github.com/goradd/gengen/pkg/maps"
	. "github.com/goradd/goradd/pkg/orm/query"
	"github.com/knq/snaker"
	"log"
	"math"
	"sort"
	"strings"
)

/*
This file contains the code that parses the data structure found in a MySQL database into
our own cross-platform internal database description object.
 */

const (
	MysqlTypeSet  = "Set"
	MysqlTypeEnum = "Enum"
)

type mysqlTable struct {
	name    string
	columns []mysqlColumn
	indexes []mysqlIndex
	fkMap   map[string]mysqlForeignKey
	comment string
}

type mysqlColumn struct {
	name            string
	defaultValue    SqlReceiver
	isNullable      string
	dataType        string
	goType          GoColumnType
	dataLen         int
	characterMaxLen sql.NullInt64
	columnType      string
	key             string
	extra           string
	comment         string
	options         *maps.SliceMap
}

type mysqlIndex struct {
	name       string
	nonUnique  bool
	columnName string
}

type mysqlForeignKey struct {
	name                 string
	columnName           string
	referencedSchema     sql.NullString
	referencedTableName  sql.NullString
	referencedColumnName sql.NullString
	updateRule           sql.NullString
	deleteRule           sql.NullString
}

/*
type Mysql5 struct {
	dbKey   string
	db      *sql.DB
	config  mysql.Config
	options DbOptions
}



func NewMysql2 (dbKey string , options DbOptions, config *mysql.Config) (*Mysql5, error) {
	source := Mysql5{}
	db, err := sql.Open("mysql", config.FormatDSN())

	source.dbKey = DbKey
	source.db = db
	source.config = *config
	source.options = options

	// Ping?

	return &source,err
}*/

func (m *Mysql5) loadDescription() {
	rawTables := m.getRawTables()
	m.description = m.descriptionFromRawTables(rawTables)
	m.description.analyze()

}

func (m *Mysql5) getRawTables() map[string]mysqlTable {
	var tableMap map[string]mysqlTable = make(map[string]mysqlTable)

	tables := m.getTables()

	for _, table := range tables {
		indexes, err := m.getIndexes(table.name)
		if err != nil {
			return nil
		}

		foreignKeys, err := m.getForeignKeys(table.name)
		if err != nil {
			return nil
		}

		// Do some processing on the foreign keys
		for _, fk := range foreignKeys {
			if fk.referencedColumnName.Valid && fk.referencedTableName.Valid {
				if _, ok := table.fkMap[fk.columnName]; ok {
					log.Printf("Warning: Column %s:%s multi-table foreign keys are not supported.", table.name, fk.columnName)
					delete(table.fkMap, fk.columnName)
				} else {
					table.fkMap[fk.columnName] = fk
				}
			}
		}

		columns, err := m.getColumns(table.name)
		if err != nil {
			return nil
		}

		table.indexes = indexes
		table.columns = columns
		tableMap[table.name] = table
	}

	return tableMap

}

// Gets some of the information for a table
func (m *Mysql5) getTables() []mysqlTable {
	var tableName, tableComment string
	tables := []mysqlTable{}

	// Use the MySQL5 Information Schema to get a list of all the tables in this database
	// (excluding views, etc.)
	dbName := m.config.DBName

	rows, err := m.db.Query(fmt.Sprintf(`
	SELECT
	table_name,
	table_comment
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
		err := rows.Scan(&tableName, &tableComment)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(tableName)
		table := mysqlTable{
			name:    tableName,
			comment: tableComment,
			columns: []mysqlColumn{},
			fkMap:   make(map[string]mysqlForeignKey),
			indexes: []mysqlIndex{},
		}
		tables = append(tables, table)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return tables
}

func (m *Mysql5) getColumns(table string) (columns []mysqlColumn, err error) {

	dbName := m.config.DBName

	rows, err := m.db.Query(fmt.Sprintf(`
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
	table_schema = '%s';
	`, table, dbName))

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var col mysqlColumn

	for rows.Next() {
		col = mysqlColumn{}
		err := rows.Scan(&col.name, &col.defaultValue.R, &col.isNullable, &col.dataType, &col.characterMaxLen, &col.columnType, &col.key, &col.extra, &col.comment)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		if col.options, err = extractOptions(col.comment); err != nil {
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

func (m *Mysql5) getIndexes(table string) (indexes []mysqlIndex, err error) {

	dbName := m.config.DBName

	rows, err := m.db.Query(fmt.Sprintf(`
	SELECT
	index_name,
	non_unique,
	column_name
	FROM
	information_schema.statistics
	WHERE
	table_name = '%s' AND
	table_schema = '%s';
	`, table, dbName))

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var index mysqlIndex

	for rows.Next() {
		index = mysqlIndex{}
		err := rows.Scan(&index.name, &index.nonUnique, &index.columnName)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		indexes = append(indexes, index)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return indexes, err
}

func (m *Mysql5) getForeignKeys(table string) (foreignKeys []mysqlForeignKey, err error) {
	dbName := m.config.DBName

	rows, err := m.db.Query(fmt.Sprintf(`
	SELECT
	k.constraint_name,
	k.column_name,
	k.referenced_table_schema,
	k.referenced_table_name,
	k.referenced_column_name,
	r.update_rule,
	r.delete_rule
	FROM
	information_schema.key_column_usage as k
	left join information_schema.referential_constraints as r on r.constraint_schema = k.table_schema AND r.constraint_name = k.constraint_name
	WHERE
	k.table_name = '%s' AND
	k.table_schema = '%s';
	`, table, dbName))

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var fk mysqlForeignKey

	for rows.Next() {
		fk = mysqlForeignKey{}
		err := rows.Scan(&fk.name, &fk.columnName, &fk.referencedSchema, &fk.referencedTableName, &fk.referencedColumnName, &fk.updateRule, &fk.deleteRule)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		foreignKeys = append(foreignKeys, fk)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return foreignKeys, err
}

// Convert the database native type to a more generic sql type, and a go table type.
func (m *Mysql5) processTypeInfo(tableName string, column mysqlColumn, cd *ColumnDescription) {
	dataLen := getDataDefLength(column.columnType)

	isUnsigned := strings.Contains(column.columnType, "unsigned")

	switch column.dataType {
	case "time":
		cd.NativeType = SqlTypeTime
		cd.ColumnType = ColTypeDateTime
	case "timestamp":
		cd.NativeType = SqlTypeTimestamp
		cd.ColumnType = ColTypeDateTime
		cd.IsTimestamp = true
	case "datetime":
		cd.NativeType = SqlTypeDatetime
		cd.ColumnType = ColTypeDateTime
	case "date":
		cd.NativeType = SqlTypeDate
		cd.ColumnType = ColTypeDateTime
	case "tinyint":
		if dataLen == 1 {
			cd.NativeType = SqlTypeBool
			cd.ColumnType = ColTypeBool
		} else {
			if isUnsigned {
				cd.NativeType = SqlTypeInteger
				cd.ColumnType = ColTypeUnsigned
				min, max := getMinMax(column.options, 0, 255, tableName, column.name)
				cd.MinValue = uint(min)
				cd.MaxValue = uint(max)
				cd.MaxCharLength = 3
			} else {
				cd.NativeType = SqlTypeInteger
				cd.ColumnType = ColTypeInteger
				min, max := getMinMax(column.options, -128, 127, tableName, column.name)
				cd.MinValue = int(min)
				cd.MaxValue = int(max)
				cd.MaxCharLength = 4 // allow for a negative sign
			}
		}

	case "int":
		if isUnsigned {
			cd.NativeType = SqlTypeInteger
			cd.ColumnType = ColTypeUnsigned
			min, max := getMinMax(column.options, 0, 4294967295, tableName, column.name)
			cd.MinValue = uint(min)
			cd.MaxValue = uint(max)
			cd.MaxCharLength = 10
		} else {
			cd.NativeType = SqlTypeInteger
			cd.ColumnType = ColTypeInteger
			min, max := getMinMax(column.options, -2147483648, 2147483647, tableName, column.name)
			cd.MinValue = int(min)
			cd.MaxValue = int(max)
			cd.MaxCharLength = 11
		}

	case "smallint":
		if isUnsigned {
			cd.NativeType = SqlTypeInteger
			cd.ColumnType = ColTypeUnsigned
			min, max := getMinMax(column.options, 0, 65535, tableName, column.name)
			cd.MinValue = uint(min)
			cd.MaxValue = uint(max)
			cd.MaxCharLength = 5
		} else {
			cd.NativeType = SqlTypeInteger
			cd.ColumnType = ColTypeInteger
			min, max := getMinMax(column.options, -32768, 32767, tableName, column.name)
			cd.MinValue = int(min)
			cd.MaxValue = int(max)
			cd.MaxCharLength = 6
		}

	case "mediumint":
		if isUnsigned {
			cd.NativeType = SqlTypeInteger
			cd.ColumnType = ColTypeUnsigned
			min, max := getMinMax(column.options, 0, 16777215, tableName, column.name)
			cd.MinValue = uint(min)
			cd.MaxValue = uint(max)
			cd.MaxCharLength = 8
		} else {
			cd.NativeType = SqlTypeInteger
			cd.ColumnType = ColTypeInteger
			min, max := getMinMax(column.options, -8388608, 8388607, tableName, column.name)
			cd.MinValue = int(min)
			cd.MaxValue = int(max)
			cd.MaxCharLength = 8
		}

	case "bigint": // We need to be explicit about this in go, since int will be whatever the OS native int size is, but go will support int64 always.
		// Also, since Json can only be decoded into float64s, we are limited in our ability to represent large min and max numbers in the json to about 2^53
		if isUnsigned {
			cd.NativeType = SqlTypeInteger
			cd.ColumnType = ColTypeUnsigned64

			if v := column.options.Get("min"); v != nil {
				if v2, ok := v.(float64); !ok {
					log.Print("Error in min value in comment for table " + tableName + ":" + column.name + ". Value is not a valid number.")
					cd.MinValue = uint64(0)
				} else if v2 < 0 {
					log.Print("Error in min value in comment for table " + tableName + ":" + column.name + ". Value cannot be less than zero.")
					cd.MinValue = uint64(0)
				} else {
					cd.MinValue = uint64(v2)
				}
			} else {
				cd.MinValue = 0
			}

			if v := column.options.Get("max"); v != nil {
				if v2, ok := v.(float64); !ok {
					log.Print("Error in max value in comment for table " + tableName + ":" + column.name + ". Value is not a valid number.")
					cd.MaxValue = uint64(math.MaxUint64)
				} else {
					cd.MaxValue = int64(v2)
				}
			} else {
				cd.MaxValue = uint64(math.MaxUint64)
			}
			cd.MaxCharLength = 20
		} else {
			cd.NativeType = SqlTypeInteger
			cd.ColumnType = ColTypeInteger64
			if v := column.options.Get("min"); v != nil {
				if v2, ok := v.(float64); !ok {
					log.Print("Error in min value in comment for table " + tableName + ":" + column.name + ". Value is not a valid number.")
					cd.MinValue = int64(math.MinInt64)
				} else {
					cd.MinValue = int64(v2)
				}
			} else {
				cd.MinValue = int64(math.MinInt64)
			}

			if v := column.options.Get("max"); v != nil {
				if v2, ok := v.(float64); !ok {
					log.Print("Error in max value in comment for table " + tableName + ":" + column.name + ". Value is not a valid number.")
					cd.MaxValue = int64(math.MaxInt64)
				} else {
					cd.MaxValue = int64(v2)
				}
			} else {
				cd.MaxValue = int64(math.MaxInt64)
			}
			cd.MaxCharLength = 20
		}

	case "float":
		cd.NativeType = SqlTypeFloat
		cd.ColumnType = ColTypeFloat
		cd.MinValue, cd.MaxValue = getMinMax(column.options, -math.MaxFloat32, math.MaxFloat32, tableName, column.name)
	case "double":
		cd.NativeType = SqlTypeDouble
		cd.ColumnType = ColTypeDouble
		cd.MinValue, cd.MaxValue = getMinMax(column.options, -math.MaxFloat64, math.MaxFloat64, tableName, column.name)

	case "varchar":
		cd.NativeType = SqlTypeVarchar
		cd.ColumnType = ColTypeString
		cd.MaxCharLength = uint64(dataLen)

	case "char":
		cd.NativeType = SqlTypeChar
		cd.ColumnType = ColTypeString
		cd.MaxCharLength = uint64(dataLen)

	case "blob":
		cd.NativeType = SqlTypeBlob
		cd.ColumnType = ColTypeBytes
		cd.MaxCharLength = 65535
	case "tinyblob":
		cd.NativeType = SqlTypeBlob
		cd.ColumnType = ColTypeBytes
		cd.MaxCharLength = 255
	case "mediumblob":
		cd.NativeType = SqlTypeBlob
		cd.ColumnType = ColTypeBytes
		cd.MaxCharLength = 16777215
	case "longblob":
		cd.NativeType = SqlTypeBlob
		cd.ColumnType = ColTypeBytes
		cd.MaxCharLength = math.MaxUint32

	case "text":
		cd.NativeType = SqlTypeText
		cd.ColumnType = ColTypeString
		cd.MaxCharLength = 65535
	case "tinytext":
		cd.NativeType = SqlTypeText
		cd.ColumnType = ColTypeString
		cd.MaxCharLength = 255
	case "mediumtext":
		cd.NativeType = SqlTypeText
		cd.ColumnType = ColTypeString
		cd.MaxCharLength = 16777215
	case "longtext":
		cd.NativeType = SqlTypeText
		cd.ColumnType = ColTypeString
		cd.MaxCharLength = math.MaxUint32

	case "decimal": // No native equivalent in Go. See the "Big" go package for support. You will need to shephard numbers into and out of string format to move data to the database
		cd.NativeType = SqlTypeDecimal
		cd.ColumnType = ColTypeString
		cd.MaxCharLength = uint64(dataLen) + 3

	case "year":
		cd.NativeType = SqlTypeInteger
		cd.ColumnType = ColTypeInteger

	case "set":
		log.Print("Note: Using association tables is preferred to using Mysql5 SET columns in table " + tableName + ":" + column.name + ".")
		cd.NativeType = MysqlTypeSet
		cd.ColumnType = ColTypeString
		cd.MaxCharLength = uint64(column.characterMaxLen.Int64)

	case "enum":
		log.Print("Note: Using type tables is preferred to using Mysql5 ENUM columns in table " + tableName + ":" + column.name + ".")
		cd.NativeType = MysqlTypeEnum
		cd.ColumnType = ColTypeString
		cd.MaxCharLength = uint64(column.characterMaxLen.Int64)

	default:
		cd.NativeType = SqlTypeUnknown
		cd.ColumnType = ColTypeString
	}

	cd.DefaultValue = column.defaultValue.Unpack(cd.ColumnType)
}

func (m *Mysql5) descriptionFromRawTables(rawTables map[string]mysqlTable) *DatabaseDescription {

	dd := NewDatabaseDescription(m.DbKey(), m.AssociatedObjectPrefix(), m.idSuffix)

tableLoop:
	for tableName, table := range rawTables {
		pkCount := 0
		td := getTableDescription(tableName, table.comment, m)
		if td == nil {
			log.Println("Skipping " + tableName)
			continue
		}

		td.DbKey = m.dbKey

		for _, column := range table.columns {
			cd := m.getColumnDescription(tableName, column, table)
			if cd.IsPk && !td.IsAssociation {
				pkCount++
				if pkCount > 1 {
					log.Println("Error, only association tables may have multiple primary keys. Skipping " + tableName)
					continue tableLoop
				}
			}
			td.Columns = append(td.Columns, cd)
			td.columnMap[cd.DbName] = cd
		}

		td.Indexes = m.getIndexDescriptions(tableName, table.indexes)

		if td.IsType {
			dd.TypeTables = append(dd.TypeTables, m.getTypeTableDescription(td))
		} else {
			dd.Tables = append(dd.Tables, td)
		}
	}

	// sort for consistent looping
	sort.Slice(dd.Tables, func(i,j int) bool {
		return dd.Tables[i].DbName < dd.Tables[j].DbName
	})
	sort.Slice(dd.TypeTables, func(i,j int) bool {
		return dd.TypeTables[i].DbName < dd.TypeTables[j].DbName
	})

	return dd
}

func (m *Mysql5) getTypeTableDescription(td *TableDescription) *TypeTableDescription {
	var pkField string
	for _, col := range td.Columns {
		if col.IsPk {
			pkField = col.GoName
			break
		}
	}

	tt := TypeTableDescription{
		DbKey:         td.DbKey,
		DbName:        td.DbName,
		EnglishName:   td.LiteralName,
		EnglishPlural: td.LiteralPlural,
		GoName:        td.GoName,
		GoPlural:      td.GoPlural,
		PkField:       pkField,
	}

	columnNames := []string{}
	columnTypes := []GoColumnType{}
	columnTypes2 := map[string]GoColumnType{}

	for _, col := range td.Columns {
		columnNames = append(columnNames, col.DbName)
		columnTypes = append(columnTypes, col.ColumnType)
		columnTypes2[col.DbName] = col.ColumnType
	}

	result, err := m.db.Query(`
	SELECT ` +
		strings.Join(columnNames, ",") +
		`
	FROM ` +
		td.DbName +
		` ORDER BY ` + columnNames[0])

	if err != nil {
		log.Fatal(err)
	}
	defer result.Close()

	tt.Values = ReceiveRows(result, columnTypes, columnNames)
	tt.FieldNames = columnNames
	tt.FieldTypes = columnTypes2

	return &tt
}

func (m *Mysql5) getColumnDescription(tableName string, column mysqlColumn, table mysqlTable) *ColumnDescription {
	options, err := extractOptions(column.comment)
	if err != nil {
		log.Print("Error in table comment for table " + tableName + ":" + column.name + ": " + err.Error())
	}

	cd := ColumnDescription{
		DbName: column.name,
	}
	var ok bool
	opt := options.Get("goName")
	if opt != nil {
		if cd.GoName, ok = opt.(string); !ok {
			log.Print("Error in table comment for table " + tableName + ":" + column.name + ": goName is not a string")
		}
	} else {
		cd.GoName = snaker.SnakeToCamel(column.name)
	}

	//cd.DefaultValue, _ = table.defaultValue.Value()

	m.processTypeInfo(tableName, column, &cd)

	cd.IsId = strings.Contains(column.extra, "auto_increment")
	cd.IsPk = (column.key == "PRI")
	cd.IsNullable = (column.isNullable == "YES")
	cd.IsUnique = (column.key == "UNI")

	// indicates that the database is handling update on modify
	// In MySQL this is detectable. In other databases, if you can set this up, but its hard to detect, you can create a comment property to spec this
	cd.IsAutoUpdateTimestamp = strings.Contains(column.extra, "CURRENT_TIMESTAMP")

	var s bool
	// indicates that we want our generated code to update the timestamp manually. This should be mutually exclusive of isAutoUpdateTimestamp
	if s, ok = options.LoadBool("shouldAutoUpdate"); options.Has("shouldAutoUpdate") && !ok {
		log.Print("Error in table comment for table " + tableName + ":" + column.name + ": shouldAutoUpdate is not a boolean")
	}
	if s {
		cd.IsTimestamp = true
	}

	if cd.IsAutoUpdateTimestamp && s {
		log.Print("Error in table comment for table " + tableName + ":" + column.name + ": shouldAutoUpdate should not be set on a table that the database is automatically updating.")
	}

	cd.Comment = column.comment

	if fk, ok := table.fkMap[cd.DbName]; ok {
		cd.ForeignKey = &ForeignKeyColumn{
			TableName:    fk.referencedTableName.String,
			ColumnName:   fk.referencedColumnName.String,
			UpdateAction: fkRuleToAction(fk.updateRule),
			DeleteAction: fkRuleToAction(fk.deleteRule),
		}
	}

	return &cd
}

// Get index description array. Must preserve the order the indexes appear in the database so that when things change, our
// generated files do not change too much. This is a convenience thing in case these files are checked in to source control.
// For this reason we use an array instead of a map.
func (m *Mysql5) getIndexDescriptions(tableName string, sqlIndexes []mysqlIndex) (indexes []IndexDescription) {
	indexMap := map[string]int{}

	for _, sqlIndex := range sqlIndexes {
		if offset, ok := indexMap[sqlIndex.name]; !ok {
			indexes = append(indexes, IndexDescription{sqlIndex.name, !sqlIndex.nonUnique, sqlIndex.name == "PRIMARY", []string{sqlIndex.columnName}})
			indexMap[sqlIndex.name] = len(indexes) - 1
		} else {
			indexes[offset].ColumnNames = append(indexes[offset].ColumnNames, sqlIndex.columnName)
		}
	}
	return indexes
}

// Process the raw list of foreign keys to return a map of foreign keys organized by table name in the current table.
// Will output any errors it finds too
/*
func (m *Mysql5) getForeignKeyDescriptions(tableName string, sqlFks []mysqlForeignKey) (foreignKeys []ForeignKeyDescription) {
	keyMap := map[string]int{}

	for _,sqlFk := range sqlFks {
		if offset, ok := keyMap[sqlFk.name]; !ok {
			if sqlFk.referencedSchema.Valid {
				foreignKeys = append(foreignKeys, ForeignKeyDescription{
						sqlFk.name,
						[]string{sqlFk.columnName},
					sqlFk.referencedSchema.String,
					sqlFk.referencedTableName.String,
					[]string{sqlFk.referencedColumnName.String},

				})
				keyMap[sqlFk.name] = len(foreignKeys) - 1
			}
		} else {
			foreignKeys[offset].Columns = append(foreignKeys[offset].Columns, sqlFk.columnName)
			foreignKeys[offset].relationColumns = append(foreignKeys[offset].relationColumns, sqlFk.referencedColumnName.String)
		}
	}
	return foreignKeys
}
*/
