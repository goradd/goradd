package db

import (
	"database/sql"
	"fmt"
	"github.com/spekary/goradd/datetime"
	. "github.com/spekary/goradd/orm/query"
	"log"
	"strconv"
)

// SqlReceiver is an encapsulation of a way of receiving data from sql queries as interface{} pointers. This allows you
// to get data without knowing the type of data you are asking for ahead of time, and is easier for dealing with NULL fields.
// Some database drivers (MySql for one), return different results in fields depending on how you call the query (using
// a prepared statement can return different results than without one), or if the data does not quite fit (UInt64 in particular
// will return a string if the returned value is bigger than MaxInt64, but smaller than MaxUint64.)
//
// Pass the address of the R member to the sql.Scan method when using an object of this type. IsRequired because there are some idiosyncracies with
// how Go treats return values that would prevent returning an address of R from a function
type SqlReceiver struct {
	R interface{}
}

func (r SqlReceiver) IntI() interface{} {
	if r.R == nil {
		return nil
	}
	switch r.R.(type) {
	case int64:
		return int(r.R.(int64))
	case int:
		return r.R
	case string:
		i, err := strconv.Atoi(r.R.(string))
		if err != nil {
			log.Panic(err)
		}
		return int(i)
	case []byte:
		i, err := strconv.Atoi(string(r.R.([]byte)[:]))
		if err != nil {
			log.Panic(err)
		}
		return int(i)

	default:
		log.Panicln("Unknown type returned from sql driver")
		return nil
	}
}

// Some drivers (like MySQL) return all integers as Int64. This converts to basic golang uint. Its up to you to make sure
// you only use this on 32-bit uints or smaller
func (r SqlReceiver) UintI() interface{} {
	if r.R == nil {
		return nil
	}
	switch r.R.(type) {
	case int64:
		return uint(r.R.(int64))
	case int:
		return uint(r.R.(int))
	case uint:
		return r.R
	case string:
		i, err := strconv.ParseUint(r.R.(string), 10, 32)
		if err != nil {
			log.Panic(err)
		}
		return uint(i)
	case []byte:
		i, err := strconv.ParseUint(string(r.R.([]byte)[:]), 10, 32)
		if err != nil {
			log.Panic(err)
		}
		return uint(i)
	default:
		log.Panicln("Unknown type returned from sql driver")
		return nil
	}
}

func (r SqlReceiver) Int64I() interface{} {
	if r.R == nil {
		return nil
	}
	switch r.R.(type) {
	case int64:
		return r.R
	case int:
		return int64(r.R.(int))
	case string:
		i, err := strconv.ParseInt(r.R.(string), 10, 64)
		if err != nil {
			log.Panic(err)
		}
		return i
	case []byte:
		i, err := strconv.ParseInt(string(r.R.([]byte)[:]), 10, 64)
		if err != nil {
			log.Panic(err)
		}
		return i
	default:
		log.Panicln("Unknown type returned from sql driver")
		return nil
	}
}

// Some drivers (like MySQL) return all integers as Int64. This converts to uint64. Its up to you to make sure
// you only use this on 64-bit uints or smaller.
func (r SqlReceiver) Uint64I() interface{} {
	if r.R == nil {
		return nil
	}
	switch r.R.(type) {
	case int64:
		return uint64(r.R.(int64))
	case int:
		return uint64(r.R.(int))
	case string: // Mysql returns this if the detected value is greater than int64 size
		i, err := strconv.ParseUint(r.R.(string), 10, 64)
		if err != nil {
			log.Panic(err)
		}
		return i
	case []byte:
		i, err := strconv.ParseUint(string(r.R.([]byte)[:]), 10, 64)
		if err != nil {
			log.Panic(err)
		}
		return i
	default:
		log.Panicln("Unknown type returned from sql driver")
		return nil
	}
}

