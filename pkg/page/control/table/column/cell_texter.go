package column

import (
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control/table"
)

type CellTexter interface {
	table.CellTexter
}

func GetCellTexter(ctrl page.ControlI, id string) CellTexter {
	return ctrl.Page().GetControl(id).(CellTexter)
}
