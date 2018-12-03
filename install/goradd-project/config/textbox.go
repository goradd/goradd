package config

import (
	"github.com/microcosm-cc/bluemonday"
	"html"
)

// Sanitizer describes an object that can sanitize user input
type Sanitizer interface {
	Sanitize(string) string
}

// GlobalSanitizer is used by textboxes to sanitize user input before it is stored
var GlobalSanitizer Sanitizer


// BlueMondaySanitizer is a sanitizer based on microcosm-cc/bluemonday. BlueMonday is designed to sanitize input
// coming from a WYSIWYG editor, so it has the annoying extra step of escaping HTML entities. We wrap the
// BlueMonday sanitizer in this structure so that we can unescape html entities before sending them to the textbox.
// We will still get all the stripping of javascript that the sanitizer normally does.
// If you want a different global sanitizer, change it here. Or, override the Sanitize function in the textbox object.
type BlueMondaySanitizer struct {
	policy *bluemonday.Policy
}

func (s BlueMondaySanitizer) Sanitize(in string) string {
	v := s.policy.Sanitize(in)
	v = html.UnescapeString(v)
	return v
}

func init() {
	GlobalSanitizer = BlueMondaySanitizer{bluemonday.StrictPolicy()}
}
