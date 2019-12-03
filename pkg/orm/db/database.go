package db

import (
	"github.com/gedex/inflector"
	. "github.com/goradd/goradd/pkg/orm/query"
	"github.com/knq/snaker"
	"log"
	"regexp"
	"strings"
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


/*
type Option struct {
	key   string
	value string
}
*/


// NewDatabase creates a new Database object from the given DatabaseDescription object.
func NewDatabase(dbKey string, foreignKeySuffix string, desc DatabaseDescription) *Database {
	d := Database {
		DbKey:                  dbKey,
		AssociatedObjectPrefix: desc.AssociatedObjectPrefix,
		ForeignKeySuffix: foreignKeySuffix,
	}
	d.analyze(desc)
	return &d
}


// Given a database description, analyze will perform an analysis of the database, and modify some of the fields to prepare
// the description for use in codegen and the orm
func (d *Database) analyze(desc DatabaseDescription) {

	d.typeTableMap = make(map[string]*TypeTable)
	d.tableMap = make(map[string]*Table)

	// deal with type tables first
	for _,table := range desc.Tables {
		if table.TypeData != nil {
			tt := d.analyzeTypeTable(table)
			d.TypeTables = append(d.TypeTables, tt)
			d.typeTableMap[tt.DbName] = tt
		}
	}

	// get the regular tables
	for _,table := range desc.Tables {
		if table.TypeData == nil {
			t := d.analyzeTable(table)
			if t != nil {
				d.Tables = append(d.Tables, t)
				d.tableMap[t.DbName] = t
			}
		}
	}

	// analyze foreign keys after the columns are in place
	for _,table := range desc.Tables {
		if table.TypeData == nil {
			d.analyzeForeignKeys(table)
		}
	}

	// analyze reverse references after the foreign keys are in place
	for _,table := range d.Tables {
		d.analyzeReverseReferences(table)
	}


	for _,assn := range desc.MM {
		d.analyzeAssociation(assn)
	}

	if d.AssociatedObjectPrefix == "" {
		d.AssociatedObjectPrefix = "o"
	}

}


// analyzeTypeTables will analyze the type tables provided by the database description
func (d *Database) analyzeTypeTable(desc TableDescription) *TypeTable {
	tt := &TypeTable {
		DbKey:         d.DbKey,
		DbName:        desc.Name,
		LiteralName:   desc.LiteralName,
		LiteralPlural: desc.LiteralPlural,
		GoName:        desc.GoName,
		GoPlural:      desc.GoPlural,
	}
	if tt.LiteralName == "" {
		tt.LiteralName = d.dbNameToEnglishName(tt.DbName)
	}
	if tt.GoName == "" {
		tt.GoName = d.dbNameToGoName(tt.DbName)
	}
	if tt.LiteralPlural == "" {
		tt.LiteralPlural = d.dbNameToEnglishPlural(tt.DbName)
	}
	if tt.GoPlural == "" {
		tt.GoPlural = d.dbNameToGoPlural(tt.DbName)
	}

	tt.LcGoName = strings.ToLower(tt.GoName[:1]) + tt.GoName[1:]

	tt.Values = desc.TypeData
	tt.FieldTypes = make(map[string]GoColumnType)

	for _,col := range desc.Columns {
		tt.FieldNames = append(tt.FieldNames, col.Name)
		tt.FieldTypes[col.Name] = ColTypeFromGoTypeString(col.GoType)
	}

	tt.Constants = make(map[uint]string, len(tt.FieldNames))
	names := tt.FieldNames
	var key uint
	var value string
	var ok bool

	if len(tt.Values) == 0 {
		log.Print("Warning: type table " + tt.DbName + " has no data entries. Specify constants by adding entries to this table.")
	}

	r := regexp.MustCompile("[^a-zA-Z0-9_]+")
	for _, m := range tt.Values {
		key, ok = m[names[0]].(uint)
		if !ok {
			key = uint(m[names[0]].(int))
		}
		value = m[names[1]].(string)
		var con string

		a := r.Split(value, -1)
		for _, word := range a {
			con += strings.Title(strings.ToLower(word))
		}
		tt.Constants[key] = con
	}

	tt.PkField = tt.FieldGoName(0)
	return tt
}

// analyzeTable will analyze the table provided by the description
func (d *Database) analyzeTable(desc TableDescription) *Table {
	t := &Table{
		DbKey:         d.DbKey,
		DbName:        desc.Name,
		LiteralName:   desc.LiteralName,
		LiteralPlural: desc.LiteralPlural,
		GoName:        desc.GoName,
		GoPlural:      desc.GoPlural,
		Comment:	   desc.Comment,
		Options:	   desc.Options,
		columnMap: make(map[string]*Column),
	}

	if t.LiteralName == "" {
		t.LiteralName = d.dbNameToEnglishName(t.DbName)
	}
	if t.GoName == "" {
		t.GoName = d.dbNameToGoName(t.DbName)
	}
	if t.LiteralPlural == "" {
		t.LiteralPlural = d.dbNameToEnglishPlural(t.DbName)
	}
	if t.GoPlural == "" {
		t.GoPlural = d.dbNameToGoPlural(t.DbName)
	}
	t.LcGoName = strings.ToLower(t.GoName[:1]) + t.GoName[1:]

	if t.GoName == t.GoPlural {
		log.Print("Error: table " + t.DbName + " is using a plural name. Change it to a singular name or assign a singular and plural go name in the comments.")
		return nil
	}

	var pkCount int
	for _,col := range desc.Columns {
		newCol := d.analyzeColumn(col)
		if newCol != nil {
			t.Columns = append(t.Columns, newCol)
			t.columnMap[newCol.DbName] = newCol
			if newCol.ColumnType == ColTypeDateTime {
				t.HasDateTime = true
			}
			if newCol.IsPk {
				pkCount++
				if pkCount > 1 {
					log.Print("Error: table " + t.DbName + " has multiple primary keys.")
					return nil
				}
			}
		}
	}

	for _,idx := range desc.Indexes {
		t.Indexes = append(t.Indexes, Index{IsUnique:idx.IsUnique, ColumnNames:idx.ColumnNames})
	}
	return t
}


func (d *Database) analyzeReverseReferences(td *Table) {
	var td2 *Table

	var cd *Column
	for _, cd = range td.Columns {
		if cd.ForeignKey != nil {
			td2 = d.Table(cd.ForeignKey.ReferencedTable)
			if td2 == nil  {
				continue // pointing to a type table
			}
			// Determine the go name, which is the name used to refer to the reverse reference.
			// This is somewhat tricky, because there is no easy way to extract an expression for this.
			// For example, a team-member  relationship would simply be called team_id from the person side, so
			// where would you get the word "member". Our strategy will be to first look for something explicit in the
			// options, and if not found, try to create a name from the table and table names.

			objName := cd.DbName
			objName = strings.TrimSuffix(objName, d.ForeignKeySuffix)
			objName = strings.Replace(objName, td2.DbName, "", 1)
			if objName != "" {
				objName = UpperCaseIdentifier(objName)
			}
			goName,_ := cd.Options["GoName"].(string)
			if goName == "" {
				goName = UpperCaseIdentifier(td.DbName)
				if objName != "" {
					goName = goName + "As" + objName
				}
			}
			goPlural,_ := cd.Options["GoPlural"].(string)
			if goPlural == "" {
				goPlural = inflector.Pluralize(td.DbName)
				goPlural = UpperCaseIdentifier(goPlural)
				if objName != "" {
					goPlural = goPlural + "As" + objName
				}
			}
			goType := td.GoName
			goTypePlural := td.GoPlural
			ref := ReverseReference{
				DbTable:              td2.DbName,
				DbColumn:             td2.PrimaryKeyColumn().DbName, // NoSQL only
				AssociatedTableName:  td.DbName,
				AssociatedColumnName: cd.DbName,
				GoName:               goName,
				GoPlural:             goPlural,
				GoType:               goType,
				GoTypePlural:         goTypePlural,
				IsUnique:             cd.IsUnique,
			}

			td2.ReverseReferences = append(td2.ReverseReferences, &ref)
			cd.ForeignKey.RR = &ref
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
		GoName:                desc.GoName,
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
		IsDateOnly:			   desc.SubType == "date",
		IsTimeOnly:			   desc.SubType == "time",
		Comment:               desc.Comment,
		Options:               desc.Options,
	}

	if c.IsId {
		c.ColumnType = ColTypeString // We treat auto-generated ids as strings for cross database compatibility.
	}

	if c.GoName == "" {
		c.GoName = d.dbNameToGoName(c.DbName)
	}


	c.ModelName = LowerCaseIdentifier(c.DbName)
	return c
}

func (d *Database) analyzeForeignKeys(desc TableDescription) {
	t := d.Table(desc.Name)
	if t != nil {
		for _,col := range desc.Columns {
			d.analyzeForeignKey(t, col)
		}
	}
}

func (d *Database) analyzeForeignKey(t *Table, cd ColumnDescription) {
	c := t.columnMap[cd.Name]
	if cd.ForeignKey != nil {
		f := &ForeignKeyInfo{
			ReferencedTable:    cd.ForeignKey.ReferencedTable,
			ReferencedColumn:   cd.ForeignKey.ReferencedColumn,
			UpdateAction: FKActionFromString(cd.ForeignKey.UpdateAction),
			DeleteAction: FKActionFromString(cd.ForeignKey.DeleteAction),
		}
		goName := c.GoName
		suffix := UpperCaseIdentifier(d.ForeignKeySuffix)
		goName = strings.TrimSuffix(goName, suffix)
		f.GoName = goName
		f.IsType = d.IsTypeTable(cd.ForeignKey.ReferencedTable)

		if f.IsType {
			tt := d.TypeTable(cd.ForeignKey.ReferencedTable)
			f.GoType = tt.GoName
			f.GoTypePlural = tt.GoPlural
		} else {
			r := d.Table(cd.ForeignKey.ReferencedTable)
			f.GoType = r.GoName
			f.GoTypePlural = r.GoPlural
			fkc := r.GetColumn(cd.ForeignKey.ReferencedColumn)
			if fkc.IsId {
				c.ColumnType = ColTypeString // Always use strings to refer to auto-generated ids for cross database compatibility
			}
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
