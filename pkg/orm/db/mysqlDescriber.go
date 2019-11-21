package db

import (
	"database/sql"
	"fmt"
	. "github.com/goradd/goradd/pkg/orm/query"
	strings2 "github.com/goradd/goradd/pkg/strings"
	"log"
	"math"
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
	options map[string]interface{}
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
	options         map[string]interface{}
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
	description := m.descriptionFromRawTables(rawTables)
	m.goraddDatabase = NewDatabase(m.dbKey, m.idSuffix, description)
}

func (m *Mysql5) getRawTables() map[string]mysqlTable {
	var tableMap = make(map[string]mysqlTable)

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
	var tables []mysqlTable

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
		err = rows.Scan(&tableName, &tableComment)
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
		if table.options, table.comment, err = extractOptions(table.comment); err != nil {
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
		err = rows.Scan(&col.name, &col.defaultValue.R, &col.isNullable, &col.dataType, &col.characterMaxLen, &col.columnType, &col.key, &col.extra, &col.comment)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		if col.options, col.comment, err = extractOptions(col.comment); err != nil {
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
		err = rows.Scan(&index.name, &index.nonUnique, &index.columnName)
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
		err = rows.Scan(&fk.name, &fk.columnName, &fk.referencedSchema, &fk.referencedTableName, &fk.referencedColumnName, &fk.updateRule, &fk.deleteRule)
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
		cd.GoType = ColTypeDateTime.GoType()
	case "timestamp":
		cd.NativeType = SqlTypeTimestamp
		cd.GoType = ColTypeDateTime.GoType()
		cd.IsTimestamp = true
	case "datetime":
		cd.NativeType = SqlTypeDatetime
		cd.GoType = ColTypeDateTime.GoType()
	case "date":
		cd.NativeType = SqlTypeDate
		cd.GoType = ColTypeDateTime.GoType()
	case "tinyint":
		if dataLen == 1 {
			cd.NativeType = SqlTypeBool
			cd.GoType = ColTypeBool.GoType()
		} else {
			if isUnsigned {
				cd.NativeType = SqlTypeInteger
				cd.GoType = ColTypeUnsigned.GoType()
				min, max := getMinMax(column.options, 0, 255, tableName, column.name)
				cd.MinValue = uint(min)
				cd.MaxValue = uint(max)
				cd.MaxCharLength = 3
			} else {
				cd.NativeType = SqlTypeInteger
				cd.GoType = ColTypeInteger.GoType()
				min, max := getMinMax(column.options, -128, 127, tableName, column.name)
				cd.MinValue = int(min)
				cd.MaxValue = int(max)
				cd.MaxCharLength = 4 // allow for a negative sign
			}
		}

	case "int":
		if isUnsigned {
			cd.NativeType = SqlTypeInteger
			cd.GoType = ColTypeUnsigned.GoType()
			min, max := getMinMax(column.options, 0, 4294967295, tableName, column.name)
			cd.MinValue = uint(min)
			cd.MaxValue = uint(max)
			cd.MaxCharLength = 10
		} else {
			cd.NativeType = SqlTypeInteger
			cd.GoType = ColTypeInteger.GoType()
			min, max := getMinMax(column.options, -2147483648, 2147483647, tableName, column.name)
			cd.MinValue = int(min)
			cd.MaxValue = int(max)
			cd.MaxCharLength = 11
		}

	case "smallint":
		if isUnsigned {
			cd.NativeType = SqlTypeInteger
			cd.GoType = ColTypeUnsigned.GoType()
			min, max := getMinMax(column.options, 0, 65535, tableName, column.name)
			cd.MinValue = uint(min)
			cd.MaxValue = uint(max)
			cd.MaxCharLength = 5
		} else {
			cd.NativeType = SqlTypeInteger
			cd.GoType = ColTypeInteger.GoType()
			min, max := getMinMax(column.options, -32768, 32767, tableName, column.name)
			cd.MinValue = int(min)
			cd.MaxValue = int(max)
			cd.MaxCharLength = 6
		}

	case "mediumint":
		if isUnsigned {
			cd.NativeType = SqlTypeInteger
			cd.GoType = ColTypeUnsigned.GoType()
			min, max := getMinMax(column.options, 0, 16777215, tableName, column.name)
			cd.MinValue = uint(min)
			cd.MaxValue = uint(max)
			cd.MaxCharLength = 8
		} else {
			cd.NativeType = SqlTypeInteger
			cd.GoType = ColTypeInteger.GoType()
			min, max := getMinMax(column.options, -8388608, 8388607, tableName, column.name)
			cd.MinValue = int(min)
			cd.MaxValue = int(max)
			cd.MaxCharLength = 8
		}

	case "bigint": // We need to be explicit about this in go, since int will be whatever the OS native int size is, but go will support int64 always.
		// Also, since Json can only be decoded into float64s, we are limited in our ability to represent large min and max numbers in the json to about 2^53
		if isUnsigned {
			cd.NativeType = SqlTypeInteger
			cd.GoType = ColTypeUnsigned64.GoType()

			if v := column.options["min"]; v != nil {
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

			if v := column.options["max"]; v != nil {
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
			cd.GoType = ColTypeInteger64.GoType()
			if v := column.options["min"]; v != nil {
				if v2, ok := v.(float64); !ok {
					log.Print("Error in min value in comment for table " + tableName + ":" + column.name + ". Value is not a valid number.")
					cd.MinValue = int64(math.MinInt64)
				} else {
					cd.MinValue = int64(v2)
				}
			} else {
				cd.MinValue = int64(math.MinInt64)
			}

			if v := column.options["max"]; v != nil {
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
		cd.GoType = ColTypeFloat.GoType()
		cd.MinValue, cd.MaxValue = getMinMax(column.options, -math.MaxFloat32, math.MaxFloat32, tableName, column.name)
	case "double":
		cd.NativeType = SqlTypeDouble
		cd.GoType = ColTypeDouble.GoType()
		cd.MinValue, cd.MaxValue = getMinMax(column.options, -math.MaxFloat64, math.MaxFloat64, tableName, column.name)

	case "varchar":
		cd.NativeType = SqlTypeVarchar
		cd.GoType = ColTypeString.GoType()
		cd.MaxCharLength = uint64(dataLen)

	case "char":
		cd.NativeType = SqlTypeChar
		cd.GoType = ColTypeString.GoType()
		cd.MaxCharLength = uint64(dataLen)

	case "blob":
		cd.NativeType = SqlTypeBlob
		cd.GoType = ColTypeBytes.GoType()
		cd.MaxCharLength = 65535
	case "tinyblob":
		cd.NativeType = SqlTypeBlob
		cd.GoType = ColTypeBytes.GoType()
		cd.MaxCharLength = 255
	case "mediumblob":
		cd.NativeType = SqlTypeBlob
		cd.GoType = ColTypeBytes.GoType()
		cd.MaxCharLength = 16777215
	case "longblob":
		cd.NativeType = SqlTypeBlob
		cd.GoType = ColTypeBytes.GoType()
		cd.MaxCharLength = math.MaxUint32

	case "text":
		cd.NativeType = SqlTypeText
		cd.GoType = ColTypeString.GoType()
		cd.MaxCharLength = 65535
	case "tinytext":
		cd.NativeType = SqlTypeText
		cd.GoType = ColTypeString.GoType()
		cd.MaxCharLength = 255
	case "mediumtext":
		cd.NativeType = SqlTypeText
		cd.GoType = ColTypeString.GoType()
		cd.MaxCharLength = 16777215
	case "longtext":
		cd.NativeType = SqlTypeText
		cd.GoType = ColTypeString.GoType()
		cd.MaxCharLength = math.MaxUint32

	case "decimal": // No native equivalent in Go. See the "Big" go package for support. You will need to shephard numbers into and out of string format to move data to the database
		cd.NativeType = SqlTypeDecimal
		cd.GoType = ColTypeString.GoType()
		cd.MaxCharLength = uint64(dataLen) + 3

	case "year":
		cd.NativeType = SqlTypeInteger
		cd.GoType = ColTypeInteger.GoType()

	case "set":
		log.Print("Note: Using association tables is preferred to using Mysql5 SET columns in table " + tableName + ":" + column.name + ".")
		cd.NativeType = MysqlTypeSet
		cd.GoType = ColTypeString.GoType()
		cd.MaxCharLength = uint64(column.characterMaxLen.Int64)

	case "enum":
		log.Print("Note: Using type tables is preferred to using Mysql5 ENUM columns in table " + tableName + ":" + column.name + ".")
		cd.NativeType = MysqlTypeEnum
		cd.GoType = ColTypeString.GoType()
		cd.MaxCharLength = uint64(column.characterMaxLen.Int64)

	default:
		cd.NativeType = SqlTypeUnknown
		cd.GoType = ColTypeString.GoType()
	}

	cd.DefaultValue = column.defaultValue.Unpack(ColTypeFromGoTypeString(cd.GoType))
}

func (m *Mysql5) descriptionFromRawTables(rawTables map[string]mysqlTable) DatabaseDescription {

	dd := DatabaseDescription{Key: m.dbKey, AssociatedObjectPrefix: m.associatedObjectPrefix}

	for tableName, table := range rawTables {
		if table.options["skip"] != nil {
			continue
		}

		if strings2.EndsWith(tableName, m.TypeTableSuffix()) {
			t := m.getTypeTableDescription(table)
			dd.Tables = append(dd.Tables, t)
		} else if strings2.EndsWith(tableName, m.AssociationTableSuffix()) {
			if mm,ok := m.getManyManyDescription(table); ok {
				dd.MM = append(dd.MM, mm)
			}
		} else {
			t := m.getTableDescription(table)
			dd.Tables = append(dd.Tables, t)
		}
	}
	return dd
}

func (m *Mysql5) getTableDescription(t mysqlTable) TableDescription {
	var columnDescriptions []ColumnDescription

	var pkCount int
	for _, col := range t.columns {
		cd := m.getColumnDescription(t, col)

		if cd.IsPk {
			// private keys go first
			// the following code does an insert after whatever previous pks have been found. Its important to do these in order.
			columnDescriptions = append(columnDescriptions, ColumnDescription{})
			copy(columnDescriptions[pkCount+1:], columnDescriptions[pkCount:])
			columnDescriptions[pkCount] = cd
			pkCount++
		} else {
			columnDescriptions = append(columnDescriptions, cd)
		}
	}

	td := TableDescription{
		Name:        t.name,
		Columns: columnDescriptions,
	}

	var ok bool
	if opt := t.options["literalName"]; opt != nil {
		if td.LiteralName, ok = opt.(string); !ok {
			log.Print("Error in table comment for table " + t.name + ": literalName is not a string")
		}
		delete(t.options, "literalName")
	}

	if opt := t.options["literalPlural"]; opt != nil {
		if td.LiteralPlural, ok = opt.(string); !ok {
			log.Print("Error in table comment for table " + t.name + ": literalPlural is not a string")
		}
		delete(t.options, "literalPlural")
	}

	if opt := t.options["goName"]; opt != nil {
		if td.GoName, ok = opt.(string); !ok {
			log.Print("Error in table comment for table " + t.name + ": goName is not a string")
		} else {
			td.GoName = strings.Title(td.GoName)
		}
		delete(t.options, "goName")
	}

	if opt := t.options["goPlural"]; opt != nil {
		if td.GoPlural, ok = opt.(string); !ok {
			log.Print("Error in table comment for table " + t.name + ": goPlural is not a string")
		} else {
			td.GoName = strings.Title(td.GoPlural)
		}
		delete(t.options, "goPlural")
	}

	td.Comment = t.comment
	td.Options = t.options

	// Build the indexes
	indexes := make(map[string]*IndexDescription)
	for _,idx := range t.indexes {
		if i,ok := indexes[idx.name]; ok {
			i.ColumnNames = append(i.ColumnNames, idx.columnName)
		} else {
			i = &IndexDescription{IsUnique:!idx.nonUnique, ColumnNames: []string{idx.columnName}}
			indexes[idx.name] = i
		}
	}
	for _,iDesc := range indexes {
		td.Indexes = append(td.Indexes, *iDesc)
	}
	return td
}


func (m *Mysql5) getTypeTableDescription(t mysqlTable) TableDescription {
	td := m.getTableDescription(t)

	var columnNames []string
	var columnTypes []GoColumnType
	columnTypes2 := map[string]GoColumnType{}

	for _, c := range t.columns {
		columnNames = append(columnNames, c.name)
		columnTypes = append(columnTypes, c.goType)
		columnTypes2[c.name] = c.goType
	}

	result, err := m.db.Query(`
	SELECT ` +
		strings.Join(columnNames, ",") +
		`
	FROM ` +
		td.Name +
		` ORDER BY ` + columnNames[0])

	if err != nil {
		log.Fatal(err)
	}
	defer result.Close()

	values := ReceiveRows(result, columnTypes, columnNames)
	td.TypeData = values
	return td
}

func (m *Mysql5) getColumnDescription(table mysqlTable, column mysqlColumn) ColumnDescription {
	cd := ColumnDescription {
		Name: column.name,
	}
	var ok bool
	var shouldAutoUpdate bool

	if opt := column.options["goName"]; opt != nil {
		if cd.GoName, ok = opt.(string); !ok {
			log.Print("Error in table comment for table " + table.name + ":" + column.name + ": goName is not a string")
		}
		delete(column.options,"goName")
	}

	if opt := column.options["shouldAutoUpdate"]; opt != nil {
		if shouldAutoUpdate, ok = opt.(bool); !ok {
			log.Print("Error in table comment for table " + table.name + ":" + column.name + ": shouldAutoUpdate is not a boolean")
		}
		delete(column.options,"shouldAutoUpdate")
	}

	m.processTypeInfo(table.name, column, &cd)

	cd.IsId = strings.Contains(column.extra, "auto_increment")
	cd.IsPk = (column.key == "PRI")
	cd.IsNullable = (column.isNullable == "YES")
	cd.IsUnique = (column.key == "UNI")

	// indicates that the database is handling update on modify
	// In MySQL this is detectable. In other databases, if you can set this up, but its hard to detect, you can create a comment property to spec this
	cd.IsAutoUpdateTimestamp = strings.Contains(column.extra, "CURRENT_TIMESTAMP")

	if shouldAutoUpdate {
		cd.IsTimestamp = true
	}

	if cd.IsAutoUpdateTimestamp && shouldAutoUpdate {
		log.Print("Error in table comment for table " + table.name + ":" + column.name + ": shouldAutoUpdate should not be set on a table that the database is automatically updating.")
	}

	cd.Comment = column.comment
	cd.Options = column.options

	if fk, ok2 := table.fkMap[cd.Name]; ok2 {
		cd.ForeignKey = &ForeignKeyDescription{
			ReferencedTable:    fk.referencedTableName.String,
			ReferencedColumn:   fk.referencedColumnName.String,
			UpdateAction: fkRuleToAction(fk.updateRule).String(),
			DeleteAction: fkRuleToAction(fk.deleteRule).String(),
		}
	}

	return cd
}


func (m *Mysql5) getManyManyDescription(t mysqlTable) (mm ManyManyDescription, ok bool) {
	td := m.getTableDescription(t)
	if len(td.Columns) != 2 {
		log.Print("Error: table " + td.Name + " must have only 2 primary key columns.")
		return
	}
	var typeIndex = -1
	for i, cd := range td.Columns {
		if !cd.IsPk {
			log.Print("Error: table " + td.Name + ":" + cd.Name + " must be a primary key.")
			return
		}

		if cd.ForeignKey == nil {
			log.Print("Error: table " + td.Name + ":" + cd.Name + " must be a foreign key.")
			return
		}

		if cd.IsNullable {
			log.Print("Error: table " + td.Name + ":" + cd.Name + " cannot be nullable.")
			return
		}

		if strings2.EndsWith(cd.ForeignKey.ReferencedTable, m.TypeTableSuffix()) {
			if typeIndex != -1 {
				log.Print("Error: table " + td.Name + ":" + " cannot have two foreign keys to type tables.")
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
	options,_,_ := extractOptions(t.columns[idx1].comment)
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

	options,_,_ = extractOptions(t.columns[idx2].comment)
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
	ok = true
	return
}

