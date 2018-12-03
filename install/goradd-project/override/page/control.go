package page

import (
	grPage "github.com/spekary/goradd/pkg/page"
	//bsPage "github.com/spekary/goradd/pkg/bootstrap/control/page"
)

// The public interface for control overrides
type ControlI interface {
	grPage.ControlI
}

// The local Control override. All controls descend from this one, and so this gives you an opportunity to affect
// how all controls work by making changes here. The only control that does not descend from this is the Form control.
// See the FormBase struct in the override directory to override things specific to the Form control.
type Control struct {
	grPage.Control
	//bsPage.BootstrapControl // implements bootstrap controls for all controls

}

/*
func (c *Control) Init(self ControlI, parent ControlI) {
	c.Control.Init(self, parent, id)

	// Put additional initializations here
}
 */

 // this() supports object-oriented features. It gives easy access to the interface to call virtual functions.
 func (c *Control) this() ControlI {
 	return c.Self.(ControlI)
 }

// You can put overrides that should apply to all your controls here.


