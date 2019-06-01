package db

import (
	"github.com/gedex/inflector"
	. "github.com/goradd/goradd/pkg/orm/query"
	"github.com/knq/snaker"
	"log"
	"regexp"
	"strings"
)

// The DatabaseDescription is the top level struct that contains a complete description of a database for purposes of
// generating and creating queries
type DatabaseDescription struct {
	// The database key corresponding to its key in the global database cluster
	DbKey string
	// Tables are the tables in the database
	Tables []*TableDescription
	// TypeTables contains a description of the enumerated types from the type tables in the database
	TypeTables []*TypeTableDescription
	// The prefix for related objects.
	AssociatedObjectPrefix string

	// Text to strip off the end of foreign key references when converting to names. Defaults to "_id"
	ForeignKeySuffix string

	// These items are filled in by analysis

	// tableMap is used to get to tables by internal name
	tableMap map[string]*TableDescription
	// typeTableMap gets to type tables by internal name
	typeTableMap map[string]*TypeTableDescription
}

// FKAction indicates how the database handles situations when one side of a relationship is deleted or the key
// is changed. These generally correspond to the options available in MySQL InnoDB databases.
type FKAction int

// The foreign key actions tell us what the database will do automatically if a foreign key object is changed. This allows
// us to do the appropriate thing when we detect in the ORM that a linked object is changing.
const (
	FKActionNone FKAction = iota // In a typical database, this is the same as Restrict. For OUR purposes, it means we should deal with it ourselves.
	// This would be the situation when we are emulating foreign key constraints for databases that don't support them.
	FKActionSetNull
	FKActionSetDefault // Not supported in MySQL!
	FKActionCascade    //
	FKActionRestrict   // The database is going to choke on this. We will try to error before something like this happens.
)

/*
type Option struct {
	key   string
	value string
}
*/

// Describer is the interface for databases to return their DatabaseDescription
type Describer interface {
	Describe() *DatabaseDescription
}

// NewDatabaseDescription creates a new database description. This is the structure returned by database analysis.
func NewDatabaseDescription(dbKey string, objectPrefix string, fkSuffix string) *DatabaseDescription {
	dd := DatabaseDescription{
		DbKey:                  dbKey,
		Tables:                 []*TableDescription{},
		TypeTables:             []*TypeTableDescription{},
		AssociatedObjectPrefix: objectPrefix,
		ForeignKeySuffix:       fkSuffix,
	}
	return &dd
}

// NewTableDescription returns a new table description
func NewTableDescription(tableName string) *TableDescription {
	td := TableDescription{
		DbName:    tableName,
		Columns:   []*ColumnDescription{},
		columnMap: map[string]*ColumnDescription{},
	}
	return &td
}

// Given a database description, analyze will perform an analysis of the database, and modify some of the fields to prepare
// the description for use in codegen and the orm
func (dd *DatabaseDescription) analyze() {
	// initialize
	dd.typeTableMap = make(map[string]*TypeTableDescription)
	dd.tableMap = make(map[string]*TableDescription)

	if dd.AssociatedObjectPrefix == "" {
		dd.AssociatedObjectPrefix = "o"
	}

	dd.analyzeTypeTables()
	dd.analyzeTables()
}

// analyzeTypeTables will analyze the type tables provided by the database description
func (dd *DatabaseDescription) analyzeTypeTables() {
	for _, tt := range dd.TypeTables {
		dd.typeTableMap[tt.DbName] = tt
		tt.Constants = make(map[uint]string, len(tt.FieldNames))
		names := tt.FieldNames
		var key uint
		var value string
		var ok bool

		tt.EnglishName = dd.dbNameToEnglishName(tt.DbName)
		tt.EnglishPlural = dd.dbNameToEnglishPlural(tt.DbName)
		tt.GoName = dd.dbNameToGoName(tt.DbName)
		tt.LcGoName = strings.ToLower(tt.GoName[:1]) + tt.GoName[1:]
		tt.GoPlural = dd.dbNameToGoPlural(tt.DbName)

		if len(tt.Values) == 0 {
			log.Print("Warning: type table " + tt.DbName + " has no data entries. Specify constants by adding entries to this table.")
		}

		for _, m := range tt.Values {
			key, ok = m[names[0]].(uint)
			if !ok {
				key = uint(m[names[0]].(int))
			}
			value = m[names[1]].(string)
			var con string

			r := regexp.MustCompile("[^a-zA-Z0-9_]+")
			a := r.Split(value, -1)
			for _, word := range a {
				con += strings.Title(strings.ToLower(word))
			}
			tt.Constants[key] = con
		}

	}
}

