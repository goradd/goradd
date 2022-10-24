package db

import (
	"context"
	"fmt"
	"github.com/gedex/inflector"
	. "github.com/goradd/goradd/pkg/orm/query"
	"github.com/kenshaw/snaker"
	"log"
	"regexp"
	"strings"
)

// These constants define the indexes used in the Options of Tables and Columns
const (
	LiteralNameOption    = "literalName"   // Used in tables only
	LiteralPluralOption  = "literalPlural" // Used in tables only
	GoNameOption         = "goName"        // Used in tables and columns
	GoPluralOption       = "goPlural"      // Used in tables and columns
	MinOption            = "min"           // Used in numeric columns
	MaxOption            = "max"           // Used in number columns
	StringerColumnOption = "stringer"      // Used in tables only to specify what column to output in the String function. Is the db name of the column.
)

// The Database is the top level struct that contains a complete description of a database for purposes of
// creating queries and doing code generation
type Database struct {
	// The database key corresponding to its key in the global database cluster
	DbKey string
	// Tables are the tables in the database
	Tables []*Table
	// TypeTables contains a description of the enumerated types from the type tables in the database
	TypeTables []*TypeTable
	// AssociatedObjectPrefix is a prefix placed in front of generated object names. Defaults to "o".
	AssociatedObjectPrefix string

	// Text to strip off the end of foreign key references when converting to names. Defaults to "_id"
	ForeignKeySuffix string

	// These items are filled in by analysis

	// tableMap is used to get to tables by internal name
	tableMap map[string]*Table
	// typeTableMap gets to type tables by internal name
	typeTableMap map[string]*TypeTable
}

// NewDatabase creates a new Database object from the given DatabaseDescription object.
func NewDatabase(dbKey string, foreignKeySuffix string, desc DatabaseDescription) *Database {
	d := Database{
		DbKey:                  dbKey,
		AssociatedObjectPrefix: desc.AssociatedObjectPrefix,
		ForeignKeySuffix:       foreignKeySuffix,
	}
	d.analyze(desc)
	return &d
}

// Given a database description, analyze will perform an analysis of the database, and modify some of the fields to prepare
// the description for use in codegen and the orm
func (d *Database) analyze(desc DatabaseDescription) {

	d.typeTableMap = make(map[string]*TypeTable)
	d.tableMap = make(map[string]*Table)
	if d.AssociatedObjectPrefix == "" {
		d.AssociatedObjectPrefix = "o"
	}

	// deal with type tables first
	for _, table := range desc.Tables {
		if table.TypeData != nil {
			tt := d.analyzeTypeTable(table)
			d.TypeTables = append(d.TypeTables, tt)
			d.typeTableMap[tt.DbName] = tt
		}
	}

	// get the regular tables
	for _, table := range desc.Tables {
		if table.TypeData == nil {
			t := d.analyzeTable(table)
			if t != nil {
				d.Tables = append(d.Tables, t)
				d.tableMap[t.DbName] = t
			}
		}
	}

	// analyze foreign keys after the columns are in place
	for _, table := range desc.Tables {
		if table.TypeData == nil {
			d.analyzeForeignKeys(table)
		}
	}

	// analyze reverse references after the foreign keys are in place
	for _, table := range d.Tables {
		d.analyzeReverseReferences(table)
	}

	for _, assn := range desc.MM {
		d.analyzeAssociation(assn)
	}
}

