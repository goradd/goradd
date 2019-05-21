package control

import (
	"bytes"
	"context"
	"github.com/goradd/gengen/pkg/maps"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/math"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/control/data"
	"github.com/goradd/goradd/pkg/page/event"
	"reflect"
	"strconv"
)

const (
	PageClick = iota + 1000
)

// PaginatedControlI is the interface that paginated controls must implement
type PaginatedControlI interface {
	data.DataManagerI
	SetTotalItems(uint)
	TotalItems() int
	SetPageSize(size int)
	PageSize() int
	PageNum() int
	SetPageNum(n int)
	AddDataPager(DataPagerI)
	CalcPageCount() int
	HasDataPagers() bool
	getDataPagers() []DataPagerI
}

// PaginatedControl is a mixin that makes a Control controllable by a data pager
type PaginatedControl struct {
	totalItems       int
	pageSize         int
	pageNum          int
	dataPagers 		 []DataPagerI
}


// DefaultPaginatorPageSize is the default number of items that a paginated control will show. You can change this in an individual control, too.
var DefaultPaginatorPageSize = 10

// DefaultMaxPagintorButtons is the default maximum number of buttons to display on the pager. You can change this in an individual control, too.
var DefaultMaxPagintorButtons = 10

// SetTotalItems sets the total number of items that the paginator keeps track of. This will be divided by
// the PageSize to determine the number of pages presented.
// You must call this each time the data size might change.
func (c *PaginatedControl) SetTotalItems(count uint) {
	c.totalItems = int(count)
	c.limitPageNumber()
}

// TotalItems returns the number of items that the paginator is aware of in the list it is managing.
func (c *PaginatedControl) TotalItems() int {
	return c.totalItems
}

// SetPageSize sets the maximum number of items that will be displayed at one time. If more than this number of items
// is being displayed, the pager will allow paging to other items.
func (c *PaginatedControl) SetPageSize(size int) {
	if size == 0 {
		size = DefaultPaginatorPageSize
	}
	c.pageSize = size
}

// PageSize returns the maximum number of items that will be allowed in a page.
func (c *PaginatedControl) PageSize() int {
	return c.pageSize
}

// PageNum returns the current page number.
func (c *PaginatedControl) PageNum() int {
	return c.pageNum
}

// SetPageNum sets the current page number. It does not redraw anything, nor does it determine if the
// page is actually visible.
func (c *PaginatedControl) SetPageNum(n int) {
	if c.pageNum != n {
		c.pageNum = n
	}
}

// AddDataPager adds a data pager to the PaginatedControl. A PaginatedControl can have multiple
// data pagers.
func (c *PaginatedControl) AddDataPager(d DataPagerI) {
	c.dataPagers = append(c.dataPagers, d)
}

func (c *PaginatedControl) limitPageNumber() {
	pageCount := c.CalcPageCount()

	if c.pageNum > pageCount {
		if pageCount <= 1 {
			c.pageNum = 1
		} else {
			c.pageNum = pageCount
		}
	}
}

// CalcPageCount will return the number of pages based on the page size and total items.
func (c *PaginatedControl) CalcPageCount() int {
	if c.pageSize == 0 || c.totalItems == 0 {
		return 0
	}
	return (c.totalItems-1)/c.pageSize + 1
}

func (c *PaginatedControl) getDataPagers() []DataPagerI{
	return c.dataPagers
}

func (c *PaginatedControl) HasDataPagers() bool {
	return len(c.dataPagers) > 0
}

// SliceOffsets returns the start and end values to use to specify a portion of a slice corresponding to the
// data the pager refers to
func (c *PaginatedControl) SliceOffsets() (start, end int) {
	start = (c.PageNum() - 1) * c.PageSize()
	_,end = math.MinInt(start+ c.PageSize(),  c.TotalItems())
	return
}

// SqlLimits returns the limits you would use in a sql database limit clause
func (c *PaginatedControl) SqlLimits() (maxRowCount, offset int) {
	offset = (c.PageNum() - 1) * c.PageSize()
	maxRowCount = c.PageSize()
	return
}


