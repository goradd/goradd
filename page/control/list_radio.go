package control

import (
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/html"
)


// RadioList is a multi-select control that presents its choices as a list of checkboxes.
// Styling is provided by divs and spans that you can provide css for in your style sheets. The
// goradd.css file has default styling to handle the basics. It wraps the whole thing in a div that can be set
// to scroll as well, so that the final structure can be styled like a multi-column table, or a single-column
// scrolling list much like a standard html select list.
type RadioList struct {
	CheckboxList
}

func NewRadioList(parent page.ControlI) *RadioList {
	l := &RadioList{}
	l.Init(l, parent)
	return l
}


func (l *RadioList) Init(self page.ControlI, parent page.ControlI) {
	l.CheckboxList.Init(self, parent)
}

// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *RadioList) DrawingAttributes() *html.Attributes {
	a := l.CheckboxList.DrawingAttributes()
	a.SetDataAttribute("grctl", "radiolist")
	return a
}

func (l *RadioList) renderItem(tag string, item ListItemI) (h string) {
	attributes := html.NewAttributes()
	attributes.SetId(item.Id())
	attributes.Set("name", l.Id())
	attributes.Set("value", item.Id())
	attributes.Set("type", "radio")
	if l.selectedIds[item.Id()] {
		attributes.Set("checked", "")
	}
	ctrl := html.RenderVoidTag("input", attributes)
	h = html.RenderLabel(html.NewAttributes().Set("for", item.Id()), item.Label(), ctrl, l.labelDrawingMode)
	attributes = item.Attributes().Clone()
	attributes.AddClass("gr-cbl-item")
	h = html.RenderTag(tag, attributes, h)
	return
}

func (l *RadioList) UpdateFormValues(ctx *page.Context) {
	controlId := l.Id()

	if ctx.RequestMode() == page.Ajax {
		if v,ok := ctx.CheckableValue(controlId); ok {
			if s, ok := v.(string); ok {
				l.selectedIds = map[string]bool{l.Id() + "_" + s:true}
			}
		}
	} else {
		if v,ok := ctx.FormValue(controlId); ok {
			l.selectedIds = map[string]bool{v:true}
		}
	}
}
