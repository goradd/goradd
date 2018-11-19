package control_base

import (
	"encoding/gob"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spekary/goradd/pkg/page"
	gr_control_base "github.com/spekary/goradd/pkg/page/control/control_base"
	html2 "html"
)

var sanitizer mySanitizer

type mySanitizer struct {
	sanitizer2 *bluemonday.Policy
}

func (s *mySanitizer) Sanitize(in string) string {
	v := s.sanitizer2.Sanitize(in)
	v = html2.UnescapeString(v)	// Fix an issue with bluemonday
	return v
}

type TextboxI interface {
	gr_control_base.TextboxI
}

// The local Textbox override. All textboxes will descend from this one. You can make changes here that will impact
// all the text fields in the system.
type Textbox struct {
	gr_control_base.Textbox
}

func NewTextbox(parent page.ControlI, id string) *Textbox {
	t := &Textbox{}
	t.Init(t, parent, id)
	return t
}

func (t *Textbox) Init(self gr_control_base.TextboxI, parent page.ControlI, id string) {
	t.Textbox.Init(self, parent, id)
	t.Textbox.SetSanitizer(&sanitizer)
}

func init() {
	sanitizer.sanitizer2 = bluemonday.StrictPolicy()
	gob.Register(&mySanitizer{})
}
