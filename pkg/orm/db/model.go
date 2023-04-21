package db

import (
	"context"
	"fmt"
	"github.com/gedex/inflector"
	. "github.com/goradd/goradd/pkg/orm/query"
	strings2 "github.com/goradd/goradd/pkg/strings"
	"github.com/kenshaw/snaker"
	"log"
	"regexp"
	"strings"
)

// These constants define the indexes used in the Options of Tables and Columns
const (
	LiteralNameOption   = "literalName"   // Used in tables only
	LiteralPluralOption = "literalPlural" // Used in tables only
	GoNameOption        = "goName"        // Used in tables and columns
	GoPluralOption      = "goPlural"      // Used in tables and columns
	MinOption           = "min"           // Used in numeric columns
	MaxOption           = "max"           // Used in number columns
)

// Model is the top level struct that contains a description of the database modeled as objects.
// It is used in code generation and query creation.
type Model struct {
	// The database key corresponding to its key in the global database cluster
	DbKey string
	// Tables are the tables in the database
	Tables []*Table
	// EnumTables contains a description of the enumerated types from the enum tables in the database
	EnumTables []*EnumTable

	// ForeignKeySuffix is the text to strip off the end of foreign key references when converting to names.
	// Defaults to "_id"
	ForeignKeySuffix string
	// EnumTableSuffix is the text to string off the end of an enum table when converting it to a type name.
	// Defaults to "_enum".
	EnumTableSuffix string

	// tableMap is used to get to tables by internal name
	tableMap map[string]*Table
	// enumTableMap gets to enum tables by internal name
	enumTableMap map[string]*EnumTable
	// ignoreSchemas indicates that the database uses table schemas, but we will ignore
	// them when generating object names. This means that the different tables in the schemas in the database
	// will not have overlapping names.
	ignoreSchemas bool
}

// NewModel creates a new Model object from the given DatabaseDescription object.
//
// dbKey is the unique key used throughout Goradd to refer to the database.
//
// foreignKeySuffix is the name ending that will be used to indicate a field is a foreign key pointer.
//
// ignoreSchemas indicates to ignore schema names when generating object names. If true and the
// database supports schemas, it will not use schema names to generate object names. If the database
// does not support schemas, this should be false
//
// desc is the description of the database.
func NewModel(dbKey string,
	foreignKeySuffix string,
	enumTableSuffix string,
	ignoreSchemas bool,
	desc DatabaseDescription) *Model {
	d := Model{
		DbKey:            dbKey,
		ForeignKeySuffix: foreignKeySuffix,
		EnumTableSuffix:  enumTableSuffix,
		ignoreSchemas:    ignoreSchemas,
	}
	d.importDescription(desc)
	return &d
}

// importDescription will convert a database description to a model which generally treats
// tables as objects and columns as member variables.
func (m *Model) importDescription(desc DatabaseDescription) {

	m.enumTableMap = make(map[string]*EnumTable)
	m.tableMap = make(map[string]*Table)

	// deal with enum tables first
	for _, table := range desc.Tables {
		if table.EnumData != nil {
			tt := m.importEnumTable(table)
			m.EnumTables = append(m.EnumTables, tt)
			m.enumTableMap[tt.DbName] = tt
		}
	}

	// get the regular tables
	for _, table := range desc.Tables {
		if table.EnumData == nil {
			t := m.importTable(table)
			if t != nil {
				m.Tables = append(m.Tables, t)
				m.tableMap[t.DbName] = t
			}
		}
	}

	// import foreign keys after the columns are in place
	for _, table := range desc.Tables {
		if table.EnumData == nil {
			m.importForeignKeys(table)
		}
	}

	// import reverse references after the foreign keys are in place
	for _, table := range m.Tables {
		m.importReverseReferences(table)
	}

	for _, assn := range desc.MM {
		m.importAssociation(assn)
	}
}