// DataPagerI is the data pager interface that allows this object to call into subclasses.
type DataPagerI interface {
	page.ControlI
	PreviousButtonsHtml() string
	NextButtonsHtml() string
	PageButtonsHtml(i int) string
}

// DataPager is a toolbar designed to aid scrolling a large set of data. It is implemented using Aria design
// best practices. It is designed to be paired with a Table or DataRepeater to aid in navigating through the data.
// It is similar to a Paginator, but a paginator is for navigating through a series of related web pages and not just for
// data on one form.
type DataPager struct {
	page.Control

	maxPageButtons   int
	ObjectName       string
	ObjectPluralName string
	LabelForNext     string
	LabelForPrevious string

	paginatedControl PaginatedControlI
	Proxy *Proxy
}

// NewDataPager creates a new DataPager
func NewDataPager(parent page.ControlI, id string, paginatedControl PaginatedControlI) *DataPager {
	d := DataPager{}
	d.Init(&d, parent, id, paginatedControl)
	return &d
}

// Init is called by subclasses of a DataPager to initialize the data pager. You do not normally need
// to call this.
func (d *DataPager) Init(self page.ControlI, parent page.ControlI, id string, paginatedControl PaginatedControlI) {
	d.Control.Init(self, parent, id)
	d.Tag = "div"
	d.LabelForNext = d.ΩT("Next")
	d.LabelForPrevious = d.ΩT("Previous")
	d.maxPageButtons = DefaultMaxPagintorButtons
	paginatedControl.AddDataPager(self.(DataPagerI))
	d.paginatedControl = paginatedControl
	d.Proxy = NewProxy(d)
	d.Proxy.On(event.Click(), action.Ajax(d.ID(), PageClick))
	d.SetAttribute("role", "tablist")
	d.paginatedControl.SetPageNum(1)
}

// ΩDrawingAttributes is called by the framework to add temporary attributes to the html.
func (d *DataPager) ΩDrawingAttributes() *html.Attributes {
	a := d.Control.ΩDrawingAttributes()
	a.SetDataAttribute("grctl", "datapager")
	return a
}

// Action is called by the framework to respond to actions.
func (d *DataPager) Action(ctx context.Context, params page.ActionParams) {
	switch params.ID {
	case PageClick:
		pageNum := params.ControlValueInt();
		if pageNum < 1 {
			pageNum = 1
		} else {
			c := d.paginatedControl.CalcPageCount()
			if pageNum > c {
				pageNum = c
			}
		}
		d.paginatedControl.SetPageNum(pageNum)
		d.paginatedControl.Refresh()
		for _,c := range d.paginatedControl.getDataPagers() {
			c.Refresh()
		}
	}
}


func (d *DataPager) refreshPaginatedControl() {
	d.paginatedControl.Refresh()
}

// SetMaxPageButtons sets the maximum number of buttons that will be displayed in the paginator.
func (d *DataPager) SetMaxPageButtons(b int) {
	d.maxPageButtons = b
}

// SetObjectNames sets the single and plural names of the objects that are represented
// in the data pager.
func (d *DataPager) SetObjectNames(singular string, plural string) {
	d.ObjectName = singular
	d.ObjectPluralName = plural
}


// SetLabels sets the previous and next labels. Translate these first.
func (d *DataPager) SetLabels(previous string, next string) {
	d.LabelForPrevious = previous
	d.LabelForNext = next
}



