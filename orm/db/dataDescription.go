package db

import (
	"log"
	"strings"
	"github.com/knq/snaker"
	"github.com/gedex/inflector"
	"github.com/spekary/goradd/util/types"
	"fmt"
	"github.com/spekary/goradd/datetime"
	. "github.com/spekary/goradd/orm/query"
)

// The DatabaseDescription is the top level struct that contains a complete description of a database for purposes of
// generating and creating queries
type DatabaseDescription struct {
	// The database key corresponding to its key in the global database cluster
	DbKey      string
	Tables     []*TableDescription
	// Type tables contain a description of an enumerated type
	TypeTables []*TypeTableDescription
	// The prefix for related objects.
	AssociatedObjectPrefix string

	// Text to strip off the end of foreign key references when converting to names. Defaults to "_id"
	ForeignKeySuffix string

	// Items filled in by analysis

	// Map to get to tables by internal name
	tableMap map[string]*TableDescription
	typeTableMap map[string]*TypeTableDescription

}

type TableDescription struct {
	// Key used to find database in the global database cluster
	DbKey string
	// The name of the database table or object
	DbName string
	// The english name of the object when describing it to the world. Use the "englishName" option in the comment to override the default.
	EnglishName string
	// The plural english name of the object. Use the "englishPlural" option in the comment to override the default.
	EnglishPlural string
	// The name of the struct when referring to it in go code. Use the "goName" option in the comment to override the default.
	GoName string
	// The name of a collection of these objects when referring to them in go code. Use the "goPlural" option in the comment to override the default.
	GoPlural      string
	// same as GoName, but with first letter lower case
	LcGoName	  string
	Columns       []*ColumnDescription
	columnMap	  map[string]*ColumnDescription
	Indexes       []IndexDescription // Creates LoadBy functions. Mapped by index name.
	Options       types.OrderedMap
	IsType        bool
	IsAssociation bool
	Comment string

	// The following items are filled in by the analyze process
	ManyManyReferences []*ManyManyReference
	ReverseReferences []*ReverseReference

	PrimaryKeyColumn *ColumnDescription

	// Do not process this table
	Skip bool
}

// Type tables essentially define enumerated types. In the SQL world, they are a table with an integer key (starting index 1) and a "name" value, though
// they can have other values associated with them too. Goradd will maintain the
// relationships in SQL, but in a No-SQL situation, it will embed all the ids and values.
type TypeTableDescription struct {
	// Key used to find database in the global database cluster
	DbKey string
	// Name in the database
	DbName string
	// The english name of the object when describing it to the world. Use the "englishName" option in the comment to override the default.
	EnglishName string
	// The plural english name of the object. Use the "englishPlural" option in the comment to override the default.
	EnglishPlural string
	// The name of the item as a go type name.
	GoName string
	GoPlural string
	LcGoName string
	FieldNames []string	// The first field name MUST be the name of the id field, and 2nd MUST be the name of the name field, the others are optional extra fields
	FieldTypes map[string]GoColumnType
	Values []map[string]interface{}

	PkField string

	// Filled in by analyzer
	Constants map[uint]string
}

type FKAction int

// The foreign key actions tell us what the database will do automatically if a foreign key object is changed. This allows
// us to do the appropriate thing when we detect in the ORM that a linked object is changing.
const (
	FK_ACTION_NONE FKAction = iota	// In a typical database, this is the same as Restrict. For OUR purposes, it means we should deal with it ourselves.
									// This would be the situation when we are emulating foreign key constraints for databases that don't support them.
	FK_ACTION_SET_NULL
	FK_ACTION_SET_DEFAULT // Not supported in MySQL!
	FK_ACTION_CASCADE	//
	FK_ACTION_RESTRICT	// The database is going to choke on this. We will try to error before something like this happens.
)

type ForeignKeyType struct {
	//DbKey string	// We don't support cross database foreign keys yet. Someday maybe.
	TableName    string
	ColumnName   string
	UpdateAction FKAction
	DeleteAction FKAction
	GoName       string	// The name we should use to refer to the related object
	GoType		 string // The type of the related object
	IsType		 bool
	RR			 *ReverseReference // Filled in by analyzer, the corresponding
}