func (r SqlReceiver) BoolI() interface{} {
	if r.R == nil {
		return nil
	}
	switch r.R.(type) {
	case bool:
		return r.R
	case int:
		return (r.R.(int) != 0)
	case int64:
		return (r.R.(int64) != 0)
	case string:
		b, err := strconv.ParseBool(r.R.(string))
		if err != nil {
			log.Panic(err)
		}
		return b
	case []byte:
		b, err := strconv.ParseBool(string(r.R.([]byte)[:]))
		if err != nil {
			log.Panic(err)
		}
		return b

	default:
		log.Panicln("Unknown type returned from sql driver")
		return nil
	}
}

func (r SqlReceiver) StringI() interface{} {
	if r.R == nil {
		return nil
	}
	switch r.R.(type) {
	case string:
		return r.R
	case []byte:
		return string(r.R.([]byte)[:])
	default:
		return fmt.Sprint(r.R)
	}
}

func (r SqlReceiver) FloatI() interface{} {
	if r.R == nil {
		return nil
	}
	switch r.R.(type) {
	case float32:
		return r.R
	case float64:
		return float32(r.R.(float64))
	case string:
		f, err := strconv.ParseFloat(r.R.(string), 32)
		if err != nil {
			log.Panic(err)
		}
		return f
	case []byte:
		f, err := strconv.ParseFloat(string(r.R.([]byte)[:]), 32)
		if err != nil {
			log.Panic(err)
		}
		return f
	default:
		log.Panicln("Unknown type returned from sql driver")
		return nil
	}
}

func (r SqlReceiver) DoubleI() interface{} {
	if r.R == nil {
		return nil
	}
	switch r.R.(type) {
	case float32:
		return float64(r.R.(float32))
	case float64:
		return r.R
	case string:
		f, err := strconv.ParseFloat(r.R.(string), 64)
		if err != nil {
			log.Panic(err)
		}
		return f
	case []byte:
		f, err := strconv.ParseFloat(string(r.R.([]byte)[:]), 64)
		if err != nil {
			log.Panic(err)
		}
		return f
	default:
		log.Panicln("Unknown type returned from sql driver")
		return nil
	}
}

func (r SqlReceiver) TimeI() interface{} {
	// TODO: We need to adjust the output to match the timezone of the server
	// Check the parseTime parameter to the database to see if it reads the server timezone
	// Otherwise we need some kind of strategy to know what timezone the data is getting returned in
	// The timezone saved may be different than that read
	// Also the server and client timezones can be different

	if r.R == nil {
		return nil
	}

	switch r.R.(type) {
	case string:
		s := string(r.R.(string))
		if s == "CURRENT_TIMESTAMP" {
			return nil // database itself is handling the setting of the time
		}
		return datetime.FromSqlDateTime(s)
	case []byte:
		s := string(r.R.([]byte)[:])
		if s == "CURRENT_TIMESTAMP" {
			return nil // database itself is handling the setting of the time
		}
		return datetime.FromSqlDateTime(s)
	default:
		log.Panicln("Unknown type returned from sql driver")
		return nil
	}

}

// Convert an SqlReceiver to a type corresponding to the given GoColumnType
func (r SqlReceiver) Unpack(typ GoColumnType) interface{} {
	switch typ {
	case ColTypeBytes:
		return r.R
	case ColTypeString:
		return r.StringI()
	case ColTypeInteger:
		return r.IntI()
	case ColTypeUnsigned:
		return r.UintI()
	case ColTypeInteger64:
		return r.Int64I()
	case ColTypeUnsigned64:
		return r.Uint64I()
	case ColTypeDateTime:
		return r.TimeI()
	case ColTypeFloat:
		return r.FloatI()
	case ColTypeDouble:
		return r.DoubleI()
	case ColTypeBool:
		return r.BoolI()
	default:
		return r.R
	}
}

// ReceiveRows gets data from a sql result set and returns it as a slice of maps. Each row is mapped to its table name.
// If you provide table names, those will be used in the map. Otherwise it will get the table names out of the
// result set provided
func ReceiveRows(rows *sql.Rows, columnTypes []GoColumnType, columnNames []string) (values []map[string]interface{}) {
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
