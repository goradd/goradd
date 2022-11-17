package pgsql

import (
	"database/sql"
	"fmt"
	"github.com/goradd/goradd/pkg/orm/db"
	sql2 "github.com/goradd/goradd/pkg/orm/db/sql"
	. "github.com/goradd/goradd/pkg/orm/query"
	"github.com/goradd/goradd/pkg/stringmap"
	strings2 "github.com/goradd/goradd/pkg/strings"
	time2 "github.com/goradd/goradd/pkg/time"
	"github.com/goradd/maps"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
)

/*
This file contains the code that parses the data structure found in a Postgresql database into
our own cross-platform internal database description object.
*/

type pgTable struct {
	name    string
	schema  string
	columns []pgColumn
	indexes []pgIndex
	fkMap   map[string]pgForeignKey
	comment string
	options map[string]interface{}
}

type pgColumn struct {
	name            string
	defaultValue    sql.NullString
	isNullable      string
	dataType        string
	charLen         int
	characterMaxLen sql.NullInt64
	comment         string
	options         map[string]interface{}
}

type pgIndex struct {
	name        string
	schema      string
	unique      bool
	primary     bool
	tableName   string
	tableSchema string
	columnName  string
}

type pgForeignKey struct {
	constraintName       string
	tableName            string
	tableSchema          string
	columnName           string
	referencedSchema     sql.NullString
	referencedTableName  sql.NullString
	referencedColumnName sql.NullString
	updateRule           sql.NullString
	deleteRule           sql.NullString
}

func (m *DB) Analyze(options Options) {
	rawTables := m.getRawTables(options)
	description := m.descriptionFromRawTables(rawTables, options)
	m.model = db.NewModel(m.DbKey(), options.ForeignKeySuffix, !options.UseQualifiedNames, description)
}

func (m *DB) getRawTables(options Options) map[string]pgTable {
	var tableMap = make(map[string]pgTable)

	tables, schemas2 := m.getTables(options.Schemas)

	indexes, err := m.getIndexes(schemas2)
	if err != nil {
		return nil
	}

	foreignKeys, err := m.getForeignKeys(schemas2)
	if err != nil {
		return nil
	}

	for _, table := range tables {
		tableIndex := table.schema + "." + table.name

		// Do some processing on the foreign keys
		for _, fk := range foreignKeys[tableIndex] {
			if fk.referencedColumnName.Valid && fk.referencedTableName.Valid {
				if _, ok := table.fkMap[fk.columnName]; ok {
					log.Printf("Warning: Column %s:%s multi-table foreign keys are not supported.", table.name, fk.columnName)
					delete(table.fkMap, fk.columnName)
				} else {
					table.fkMap[fk.columnName] = fk
				}
			}
		}

		columns, err2 := m.getColumns(table.name, table.schema)
		if err2 != nil {
			return nil
		}

		table.indexes = indexes[tableIndex]
		table.columns = columns
		tableMap[tableIndex] = table
	}

	return tableMap

}

// Gets information for a table
func (m *DB) getTables(schemas []string) ([]pgTable, []string) {
	var tableName, tableSchema, tableComment string
	var tables []pgTable
	var schemaMap maps.Set[string]

	stmt := `
	SELECT
	t.table_name,
	t.table_schema,
	COALESCE(obj_description((table_schema||'.'||quote_ident(table_name))::regclass), '')
	FROM
	information_schema.tables t
	WHERE
	table_type <> 'VIEW'`

	if schemas != nil {
		stmt += fmt.Sprintf(` AND table_schema IN ('%s')`, strings.Join(schemas, `','`))
	} else {
		stmt += `AND table_schema NOT IN ('pg_catalog', 'information_schema')`
	}

	rows, err := m.SqlDb().Query(stmt)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&tableName, &tableSchema, &tableComment)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(tableSchema + "." + tableName)
		schemaMap.Add(tableSchema)
		table := pgTable{
			name:    tableName,
			schema:  tableSchema,
			comment: tableComment,
			columns: []pgColumn{},
			fkMap:   make(map[string]pgForeignKey),
			indexes: []pgIndex{},
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

	return tables, schemaMap.Values()
}

