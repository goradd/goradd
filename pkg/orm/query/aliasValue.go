package query

import (
	"fmt"
	"github.com/goradd/goradd/pkg/datetime"
	"log"
	"strconv"
)

// An AliasValue is returned by the GetAlias function that is generated for each type. You then convert the alias to a
// particular type to use it.
type AliasValue struct {
	value string
	isNil bool
}

// NewAliasValue is used by the ORM to wrap an aliased operation or computed value that was returned by a query.
// You would not normally call this function.
func NewAliasValue(a interface{}) AliasValue {
	switch v := a.(type) {
	case nil:
		return AliasValue{"", true}
	case []byte:
		return AliasValue{string(v[:]), false}
	default:
		return AliasValue{fmt.Sprint(a), false}
	}
}

// IsNil returns true if the value returned was a NULL value from the database.
func (a AliasValue) IsNil() bool {
	return a.isNil
}

// IsNull is the same as IsNil, and returns true if the operation return a NULL value from the database.
func (a AliasValue) IsNull() bool {
	return a.isNil
}

// String returns the value as a string. A NULL value will be an empty string.
func (a AliasValue) String() string {
	return string(a.value)
}

// Int returns the value as an integer.
func (a AliasValue) Int() int {
	i, err := strconv.ParseInt(a.String(), 10, 64)
	if err != nil {
		log.Panic(err)
	}
	return int(i)
}

// DateTime returns the value as a datetime.DateTime value.
func (a AliasValue) DateTime() datetime.DateTime {
	t, err := datetime.FromSqlDateTime(a.String())
	if err != nil {
		panic("Alias DateTime returned unparsable value: " + a.String() + " : " + err.Error())
	}
	return t
}

// Float returns the value as a float64
func (a AliasValue) Float() float64 {
	f, err := strconv.ParseFloat(a.String(), 64)
	if err != nil {
		log.Panic(err)
	}
	return f
}

// Bool returns the value as a bool.
func (a AliasValue) Bool() bool {
	b, err := strconv.ParseBool(a.String())
	if err != nil {
		log.Panic(err)
	}
	return b
}
