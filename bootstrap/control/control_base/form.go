package control_base

import (
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
)

// Support routines for form objects

// FormControlAddAttributes should be called from the DataAttributes function of any form object that should become
// a boostrap form-control type control
func FormControlAddAttributes(ctrl page.ControlI, attr *html.Attributes) {
	switch ctrl.ValidationState() {
	case page.Valid:
		attr.AddClass("is-valid")
	case page.Invalid:
		attr.AddClass("is-invalid")
	}
	attr.AddClass("form-control")
}