// importEnumTable will import the enum table provided by the database description
func (m *Model) importEnumTable(desc TableDescription) *EnumTable {
	typeName := desc.Name
	if m.ignoreSchemas {
		parts := strings.Split(typeName, ".")
		if len(parts) == 2 {
			typeName = parts[1]
		}
	}
	typeName = strings.TrimSuffix(typeName, m.EnumTableSuffix)

	t := &EnumTable{
		DbKey:         m.DbKey,
		DbName:        desc.Name,
		LiteralName:   m.dbNameToEnglishName(typeName),
		LiteralPlural: m.dbNameToEnglishPlural(typeName),
		GoName:        m.dbNameToGoName(typeName),
		GoPlural:      m.dbNameToGoPlural(typeName),
	}

	var ok bool
	if opt := desc.Options[LiteralNameOption]; opt != nil {
		if t.LiteralName, ok = opt.(string); !ok {
			log.Print("Error in option for table " + desc.Name + ": literalName is not a string")
		}
	}

	if opt := desc.Options[LiteralPluralOption]; opt != nil {
		if t.LiteralPlural, ok = opt.(string); !ok {
			log.Print("Error in option for table " + desc.Name + ": literalPlural is not a string")
		}
	}

	if opt := desc.Options[GoNameOption]; opt != nil {
		if t.GoName, ok = opt.(string); !ok {
			log.Print("Error in option for table " + desc.Name + ": goName is not a string")
		}
	}

	if opt := desc.Options[GoPluralOption]; opt != nil {
		if t.GoPlural, ok = opt.(string); !ok {
			log.Print("Error in option for table " + desc.Name + ": goPlural is not a string")
		}
	}

	t.LcGoName = strings.ToLower(t.GoName[:1]) + t.GoName[1:]

	t.Values = desc.EnumData
	t.FieldTypes = make(map[string]GoColumnType)

	for _, col := range desc.Columns {
		t.FieldNames = append(t.FieldNames, col.Name)
		t.FieldTypes[col.Name] = ColTypeFromGoTypeString(col.GoType)
	}

	t.Constants = make(map[int]string, len(t.FieldNames))
	names := t.FieldNames

	if len(t.Values) == 0 {
		log.Print("Warning: enum table " + t.DbName + " has no data entries. Specify constants by adding entries to this table.")
	}

	r := regexp.MustCompile("[^a-zA-Z0-9_]+")
	for _, val := range t.Values {
		key, ok := val[names[0]].(int)
		if !ok {
			panic("first column of enum table must be an integer")
		}
		value := val[names[1]].(string)
		var con string

		a := r.Split(value, -1)
		for _, word := range a {
			con += strings2.Title(strings.ToLower(word))
		}
		t.Constants[key] = con
	}

	t.PkField = t.FieldGoName(0)
	return t
}

// importTable will import the table provided by the description
func (m *Model) importTable(desc TableDescription) *Table {
	tableName := desc.Name
	if m.ignoreSchemas {
		parts := strings.Split(tableName, ".")
		if len(parts) == 2 {
			tableName = parts[1]
		}
	}

	t := &Table{
		DbKey:         m.DbKey,
		DbName:        desc.Name,
		LiteralName:   m.dbNameToEnglishName(tableName),
		LiteralPlural: m.dbNameToEnglishPlural(tableName),
		GoName:        m.dbNameToGoName(tableName),
		GoPlural:      m.dbNameToGoPlural(tableName),
		Comment:       desc.Comment,
		Options:       desc.Options,
		columnMap:     make(map[string]*Column),
	}

	var ok bool
	if opt := desc.Options[LiteralNameOption]; opt != nil {
		if t.LiteralName, ok = opt.(string); !ok {
			log.Print("Error in option for table " + desc.Name + ": literalName is not a string")
		}
	}

	if opt := desc.Options[LiteralPluralOption]; opt != nil {
		if t.LiteralPlural, ok = opt.(string); !ok {
			log.Print("Error in option for table " + desc.Name + ": literalPlural is not a string")
		}
	}

	if opt := desc.Options[GoNameOption]; opt != nil {
		if t.GoName, ok = opt.(string); !ok {
			log.Print("Error in option for table " + desc.Name + ": goName is not a string")
		}
	}

	if opt := desc.Options[GoPluralOption]; opt != nil {
		if t.GoPlural, ok = opt.(string); !ok {
			log.Print("Error in option for table " + desc.Name + ": goPlural is not a string")
		}
	}

	t.LcGoName = strings.ToLower(t.GoName[:1]) + t.GoName[1:]

	if t.GoName == t.GoPlural {
		log.Print("Error: table " + t.DbName + " is using a plural name. Change it to a singular name or assign a singular and plural go name in the comments.")
		return nil
	}

	var pkCount int
	for _, col := range desc.Columns {
		newCol := m.importColumn(col)
		if newCol != nil {
			t.Columns = append(t.Columns, newCol)
			t.columnMap[newCol.DbName] = newCol
			if newCol.IsPk {
				pkCount++
				if pkCount > 1 {
					log.Print("Error: table " + t.DbName + " has multiple primary keys.")
					return nil
				}
			}
		}
	}

	for _, idx := range desc.Indexes {
		var columns []*Column
		for _, name := range idx.ColumnNames {
			col := t.GetColumn(name)
			if col == nil {
				panic("Cannot find column " + name + " of table " + t.DbName)
			}
			columns = append(columns, col)
		}
		t.Indexes = append(t.Indexes, Index{IsUnique: idx.IsUnique, Columns: columns})
	}
	return t
}

