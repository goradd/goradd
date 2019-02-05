package db

import (
	"fmt"
	"github.com/gedex/inflector"
	"github.com/knq/snaker"
	"github.com/goradd/gengen/pkg/maps"
	"github.com/goradd/goradd/pkg/datetime"
	. "github.com/goradd/goradd/pkg/orm/query"
	"log"
	"strings"
	"regexp"
)

// The DatabaseDescription is the top level struct that contains a complete description of a database for purposes of
// generating and creating queries
type DatabaseDescription struct {
	// The database key corresponding to its key in the global database cluster
	DbKey  string
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
	tableMap     map[string]*TableDescription
	// typeTableMap gets to type tables by internal name
	typeTableMap map[string]*TypeTableDescription
}

type TableDescription struct {
	// DbKey is the key used to find database in the global database cluster
	DbKey string
	// DbName is the name of the database table or object in the database.
	DbName string
	// LiteralName is the name of the object when describing it to the world. Use the "literalName" option in the comment to override the default. Should be lower case.
	LiteralName string
	// LiteralPlural is the plural name of the object. Use the "literalPlural" option in the comment to override the default. Should be lower case.
	LiteralPlural string
	// GoName is name of the struct when referring to it in go code. Use the "goName" option in the comment to override the default.
	GoName string
	// GoPlural is the name of a collection of these objects when referring to them in go code. Use the "goPlural" option in the comment to override the default.
	GoPlural string
	// LcGoName is the same as GoName, but with first letter lower case.
	LcGoName      string
	// Columns is a list of ColumnDescriptions, one for each column in the table.
	Columns       []*ColumnDescription
	// columnMap is an internal map of the columns
	columnMap     map[string]*ColumnDescription
	// Indexes are the indexes defined in the database. Unique indexes will result in LoadBy* functions.
	Indexes       []IndexDescription
	// Options are key-value pairs of values that can be used to customize how code generation is performed
	Options       maps.SliceMap
	// IsType is true if this is a type table
	IsType        bool
	// IsAssociation is true if this is an association table, which is used to create a many-to-many relationship between two tables.
	IsAssociation bool
	// Comment is the general comment included in the database
	Comment       string

	// The following items are filled in by the analyze process

	// ManyManyReferences describe the many-to-many references pointing to this table
	ManyManyReferences []*ManyManyReference
	// ReverseReferences describes the many-to-one references pointing to this table
	ReverseReferences  []*ReverseReference
	// HasDateTime is true if the table contains a DateTime column.
	HasDateTime bool
	// PrimaryKeyColumn points to the column that contains the primary key of the table.
	PrimaryKeyColumn *ColumnDescription

	// Skip will cause the table to be skipped in code generation
	Skip bool
}

// TypeTableDescription describes a type table, which essentially defines an enumerated type.
// In the SQL world, they are a table with an integer key (starting index 1) and a "name" value, though
// they can have other values associated with them too. Goradd will maintain the
// relationships in SQL, but in a No-SQL situation, it will embed all the ids and values.
type TypeTableDescription struct {
	// DbKey is the key used to find the database in the global database cluster
	DbKey string
	// DbName is the name of the table in the database
	DbName string
	// EnglishName is the english name of the object when describing it to the world. Use the "literalName" option in the comment to override the default.
	EnglishName string
	// EnglishPlural is the plural english name of the object. Use the "literalPlural" option in the comment to override the default.
	EnglishPlural string
	// GoName is the name of the item as a go type name.
	GoName     string
	// GoPlural is the plural of the go type
	GoPlural   string
	// LcGoName is the lower case version of the go name
	LcGoName   string
	// FieldNames are the names of the fields defined in the table. The first field name MUST be the name of the id field, and 2nd MUST be the name of the name field, the others are optional extra fields.
	FieldNames []string
	// FieldTypes are the go column types of the fields, indexed by field name
	FieldTypes map[string]GoColumnType
	// Values are the constant values themselves as defined in the table, mapped to field names in each row.
	Values     []map[string]interface{}
	// PkField is the name of the private key field
	PkField string

	// Filled in by analyzer
	Constants map[uint]string
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
	FKActionSetDefault  // Not supported in MySQL!
	FKActionCascade     //
	FKActionRestrict    // The database is going to choke on this. We will try to error before something like this happens.
)


