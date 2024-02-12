package config

import (
	"github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/strings"
	"github.com/microcosm-cc/bluemonday"
	"html"
)

// Sanitizer describes an object that can sanitize user input
type Sanitizer interface {
	Sanitize(string) string
}

// GlobalSanitizer is used by control.Textbox controls to sanitize user input before it is stored
var GlobalSanitizer Sanitizer

// BlueMondaySanitizer is a sanitizer based on microcosm-cc/bluemonday. BlueMonday is designed to sanitize input
// coming from a WYSIWYG HTML editor, so it has the annoying extra step of escaping HTML entities. We wrap the
// BlueMonday sanitizer in this structure so that we can unescape html entities before sending them to the textbox.
// We will still get all the stripping of javascript that the sanitizer normally does.
// If you want a different global sanitizer, change it here.
// Or, override the Sanitize function in the textbox object.
//
// This sanitizer is no longer used by default, because it removes too much valid text. For
// example, "a<b" is changed.
type BlueMondaySanitizer struct {
	policy *bluemonday.Policy
}

func (s BlueMondaySanitizer) Sanitize(in string) string {
	if !strings.IsUTF8(in) {
		log.FrameworkInfo("input string is not UTF-8")
		return ""
	}
	v := s.policy.Sanitize(in)
	v = html.UnescapeString(v)
	return v
}

type DefaultSanitizer struct {
}

func (s DefaultSanitizer) Sanitize(in string) string {
	if !strings.IsUTF8(in) {
		log.FrameworkInfo("input string is not UTF-8")
		return ""
	}
	if strings.HasNull(in) {
		log.FrameworkInfo("input string has a null value")
		return "" // completely reject any attempt to insert a null inside an input
	}
	return in
}

func init() {
	GlobalSanitizer = DefaultSanitizer{}
	//GlobalSanitizer = BlueMondaySanitizer{bluemonday.UGCPolicy()}
}