type ColumnDescription struct {
	DbName                string // name in database. Blank if this is a "virtual" column for sql tables. i.e. an association or virtual attribute query
	GoName                string
	NativeType            string       // type of column based on native description of column
	GoType                GoColumnType // represents a basic go type
	MaxCharLength         uint64       // To help fields limit the number of characters or bytes they will accept, when its known
	DefaultValue          interface{}  // gets cast to the goType
	MaxValue              interface{}
	MinValue              interface{}
	IsId                  bool // Is this a unique identifier that is generated by the database system?
	IsPk                  bool // A primary key. Could be generated by us.
	IsNullable            bool
	IsIndexed             bool // Does it have a single index on itself. Will generate a LoadBy function
	IsUnique              bool
	IsAutoUpdateTimestamp bool // Database is auto updating this timestamp, or
	ShouldUpdateTimestamp bool // We should generate code to auto update this timestamp.
	Comment          	  string

	// Filled in by analyzer
	Options			 *types.OrderedMap
	ForeignKey       *ForeignKeyType
	VarName			 string // code generating convenience. The name to use for the internal variable name corresponding to this column.
}

// The ManyManyReference structure is used by the templates during the codegen process to describe a many-to-any relationship.
type ManyManyReference struct {
	// NoSQL: The originating table. SQL: The association table
	AssnTableName string
	// NoSQL: The column storing the array of ids on the other end. SQL: the column in the association table pointing towards us.
	AssnColumnName string

	// NoSQL & SQL: The table we are joining to
	AssociatedTableName string
	// NoSQL: column point backwards to us. SQL: Column in association table pointing forwards to refTable
	AssociatedColumnName string

	AssociatedObjectName string

	// The virtual column names used to describe the objects on the other end of the association
	GoName string
	GoPlural string

	IsTypeAssociation bool
	Options types.OrderedMap

	MM *ManyManyReference // ManyManyReference pointing back towards this one
}

type ReverseReference struct {
	DbTable string
	DbColumn string
	AssociatedTableName string
	AssociatedColumnName string
	GoName string
	GoPlural string
	GoType string
	IsUnique bool
	//Options types.OrderedMap
}


type IndexDescription struct {
	keyName string
	isUnique bool
	isPrimary bool
	columnNames []string
}

// Describes a foreign key relationship between columns in one table and columns in a different table
// We currently allow the collection of multi-column and cross-database fk data, but we don't currently support them in codegen.
type ForeignKeyDescription struct {
	KeyName         string
	Columns         []string
	RelationSchema  string
	RelationTable   string
	relationColumns []string // must be ordered that same as columns
}

type Option struct {
	key string
	value string
}

type  Describer interface {
	Describe() *DatabaseDescription
}

func NewDatabaseDescription(dbKey string, objectPrefix string, fkSuffix string) *DatabaseDescription {
	dd := DatabaseDescription {
		DbKey:      dbKey,
		Tables:     []*TableDescription{},
		TypeTables: []*TypeTableDescription{},
		AssociatedObjectPrefix:objectPrefix,
		ForeignKeySuffix:fkSuffix,
	}
	return &dd
}

func NewTableDescription(tableName string) *TableDescription {
	td := TableDescription{
		DbName: tableName,
		Columns:[]*ColumnDescription{},
		columnMap:map[string]*ColumnDescription{},

	}
	return &td
}

// Given a database description, will perform an analysis of the database, and modify some of the fields to prepare
// the description for use in codegen and the orm
func (dd *DatabaseDescription) analyze()  {
	// initialize
	dd.typeTableMap = make(map[string]*TypeTableDescription)
	dd.tableMap = make(map[string]*TableDescription)

	if dd.AssociatedObjectPrefix == "" {
		dd.AssociatedObjectPrefix = "o"
	}

	dd.analyzeTypeTables()
	dd.analyzeTables()
}

