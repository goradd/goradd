package control_base

import (
	gr_control_base "github.com/spekary/goradd/page/control/control_base"
)


type CheckboxI interface {
	gr_control_base.CheckboxI
}

// Checkbox is the local override for the Checkbox base class, which is the foundation for both html checkboxes and
// radio buttons. You can make any local changes here that affects both of these controls throughout your app.
type Checkbox struct {
	gr_control_base.Checkbox
}
