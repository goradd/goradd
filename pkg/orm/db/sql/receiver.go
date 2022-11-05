package sql

import (
	"fmt"
	. "github.com/goradd/goradd/pkg/orm/query"
	strings2 "github.com/goradd/goradd/pkg/strings"
	time2 "github.com/goradd/goradd/pkg/time"
	"log"
	"strconv"
	"strings"
	"time"
)

// SqlReceiver is an encapsulation of a way of receiving data from sql queries as interface{} pointers. This allows you
// to get data without knowing the type of data you are asking for ahead of time, and is easier for dealing with NULL fields.
// Some database drivers (MySql for one) return different results in fields depending on how you call the query (using
// a prepared statement can return different results than without one), or if the data does not quite fit (UInt64 in particular
// will return a string if the returned value is bigger than MaxInt64, but smaller than MaxUint64.)
//
// Pass the address of the R member to the sql.Scan method when using an object of this type,
// because there are some idiosyncrasies with
// how Go treats return values that prevents returning an address of R from a function
type SqlReceiver struct {
	R interface{}
}

// IntI returns the receiver as an interface to an int.
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
		return i
	case []byte:
		v := string(r.R.([]byte)[:])
		if v == "NULL" {
			return nil
		} // MariaDB does this for default values
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Panic(err)
		}
		return i

	default:
		log.Panicln("Unknown type returned from sql driver")
		return nil
	}
}

// UintI converts the value to an interface to a GO uint.
//
// Some drivers (like MySQL) return all integers as Int64. Its up to you to make sure
// you only use this on 32-bit uints or smaller.
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
		v := string(r.R.([]byte)[:])
		if v == "NULL" {
			return nil
		} // MariaDB does this for default values
		i, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			log.Panic(err)
		}
		return uint(i)
	default:
		log.Panicln("Unknown type returned from sql driver")
		return nil
	}
}

// Int64I returns the given value as an interface to an Int64
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
		v := string(r.R.([]byte)[:])
		if v == "NULL" {
			return nil
		} // MariaDB does this for default values
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Panic(err)
		}
		return i
	default:
		log.Panicln("Unknown type returned from sql driver")
		return nil
	}
}

// Uint64I returns a value as an interface to a UInt64.
//
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
	case uint64:
		return r.R
	case string: // Mysql returns this if the detected value is greater than int64 size
		i, err := strconv.ParseUint(r.R.(string), 10, 64)
		if err != nil {
			log.Panic(err)
		}
		return i
	case []byte:
		v := string(r.R.([]byte)[:])
		if v == "NULL" {
			return nil
		} // MariaDB does this for default values
		i, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			log.Panic(err)
		}
		return i
	default:
		log.Panicln("Unknown type returned from sql driver")
		return nil
	}
}

// BoolI returns the value as an interface to a boolean
func (r SqlReceiver) BoolI() interface{} {
	if r.R == nil {
		return nil
	}
	switch r.R.(type) {
	case bool:
		return r.R
	case int:
		return r.R.(int) != 0
	case int64:
		return r.R.(int64) != 0
	case string:
		b, err := strconv.ParseBool(r.R.(string))
		if err != nil {
			log.Panic(err)
		}
		return b
	case []byte:
		v := string(r.R.([]byte)[:])
		if v == "NULL" {
			return nil
		} // MariaDB does this for default values
		b, err := strconv.ParseBool(v)
		if err != nil {
			log.Panic(err)
		}
		return b

	default:
		log.Panicln("Unknown type returned from sql driver")
		return nil
	}
}

// StringI returns the value as an interface to a string
func (r SqlReceiver) StringI() interface{} {
	if r.R == nil {
		return nil
	}
	switch r.R.(type) {
	case string:
		return r.R
	case []byte:
		v := string(r.R.([]byte)[:])
		return v
	default:
		return fmt.Sprint(r.R)
	}
}

// FloatI returns the value as an interface to a float32 value.
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
		v := string(r.R.([]byte)[:])
		if v == "NULL" {
			return nil
		} // MariaDB does this for default values
		f, err := strconv.ParseFloat(v, 32)
		if err != nil {
			log.Panic(err)
		}
		return float32(f)
	default:
		log.Panicln("Unknown type returned from sql driver")
		return nil
	}
}

// DoubleI returns the value as a float64 interface
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
		v := string(r.R.([]byte)[:])
		if v == "NULL" {
			return nil
		} // MariaDB does this for default values
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			log.Panic(err)
		}
		return f
	default:
		log.Panicln("Unknown type returned from sql driver")
		return nil
	}
}

// TimeI returns the value as a time.Time value in UTC, or in the case of CURRENT_TIME, a string "now".
func (r SqlReceiver) TimeI() interface{} {
	if r.R == nil {
		return nil
	}

	var t time.Time
	var err error
	switch v := r.R.(type) {
	case time.Time:
		t = v
	case string:
		t = time2.FromSqlDateTime(v) // Note that this must always include timezone information if coming from a timestamp with timezone column
	case []byte:
		s := string(v)
		if s == "NULL" {
			return nil
		}
		u := strings.ToUpper(s)
		if strings2.StartsWith(u, "CURRENT_TIMESTAMP") {
			// Mysql version of now. This would only be asked for if we were looking for a default value.
			return "now"
		}
		t = time2.FromSqlDateTime(s)
		if err != nil {
			return nil
		}
		// TODO: SQL Lite, may return an int or float. Not sure we can support these.
	default:
		log.Panicln("Unknown type returned from sql driver")
		return nil
	}

	return t
}

// Unpack converts a SqlReceiver to a type corresponding to the given GoColumnType
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
	case ColTypeTime:
		return r.TimeI()
	case ColTypeFloat32:
		return r.FloatI()
	case ColTypeFloat64:
		return r.DoubleI()
	case ColTypeBool:
		return r.BoolI()
	default:
		return r.R
	}
}

// UnpackDefaultValue converts a SqlReceiver used to get the default value
// to a type corresponding to the given GoColumnType.
func (r SqlReceiver) UnpackDefaultValue(typ GoColumnType) interface{} {
	switch typ {
	case ColTypeBytes:
		return r.R
	case ColTypeString:
		s := r.StringI()
		if s == nil {
			return s
		}
		if s.(string) == "NULL" {
			return nil
		}
		// Unwrap single quotes coming from mariadb
		s = strings.Trim(s.(string), `"'`)
		return s

	case ColTypeInteger:
		return r.IntI()
	case ColTypeUnsigned:
		return r.UintI()
	case ColTypeInteger64:
		return r.Int64I()
	case ColTypeUnsigned64:
		return r.Uint64I()
	case ColTypeTime:
		return r.TimeI()
	case ColTypeFloat32:
		return r.FloatI()
	case ColTypeFloat64:
		return r.DoubleI()
	case ColTypeBool:
		return r.BoolI()
	default:
		return r.R
	}
}