// ColumnDescription describes a database column. Most of the information is either
// gleaned from the structure of the database, or is taken from a file that describes the relationships between
// different record types. Some of the information is filled in after analysis. Some of the information can be
// provided through information embedded in database comments.
type ColumnDescription struct {
	// DbName is the name of the column in the database. This is blank if this is a "virtual" table for sql tables like an association or virtual attribute query.
	DbName        string
	// GoName is the name of the column in go code
	GoName        string
	// NativeType is the type of the column as described by the database itself.
	NativeType    string
	//  ColumnType is the goradd defined column type
	ColumnType    GoColumnType
	// MaxCharLength is the maximum length of characters to allow in the column if a string type column.
	// If the database has the ability to specify this, this will correspond to what is specified.
	// In any case, we will generate code to prevent fields from getting bigger than this.
	MaxCharLength uint64
	// DefaultValue is the default value as specified by the database. We will initialize new ORM objects
	// with this value. It will be case to the corresponding GO type.
	DefaultValue  interface{}
	// MaxValue is the maximum value allowed for numeric values. This can be used by UI objects to tell the user what the limits are.
	MaxValue      interface{}
	// MinValue is the minimum value allowed for numeric values. This can be used by UI objects to tell the user what the limits are.
	MinValue      interface{}
	// IsId is true if this column represents a unique identifier generated by the database
	IsId          bool
	// IsPk is true if this is the primary key column. PK's do not necessarily need to be ID columns, and if not, we will need to do our own work to generate unique PKs.
	IsPk          bool
	// IsNullable is true if the column can be given a NULL value
	IsNullable    bool
	// IsIndexed is true if the column's table has a single index on the column, which will generate a LoadArrayBy function.
	IsIndexed             bool
	// IsUnique is true if the column's table has a single unique index on the column.
	IsUnique              bool
	// IsTimestamp is true if the field is a timestamp. Timestamps represent a specific point in world time.
	IsTimestamp 		  bool
	// IsAutoUpdateTimestamp is true if the database is updating the timestamp. Otherwise we will do it manually.
	IsAutoUpdateTimestamp bool
	// Comment is the contents of the comment associated with this field
	Comment               string

	// Filled in by analyzer
	// Options are the options extracted from the comments string
	Options    *maps.SliceMap
	// ForeignKey is additional information describing a foreign key relationship
	ForeignKey *ForeignKeyColumn
	// ModelName is a cache for the internal model name of this column.
	ModelName    string
}

// ForeignKeyColumn is additional information to describe what a foreign key points to
type ForeignKeyColumn struct {
	//DbKey string	// We don't support cross database foreign keys yet. Someday maybe.
	// TableName is the name of the table on the other end of the foreign key
	TableName    string
	// ColumnName is the database column name in the linked table that matches this column name. Often that is the primary key of the other table.
	ColumnName   string
	// UpdateAction indicates how the database will react when the other end of the relationship's value changes.
	UpdateAction FKAction
	// DeleteAction indicates how the database will react when the other end of the relationship's record is deleted.
	DeleteAction FKAction
	// GoName is the name we should use to refer to the related object
	GoName       string
	// GoType is the type of the related object
	GoType       string
	// IsType is true if this is a related type
	IsType       bool
	// RR is filled in by the analyzer and represent a reverse reference relationship
	RR           *ReverseReference
}


// The ManyManyReference structure is used by the templates during the codegen process to describe a many-to-any relationship.
type ManyManyReference struct {
	// AssnTableName is the database table creating the association. NoSQL: The originating table. SQL: The association table
	AssnTableName string
	// AssnColumnName is the column creating the association. NoSQL: The table storing the array of ids on the other end. SQL: the table in the association table pointing towards us.
	AssnColumnName string

	// AssociatedTableName is the database table being linked. NoSQL & SQL: The table we are joining to
	AssociatedTableName string
	// AssociatedColumnName is the database column being linked. NoSQL: table point backwards to us. SQL: Column in association table pointing forwards to refTable
	AssociatedColumnName string
	// AssociatedObjectName is the go name of the object created by this reference
	AssociatedObjectName string

	// GoName is the name used to refer to an object on the other end of the reference.
	GoName   string
	// GoPlural is the name used to refer to the group of objects on the other end of the reference.
	GoPlural string

	// IsTypeAssociation is true if this is a many-many relationship with a type table
	IsTypeAssociation bool
	// Options are the key-value options taken from the Comments in the association table, if there is one.
	Options           maps.SliceMap

	// MM is the many-many reference on the other end of the relationship that points back to this one.
	MM *ManyManyReference
}

// ReverseReference represents a kind of virtual column that is a result of a foreign-key
// pointing back to this column. This is the "one" side of a one-to-many relationship. Or, if
// the relationship is unique, this creates a one-to-one relationship.
// In SQL, since there is only a one-way foreign key, the side being pointed at does not have any direct
// data in a table indicating the relationship. We create a ReverseReference during data analysis and include
// it with the table description so that the table can know about the relationship and use it when doing queries.
type ReverseReference struct {
	// DbTable is the database table on the "one" end of the relationship. Its the table that the ReverseReference belongs to.
	DbTable              string
	// DbColumn is the database column that is referred to from the many end. This is most likely the primary key of DbTable.
	DbColumn             string
	// AssociatedTableName is the table on the "many" end that is pointing to the table containing the ReverseReference.
	AssociatedTableName  string
	// AssociatedColumnName is the column on the "many" end that is pointing to the table containing the ReverseReference. It is a foreign-key.
	AssociatedColumnName string
	// GoName is the name used to represent an object in the reverse relationship
	GoName               string
	// GoPlural is the name used to represent the group of objects in the reverse relationship
	GoPlural             string
	// GoType is the type of object in the collection of "many" objects, which corresponds to the name of the struct corresponding to the table
	GoType               string
	// IsUnique is true if the ReverseReference is unique. A unique reference creates a one-to-one relationship rather than one-to-many
	IsUnique             bool
	//Options maps.SliceMap
}