// analyzeTables will analyze the table provided by the description
func (dd *DatabaseDescription) analyzeTables() {
	for _, td := range dd.Tables {
		td.ManyManyReferences = []*ManyManyReference{}
		td.ReverseReferences = []*ReverseReference{}
	}

	for _, td := range dd.Tables {
		if !td.IsAssociation {
			dd.analyzeTable(td)
		}
		if !td.Skip {
			dd.tableMap[td.DbName] = td
		}
	}

	for _, td := range dd.tableMap {
		if !td.IsAssociation {
			dd.analyzeColumns(td)
		}
	}

	for _, td := range dd.tableMap {
		if td.Skip {
			continue
		}
		if !td.IsAssociation {
			dd.analyzeReverseReferences(td)
		}
	}

	for _, td := range dd.Tables {
		if td.Skip {
			continue
		}
		if td.IsAssociation {
			dd.analyzeAssociation(td)
		}
	}

}

func (dd *DatabaseDescription) analyzeTable(td *TableDescription) {

	if td.LiteralName == "" {
		td.LiteralName = dd.dbNameToEnglishName(td.DbName)
	}

	if td.LiteralName == td.LiteralPlural {
		log.Print("Error: table " + td.DbName + " is a plural name. Change it to a singular name.")
		td.Skip = true
	}

	if td.GoName == "" {
		td.GoName = dd.dbNameToGoName(td.DbName)
		td.LcGoName = strings.ToLower(td.GoName[:1]) + td.GoName[1:]
	}
	if td.GoPlural == "" {
		td.GoPlural = dd.dbNameToGoPlural(td.DbName)
	}
}

func (dd *DatabaseDescription) analyzeColumns(td *TableDescription) {
	var pkCount = 0

	for _, cd := range td.Columns {
		dd.analyzeColumn(td, cd)
		if cd.IsPk {
			pkCount++
			if td.PrimaryKeyColumn == nil {
				td.PrimaryKeyColumn = cd
			} else {
				log.Print("Error: table " + td.DbName + " has multiple primary keys.")
				td.Skip = true
			}
		}
		if cd.ColumnType == ColTypeDateTime {
			td.HasDateTime = true
		}
	}
}