/*
CalcBunch is called by the framework to lay out the data pager based on the number of pages
in the pager. It should try to represent an easy to navigate interface that can manage 2 or
2000 pages.

A "Bunch" is defined as the collection of numbers that lies in between the pair of Ellipsis ("...")

Layout

For an IndexCount of 10
2   213   2 (two items to the left of the bunch, and then 2 indexes, selected index, 3 indexes, and then two items to the right of the bunch)
e.g. 1 ... 5 6 *7* 8 9 10 ... 100

For IndexCount of 11
2   313   2

For IndexCount of 12
2   314   2

For IndexCount of 13
2   414   2

For IndexCount of 14
2   415   2

Start/end page numbers for the bunch

For IndexCount of 10
1 2 3 4 5 6 7 8 .. 100
1 .. 4 5 *6* 7 8 9 .. 100
1 .. 92 93 *94* 95 96 97 .. 100
1 .. 93 94 95 96 97 98 99 100

For IndexCount of 11
1 2 3 4 5 6 7 8 9 .. 100
1 .. 4 5 6 *7* 8 9 10 .. 100
1 .. 91 92 93 *94* 95 96 97 .. 100
1 .. 92 93 94 95 96 97 98 99 100

For IndexCount of 12
1 2 3 4 5 6 7 8 9 10 .. 100
1 .. 4 5 6 *7* 8 9 10 11 .. 100
1 .. 90 91 92 *93* 94 95 96 97 .. 100
1 .. 91 92 93 94 95 96 97 98 99 100

For IndexCount of 13
1 2 3 4 5 6 7 8 9 11 .. 100
1 .. 4 5 6 7 *8* 9 10 11 12 .. 100
1 .. 89 90 91 92 *93* 94 95 96 97 .. 100
1 .. 90 91 92 93 94 95 96 97 98 99 100

Note: there are likely better ways to do this. Some innovative ones are to have groups of 10s, and then 100s etc.
Or, use the ellipsis as a dropdown menu for more selections
*/
func (d *DataPager) CalcBunch() (pageStart, pageEnd int) {

	pageCount := d.paginatedControl.CalcPageCount()
	pageNum := d.paginatedControl.PageNum()

	if pageCount <= d.maxPageButtons {
		return 1, pageCount
	} else {
		_,minEndOfBunch := math.MinInt(d.maxPageButtons-2, pageCount)
		_,maxStartOfBunch := math.MaxInt(pageCount-d.maxPageButtons+3, 1)

		leftOfBunchCount := (d.maxPageButtons - 5) / 2
		rightOfBunchCount := (d.maxPageButtons - 4) / 2

		leftBunchTrigger := leftOfBunchCount + 4
		rightBunchTrigger := maxStartOfBunch + (d.maxPageButtons-7)/2

		if pageNum < leftBunchTrigger {
			pageStart = 1
		} else {
			_,pageStart = math.MinInt(maxStartOfBunch, pageNum-leftOfBunchCount)
		}

		if pageNum > rightBunchTrigger {
			pageEnd = pageCount
		} else {
			_,pageEnd = math.MaxInt(minEndOfBunch, pageNum+rightOfBunchCount)
		}
		return
	}
}

// ΩPreRender is called by the framework to load data into the paginated control just before drawing.
func (d *DataPager) ΩPreRender(ctx context.Context, buf *bytes.Buffer) (err error) {
	err = d.Control.ΩPreRender(ctx, buf)

	if err == nil {
		// If we are being drawn before the paginated control, we must tell the paginated control to load up its
		// data so that we can figure out what to do
		if !d.paginatedControl.WasRendered() &&
			!d.paginatedControl.IsRendering() { // not a child control
			d.paginatedControl.GetData(ctx, d.paginatedControl)
		}
	}
	return
}

// ΩDrawInnerHtml is called by the framework to draw the control's inner html.
func (d *DataPager) ΩDrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	h := d.Self.(DataPagerI).PreviousButtonsHtml()
	pageStart, pageEnd := d.CalcBunch()
	for i := pageStart; i <= pageEnd; i++ {
		h += d.Self.(DataPagerI).PageButtonsHtml(i)
	}

	h += d.Self.(DataPagerI).NextButtonsHtml()
	_, err = buf.WriteString(h)
	return
}

// PreviousButtonsHtml returns the html to draw the previous buttons. Subclasses can override this to
// change how the Previous buttons are drawn.
func (d *DataPager) PreviousButtonsHtml() string {
	var prev string
	var actionValue string

	pageNum := d.paginatedControl.PageNum()

	actionValue = strconv.Itoa(pageNum - 1)

	attr := html.NewAttributes().
		Set("id", d.ID()+"_arrow_"+actionValue).
		SetClass("arrow previous")

	if pageNum <= 1 {
		attr.SetDisabled(true)
		attr.SetStyle("cursor", "not-allowed")
	}

	prev = d.Proxy.ButtonHtml(d.LabelForPrevious, actionValue, attr, false)

	h := prev
	pageStart, _ := d.CalcBunch()
	if pageStart != 1 {
		h += d.Self.(DataPagerI).PageButtonsHtml(1)
		h += `<span class="ellipsis">&hellip;</span>`
	}
	return h
}