func (m *Model) importReverseReferences(td *Table) {
	var td2 *Table

	var col *Column
	for _, col = range td.Columns {
		if col.ForeignKey != nil {
			td2 = m.Table(col.ForeignKey.ReferencedTable)
			if td2 == nil {
				continue // pointing to a enum table
			}
			// Determine the go name, which is the name used to refer to the reverse reference.
			// This is somewhat tricky, because there is no easy way to extract an expression for this.
			// For example, a team-member  relationship would simply be called team_id from the person side, so
			// where would you get the word "member". Our strategy will be to first look for something explicit in the
			// options, and if not found, try to create a name from the table and table names.

			objName := col.DbName
			objName = strings.TrimSuffix(objName, m.ForeignKeySuffix)
			objName = UpperCaseIdentifier(objName)
			objName = strings.Replace(objName, td2.GoName, "", 1)

			goName, _ := col.Options[GoNameOption].(string)
			if goName == "" {
				goName = td.GoName
				if objName != "" {
					goName = goName + "As" + objName
				}
			}
			goPlural, _ := col.Options[GoPluralOption].(string)
			if goPlural == "" {
				goPlural = td.GoPlural
				if objName != "" {
					goPlural = goPlural + "As" + objName
				}
			}
			goType := td.GoName
			goTypePlural := td.GoPlural

			// Check for name conflicts
			for _, col2 := range td2.Columns {
				if goName == col2.GoName {
					log.Printf("Error: table %s has a field name %s that is the same as the %s table that is referring to it. Either change these names, or provide an alternate goName in the options.", td2.GoName, goName, td.GoName)
				}
			}

			var colName string
			if col.IsUnique {
				colName = snaker.CamelToSnake(goName)
			} else {
				colName = snaker.CamelToSnake(goPlural)
			}

			ref := ReverseReference{
				DbColumn:         colName,
				AssociatedTable:  td,
				AssociatedColumn: col,
				GoName:           goName,
				GoPlural:         goPlural,
				GoType:           goType,
				GoTypePlural:     goTypePlural,
			}

			td2.ReverseReferences = append(td2.ReverseReferences, &ref)
			col.ForeignKey.RR = &ref
		}
	}
}

// Analyzes an association table and creates special virtual columns in the corresponding tables it points to.
// Association tables are used by SQL databases to create many-many relationships. NoSQL databases can define their
// association columns directly and store an array of records numbers on either end of the association.
func (m *Model) importAssociation(mm ManyManyDescription) {
	if m.EnumTable(mm.Table2) == nil {
		ref1 := m.makeManyManyRef(mm.Table1, mm.Column1, mm.Table2, mm.Column2, mm.GoName2, mm.GoPlural2, mm.AssnTableName, false, mm.SupportsForeignKeys)
		ref2 := m.makeManyManyRef(mm.Table2, mm.Column2, mm.Table1, mm.Column1, mm.GoName1, mm.GoPlural1, mm.AssnTableName, false, mm.SupportsForeignKeys)
		ref1.MM = ref2
		ref2.MM = ref1
	} else {
		// enum table
		m.makeManyManyRef(mm.Table1, mm.Column1, mm.Table2, mm.Column2, mm.GoName2, mm.GoPlural2, mm.AssnTableName, true, mm.SupportsForeignKeys)
	}
}

