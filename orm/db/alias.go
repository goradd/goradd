package db

import (
	"github.com/spekary/goradd/datetime"
	"strconv"
	"log"
	"fmt"
)

// Aliases are returned by the GetAlias function that is generated for each type. You then convert the alias to a
// particular type to use it.


type Alias struct {
	value string
	isNil bool
}

func NewAlias(a interface{}) Alias {
	switch v := a.(type) {
	case nil:
		return  Alias{"", true}
	case []byte:
		return Alias{string(v[:]), false}
	default:
		return Alias{fmt.Sprint(a), false}
	}
}

func (a Alias) IsNil() bool {
	return a.isNil
}

func (a Alias) IsNull() bool {
	return a.isNil
}


func (a Alias) String() string {
	return string(a.value[:])
}

func (a Alias) Int() int {
	i, err := strconv.ParseInt(a.String(), 10, 64)
	if err != nil {
		log.Panic(err)
	}
	return int(i)
}

func (a Alias) DateTime() datetime.DateTime {
	return datetime.FromSqlDateTime(a.String())
}

func (a Alias) Float() float64 {
	f,err := strconv.ParseFloat(a.String(), 64)
	if err != nil {
		log.Panic(err)
	}
	return f
}

func (a Alias) Bool() bool {
	b, err := strconv.ParseBool(a.String())
	if err != nil {
		log.Panic(err)
	}
	return b
}