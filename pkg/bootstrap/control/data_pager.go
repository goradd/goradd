package control

import (
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
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

func NewDataPager(parent page.ControlI, id string, pagedControl control.PagedControlI) *DataPager {
	d := DataPager{}
	d.Init(&d, parent, id, pagedControl)
	return &d
}

func (d *DataPager) Init(self page.ControlI, parent page.ControlI, id string, pagedControl control.PagedControlI) {
	d.DataPager.Init(self, parent, id, pagedControl)
	d.SetLabels(`<span aria-hidden="true">&laquo;</span><span class="sr-only">Previous</span>`,
		`<span aria-hidden="true">&raquo;</span> <span class="sr-only">Next</span>`)
	d.ButtonStyle = ButtonStyleOutlineSecondary
	d.HighlightStyle = ButtonStylePrimary
	d.SetAttribute("aria-label", "Data pager")
}

func (d *DataPager) this() DataPagerI {
	return d.Self.(DataPagerI)
}

func (l *DataPager) ΩDrawingAttributes(ctx context.Context) html.Attributes {
	a := l.DataPager.ΩDrawingAttributes(ctx)
	a.AddClass("btn-group")
	return a
}

func (d *DataPager) PreviousButtonsHtml() string {
	var prev string
	var actionValue string

	pageNum := d.PagedControl().PageNum()
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
	pageNum := d.PagedControl().PageNum()

	actionValue = strconv.Itoa(pageNum + 1)

	attr := html.NewAttributes().
		Set("id", d.ID()+"_arrow_"+actionValue).
		SetClass("btn " + string(d.ButtonStyle))

	_, pageEnd := d.CalcBunch()
	pageCount := d.PagedControl().CalcPageCount()

	if pageNum >= pageCount {
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
	pageNum := d.PagedControl().PageNum()

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


func (d *DataPager) Deserialize(dec page.Decoder) (err error) {
	if err = d.DataPager.Deserialize(dec); err != nil {
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

// DataPagerCreator is the initialization structure for declarative creation of data pagers
type DataPagerCreator struct {
	// ID is the control id
	ID string
	// MaxPageButtons is the maximum number of page buttons to display in the pager
	MaxPageButtons int
	// ObjectName is the name of the object being displayed in the table
	ObjectName string
	// ObjectPluralName is the plural name of the object being displayed
	ObjectPluralName string
	// LabelForNext is the text to use in the Next button
	LabelForNext string
	// LabelForPrevious is the text to use in the Previous button
	LabelForPrevious string
	// PagedControl is the id of the control that will be paged by the pager
	PagedControl string
	page.ControlOptions
	// ButtonStyle is the style that will be used to draw the standard buttons
	ButtonStyle    ButtonStyle
	// HighlightStyle is the style that will be used to draw the highlighted buttons
	HighlightStyle ButtonStyle

}


// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c DataPagerCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	if !parent.Page().HasControl(c.PagedControl) {
		panic ("you must declare the paged control before the data pager")
	}
	p := parent.Page().GetControl(c.PagedControl).(control.PagedControlI)
	ctrl := NewDataPager(parent, c.ID, p)
	c.Init(ctx, ctrl)
	return ctrl
}

// Init is called by implementations of Buttons to initialize a control with the
// creator. You do not normally need to call this.
func (c DataPagerCreator) Init(ctx context.Context, ctrl DataPagerI) {
	ctrl.(*DataPager).ButtonStyle = c.ButtonStyle
	ctrl.(*DataPager).HighlightStyle = c.HighlightStyle

	sub := control.DataPagerCreator{
		MaxPageButtons: c.MaxPageButtons,
		ObjectName: c.ObjectName,
		ObjectPluralName: c.ObjectPluralName,
		LabelForNext: c.LabelForNext,
		LabelForPrevious: c.LabelForPrevious,
		PagedControl:  c.PagedControl,
		ControlOptions: c.ControlOptions,
	}
	sub.Init(ctx, ctrl)
}

func init() {
	page.RegisterControl(DataPager{})
}