func (m *Model) makeManyManyRef(
	t1, c1, t2, c2, g2, g2p, t string,
	isEnum, supportsForeignKeys bool,
) *ManyManyReference {
	sourceTableName := t1
	destTableName := t2
	sourceObjName := strings.TrimSuffix(c1, m.ForeignKeySuffix)
	destObjName := strings.TrimSuffix(c2, m.ForeignKeySuffix)
	sourceTable := m.Table(sourceTableName)

	if m.ignoreSchemas {
		parts := strings.Split(sourceTableName, ".")
		if len(parts) == 2 {
			sourceTableName = parts[1]
		}
	}
	var objName string
	var objPlural string
	var pkType string

	if isEnum {
		destTable := m.EnumTable(destTableName)
		objName = destTable.GoName
		objPlural = destTable.GoPlural
		pkType = "int"
	} else {
		destTable := m.Table(destTableName)
		objName = destTable.GoName
		objPlural = destTable.GoPlural
		pkType = destTable.PrimaryKeyGoType()
	}

	goName := g2
	goPlural := g2p

	if goName != "" {
		// use it
	} else if sourceObjName != sourceTableName {
		goName = UpperCaseIdentifier(destObjName) + "As" + UpperCaseIdentifier(sourceObjName)
	} else {
		goName = UpperCaseIdentifier(destObjName)
	}

	if goPlural != "" {
		// use it
	} else if sourceObjName != sourceTableName {
		goPlural = inflector.Pluralize(UpperCaseIdentifier(destObjName)) + "As" + UpperCaseIdentifier(sourceObjName)
	} else {
		goPlural = inflector.Pluralize(goName)
	}

	ref := ManyManyReference{
		AssnTableName:         t,
		AssnColumnName:        c1,
		AssociatedTableName:   t2,
		AssociatedTablePkType: pkType,
		AssociatedColumnName:  c2,
		AssociatedObjectType:  objName,
		AssociatedObjectTypes: objPlural,
		GoName:                goName,
		GoPlural:              goPlural,
		IsEnumAssociation:     isEnum,
		SupportsForeignKeys:   supportsForeignKeys,
	}
	sourceTable.ManyManyReferences = append(sourceTable.ManyManyReferences, &ref)
	return &ref
}

func (m *Model) importColumn(desc ColumnDescription) *Column {
	c := &Column{
		DbName:        desc.Name,
		GoName:        m.dbNameToGoName(desc.Name),
		NativeType:    desc.NativeType,
		ColumnType:    ColTypeFromGoTypeString(desc.GoType),
		MaxCharLength: desc.MaxCharLength,
		DefaultValue:  desc.DefaultValue,
		MaxValue:      desc.MaxValue,
		MinValue:      desc.MinValue,
		IsId:          desc.IsId,
		IsPk:          desc.IsPk,
		IsNullable:    desc.IsNullable,
		IsUnique:      desc.IsUnique,
		IsTimestamp:   desc.SubType == "timestamp",
		IsDateOnly:    desc.SubType == "date",
		IsTimeOnly:    desc.SubType == "time",
		Comment:       desc.Comment,
		Options:       desc.Options,
	}

	if c.IsId {
		c.ColumnType = ColTypeString // We treat auto-generated ids as strings for cross database compatibility.
	}

	var ok bool
	if opt := desc.Options[GoNameOption]; opt != nil {
		if c.GoName, ok = opt.(string); !ok {
			log.Printf("Error in option for column " + desc.Name + ": goName is not a string")
		}
	}

	c.modelName = LowerCaseIdentifier(c.DbName)

	var err error
	if opt := desc.Options[MinOption]; opt != nil {
		if c.MinValue, err = getMinOption(c.MinValue, opt); err != nil {
			log.Printf("Error in 'min' option for column %s: %s", desc.Name, err.Error())
		}
	}

	if opt := desc.Options[MaxOption]; opt != nil {
		if c.MaxValue, err = getMaxOption(c.MaxValue, opt); err != nil {
			log.Printf("Error in 'max' option for column %s: %s", desc.Name, err.Error())
		}
	}

	return c
}

func (m *Model) importForeignKeys(desc TableDescription) {
	t := m.Table(desc.Name)
	if t != nil {
		if t.PrimaryKeyColumn() == nil {
			log.Printf("*** Error: table %s must have a primary key if it has a foreign key.", desc.Name)
		} else {
			for _, col := range desc.Columns {
				m.importForeignKey(t, col)
			}
		}
	}
}

