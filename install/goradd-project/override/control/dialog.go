package control

import (
	"github.com/spekary/goradd/pkg/page"
	"github.com/spekary/goradd/pkg/page/control"
	bs "github.com/spekary/goradd/pkg/bootstrap/control"
)


// TODO: Create ErrorAlert, WarningAlert, InfoAlert, YesNo Alert and Alert functions and put them in an interface.
func ErrorAlert(form page.FormI, msg string) {
	d := control.Alert(form, msg, "OK")
	d.SetTitle(form.T("Error"))
	d.(*bs.Modal).AddTitlebarClass("bg-error text-light")
}