// NextButtonsHtml returns the html for the next buttons. Subclasses can override this to change how the
// next buttons are drawn.
func (d *DataPager) NextButtonsHtml() string {
	var next string
	var actionValue string

	pageNum := d.paginatedControl.PageNum()

	attr := html.NewAttributes().
		Set("id", d.ID()+"_arrow_"+actionValue).
		SetClass("arrow next")

	actionValue = strconv.Itoa(pageNum + 1)

	_, pageEnd := d.CalcBunch()
	pageCount := d.paginatedControl.CalcPageCount()

	if pageNum >= pageCount {
		attr.SetDisabled(true)
		attr.SetStyle("cursor", "not-allowed")
	}

	next = d.Proxy.ButtonHtml(d.LabelForNext, actionValue, attr, false)

	h := next
	if pageEnd != pageCount {
		h += d.Self.(DataPagerI).PageButtonsHtml(pageCount) + h
		h = `<span class="ellipsis">&hellip;</span>` + h
	}
	return h
}

// PageButtonsHtml returns the html for the page buttons. Subclasses can override this to change how
// the page buttons are drawn.
func (d *DataPager) PageButtonsHtml(i int) string {
	actionValue := strconv.Itoa(i)
	buttonId := d.ID() + "_page_" + actionValue
	attr := html.NewAttributes().Set("id", buttonId).Set("role", "tab").AddClass("page")
	pageNum := d.paginatedControl.PageNum()

	if pageNum == i {
		attr.AddClass("selected")
		attr.Set("aria-selected", "true")
		attr.Set("tabindex", "0")
	} else {
		attr.Set("aria-selected", "false")
		attr.Set("tabindex", "-1")
		// TODO: We need javascript to respond to arrow keys to set the focus on the buttons. User could then press space to click on button.
	}
	return d.Proxy.ButtonHtml(actionValue, actionValue, attr, false)
}

// ΩMarshalState is an internal function to save the state of the control
func (d *DataPager) ΩMarshalState(m maps.Setter) {
	m.Set("pageNum", d.paginatedControl.PageNum())
}

// ΩUnmarshalState is an internal function to restore the state of the control
func (d *DataPager) ΩUnmarshalState(m maps.Loader) {
	if v,ok := m.Load("pageNum"); ok {
		if i, ok := v.(int); ok {
			d.paginatedControl.SetPageNum (i) // admittedly, multiple pagers will repeat the same call, but not likely to effect performance
		}
	}
}

func (d *DataPager) PaginatedControl() PaginatedControlI {
	return d.paginatedControl
}

func (d *DataPager) Serialize(e page.Encoder) (err error) {
	if err = d.Control.Serialize(e); err != nil {
		return
	}

	if err = e.Encode(d.maxPageButtons); err != nil {
		return
	}

	if err = e.Encode(d.ObjectName); err != nil {
		return
	}

	if err = e.Encode(d.ObjectPluralName); err != nil {
		return
	}

	if err = e.Encode(d.LabelForNext); err != nil {
		return
	}

	if err = e.Encode(d.LabelForPrevious); err != nil {
		return
	}

	if err = e.EncodeControl(d.paginatedControl); err != nil {
		return
	}

	if err = e.EncodeControl(d.Proxy); err != nil {
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
	if err = d.Control.Deserialize(dec, p); err != nil {
		return
	}

	if err = dec.Decode(&d.maxPageButtons); err != nil {
		return
	}

	if err = dec.Decode(&d.ObjectName); err != nil {
		return
	}

	if err = dec.Decode(&d.ObjectPluralName); err != nil {
		return
	}

	if err = dec.Decode(&d.LabelForNext); err != nil {
		return
	}

	if ci,err := dec.DecodeControl(p); err != nil {
		return err
	} else {
		d.paginatedControl = ci.(PaginatedControlI)
	}

	if ci,err := dec.DecodeControl(p); err != nil {
		return err
	} else {
		d.Proxy = ci.(*Proxy)
	}
	
	return
}
