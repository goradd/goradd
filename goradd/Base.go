package goradd

import (
	"reflect"
)

type BaseI interface {
	Type() string
}

type Base struct {
	Self interface{}
}

func (b *Base) Init(self interface{}) {
	b.Self = self
}

func (b *Base) Type() string {
	return reflect.TypeOf(b.Self).String()
}

func (b *Base) GobEncode() (data []byte, err error) {
	panic("the Base class is not encodable")
}

func (b *Base) GobDecode(data []byte) (err error) {
	panic("the Base class is not decodable")
}