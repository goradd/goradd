package column

import (
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
)

type CellTexter interface {
	control.CellTexter
}

func GetCellTexter(ctrl page.ControlI, id string) CellTexter {
	return ctrl.Page().GetControl(id).(CellTexter)
}