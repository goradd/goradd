package control_base

import "github.com/spekary/goradd/pkg/page/control/control_base"


type ButtonI interface {
	control_base.ButtonI
}

// Button is the local override for the Button control. Buttons are created by the framework in the CRUD forms,
// and you can override that process here.
type Button struct {
	control_base.Button
}