// Currently SQL only
func (dd *DatabaseDescription) analyzeReverseReferences(td *TableDescription) {
	var td2 *TableDescription

	if td.Skip {
		return
	}
	var cd *ColumnDescription
	for _, cd = range td.Columns {
		if cd.ForeignKey != nil {
			td2 = dd.TableDescription(cd.ForeignKey.TableName)
			if td2 == nil || td2.Skip {
				continue
			}
			if td2 == nil {
				// This is pointing to a type table?
				if _, ok := dd.typeTableMap[cd.ForeignKey.TableName]; !ok {
					log.Println("Error: Could not find foreign key table for table " + cd.ForeignKey.TableName)
					return
				}
			} else {
				// Determine the go name, which is the name used to refer to the reverse reference.
				// This is somewhat tricky, because there is no easy way to extract an expression for this.
				// For example, a team-member  relationship would simply be called team_id from the person side, so
				// where would you get the word "member". Our strategy will be to first look for something explicit in the
				// options, and if not found, try to create a name from the table and table names.

				goName, _ := cd.Options.LoadString("GoName")
				goPlural, _ := cd.Options.LoadString("GoPlural")
				objName := cd.DbName
				objName = strings.TrimSuffix(objName, dd.ForeignKeySuffix)
				objName = strings.Replace(objName, td2.DbName, "", 1)
				if objName != "" {
					objName = UpperCaseIdentifier(objName)
				}
				if goName == "" {
					goName = UpperCaseIdentifier(td.DbName)
					if objName != "" {
						goName = goName + "As" + objName
					}
				}
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
					DbColumn:             td2.PrimaryKeyColumn.DbName, // NoSQL only
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
}

// Analyzes an association table and creates special virtual columns in the corresponding tables it points to.
// Association tables are used by SQL databases to create many-many relationships. NoSQL databases can define their
// association columns directly and store an array of records numbers on either end of the association.
func (dd *DatabaseDescription) analyzeAssociation(td *TableDescription) {
	if len(td.Columns) != 2 {
		log.Print("Error: table " + td.DbName + " must have only 2 primary key columns.")
		return
	}
	var typeIndex = -1
	for i, cd := range td.Columns {
		if !cd.IsPk {
			log.Print("Error: table " + td.DbName + ":" + cd.DbName + " must be a primary key.")
			return
		}

		if cd.ForeignKey == nil {
			log.Print("Error: table " + td.DbName + ":" + cd.DbName + " must be a foreign key.")
			return
		}

		if cd.IsNullable {
			log.Print("Error: table " + td.DbName + ":" + cd.DbName + " cannot be nullable.")
			return
		}

		if _, ok := dd.typeTableMap[cd.ForeignKey.TableName]; ok {
			typeIndex = i
		}
	}

	if typeIndex == -1 { // normal table-table link
		ref1 := dd.makeManyManyRef(td, td.Columns[0], td.Columns[1])
		ref2 := dd.makeManyManyRef(td, td.Columns[1], td.Columns[0])
		ref1.MM = ref2
		ref2.MM = ref1

	} else {
		dd.makeManyManyRef(td, td.Columns[1-typeIndex], td.Columns[typeIndex])
	}
}

func (dd *DatabaseDescription) makeManyManyRef(td *TableDescription, cdSource *ColumnDescription, cdDest *ColumnDescription) *ManyManyReference {
	sourceTableName := cdSource.ForeignKey.TableName
	destTableName := cdDest.ForeignKey.TableName
	sourceObjName := strings.TrimSuffix(cdSource.DbName, dd.ForeignKeySuffix)
	destObjName := strings.TrimSuffix(cdDest.DbName, dd.ForeignKeySuffix)
	sourceTable := dd.TableDescription(sourceTableName)

	_, isType := dd.typeTableMap[cdDest.ForeignKey.TableName]

	var objName string
	if isType {
		destTable := dd.TypeTableDescription(destTableName)
		objName = destTable.GoName
	} else {
		destTable := dd.TableDescription(destTableName)
		objName = destTable.GoName
	}

	var goName, goPlural string

	if sourceObjName != sourceTableName {
		goName = UpperCaseIdentifier(destObjName) + "As" + UpperCaseIdentifier(sourceObjName)
		goPlural = inflector.Pluralize(UpperCaseIdentifier(destObjName)) + "As" + UpperCaseIdentifier(sourceObjName)
	} else {
		goName = UpperCaseIdentifier(destObjName)
		goPlural = inflector.Pluralize(goName)
	}

	ref := ManyManyReference{
		AssnTableName:        td.DbName,
		AssnColumnName:       cdSource.DbName,
		AssociatedTableName:  destTableName,
		AssociatedColumnName: cdDest.DbName,
		AssociatedObjectName: objName,
		GoName:               goName,
		GoPlural:             goPlural,
		Options:              td.Options,
		IsTypeAssociation:    isType,
	}
	sourceTable.ManyManyReferences = append(sourceTable.ManyManyReferences, &ref)
	return &ref
}

func (dd *DatabaseDescription) analyzeColumn(td *TableDescription, cd *ColumnDescription) {
	var err error
	cd.Options, err = extractOptions(cd.Comment)
	if err != nil {
		log.Println(err)
	}

	if cd.IsId {
		cd.ColumnType = ColTypeString // We treat auto-generated ids as strings for cross database compatibility.
	}
	if cd.ForeignKey != nil {
		goName := cd.GoName
		suffix := UpperCaseIdentifier(dd.ForeignKeySuffix)
		goName = strings.TrimSuffix(goName, suffix)
		cd.ForeignKey.GoName = goName
		cd.ForeignKey.IsType = dd.IsTypeTable(cd.ForeignKey.TableName)
		if cd.ForeignKey.IsType {
			tt := dd.TypeTableDescription(cd.ForeignKey.TableName)
			cd.ForeignKey.GoType = tt.GoName
		} else {
			td := dd.TableDescription(cd.ForeignKey.TableName)
			cd.ForeignKey.GoType = td.GoName
			fkc := td.GetColumn(cd.ForeignKey.ColumnName)
			if fkc.IsId {
				cd.ColumnType = ColTypeString // Always use strings to refer to auto-generated ids for cross database compatibility
			}
		}
	}

	cd.ModelName = LowerCaseIdentifier(cd.DbName)
}

func (dd *DatabaseDescription) dbNameToEnglishName(name string) string {
	return strings.Title(strings.Replace(name, "_", " ", -1))
}

func (dd *DatabaseDescription) dbNameToEnglishPlural(name string) string {
	return inflector.Pluralize(dd.dbNameToEnglishName(name))
}

func (dd *DatabaseDescription) dbNameToGoName(name string) string {
	return UpperCaseIdentifier(name)
}

func (dd *DatabaseDescription) dbNameToGoPlural(name string) string {
	return inflector.Pluralize(dd.dbNameToGoName(name))
}

// TableDescription returns a TableDescription from the database given the table name.
func (dd *DatabaseDescription) TableDescription(name string) *TableDescription {
	if v, ok := dd.tableMap[name]; ok {
		return v
	} else {
		return nil
	}
}

// TypeTableDescription returns a TypeTableDescription from the database given the table name.
func (dd *DatabaseDescription) TypeTableDescription(name string) *TypeTableDescription {
	return dd.typeTableMap[name]
}

// IsTypeTable returns true if the given name is the name of a type table in the database
func (dd *DatabaseDescription) IsTypeTable(name string) bool {
	_, ok := dd.typeTableMap[name]
	return ok
}

// GetColumn returns a ColumnDescription given the name of a column
func (td *TableDescription) GetColumn(name string) *ColumnDescription {
	return td.columnMap[name]
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
