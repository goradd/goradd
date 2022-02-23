package sql

// Helper utilities for extracting a description out of a database

import (
	"database/sql"
	"encoding/json"
	"github.com/goradd/gengen/pkg/maps"
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/orm/query"
	"log"
	"strconv"
	"strings"
)

const (
	SqlTypeUnknown   = "Unknown"
	SqlTypeBlob      = "Blob"
	SqlTypeVarchar   = "VarChar"
	SqlTypeChar      = "Char"
	SqlTypeText      = "Text"
	SqlTypeInteger   = "Int"
	SqlTypeTimestamp = "Timestamp"
	SqlTypeDatetime  = "DateTime"
	SqlTypeDate      = "Date"
	SqlTypeTime      = "Time"
	SqlTypeFloat     = "Float"
	SqlTypeDouble    = "Double"
	SqlTypeBool      = "Bool"
	SqlTypeDecimal   = "Decimal" // a fixed point type
)


// Find the json encoded list of options in the given string
func ExtractOptions(comment string) (options map[string]interface{}, remainingComment string, err error) {
	var optionString string
	firstIndex := strings.Index(comment, "{")
	lastIndex := strings.LastIndex(comment, "}")
	options = make(map[string]interface{})

	if firstIndex != -1 &&
		lastIndex != -1 &&
		lastIndex > firstIndex {

		optionString = comment[firstIndex : lastIndex+1]
		remainingComment = strings.TrimSpace(comment[:firstIndex] + comment[lastIndex+1:])

		err = json.Unmarshal([]byte(optionString), &options)
	}
	return
}

// GetDataDefLength will extract the length from the definition given a data definition description of the table.
// If more than one number, returns the first number
// Example:
//	bigint(21) -> 21
// varchar(50) -> 50
// decimal(10,2) -> 10
func GetDataDefLength(description string) int {
	var lastPos, lenPos int
	var size string
	if lenPos = strings.Index(description, "("); lenPos != -1 {
		lastPos = strings.LastIndex(description, ")")
		size = description[lenPos+1 : lastPos]
		sizes := strings.Split(size, ",")
		i, _ := strconv.Atoi(sizes[0])
		return i
	}
	return 0
}

// Retrieves a numeric value from the options, which is always going to return a float64
func getNumericOption(o map[string]interface{}, option string, defaultValue float64) (float64, bool) {
	if v := o[option]; v != nil {
		if v2, ok := v.(float64); !ok {
			return defaultValue, false
		} else {
			return v2, true
		}
	} else {
		return defaultValue, true
	}
}

// Retrieves a boolean value from the options
func getBooleanOption(o *maps.SliceMap, option string) (val bool, ok bool) {
	val, ok = o.LoadBool(option)
	return
}

// Extracts a minimum and maximum value from the option map, returning defaults if none was found, and making sure
// the boundaries of anything found are not exceeded
func GetMinMax(o map[string]interface{}, defaultMin float64, defaultMax float64, tableName string, columnName string) (min float64, max float64) {
	var errString string

	if columnName == "" {
		errString = "table " + tableName
	} else {
		errString = "table " + tableName + ":" + columnName
	}

	v, ok := getNumericOption(o, "min", defaultMin)
	if !ok {
		log.Print("Error in min value in comment for " + errString + ". Value is not a valid number.")
		min = defaultMin
	} else {
		if v < defaultMin {
			log.Print("Error in min value in comment for " + errString + ". Value is less than the allowed minimum.")
			min = defaultMin
		} else {
			min = v
		}
	}
	delete(o, "min")

	v, ok = getNumericOption(o, "max", defaultMax)
	if !ok {
		log.Print("Error in max value in comment for " + errString + ". Value is not a valid number.")
		max = defaultMax
	} else {
		if v > defaultMax {
			log.Print("Error in max value in comment for " + errString + ". Value is more than the allowed maximum.")
			max = defaultMax
		} else {
			max = v
		}
	}
	delete(o, "max")

	return
}

func FkRuleToAction(rule sql.NullString) db.FKAction {

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


// SqlReceiveRows gets data from a sql result set and returns it as a slice of maps.
//
// Each column is mapped to its column name.
// If you provide columnNames, those will be used in the map. Otherwise it will get the column names out of the
// result set provided.
func SqlReceiveRows(rows *sql.Rows,
	columnTypes []query.GoColumnType,
	columnNames []string,
	builder *Builder,
) []map[string]interface{} {

	var values []map[string]interface{}

	cursor := NewSqlCursor(rows, columnTypes, columnNames, nil)
	defer cursor.Close()
	for v := cursor.Next();v != nil;v = cursor.Next() {
		values = append(values, v)
	}
	if builder != nil {
		values = builder.unpackResult(values)
	}

	return values
}


// ReceiveRows gets data from a sql result set and returns it as a slice of maps. Each column is mapped to its column name.
// If you provide column names, those will be used in the map. Otherwise it will get the column names out of the
// result set provided
func sqlReceiveRows2(rows *sql.Rows, columnTypes []query.GoColumnType, columnNames []string) (values []map[string]interface{}) {
	var err error

	values = []map[string]interface{}{}

	columnReceivers := make([]SqlReceiver, len(columnTypes))
	columnValueReceivers := make([]interface{}, len(columnTypes))

	if columnNames == nil {
		columnNames, err = rows.Columns()
		if err != nil {
			log.Panic(err)
		}
	}

	for i, _ := range columnReceivers {
		columnValueReceivers[i] = &(columnReceivers[i].R)
	}

	for rows.Next() {
		err = rows.Scan(columnValueReceivers...)

		if err != nil {
			log.Panic(err)
		}

		v1 := make(map[string]interface{}, len(columnReceivers))
		for j, vr := range columnReceivers {
			v1[columnNames[j]] = vr.Unpack(columnTypes[j])
		}
		values = append(values, v1)

	}
	err = rows.Err()
	if err != nil {
		log.Panic(err)
	}
	return
}