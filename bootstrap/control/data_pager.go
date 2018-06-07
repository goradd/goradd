package control

import (
	"fmt"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/page/control"
	"strconv"
)

// DataPager is a toolbar designed to aid scrolling through a large set of data. It is implemented using Aria design
// best practices. It is designed to be paired with a Table or DataRepeater to aid in navigating through the data.
// It is similar to a Paginator, but a paginator is for navigating through a whole series of pages and not just for
// data on one page.
type DataPager struct {
	control.DataPager
	ButtonStyle    ButtonStyle
	HighlightStyle ButtonStyle
}

func NewDataPager(parent page.ControlI, paginatedControl page.ControlI) *DataPager {
	d := DataPager{}
	d.Init(&d, parent, paginatedControl)
	d.Tag = "div"
	d.SetLabels(`<span aria-hidden="true">&laquo;</span><span class="sr-only">Previous</span>`,
		`<span aria-hidden="true">&raquo;</span> <span class="sr-only">Next</span>`)
	d.SetMaxPageButtons(10)
	d.ButtonStyle = ButtonStyleOutlineSecondary
	d.HighlightStyle = ButtonStyleSecondary
	return &d
}

func (d *DataPager) Init(self page.ControlI, parent page.ControlI, paginatedControl page.ControlI) {
	d.DataPager.Init(self, parent, paginatedControl)
	d.SetAttribute("aria-label", "Data pager")
}

func (l *DataPager) DrawingAttributes() *html.Attributes {
	a := l.DataPager.DrawingAttributes()
	a.AddClass("btn-group")
	return a
}

func (d *DataPager) PreviousButtonsHtml() string {
	var prev string
	var actionValue string
	actionValue = strconv.Itoa(d.PageNum() - 1)

	attr := html.NewAttributes().
		Set("id", d.ID()+"_arrow_"+actionValue).
		SetClass("btn " + string(d.ButtonStyle))

	if d.PageNum() <= 1 {
		attr.SetDisabled(true)
		attr.SetStyle("cursor", "not-allowed")
	}

	prev = d.Proxy.ButtonHtml(d.LabelForPrevious, actionValue, attr, true)

	h := prev
	pageStart, _ := d.CalcBunch()
	if pageStart != 1 {
		h += d.PageButtonsHtml(1)
		h += fmt.Sprintf(`<button disabled class="btn %s" style="cursor: not-allowed">&hellip;</button>`, d.ButtonStyle)
	}
	return h
}

func (d *DataPager) NextButtonsHtml() string {
	var next string
	var actionValue string

	attr := html.NewAttributes().
		Set("id", d.ID()+"_arrow_"+actionValue).
		SetClass("btn " + string(d.ButtonStyle))

	actionValue = strconv.Itoa(d.PageNum() + 1)

	_, pageEnd := d.CalcBunch()
	pageCount := d.CalcPageCount()

	if d.PageNum() >= pageCount-1 {
		attr.SetDisabled(true)
		attr.SetStyle("cursor", "not-allowed")
	}

	next = d.Proxy.ButtonHtml(d.LabelForNext, actionValue, attr, true)

	h := next
	if pageEnd != pageCount {
		h += d.PageButtonsHtml(pageCount) + h
		h = fmt.Sprintf(`<button disabled class="btn %s" style="cursor: not-allowed">&hellip;</button>`, d.ButtonStyle) + h
	}
	return h
}

func (d *DataPager) PageButtonsHtml(i int) string {
	actionValue := strconv.Itoa(i)
	attr := html.NewAttributes().Set("id", d.ID()+"_page_"+actionValue).
		Set("role", "tab").
		AddClass("btn")
	if d.PageNum() == i {
		attr.AddClass(string(d.HighlightStyle))
		attr.Set("aria-selected", "true")
		attr.Set("tabindex", "0")
	} else {
		attr.AddClass(string(d.ButtonStyle))
		attr.Set("aria-selected", "false")
		attr.Set("tabindex", "-1")
		// TODO: We need javascript to respond to arrow keys to set the focus on the buttons. User could then press space to click on button.
	}
	return d.Proxy.ButtonHtml(actionValue, actionValue, attr, false)
}
