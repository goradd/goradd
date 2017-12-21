package goradd

type BaseI interface {
}

type Base struct {
	Self interface{}
}

func (b *Base) Init (self interface{}) {
	b.Self = self
}