func (m *Model) importForeignKey(t *Table, cd ColumnDescription) {
	c := t.columnMap[cd.Name]
	if cd.ForeignKey != nil {
		f := &ForeignKeyInfo{
			ReferencedTable:  cd.ForeignKey.ReferencedTable,
			ReferencedColumn: cd.ForeignKey.ReferencedColumn,
			UpdateAction:     cd.ForeignKey.UpdateAction,
			DeleteAction:     cd.ForeignKey.DeleteAction,
		}

		if (f.UpdateAction == FKActionSetNull || f.DeleteAction == FKActionSetNull) &&
			!cd.IsNullable {
			panic(fmt.Sprintf("a foreign key cannot have an action of Null if the column is not nullable. Table: %s, Col: %s", t.DbName, cd.Name))
		}

		goName := c.GoName
		suffix := UpperCaseIdentifier(m.ForeignKeySuffix)
		goName = strings.TrimSuffix(goName, suffix)
		if goName == "" || c.IsPk {
			// Either:
			// - the primary key is also a foreign key, which would be a 1-1 relationship, or
			// - the name of column is just a foreign key suffix, as in "id"
			// So we use the name of the foreign key table as the object name
			goName = f.ReferencedTable
			goName = UpperCaseIdentifier(goName)
		}
		f.GoName = goName
		f.IsEnum = m.IsEnumTable(cd.ForeignKey.ReferencedTable)

		if f.IsEnum {
			tt := m.EnumTable(cd.ForeignKey.ReferencedTable)
			f.GoType = tt.GoName
			f.GoTypePlural = tt.GoPlural
			suf := UpperCaseIdentifier(m.ForeignKeySuffix)
			c.referenceFunction = strings.TrimSuffix(f.GoName, suf)
		} else {
			r := m.Table(cd.ForeignKey.ReferencedTable)
			f.GoType = r.GoName
			f.GoTypePlural = r.GoPlural
			fkc := r.GetColumn(cd.ForeignKey.ReferencedColumn)
			if fkc.IsId {
				c.ColumnType = ColTypeString // Always use strings to refer to auto-generated ids for cross database compatibility
			}
			c.referenceFunction = f.GoName
		}
		c.ForeignKey = f
	}
}

func (m *Model) dbNameToEnglishName(name string) string {
	return strings2.Title(name)
}

func (m *Model) dbNameToEnglishPlural(name string) string {
	return inflector.Pluralize(m.dbNameToEnglishName(name))
}

func (m *Model) dbNameToGoName(name string) string {
	return UpperCaseIdentifier(name)
}

func (m *Model) dbNameToGoPlural(name string) string {
	return inflector.Pluralize(m.dbNameToGoName(name))
}

// Table returns a Table from the database given the table name.
func (m *Model) Table(name string) *Table {
	if v, ok := m.tableMap[name]; ok {
		return v
	} else {
		return nil
	}
}

// EnumTable returns a EnumTable from the database given the table name.
func (m *Model) EnumTable(name string) *EnumTable {
	return m.enumTableMap[name]
}

// IsEnumTable returns true if the given name is the name of a enum table in the database
func (m *Model) IsEnumTable(name string) bool {
	_, ok := m.enumTableMap[name]
	return ok
}

func isReservedIdentifier(s string) bool {
	switch s {
	case "break":
		return true
	case "case":
		return true
	case "chan":
		return true
	case "const":
		return true
	case "continue":
		return true
	case "default":
		return true
	case "defer":
		return true
	case "else":
		return true
	case "fallthrough":
		return true
	case "for":
		return true
	case "func":
		return true
	case "go":
		return true
	case "goto":
		return true
	case "if":
		return true
	case "import":
		return true
	case "interface":
		return true
	case "map":
		return true
	case "package":
		return true
	case "range":
		return true
	case "return":
		return true
	case "select":
		return true
	case "struct":
		return true
	case "switch":
		return true
	case "type":
		return true
	case "var":
		return true
	}
	return false
}

func LowerCaseIdentifier(s string) (i string) {
	if strings.Contains(s, "_") {
		i = snaker.ForceLowerCamelIdentifier(snaker.SnakeToCamelIdentifier(s))
	} else {
		// Not a snake string, but still might need some fixing up
		i = snaker.ForceLowerCamelIdentifier(s)
	}
	i = strings.TrimSpace(i)
	if isReservedIdentifier(i) {
		panic("Cannot use '" + i + "' as an identifier.")
	}
	if i == "" {
		panic("Cannot use blank as an identifier.")
	}
	return
}

func UpperCaseIdentifier(s string) (i string) {
	if strings.Contains(s, "_") {
		i = snaker.ForceCamelIdentifier(snaker.SnakeToCamelIdentifier(s))
	} else {
		// Not a snake string, but still might need some fixing up
		i = snaker.ForceCamelIdentifier(s)
	}
	i = strings.TrimSpace(i)
	if i == "" {
		panic("Cannot use blank as an identifier.")
	}
	return
}

// ExecuteTransaction wraps the function in a database transaction
func ExecuteTransaction(ctx context.Context, d DatabaseI, f func()) {
	txid := d.Begin(ctx)
	defer d.Rollback(ctx, txid)
	f()
	d.Commit(ctx, txid)
}
