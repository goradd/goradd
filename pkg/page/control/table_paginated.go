package control

import (
	"github.com/spekary/goradd/pkg/page"
	"goradd-project/config"
)

type PaginatedTable struct {
	Table
	PaginatedControl
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