func (dd *DatabaseDescription) analyzeTypeTables()  {
	for _, tt := range dd.TypeTables {
		dd.typeTableMap[tt.DbName] = tt
		tt.Constants = make(map[uint]string, len(tt.FieldNames))
		names := tt.FieldNames
		var key uint
		var value string
		var ok bool
		for _,m := range tt.Values {
			key,ok = m[names[0]].(uint)
			if !ok {
				key = uint(m[names[0]].(int))
			}
			value = m[names[1]].(string)
			var con string
			if strings.ToUpper(value) == value {
				// All upper case values
				con = strings.Replace(value, " ", "_", -1)
			} else {
				con = strings.Replace(value, " ", "", -1)
				con = snaker.CamelToSnake(con)
				con = strings.ToUpper(con)
			}
			tt.Constants[key] = con

			if tt.EnglishName == "" {
				tt.EnglishName = dd.dbNameToEnglishName(tt.DbName)
			}
			if tt.EnglishPlural == "" {
				tt.EnglishPlural = dd.dbNameToEnglishPlural(tt.DbName)
			}
			if tt.GoName == "" {
				tt.GoName = dd.dbNameToGoName(tt.DbName)
				tt.LcGoName = strings.ToLower(tt.GoName[:1]) + tt.GoName[1:]
			}
			if tt.GoPlural == "" {
				tt.GoPlural = dd.dbNameToGoPlural(tt.DbName)
			}
		}

	}
}

func (dd *DatabaseDescription) analyzeTables() {
	for _,td := range dd.Tables {
		td.ManyManyReferences = []*ManyManyReference{}
		td.ReverseReferences = []*ReverseReference{}
	}

	for _,td := range dd.Tables {
		if !td.IsAssociation {
			dd.analyzeTable(td)
		}
		dd.tableMap[td.DbName] = td
	}

	for _,td := range dd.tableMap {
		if !td.IsAssociation {
			dd.analyzeColumns(td)
		}
	}

	for _,td := range dd.tableMap {
		if td.Skip {
			continue
		}
		if !td.IsAssociation {
			dd.analyzeReverseReferences(td)
		}
	}

	for _,td := range dd.Tables {
		if td.Skip {
			continue
		}
		if td.IsAssociation {
			dd.analyzeAssociation(td)
		}
	}

}

func (dd *DatabaseDescription) analyzeTable(td *TableDescription)  {

	if td.EnglishName == "" {
		td.EnglishName = dd.dbNameToEnglishName(td.DbName)
	}
	if td.EnglishPlural == "" {
		td.EnglishPlural = dd.dbNameToEnglishPlural(td.DbName)
	}
	if td.GoName == "" {
		td.GoName = dd.dbNameToGoName(td.DbName)
		td.LcGoName = strings.ToLower(td.GoName[:1]) + td.GoName[1:]
	}
	if td.GoPlural == "" {
		td.GoPlural = dd.dbNameToGoPlural(td.DbName)
	}
}

