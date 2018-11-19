package control_base

import (
	"github.com/spekary/goradd/pkg/page/control/control_base/table"
)


type TableI interface {
	table.TableI
}

// Table is the local override for the Button control. Tables are created by the framework in list forms.
type Table struct {
	table.Table
}