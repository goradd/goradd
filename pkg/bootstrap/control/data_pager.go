package control

import (
	"fmt"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
	"reflect"
	"strconv"
)

type DataPagerI interface {
	control.DataPagerI
}

// DataPager is a toolbar designed to aid scrolling through a large set of data. It is implemented using Aria design
// best practices. It is designed to be paired with a Table or DataRepeater to aid in navigating through the data.
// It is similar to a Paginator, but a paginator is for navigating through a whole series of pages and not just for
// data on one override.
type DataPager struct {
	control.DataPager
	ButtonStyle    ButtonStyle
	HighlightStyle ButtonStyle
}

func NewDataPager(parent page.ControlI, id string, paginatedControl control.PaginatedControlI) *DataPager {
	d := DataPager{}
	d.Init(&d, parent, id, paginatedControl)
	return &d
}

func (d *DataPager) Init(self page.ControlI, parent page.ControlI, id string, paginatedControl control.PaginatedControlI) {
	d.DataPager.Init(self, parent, id, paginatedControl)
	d.SetLabels(`<span aria-hidden="true">&laquo;</span><span class="sr-only">Previous</span>`,
		`<span aria-hidden="true">&raquo;</span> <span class="sr-only">Next</span>`)
	d.ButtonStyle = ButtonStyleOutlineSecondary
	d.HighlightStyle = ButtonStylePrimary
	d.SetAttribute("aria-label", "Data pager")
}

func (d *DataPager) this() DataPagerI {
	return d.Self.(DataPagerI)
}

func (l *DataPager) ΩDrawingAttributes() *html.Attributes {
	a := l.DataPager.ΩDrawingAttributes()
	a.AddClass("btn-group")
	return a
}

func (d *DataPager) PreviousButtonsHtml() string {
	var prev string
	var actionValue string

	pageNum := d.PaginatedControl().PageNum()
	actionValue = strconv.Itoa(pageNum - 1)

	attr := html.NewAttributes().
		Set("id", d.ID()+"_arrow_"+actionValue).
		SetClass("btn " + string(d.ButtonStyle))

	if pageNum <= 1 {
		attr.SetDisabled(true)
		attr.SetStyle("cursor", "not-allowed")
	}

	prev = d.ButtonProxy().ButtonHtml(d.LabelForPrevious, actionValue, attr, true)

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
	pageNum := d.PaginatedControl().PageNum()

	attr := html.NewAttributes().
		Set("id", d.ID()+"_arrow_"+actionValue).
		SetClass("btn " + string(d.ButtonStyle))

	actionValue = strconv.Itoa(pageNum + 1)

	_, pageEnd := d.CalcBunch()
	pageCount := d.PaginatedControl().CalcPageCount()

	if pageNum >= pageCount-1 {
		attr.SetDisabled(true)
		attr.SetStyle("cursor", "not-allowed")
	}

	next = d.ButtonProxy().ButtonHtml(d.LabelForNext, actionValue, attr, true)

	h := next
	if pageEnd != pageCount {
		h += d.PageButtonsHtml(pageCount) + h
		h = fmt.Sprintf(`<button disabled class="btn %s" style="cursor: not-allowed">&hellip;</button>`, d.ButtonStyle) + h
	}
	return h
}

func (d *DataPager) PageButtonsHtml(i int) string {
	pageNum := d.PaginatedControl().PageNum()

	actionValue := strconv.Itoa(i)
	attr := html.NewAttributes().Set("id", d.ID()+"_page_"+actionValue).
		Set("role", "tab").
		AddClass("btn")
	if pageNum == i {
		attr.AddClass(string(d.HighlightStyle))
		attr.Set("aria-selected", "true")
		attr.Set("tabindex", "0")
	} else {
		attr.AddClass(string(d.ButtonStyle))
		attr.Set("aria-selected", "false")
		attr.Set("tabindex", "-1")
		// TODO: We need javascript to respond to arrow keys to set the focus on the buttons. User could then press space to click on button.
	}
	return d.ButtonProxy().ButtonHtml(actionValue, actionValue, attr, false)
}

func (d *DataPager) Serialize(e page.Encoder) (err error) {
	if err = d.DataPager.Serialize(e); err != nil {
		return
	}

	if err = e.Encode(d.ButtonStyle); err != nil {
		return
	}

	if err = e.Encode(d.HighlightStyle); err != nil {
		return
	}

	return
}

// ΩisSerializer is used by the automated control serializer to determine how far down the control chain the control
// has to go before just calling serialize and deserialize
func (d *DataPager) ΩisSerializer(i page.ControlI) bool {
	return reflect.TypeOf(d) == reflect.TypeOf(i)
}

func (d *DataPager) Deserialize(dec page.Decoder, p *page.Page) (err error) {
	if err = d.DataPager.Deserialize(dec, p); err != nil {
		return
	}

	if err = dec.Decode(&d.ButtonStyle); err != nil {
		return
	}

	if err = dec.Decode(&d.HighlightStyle); err != nil {
		return
	}

	return
}