func (dd *DatabaseDescription) analyzeColumns(td *TableDescription)  {
	var pkCount int = 0

	for _,cd := range td.Columns {
		dd.analyzeColumn(td, cd)
		if cd.IsPk {
			pkCount ++
			if td.PrimaryKeyColumn == nil {
				td.PrimaryKeyColumn = cd
			} else {
				log.Print("Error: table " + td.DbName + " has multiple primary keys.")
				td.Skip = true
			}
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
	for _,cd = range td.Columns {
		if cd.ForeignKey != nil {
			td2 = dd.TableDescription(cd.ForeignKey.TableName)
			if td2 == nil || td2.Skip {
				continue
			}
			if td2 == nil {
				// This is pointing to a type table?
				if _,ok := dd.typeTableMap[cd.ForeignKey.TableName]; !ok {
					log.Println("Error: Could not find foreign key table for table " + cd.ForeignKey.TableName)
					return
				}
			} else {
				// Determine the go name, which is the name used to refer to the reverse reference.
				// This is somewhat tricky, because there is no easy way to extract an expression for this.
				// For example, a team-member  relationship would simply be called team_id from the person side, so
				// where would you get the word "member". Our strategy will be to first look for something explicit in the
				// options, and if not found, try to create a name from the column and table names.

				goName,_ := cd.Options.GetString("GoName")
				goPlural,_ := cd.Options.GetString("GoPlural")
				objName := cd.DbName
				objName = strings.TrimSuffix(objName, dd.ForeignKeySuffix)
				objName = strings.Replace(objName, td2.DbName, "",1)
				if goName == "" {
					goName = td.DbName
					if objName != "" {
						goName = goName + "_as_" + objName
					}
					goName = snaker.SnakeToCamel(goName)
				}
				if goPlural == "" {
					goPlural = inflector.Pluralize(td.DbName)
					if objName != "" {
						goPlural = goPlural + "_as_" + objName
					}
					goPlural = snaker.SnakeToCamel(goPlural)
				}
				goType := td.GoName
				ref := ReverseReference {
					DbTable: td2.DbName,
					DbColumn: td2.PrimaryKeyColumn.DbName, // NoSQL only
					AssociatedTableName: td.DbName,
					AssociatedColumnName: cd.DbName,
					GoName: goName,
					GoPlural: goPlural,
					GoType: goType,
					IsUnique: cd.IsUnique,
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
	var typeIndex int = -1
	for i,cd := range td.Columns {
		if !cd.IsPk {
			log.Print("Error: column " + td.DbName + ":" + cd.DbName + " must be a primary key.")
			return
		}

		if cd.ForeignKey == nil {
			log.Print("Error: column " + td.DbName + ":" + cd.DbName + " must be a foreign key.")
			return
		}

		if cd.IsNullable {
			log.Print("Error: column " + td.DbName + ":" + cd.DbName + " cannot be nullable.")
			return
		}

		if _, ok := dd.typeTableMap[cd.ForeignKey.TableName]; ok {
			typeIndex = i
		}
	}

	if typeIndex == -1 {	// normal column-column link
		ref1 := dd.makeManyManyRef(td, td.Columns[0], td.Columns[1])
		ref2 := dd.makeManyManyRef(td, td.Columns[1], td.Columns[0])
		ref1.MM = ref2
		ref2.MM = ref1

	} else {
		dd.makeManyManyRef(td, td.Columns[1-typeIndex], td.Columns[typeIndex])
	}
}

func (dd *DatabaseDescription) makeManyManyRef(td *TableDescription, cdSource *ColumnDescription, cdDest *ColumnDescription)  *ManyManyReference {
	sourceTableName := cdSource.ForeignKey.TableName
	destTableName := cdDest.ForeignKey.TableName
	sourceObjName := strings.TrimSuffix(cdSource.DbName, dd.ForeignKeySuffix)
	destObjName := strings.TrimSuffix(cdDest.DbName, dd.ForeignKeySuffix)
	sourceTable := dd.TableDescription(sourceTableName)

	_,isType := dd.typeTableMap[cdDest.ForeignKey.TableName]

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
		goName = snaker.SnakeToCamel(destObjName) + "As" + snaker.SnakeToCamel(sourceObjName)
		goPlural = inflector.Pluralize(snaker.SnakeToCamel(destObjName)) + "As" + snaker.SnakeToCamel(sourceObjName)
	} else {
		goName = snaker.SnakeToCamel(destObjName)
		goPlural = inflector.Pluralize(goName)
	}

	ref := ManyManyReference {
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
	sourceTable.ManyManyReferences = append (sourceTable.ManyManyReferences, &ref)
	return &ref
}


func (dd *DatabaseDescription) analyzeColumn(td *TableDescription, cd *ColumnDescription) {
	var err error
	cd.Options, err = extractOptions(cd.Comment)
	if err != nil {
		log.Println(err)
	}

	if cd.IsId {
		cd.GoType = COL_TYPE_STRING	// We treat auto-generated ids as strings for cross database compatibility.
	}
	if cd.ForeignKey != nil {
		goName := cd.GoName
		suffix := snaker.SnakeToCamel(dd.ForeignKeySuffix)
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
				cd.GoType = COL_TYPE_STRING // Always use strings to refer to auto-generated ids for cross database compatibility
			}
		}
	}

	if strings.Contains(cd.DbName, "_") {
		camel := snaker.SnakeToCamel(cd.DbName)
		cd.VarName = strings.ToLower(camel[0:1]) + camel[1:]
	} else {
		cd.VarName = cd.DbName
	}
}

func (dd *DatabaseDescription) dbNameToEnglishName(name string) string {
	return strings.Title(strings.Replace(name, "_", " ", -1))
}

func (dd *DatabaseDescription) dbNameToEnglishPlural(name string) string {
	return inflector.Pluralize(dd.dbNameToEnglishName(name))
}

func (dd *DatabaseDescription) dbNameToGoName(name string) string {
	return snaker.SnakeToCamel(name)
}

func (dd *DatabaseDescription) dbNameToGoPlural(name string) string {
	return inflector.Pluralize(dd.dbNameToGoName(name))
}


func (dd *DatabaseDescription) TableDescription(name string) *TableDescription {
	if v,ok := dd.tableMap[name]; ok {
		return v
	} else {
		return nil
	}
}

func (dd *DatabaseDescription) TypeTableDescription(name string) *TypeTableDescription {
	return dd.typeTableMap[name]
}

func (dd *DatabaseDescription) IsTypeTable(name string) bool {
	_,ok := dd.typeTableMap[name]
	return ok
}

func (td *TableDescription) GetColumn(name string) *ColumnDescription {
	return td.columnMap[name]
}

// Return the go name corresponding to the given field offset
func (tt *TypeTableDescription) FieldGoName(i int) string {
	if i >= len(tt.FieldNames) {
		return ""
	}
	fn := tt.FieldNames[i]
	fn = snaker.SnakeToCamel(fn)
	return fn
}

// Return the go type corresponding to the given field offset
func (tt *TypeTableDescription) FieldGoType(i int) string {
	if i >= len(tt.FieldNames) {
		return ""
	}
	fn := tt.FieldNames[i]
	ft := tt.FieldTypes[fn]
	return ft.String()
}


func (cd *ColumnDescription) IsReference() bool {
	return cd.ForeignKey != nil
}

// Returns the default
func (cd *ColumnDescription) DefaultValueAsConstant() string {
	if cd.GoType == COL_TYPE_DATETIME {
		if cd.DefaultValue == nil {
			return "datetime.Zero"	// pass this to datetime.NewDateTime()
		} else {
			d := cd.DefaultValue.(datetime.DateTime)
			if b,_ := d.MarshalText(); b == nil {
				return "datetime.Zero"
			} else {
				s := string(b[:])
				return s
			}
		}
	} else if cd.DefaultValue == nil {
		return cd.GoType.DefaultValue()
	} else {
		return fmt.Sprintf("%#v", cd.DefaultValue)
	}
}

func (cd *ColumnDescription) DefaultValueAsValue() string {
	if cd.DefaultValue == nil {
		v := cd.GoType.DefaultValue()
		if v == "" {
			return "nil"
		} else {
			return v
		}
	} else if cd.GoType == COL_TYPE_DATETIME {
		if b, _ := cd.DefaultValue.(datetime.DateTime).MarshalText(); b == nil {
			return cd.GoType.DefaultValue()
		} else {
			s := string(b[:])
			return fmt.Sprintf("datetime.NewDateTime(%#v)", s)
		}

	} else {
		return fmt.Sprintf("%#v", cd.DefaultValue)
	}
}


func (cd *ColumnDescription) DefaultConstantName(td *TableDescription) string {
	title := td.GoName + cd.GoName + "Default"
	title = snaker.CamelToSnake(title)
	title = strings.ToUpper(title)
	return title
}

func (fk *ForeignKeyType) GoVarName() string {
	return "obj" + fk.GoName
}