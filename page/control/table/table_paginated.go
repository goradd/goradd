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

func NewPaginatedTable(parent page.ControlI, id string) *PaginatedTable {
	t := &PaginatedTable{}
	t.Init(t, parent, id)
	return t
}

func (t *PaginatedTable) Init(self page.ControlI, parent page.ControlI, id string) {
	t.Table.Init(self, parent, id)
	t.PaginatedControl.SetPageSize(config.DefaultPageSize)
}