// analyzeTypeTables will analyze the type tables provided by the database description
func (d *Database) analyzeTypeTable(desc TableDescription) *TypeTable {
	t := &TypeTable{
		DbKey:         d.DbKey,
		DbName:        desc.Name,
		LiteralName:   d.dbNameToEnglishName(desc.Name),
		LiteralPlural: d.dbNameToEnglishPlural(desc.Name),
		GoName:        d.dbNameToGoName(desc.Name),
		GoPlural:      d.dbNameToGoPlural(desc.Name),
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

	t.Values = desc.TypeData
	t.FieldTypes = make(map[string]GoColumnType)

	for _, col := range desc.Columns {
		t.FieldNames = append(t.FieldNames, col.Name)
		t.FieldTypes[col.Name] = ColTypeFromGoTypeString(col.GoType)
	}

	t.Constants = make(map[int]string, len(t.FieldNames))
	names := t.FieldNames
	var key int
	var value string

	if len(t.Values) == 0 {
		log.Print("Warning: type table " + t.DbName + " has no data entries. Specify constants by adding entries to this table.")
	}

	r := regexp.MustCompile("[^a-zA-Z0-9_]+")
	for _, m := range t.Values {
		key, ok = m[names[0]].(int)
		if !ok {
			key = int(m[names[0]].(uint))
		}
		value = m[names[1]].(string)
		var con string

		a := r.Split(value, -1)
		for _, word := range a {
			con += strings.Title(strings.ToLower(word))
		}
		t.Constants[key] = con
	}

	t.PkField = t.FieldGoName(0)
	return t
}

// analyzeTable will analyze the table provided by the description
func (d *Database) analyzeTable(desc TableDescription) *Table {
	t := &Table{
		DbKey:         d.DbKey,
		DbName:        desc.Name,
		LiteralName:   d.dbNameToEnglishName(desc.Name),
		LiteralPlural: d.dbNameToEnglishPlural(desc.Name),
		GoName:        d.dbNameToGoName(desc.Name),
		GoPlural:      d.dbNameToGoPlural(desc.Name),
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
		newCol := d.analyzeColumn(col)
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

func (d *Database) analyzeReverseReferences(td *Table) {
	var td2 *Table

	var col *Column
	for _, col = range td.Columns {
		if col.ForeignKey != nil {
			td2 = d.Table(col.ForeignKey.ReferencedTable)
			if td2 == nil {
				continue // pointing to a type table
			}
			// Determine the go name, which is the name used to refer to the reverse reference.
			// This is somewhat tricky, because there is no easy way to extract an expression for this.
			// For example, a team-member  relationship would simply be called team_id from the person side, so
			// where would you get the word "member". Our strategy will be to first look for something explicit in the
			// options, and if not found, try to create a name from the table and table names.

			objName := col.DbName
			objName = strings.TrimSuffix(objName, d.ForeignKeySuffix)
			objName = strings.Replace(objName, td2.DbName, "", 1)
			if objName != "" {
				objName = UpperCaseIdentifier(objName)
			}
			goName, _ := col.Options[GoNameOption].(string)
			if goName == "" {
				goName = UpperCaseIdentifier(td.DbName)
				if objName != "" {
					goName = goName + "As" + objName
				}
			}
			goPlural, _ := col.Options[GoPluralOption].(string)
			if goPlural == "" {
				goPlural = inflector.Pluralize(td.DbName)
				goPlural = UpperCaseIdentifier(goPlural)
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
				Values:           make(map[string]string),
			}

			td2.ReverseReferences = append(td2.ReverseReferences, &ref)
			col.ForeignKey.RR = &ref
		}
	}
}

// Analyzes an association table and creates special virtual columns in the corresponding tables it points to.
// Association tables are used by SQL databases to create many-many relationships. NoSQL databases can define their
// association columns directly and store an array of records numbers on either end of the association.
func (d *Database) analyzeAssociation(mm ManyManyDescription) {
	if d.TypeTable(mm.Table2) == nil {
		ref1 := d.makeManyManyRef(mm.Table1, mm.Column1, mm.Table2, mm.Column2, mm.GoName2, mm.GoPlural2, mm.AssnTableName, false)
		ref2 := d.makeManyManyRef(mm.Table2, mm.Column2, mm.Table1, mm.Column1, mm.GoName1, mm.GoPlural1, mm.AssnTableName, false)
		ref1.MM = ref2
		ref2.MM = ref1

	} else {
		// type table
		d.makeManyManyRef(mm.Table1, mm.Column1, mm.Table2, mm.Column2, mm.GoName2, mm.GoPlural2, mm.AssnTableName, true)
	}
}

func (d *Database) makeManyManyRef(t1, c1, t2, c2, g2, g2p, t string, isType bool) *ManyManyReference {
	sourceTableName := t1
	destTableName := t2
	sourceObjName := strings.TrimSuffix(c1, d.ForeignKeySuffix)
	destObjName := strings.TrimSuffix(c2, d.ForeignKeySuffix)
	sourceTable := d.Table(sourceTableName)

	var objName string
	if isType {
		destTable := d.TypeTable(destTableName)
		objName = destTable.GoName
	} else {
		destTable := d.Table(destTableName)
		objName = destTable.GoName
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
		AssnTableName:        t,
		AssnColumnName:       c1,
		AssociatedTableName:  t2,
		AssociatedColumnName: c2,
		AssociatedObjectName: objName,
		GoName:               goName,
		GoPlural:             goPlural,
		IsTypeAssociation:    isType,
	}
	sourceTable.ManyManyReferences = append(sourceTable.ManyManyReferences, &ref)
	return &ref
}

func (d *Database) analyzeColumn(desc ColumnDescription) *Column {
	c := &Column{
		DbName:                desc.Name,
		GoName:                d.dbNameToGoName(desc.Name),
		NativeType:            desc.NativeType,
		ColumnType:            ColTypeFromGoTypeString(desc.GoType),
		MaxCharLength:         desc.MaxCharLength,
		DefaultValue:          desc.DefaultValue,
		MaxValue:              desc.MaxValue,
		MinValue:              desc.MinValue,
		IsId:                  desc.IsId,
		IsPk:                  desc.IsPk,
		IsNullable:            desc.IsNullable,
		IsUnique:              desc.IsUnique,
		IsTimestamp:           desc.SubType == "timestamp" || desc.SubType == "auto timestamp",
		IsAutoUpdateTimestamp: desc.SubType == "auto timestamp",
		IsDateOnly:            desc.SubType == "date",
		IsTimeOnly:            desc.SubType == "time",
		Comment:               desc.Comment,
		Options:               desc.Options,
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

func (d *Database) analyzeForeignKeys(desc TableDescription) {
	t := d.Table(desc.Name)
	if t != nil {
		for _, col := range desc.Columns {
			d.analyzeForeignKey(t, col)
		}
	}
}

func (d *Database) analyzeForeignKey(t *Table, cd ColumnDescription) {
	c := t.columnMap[cd.Name]
	if cd.ForeignKey != nil {
		f := &ForeignKeyInfo{
			ReferencedTable:  cd.ForeignKey.ReferencedTable,
			ReferencedColumn: cd.ForeignKey.ReferencedColumn,
			UpdateAction:     FKActionFromString(cd.ForeignKey.UpdateAction),
			DeleteAction:     FKActionFromString(cd.ForeignKey.DeleteAction),
		}

		if (f.UpdateAction == FKActionSetNull || f.DeleteAction == FKActionSetNull) &&
			!cd.IsNullable {
			panic(fmt.Sprintf("a foreign key cannot have an action of Null if the column is not nullable. Table: %s, Col: %s", t.DbName, cd.Name))
		}

		goName := c.GoName
		suffix := UpperCaseIdentifier(d.ForeignKeySuffix)
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
		f.IsType = d.IsTypeTable(cd.ForeignKey.ReferencedTable)

		if f.IsType {
			tt := d.TypeTable(cd.ForeignKey.ReferencedTable)
			f.GoType = tt.GoName
			f.GoTypePlural = tt.GoPlural
			suf := UpperCaseIdentifier(d.ForeignKeySuffix)
			c.referenceFunction = strings.TrimSuffix(f.GoName, suf)
		} else {
			r := d.Table(cd.ForeignKey.ReferencedTable)
			f.GoType = r.GoName
			f.GoTypePlural = r.GoPlural
			fkc := r.GetColumn(cd.ForeignKey.ReferencedColumn)
			if fkc.IsId {
				c.ColumnType = ColTypeString // Always use strings to refer to auto-generated ids for cross database compatibility
			}
			c.referenceName = d.AssociatedObjectPrefix + f.GoName
			c.referenceFunction = f.GoName
		}
		c.ForeignKey = f
	}
}

func (d *Database) dbNameToEnglishName(name string) string {
	return strings.Title(strings.Replace(name, "_", " ", -1))
}

func (d *Database) dbNameToEnglishPlural(name string) string {
	return inflector.Pluralize(d.dbNameToEnglishName(name))
}

func (d *Database) dbNameToGoName(name string) string {
	return UpperCaseIdentifier(name)
}

func (d *Database) dbNameToGoPlural(name string) string {
	return inflector.Pluralize(d.dbNameToGoName(name))
}

// Table returns a Table from the database given the table name.
func (d *Database) Table(name string) *Table {
	if v, ok := d.tableMap[name]; ok {
		return v
	} else {
		return nil
	}
}

// TypeTable returns a TypeTable from the database given the table name.
func (d *Database) TypeTable(name string) *TypeTable {
	return d.typeTableMap[name]
}

// IsTypeTable returns true if the given name is the name of a type table in the database
func (d *Database) IsTypeTable(name string) bool {
	_, ok := d.typeTableMap[name]
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
		panic("Cannot use " + i + " as an identifier.")
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