// IndexDescription is used by SQL analysis to extract details about an Index in the database. We can use indexes
// to know how to get to sorted data easily.
type IndexDescription struct {
	// KeyName is the name of the index
	KeyName     string
	// IsUnique indicates whether the index is for a unique index
	IsUnique    bool
	// IsPrimary indicates whether this is the index for the primary key
	IsPrimary   bool
	// ColumnNames are the columns that are part of the index
	ColumnNames []string
}

/*
// ForeignKeyDescription describes a foreign key relationship between columns in one table and columns in a different table.
// We currently allow the collection of multi-table and cross-database fk data, but we don't currently support them in codegen.
// ForeignKeyDescription is primarily for SQL databases that have specific support for foreign keys, like MySQL InnoDB tables.
type ForeignKeyDescription struct {
	KeyName         string
	Columns         []string
	RelationSchema  string
	RelationTable   string
	relationColumns []string // must be ordered that same as columns
}
*/

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
			for _,word := range a {
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
	var pkCount int = 0

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
				ref := ReverseReference{
					DbTable:              td2.DbName,
					DbColumn:             td2.PrimaryKeyColumn.DbName, // NoSQL only
					AssociatedTableName:  td.DbName,
					AssociatedColumnName: cd.DbName,
					GoName:               goName,
					GoPlural:             goPlural,
					GoType:               goType,
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
	var typeIndex int = -1
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
		goName = snaker.SnakeToCamel(destObjName) + "As" + snaker.SnakeToCamel(sourceObjName)
		goPlural = inflector.Pluralize(snaker.SnakeToCamel(destObjName)) + "As" + snaker.SnakeToCamel(sourceObjName)
	} else {
		goName = snaker.SnakeToCamel(destObjName)
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
				cd.ColumnType = ColTypeString // Always use strings to refer to auto-generated ids for cross database compatibility
			}
		}
	}

	if strings.Contains(cd.DbName, "_") {
		camel := snaker.SnakeToCamel(cd.DbName)
		cd.ModelName = strings.ToLower(camel[0:1]) + camel[1:]
	} else {
		cd.ModelName = cd.DbName
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

// FieldGoName returns the go name corresponding to the given field offset
func (tt *TypeTableDescription) FieldGoName(i int) string {
	if i >= len(tt.FieldNames) {
		return ""
	}
	fn := tt.FieldNames[i]
	fn = snaker.SnakeToCamel(fn)
	return fn
}

// FieldGoColumnType returns the GoColumnType corresponding to the given field offset
func (tt *TypeTableDescription) FieldGoColumnType(i int) GoColumnType {
	if i >= len(tt.FieldNames) {
		return ColTypeUnknown
	}
	fn := tt.FieldNames[i]
	ft := tt.FieldTypes[fn]
	return ft
}

// IsReference returns true if the column is a foreign key pointing to another table
func (cd *ColumnDescription) IsReference() bool {
	return cd.ForeignKey != nil
}

// DefaultValueAsConstant returns the default value of the column as a GO constant
func (cd *ColumnDescription) DefaultValueAsConstant() string {
	if cd.ColumnType == ColTypeDateTime {
		if cd.DefaultValue == nil {
			return "datetime.Zero" // pass this to datetime.NewDateTime()
		} else {
			d := cd.DefaultValue.(datetime.DateTime)
			if b, _ := d.MarshalText(); b == nil {
				return "datetime.Zero"
			} else {
				s := string(b[:])
				return fmt.Sprintf("%#v", s)
			}
		}
	} else if cd.DefaultValue == nil {
		return cd.ColumnType.DefaultValue()
	} else {
		return fmt.Sprintf("%#v", cd.DefaultValue)
	}
}

// DefaultValueAsValue returns the default value of the column as a GO value
func (cd *ColumnDescription) DefaultValueAsValue() string {
	if cd.DefaultValue == nil {
		v := cd.ColumnType.DefaultValue()
		if v == "" {
			return "nil"
		} else {
			return v
		}
	} else if cd.ColumnType == ColTypeDateTime {
		if b, _ := cd.DefaultValue.(datetime.DateTime).MarshalText(); b == nil {
			return cd.ColumnType.DefaultValue()
		} else {
			s := string(b[:])
			if cd.IsTimestamp {
				return fmt.Sprintf("datetime.NewTimestamp(%#v)", s)
			} else {
				return fmt.Sprintf("datetime.NewDateTime(%#v)", s)
			}
		}

	} else {
		return fmt.Sprintf("%#v", cd.DefaultValue)
	}
}

// DefaultConstantName returns the name of the default value constant that will be used to refer to the default value
func (cd *ColumnDescription) DefaultConstantName(tableName string) string {
	title := tableName + cd.GoName + "Default"
	return title
}

// GoVarName returns the name of the go object used to refer to the kind of object the foreign key points to.
func (fk *ForeignKeyColumn) GoVarName() string {
	return "obj" + fk.GoName
}
