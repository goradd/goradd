package query

import (
	"fmt"
	"github.com/goradd/goradd/pkg/datetime"
	"log"
	"strconv"
)

// AliasValues are returned by the GetAlias function that is generated for each type. You then convert the alias to a
// particular type to use it.

type AliasValue struct {
	value string
	isNil bool
}

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

func (a AliasValue) IsNil() bool {
	return a.isNil
}

func (a AliasValue) IsNull() bool {
	return a.isNil
}

func (a AliasValue) String() string {
	return string(a.value[:])
}

func (a AliasValue) Int() int {
	i, err := strconv.ParseInt(a.String(), 10, 64)
	if err != nil {
		log.Panic(err)
	}
	return int(i)
}

func (a AliasValue) DateTime() datetime.DateTime {
	return datetime.FromSqlDateTime(a.String())
}

func (a AliasValue) Float() float64 {
	f, err := strconv.ParseFloat(a.String(), 64)
	if err != nil {
		log.Panic(err)
	}
	return f
}

func (a AliasValue) Bool() bool {
	b, err := strconv.ParseBool(a.String())
	if err != nil {
		log.Panic(err)
	}
	return b
}
