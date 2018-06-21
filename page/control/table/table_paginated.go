package table

import (
	"github.com/spekary/goradd/page/control"
	"github.com/spekary/goradd/page"
	"goradd/config"
)

type PaginatedTable struct {
	Table
	control.PaginatedControl
}

func NewPaginatedTable(parent page.ControlI) *PaginatedTable {
	t := &PaginatedTable{}
	t.Init(t, parent)
	return t
}

func (t *PaginatedTable) Init(self page.ControlI, parent page.ControlI) {
	t.Table.Init(self, parent)
	t.PaginatedControl.SetPageSize(config.DefaultPageSize)
}