func (m *DB) getColumns(table string, schema string) (columns []pgColumn, err error) {

	s := fmt.Sprintf(`
	SELECT
	c.column_name,
	c.column_default,
	c.is_nullable,
	c.data_type,
	c.character_maximum_length,
	pgd.description
FROM
	information_schema.columns as c
JOIN 
	pg_catalog.pg_statio_all_tables as st
	on c.table_schema = st.schemaname
	and c.table_name = st.relname
LEFT JOIN 
	pg_catalog.pg_description pgd
	on pgd.objoid=st.relid
	and pgd.objsubid=c.ordinal_position
WHERE
	c.table_name = '%s' AND
	c.table_schema = '%s'
ORDER BY
	c.ordinal_position;
	`, table, schema)

	rows, err := m.SqlDb().Query(s)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var col pgColumn

	for rows.Next() {
		col = pgColumn{}
		var descr sql.NullString
		err = rows.Scan(&col.name, &col.defaultValue, &col.isNullable, &col.dataType, &col.characterMaxLen, &descr)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		if descr.Valid {
			if col.options, col.comment, err = sql2.ExtractOptions(descr.String); err != nil {
				log.Print("Error in table comment options for table " + table + ":" + col.name + " - " + err.Error())
			}
		}
		columns = append(columns, col)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return columns, err
}

func (m *DB) getIndexes(schemas []string) (indexes map[string][]pgIndex, err2 error) {

	indexes = make(map[string][]pgIndex)

	sql := fmt.Sprintf(`
	select idx.relname as index_name, 
       insp.nspname as index_schema,
       tbl.relname as table_name,
       tnsp.nspname as table_schema,
	   pgi.indisunique,
	   pgi.indisprimary,
	   a.attname as column_name
from pg_index pgi
  join pg_class idx on idx.oid = pgi.indexrelid
  join pg_namespace insp on insp.oid = idx.relnamespace
  join pg_class tbl on tbl.oid = pgi.indrelid
  join pg_namespace tnsp on tnsp.oid = tbl.relnamespace
  join pg_attribute a on a.attrelid = idx.oid
where
  tnsp.nspname IN ('%s')
	`, strings.Join(schemas, "','"))

	rows, err := m.SqlDb().Query(sql)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var index pgIndex

	for rows.Next() {
		index = pgIndex{}
		err = rows.Scan(&index.name, &index.schema, &index.tableName, &index.tableSchema, &index.unique, &index.primary, &index.columnName)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		indexKey := index.schema + "." + index.tableName
		tableIndexes := indexes[indexKey]
		tableIndexes = append(tableIndexes, index)
		indexes[indexKey] = tableIndexes
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
func (m *DB) getForeignKeys(schemas []string) (foreignKeys map[string][]pgForeignKey, err error) {
	fkMap := make(map[string]pgForeignKey)

	stmt := fmt.Sprintf(`
SELECT
    tc.constraint_name, 
    tc.table_name, 
    tc.table_schema, 
    kcu.column_name, 
    ccu.table_schema AS foreign_table_schema,
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name, 
    pgc.confdeltype,
    pgc.confupdtype
FROM 
    information_schema.table_constraints AS tc 
    JOIN information_schema.key_column_usage AS kcu
      ON tc.constraint_name = kcu.constraint_name
      AND tc.table_schema = kcu.table_schema
    JOIN information_schema.constraint_column_usage AS ccu
      ON ccu.constraint_name = tc.constraint_name
      AND ccu.table_schema = tc.table_schema
    JOIN pg_constraint as pgc
      ON tc.constraint_name = pgc.conname
WHERE tc.constraint_type = 'FOREIGN KEY' AND
      tc.table_schema IN ('%s')
	`, strings.Join(schemas, "','"))

	rows, err := m.SqlDb().Query(stmt)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		fk := pgForeignKey{}
		err = rows.Scan(&fk.constraintName,
			&fk.tableName,
			&fk.tableSchema,
			&fk.columnName,
			&fk.referencedSchema,
			&fk.referencedTableName,
			&fk.referencedColumnName,
			&fk.updateRule,
			&fk.deleteRule)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		if fk.referencedColumnName.Valid {
			fkMap[fk.constraintName] = fk
		}
	}

	rows.Close()

	foreignKeys = make(map[string][]pgForeignKey)
	stringmap.Range(fkMap, func(_ string, val interface{}) bool {
		fk := val.(pgForeignKey)
		i := fk.tableSchema + "." + fk.tableName
		tableKeys := foreignKeys[i]
		tableKeys = append(tableKeys, fk)
		foreignKeys[i] = tableKeys
		return true
	})
	return foreignKeys, err
}

// Convert the database native type to a more generic sql type, and a go table type.
func (m *DB) processTypeInfo(tableName string, column pgColumn, cd *db.ColumnDescription) {

	switch column.dataType {
	case "time without time zone":
		fallthrough
	case "time":
		cd.GoType = ColTypeTime.GoType()
		cd.SubType = "time"
	case "timestamp":
		fallthrough
	case "timestamp with time zone":
		cd.GoType = ColTypeTime.GoType()
		cd.SubType = "timestamp"
	case "datetime":
		fallthrough
	case "timestamp without time zone":
		cd.GoType = ColTypeTime.GoType()
	case "date":
		cd.GoType = ColTypeTime.GoType()
		cd.SubType = "date"

	case "boolean":
		cd.GoType = ColTypeBool.GoType()

	case "integer":
		fallthrough
	case "int":
		cd.GoType = ColTypeInteger.GoType()
		cd.MinValue = int64(-2147483648)
		cd.MaxValue = int64(2147483647)
		cd.MaxCharLength = 11

	case "smallint":
		cd.GoType = ColTypeInteger.GoType()
		cd.MinValue = int64(-32768)
		cd.MaxValue = int64(32767)
		cd.MaxCharLength = 6

	case "bigint": // We need to be explicit about this in go, since int will be whatever the OS native int size is, but go will support int64 always.
		// Also, since Json can only be decoded into float64s, we are limited in our ability to represent large min and max numbers in the json to about 2^53
		cd.GoType = ColTypeInteger64.GoType()
		cd.MinValue = int64(math.MinInt64)
		cd.MaxValue = int64(math.MaxInt64)
		cd.MaxCharLength = 20

	case "real":
		cd.GoType = ColTypeFloat32.GoType()
		cd.MinValue = -math.MaxFloat32 // float64 type
		cd.MaxValue = math.MaxFloat32

	case "double precision":
		cd.GoType = ColTypeFloat64.GoType()
		cd.MinValue = -math.MaxFloat64
		cd.MaxValue = math.MaxFloat64

	case "character varying":
		cd.GoType = ColTypeString.GoType()
		cd.MaxCharLength = uint64(column.characterMaxLen.Int64)

	case "char":
		cd.GoType = ColTypeString.GoType()
		cd.MaxCharLength = uint64(column.characterMaxLen.Int64)

	case "bytea":
		cd.GoType = ColTypeBytes.GoType()
		cd.MaxCharLength = 65535

	case "text":
		cd.GoType = ColTypeString.GoType()
		cd.MaxCharLength = 65535

	case "numeric":
		// No native equivalent in Go.
		// See the shopspring/decimal package for support.
		// You will need to shepherd numbers into and out of string format to move data to the database.
		cd.GoType = ColTypeString.GoType()
		cd.MaxCharLength = uint64(column.characterMaxLen.Int64) + 3

	case "year":
		cd.GoType = ColTypeInteger.GoType()

	default:
		cd.GoType = ColTypeString.GoType()
	}

	cd.NativeType = column.dataType
	cd.DefaultValue = getDefaultValue(column.defaultValue, ColTypeFromGoTypeString(cd.GoType))
}

func (m *DB) descriptionFromRawTables(rawTables map[string]pgTable, options Options) db.DatabaseDescription {

	dd := db.DatabaseDescription{}

	keys := stringmap.SortedKeys(rawTables)
	for _, tableName := range keys {
		table := rawTables[tableName]
		if table.options["skip"] != nil {
			continue
		}
		if strings.Contains(table.name, ".") {
			log.Print("Error: Table " + table.schema + "." + table.name + " cannot contain a period in its name. Skipping.")
			continue
		}
		if strings.Contains(table.schema, ".") {
			log.Print("Error: Schema " + table.schema + "." + table.name + " cannot contain a period in its schema name. Skipping.")
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

func (m *DB) getTableDescription(t pgTable) db.TableDescription {
	var columnDescriptions []db.ColumnDescription

	// Build the indexes
	pkColumns := make(map[string]bool)
	indexes := make(map[string]*db.IndexDescription)
	uniqueColumns := make(map[string]bool)

	// Fill pkColumns map with the column names of all the pk columns
	// Also file the indexes map with a list of columns for each index
	for _, idx := range t.indexes {
		if idx.primary {
			pkColumns[idx.columnName] = true
		} else if i, ok2 := indexes[idx.name]; ok2 {
			i.ColumnNames = append(i.ColumnNames, idx.columnName)
			sort.Strings(i.ColumnNames) // make sure this list stays in a predictable order each time
		} else {
			i = &db.IndexDescription{IsUnique: idx.unique, ColumnNames: []string{idx.columnName}}
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
		for k := range pkColumns {
			uniqueColumns[k] = true
		}
	}

	var pkCount int
	for _, col := range t.columns {
		if strings.Contains(col.name, ".") {
			log.Print(`Error: Column "` + col.name + `" cannot contain a period in its name. Skipping.`)
			continue
		}

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

	tableName := t.name
	if t.schema != "" {
		tableName = t.schema + "." + tableName
	}

	td := db.TableDescription{
		Name:                tableName,
		Columns:             columnDescriptions,
		SupportsForeignKeys: true, // Postgres supports foreign keys in all tables
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

func (m *DB) getTypeTableDescription(t pgTable) db.TableDescription {
	td := m.getTableDescription(t)

	var columnNames []string
	var quotedNames []string
	var columnTypes []GoColumnType

	for i, c := range td.Columns {
		columnNames = append(columnNames, c.Name)
		quotedNames = append(quotedNames, iq(c.Name))
		colType := ColTypeFromGoTypeString(c.GoType)
		if i == 0 {
			colType = ColTypeInteger // Force first value to be treated like an integer
		}
		columnTypes = append(columnTypes, colType)
	}

	stmt := fmt.Sprintf(`
SELECT
	%s
FROM
    %s
ORDER BY
    %s
`,
		strings.Join(quotedNames, `,`),
		iq(td.Name),
		quotedNames[0])

	result, err := m.SqlDb().Query(stmt)

	if err != nil {
		log.Fatal(err)
	}

	values := sql2.SqlReceiveRows(result, columnTypes, columnNames, nil)
	td.TypeData = values
	return td
}

func (m *DB) getColumnDescription(table pgTable, column pgColumn, isPk bool, isUnique bool) db.ColumnDescription {
	cd := db.ColumnDescription{
		Name: column.name,
	}
	m.processTypeInfo(table.name, column, &cd)

	// treat auto incrementing values as id values
	cd.IsId = column.defaultValue.Valid && strings.Contains(column.defaultValue.String, "nextval")
	cd.IsPk = isPk
	cd.IsNullable = column.isNullable == "YES"
	cd.IsUnique = isUnique

	cd.Comment = column.comment
	cd.Options = column.options

	if fk, ok2 := table.fkMap[cd.Name]; ok2 {
		tableName := fk.referencedTableName.String
		if fk.referencedSchema.Valid && fk.referencedSchema.String != "" {
			tableName = fk.referencedSchema.String + "." + tableName
		}

		cd.ForeignKey = &db.ForeignKeyDescription{
			ReferencedTable:  tableName,
			ReferencedColumn: fk.referencedColumnName.String,
			UpdateAction:     fkRuleToAction(fk.updateRule),
			DeleteAction:     fkRuleToAction(fk.deleteRule),
		}
	}

	return cd
}

func (m *DB) getManyManyDescription(t pgTable, typeTableSuffix string) (mm db.ManyManyDescription, ok bool) {
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

		if cd.ForeignKey.DeleteAction != db.FKActionCascade {
			log.Print("Warning: column " + td.Name + ":" + cd.Name + " has a DELETE action that is not CASCADE. You will need to manually delete the relationship before the associated object is deleted.")
		}

		if cd.IsNullable {
			log.Print("Error: table " + td.Name + ":" + cd.Name + " cannot be nullable.")
			return
		}

		if strings2.EndsWith(cd.ForeignKey.ReferencedTable, typeTableSuffix) {
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
	mm.SupportsForeignKeys = true
	ok = true
	return
}

func getDefaultValue(sqlVal sql.NullString, typ GoColumnType) interface{} {
	if !sqlVal.Valid {
		return nil
	}
	v := sqlVal.String

	if strings2.StartsWith(v, "NULL") {
		return nil
	}

	switch typ {
	case ColTypeBytes:
		return nil
	case ColTypeString:
		return v
	case ColTypeInteger:
		i, _ := strconv.Atoi(v)
		return i
	case ColTypeInteger64:
		i, _ := strconv.Atoi(v)
		return int64(i)
	case ColTypeTime:
		if v == "CURRENT_TIMESTAMP" {
			return "now"
		}
		return time2.FromSqlDateTime(v).UTC()
	case ColTypeFloat32:
		i, _ := strconv.ParseFloat(v, 32)
		return float32(i)
	case ColTypeFloat64:
		i, _ := strconv.ParseFloat(v, 64)
		return i
	case ColTypeBool:
		return v == "TRUE"
	default:
		return nil
	}
}

func fkRuleToAction(rule sql.NullString) db.FKAction {

	if !rule.Valid {
		return db.FKActionNone // This means we will emulate foreign key actions
	}
	switch rule.String {
	case "":
		fallthrough
	case "r":
		return db.FKActionRestrict
	case "c":
		return db.FKActionCascade
	case "d":
		return db.FKActionSetDefault
	case "n":
		return db.FKActionSetNull

	}
	return db.FKActionNone
}
