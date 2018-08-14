package control_base

import (
	gr_control_base "github.com/spekary/goradd/page/control/control_base"
)


type PanelI interface {
	gr_control_base.PanelI
}

// Panel is a Goradd control that is a basic "div" wrapper and is also the bases for many kinds of custom java controls.
// This local override gives you the ability to affect all panels and controls based on panels.
type Panel struct {
	gr_control_base.Panel
}

