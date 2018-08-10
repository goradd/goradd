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
