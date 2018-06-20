package html

import (
	"fmt"
	"goradd/config"
	html2 "html"
	"strings"
)

type LabelDrawingMode int

// The label drawing mode describes how to draw a label when it is drawn. Various CSS frameworks expect it a certain way. Many
// are not very forgiving when you don't do it the way they expect.
const (
	LABEL_DEFAULT     LabelDrawingMode = iota // Label mode is defined elsewhere, like in a config setting
	LABEL_BEFORE                              // Label tag is before the control's tag, and terminates before the control
	LABEL_AFTER                               // Label tag is after the control's tag, and start after the control
	WRAP_LABEL_BEFORE                         // Label tag is before the control's tag, and wraps the control tag
	WRAP_LABEL_AFTER                          // Label tag is after the control's tag, and wraps the control tag
)

type VoidTag struct {
	Tag  string
	Attr *Attributes
}

func (t VoidTag) Render() string {
	return RenderVoidTag(t.Tag, t.Attr)
}

// Render a void tag that has no closing tag
func RenderVoidTag(tag string, attr *Attributes) (s string) {
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

// RenderTag renders a standard html tag with a closing tag
// innerHtml is html, and must already be escaped if needed
// The tag will be surrounded with newlines to force general formatting consistency. This will cause the tag to be rendered
// with a space between it and its neighbors if the tag is not a block tag. In the few situations where you would want to
// get rid of this space, call RenderTagNoSpace()
func RenderTag(tag string, attr *Attributes, innerHtml string) string {
	var attrString string

	if attr != nil {
		attrString = " " + attr.String()
	}
	ret := "<" + tag + attrString + ">"

	if innerHtml == "" {
		ret += "</" + tag + ">"
	} else {
		if innerHtml[len(innerHtml)-1:] != "\n" {
			innerHtml += "\n"
		}
		if !config.Minify {
			innerHtml = Indent(innerHtml)
		}

		ret += "\n" + // required here for consistency, will force a space between itself and its neighbors in certain situations
			innerHtml +
			"</" + tag + ">\n"
	}
	return ret
}

// RenderTagNoSpace is similar to RenderTag, but should be used in situations where the tag is an inline tag that you want
// to visually be right next to its neighbors with no space.
func RenderTagNoSpace(tag string, attr *Attributes, innerHtml string) string {
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
		ret += "\n" + // innerhtml is a tag, and so spacing will not matter, so make it look good
			innerHtml +
			"</" + tag + ">\n"
	}
	return ret
}

// A utility function to render a label, together with its text. Various CSS frameworks require labels to be rendered
// a certain way.
func RenderLabel(labelAttributes *Attributes, label string, ctrlHtml string, mode LabelDrawingMode) string {
	tag := "label"
	label = html2.EscapeString(label)
	switch mode {
	case LABEL_BEFORE:
		return RenderTagNoSpace(tag, labelAttributes, label) + " " + ctrlHtml
	case LABEL_AFTER:
		return ctrlHtml + " " + RenderTagNoSpace(tag, labelAttributes, label)
	case WRAP_LABEL_BEFORE:
		return RenderTag(tag, labelAttributes, label + " " + ctrlHtml)
	case WRAP_LABEL_AFTER:
		return RenderTag(tag, labelAttributes, ctrlHtml + " " + label)
	}
	panic("Unknown label mode")
}

func RenderImage(src string, alt string, attributes *Attributes) string {
	var a *Attributes

	if attributes != nil {
		a = attributes.Clone()
	} else {
		a = NewAttributes()
	}
	a.Set("src", src)
	a.Set("alt", alt)
	return RenderVoidTag("img", a)
}

// Indent will add space to the front of every line in the string. Since indent is used to format code for reading
// while we are in development mode, we do not need it to be particularly efficient.
func Indent(s string) string {
	if config.Minify {
		return s
	}

	in := "  "
	r := strings.NewReplacer("\n", "\n"+in)
	s = r.Replace(s)
	return in + strings.TrimSuffix(s, in)
}

// Comment turns the given text into an html comment and returns the comment
func Comment(s string) string {
	return fmt.Sprintf("<!-- %s -->", s)
}
