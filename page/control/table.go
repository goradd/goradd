package control

import (
	"goradd-project/override/control_base"
	"github.com/spekary/goradd/page"
)

const (
	ColumnAction = iota + 2000
	SortClick
)

type TableI interface {
	control_base.TableI
}

type Table struct {
	control_base.Table
}

func NewTable(parent page.ControlI, id string) *Table {
	t := &Table{}
	t.Init(t, parent, id)
	return t
}

