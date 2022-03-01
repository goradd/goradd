package html

import (
	"fmt"
	"github.com/goradd/goradd/pkg/config" // TODO: remove reliance on outside package
	html2 "html"
	"strings"
)

// The LabelDrawingMode describes how to draw a label when it is drawn.
// Various CSS frameworks expect it a certain way. Many are not very forgiving when
// you don't do it the way they expect.
type LabelDrawingMode int

const (
	// LabelDefault means the mode is defined elsewhere, like in a config setting
	LabelDefault LabelDrawingMode = iota
	// LabelBefore indicates the label is in front of the control.
	// Example: <label>MyLabel</label><input ... />
	LabelBefore
	// LabelAfter indicates the label is after the control.
	// Example: <input ... /><label>MyLabel</label>
	LabelAfter
	// LabelWrapBefore indicates the label is before the control's tag, and wraps the control tag.
	// Example: <label>MyLabel<input ... /></label>
	LabelWrapBefore
	// LabelWrapAfter indicates the label is after the control's tag, and wraps the control tag.
	// Example: <label><input ... />MyLabel</label>
	LabelWrapAfter
)

// VoidTag represents a void tag, which is a tag that does not need a matching closing tag.
type VoidTag struct {
	Tag  string
	Attr Attributes
}

// Render returns the rendered version of the tag.
func (t VoidTag) Render() string {
	return RenderVoidTag(t.Tag, t.Attr)
}

// RenderVoidTag renders a void tag using the given tag name and attributes.
func RenderVoidTag(tag string, attr Attributes) (s string) {
	if attr == nil {
		s = "<" + tag + " />"
	} else {
		s = "<" + tag + " " + attr.String() + " />"
	}
	if !config.Minify {
		s += "\n"
	}
	return s
}

// RenderTag renders a standard html tag with a closing tag.
//
// innerHtml is html, and must already be escaped if needed.
//
// The tag will be surrounded with newlines to force general formatting consistency.
// This will cause the tag to be rendered with a space between it and its neighbors if the tag is
// not a block tag.
//
// In the few situations where you would want to
// get rid of this space, call RenderTagNoSpace()
func RenderTag(tag string, attr Attributes, innerHtml string) string {
	b := strings.Builder{}

	_,_ = fmt.Fprint(&b, "<", tag)
	if attr != nil {
		_,_ = fmt.Fprint(&b, " ", attr.String())
	}
	_,_ = fmt.Fprint(&b, ">")

	if innerHtml != "" {
		if !config.Minify {
			innerHtml = Indent(innerHtml)
		}
		b.WriteString("\n") // required for consistency, will force a space between itself and its neighbors in certain situations
		b.WriteString(innerHtml)
		// if the innerHtml does not already have a newline, add a newline
		if innerHtml[len(innerHtml)-1:] != "\n" {
			b.WriteString("\n")
		}
	}
	_,_ = fmt.Fprint(&b, "</", tag, ">\n")
	return b.String()
}


// RenderTagNoSpace is similar to RenderTag, but should be used in situations where the tag is an
// inline tag that you want to visually be right next to its neighbors with no space.
func RenderTagNoSpace(tag string, attr Attributes, innerHtml string) string {
	innerHtml = strings.TrimSpace(innerHtml)
	var attrString string

	if attr != nil {
		attrString = " " + attr.String()
	}
	ret := "<" + tag + attrString + ">"

	if innerHtml == "" || innerHtml[:1] != "<" {
		// either innerHtml is blank, or it is text and not a tag, so reproduce it verbatim
		ret += innerHtml + "</" + tag + ">"
	} else {
		if !config.Minify {
			innerHtml = Indent(innerHtml)
		}
		if innerHtml[len(innerHtml)-1:] != "\n" {
			innerHtml += "\n"
		}
		ret += "\n" + // innerhtml is a tag, and so spacing inside will not matter, so make it look good
			innerHtml +
			"</" + tag + ">"
	}
	return ret
}
/* TODO:
func renderTag(w io.Writer, tag string, attr Attributes, innerHtml io.WriterTo, minify bool, noSpace bool) {
	_,_ = fmt.Fprint(w, "<", tag)
	if attr != nil {
		_,_ = fmt.Fprint(w, " ", attr.String())
	}
	_,_ = fmt.Fprint(w, ">")

	if innerHtml != nil {
		if !minify {
			innerHtml = Indent(innerHtml)
		}
		b.WriteString("\n") // required for consistency, will force a space between itself and its neighbors in certain situations
		b.WriteString(innerHtml)
		// if the innerHtml does not already have a newline, add a newline
		if innerHtml[len(innerHtml)-1:] != "\n" {
			b.WriteString("\n")
		}
	}
	_,_ = fmt.Fprint(&b, "</", tag, ">\n")
	return b.String()


}
*/

// RenderLabel is a utility function to render a label, together with its text.
// Various CSS frameworks require labels to be rendered a certain way.
func RenderLabel(labelAttributes Attributes, label string, ctrlHtml string, mode LabelDrawingMode) string {
	tag := "label"
	label = html2.EscapeString(label)
	switch mode {
	case LabelBefore:
		return RenderTagNoSpace(tag, labelAttributes, label) + " " + ctrlHtml
	case LabelAfter:
		return ctrlHtml + " " + RenderTagNoSpace(tag, labelAttributes, label)
	case LabelWrapBefore:
		return RenderTag(tag, labelAttributes, label+" "+ctrlHtml)
	case LabelWrapAfter:
		return RenderTag(tag, labelAttributes, ctrlHtml+" "+label)
	}
	panic("Unknown label mode")
}

// RenderImage renders an image tag with the given source, alt and attribute values.
func RenderImage(src string, alt string, attributes Attributes) string {
	var a Attributes
	a = NewAttributes()
	a.Merge(attributes)

	a.Set("src", src)
	a.Set("alt", alt)
	return RenderVoidTag("img", a)
}

// Indent will add space to the front of every line in the string. Since indent is used to format code for reading
// while we are in development mode, we do not need it to be particularly efficient.
// It will not do this for textarea tags, since that would change the text in the tag.
func Indent(s string) string {
	var out string
	if config.Minify {
		return s
	}
	var taOffset int
	for {
		taOffset = strings.Index(s, "<textarea")
		if taOffset == -1 {
			out += indent(s)
			return out
		}
		out += indent(s[:taOffset])
		s = s[taOffset:]
		taOffset = strings.Index(s, "</textarea>")
		if taOffset == -1 {
			// This is an error in the html, so just return the original
			return s
		}
		out += s[:taOffset+11] // skip textarea close tag
		s = s[taOffset+11:]
	}
}

// indents the string unsafely, in that it does not check for allowable tags to indent
func indent(s string) string {
	if s == "" {
		return "" // don't indent empty strings
	}
	in := "  "
	r := strings.NewReplacer("\n", "\n"+in)
	s = r.Replace(s)
	return in + strings.TrimSuffix(s, in)
}

// Comment turns the given text into an HTML comment and returns the rendered comment
func Comment(s string) string {
	return fmt.Sprintf("<!-- %s -->", s)
}
