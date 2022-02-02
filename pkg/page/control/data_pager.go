package control

import (
	"context"
	"github.com/goradd/gengen/pkg/maps"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/math"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/event"
	"io"
	"strconv"
)

const (
	PageClick = iota + 1000
)

// PagedControlI is the interface that paged controls must implement
type PagedControlI interface {
	DataManagerI
	SetTotalItems(uint)
	TotalItems() int
	SetPageSize(size int)
	PageSize() int
	PageNum() int
	SetPageNum(n int)
	AddDataPager(DataPagerI)
	CalcPageCount() int
	HasDataPagers() bool
	GetDataPagerIDs() []string
	SliceOffsets() (start, end int)
}

// PagedControl is a mixin that makes a ControlBase controllable by a data pager. All embedders of a
// PagedControl MUST implement the Serialize and Deserialize methods so that the base ControlBase versions
// of these functions will get called.
type PagedControl struct {
	totalItems int
	pageSize   int
	pageNum    int
	dataPagers []string
}

// DefaultPagerPageSize is the default number of items that a paged control will show. You can change this in an individual control, too.
var DefaultPagerPageSize = 10

// DefaultMaxPagerButtons is the default maximum number of buttons to display on the pager. You can change this in an individual control, too.
var DefaultMaxPagerButtons = 10

// SetTotalItems sets the total number of items that the paginator keeps track of. This will be divided by
// the PageSize to determine the number of pages presented.
// You must call this each time the data size might change.
func (c *PagedControl) SetTotalItems(count uint) {
	c.totalItems = int(count)
	c.limitPageNumber()
}

// TotalItems returns the number of items that the paginator is aware of in the list it is managing.
func (c *PagedControl) TotalItems() int {
	return c.totalItems
}

// SetPageSize sets the maximum number of items that will be displayed at one time. If more than this number of items
// is being displayed, the pager will allow paging to other items.
func (c *PagedControl) SetPageSize(size int) {
	if size == 0 {
		size = DefaultPagerPageSize
	}
	c.pageSize = size
}

// PageSize returns the maximum number of items that will be allowed in a page.
func (c *PagedControl) PageSize() int {
	return c.pageSize
}

// PageNum returns the current page number.
func (c *PagedControl) PageNum() int {
	return c.pageNum
}

// SetPageNum sets the current page number. It does not redraw anything, nor does it determine if the
// page is actually visible.
func (c *PagedControl) SetPageNum(n int) {
	if c.pageNum != n {
		c.pageNum = n
	}
}

func (c *PagedControl) GetDataPagerIDs() []string {
	return c.dataPagers
}

// AddDataPager adds a data pager to the PagedControl. A PagedControl can have multiple
// data pagers.
func (c *PagedControl) AddDataPager(d DataPagerI) {
	c.dataPagers = append(c.dataPagers, d.ID())
}

func (c *PagedControl) limitPageNumber() {
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
func (c *PagedControl) CalcPageCount() int {
	if c.pageSize == 0 || c.totalItems == 0 {
		return 0
	}
	return (c.totalItems-1)/c.pageSize + 1
}

func (c *PagedControl) HasDataPagers() bool {
	return len(c.dataPagers) > 0
}

// SliceOffsets returns the start and end values to use to specify a portion of a slice corresponding to the
// data the pager refers to
func (c *PagedControl) SliceOffsets() (start, end int) {
	start = (c.PageNum() - 1) * c.PageSize()
	_, end = math.MinInt(start+c.PageSize(), c.TotalItems())
	return
}

// SqlLimits returns the limits you would use in a sql database limit clause
func (c *PagedControl) SqlLimits() (maxRowCount, offset int) {
	offset = (c.PageNum() - 1) * c.PageSize()
	maxRowCount = c.PageSize()
	return
}

// MarshalState is an internal function to save the state of the control
func (c *PagedControl) MarshalState(m maps.Setter) {
	m.Set("pn", c.pageNum)
}

// UnmarshalState is an internal function to restore the state of the control
func (c *PagedControl) UnmarshalState(m maps.Loader) {
	if v, ok := m.Load("pn"); ok {
		if pn, ok := v.(int); ok {
			c.pageNum = pn
		}
	}
}


// Serialize encodes the PagedControl data for serialization. Note that all control implementations
// that use a PagedControl MUST create their own Serialize method, call the base ControlBase's version first,
// and then call this Serialize method.
func (c *PagedControl) Serialize(e page.Encoder) (err error) {
	if err = e.Encode(c.totalItems); err != nil {
		return
	}
	if err = e.Encode(c.pageSize); err != nil {
		return
	}
	if err = e.Encode(c.pageNum); err != nil {
		return
	}
	if err = e.Encode(c.dataPagers); err != nil {
		return
	}

	return
}

func (c *PagedControl) Deserialize(dec page.Decoder) (err error) {
	if err = dec.Decode(&c.totalItems); err != nil {
		panic(err)
	}

	if err = dec.Decode(&c.pageSize); err != nil {
		panic(err)
	}

	if err = dec.Decode(&c.pageNum); err != nil {
		panic(err)
	}

	if err = dec.Decode(&c.dataPagers); err != nil {
		panic(err)
	}

	return
}


// DataPagerI is the data pager interface that allows this object to call into subclasses.
type DataPagerI interface {
	page.ControlI
	PreviousButtonsHtml() string
	NextButtonsHtml() string
	PageButtonsHtml(i int) string
	SetMaxPageButtons(b int)
	SetObjectNames(singular string, plural string)
	SetLabels(previous string, next string)
}

// DataPager is a toolbar designed to aid scrolling a large set of data. It is implemented using Aria design
// best practices. It is designed to be paired with a Table or DataRepeater to aid in navigating through the data.
// It is similar to a Paginator, but a paginator is for navigating through a series of related web pages and not just for
// data on one form.
type DataPager struct {
	page.ControlBase

	maxPageButtons   int
	ObjectName       string
	ObjectPluralName string
	LabelForNext     string
	LabelForPrevious string

	pagedControlID string
}

// NewDataPager creates a new DataPager
func NewDataPager(parent page.ControlI, id string, pagedControl PagedControlI) *DataPager {
	d := &DataPager{}
	d.Self = d
	d.Init(parent, id, pagedControl)
	return d
}

// Init is called by subclasses of a DataPager to initialize the data pager. You do not normally need
// to call this.
func (d *DataPager) Init(parent page.ControlI, id string, pagedControl PagedControlI) {
	d.ControlBase.Init(parent, id)
	d.Tag = "div"
	d.LabelForNext = d.GT("Next")
	d.LabelForPrevious = d.GT("Previous")
	d.maxPageButtons = DefaultMaxPagerButtons
	pagedControl.AddDataPager(d.Self.(DataPagerI))
	d.pagedControlID = pagedControl.ID()
	pxy := NewProxy(d, d.proxyID())
	pxy.On(event.Click().Bubbles(), action.Ajax(d.ID(), PageClick))
	d.SetAttribute("role", "tablist")
	d.PagedControl().SetPageNum(1)
}

func (d *DataPager) proxyID() string {
	return d.ID() + "-pxy"
}

func (d *DataPager) ButtonProxy() *Proxy {
	return GetProxy(d, d.proxyID())
}


// DrawingAttributes is called by the framework to add temporary attributes to the html.
func (d *DataPager) DrawingAttributes(ctx context.Context) html.Attributes {
	a := d.ControlBase.DrawingAttributes(ctx)
	a.SetDataAttribute("grctl", "datapager")
	return a
}

// Action is called by the framework to respond to actions.
func (d *DataPager) Action(ctx context.Context, params page.ActionParams) {
	switch params.ID {
	case PageClick:
		pageNum := params.ControlValueInt();
		p := d.PagedControl()
		if pageNum < 1 {
			pageNum = 1
		} else {
			c := p.CalcPageCount()
			if pageNum > c {
				pageNum = c
			}
		}
		p.SetPageNum(pageNum)
		p.Refresh()
		for _, c := range p.GetDataPagerIDs() {
			GetDataPager(d, c).Refresh()
		}
	}
}

func (d *DataPager) refreshPagedControl() {
	d.PagedControl().Refresh()
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
	p := d.PagedControl()
	pageCount := p.CalcPageCount()
	pageNum := p.PageNum()

	if pageCount <= d.maxPageButtons {
		return 1, pageCount
	} else {
		_, minEndOfBunch := math.MinInt(d.maxPageButtons-2, pageCount)
		_, maxStartOfBunch := math.MaxInt(pageCount-d.maxPageButtons+3, 1)

		leftOfBunchCount := (d.maxPageButtons - 5) / 2
		rightOfBunchCount := (d.maxPageButtons - 4) / 2

		leftBunchTrigger := leftOfBunchCount + 4
		rightBunchTrigger := maxStartOfBunch + (d.maxPageButtons-7)/2

		if pageNum < leftBunchTrigger {
			pageStart = 1
		} else {
			_, pageStart = math.MinInt(maxStartOfBunch, pageNum-leftOfBunchCount)
		}

		if pageNum > rightBunchTrigger {
			pageEnd = pageCount
		} else {
			_, pageEnd = math.MaxInt(minEndOfBunch, pageNum+rightOfBunchCount)
		}
		return
	}
}

// PreRender is called by the framework to load data into the paged control just before drawing.
func (d *DataPager) PreRender(ctx context.Context, w io.Writer) (err error) {
	err = d.ControlBase.PreRender(ctx, w)
	p := d.PagedControl()

	if err == nil {
		// If we are being drawn before the paged control, we must tell the paged control to load up its
		// data so that we can figure out what to do
		if !p.WasRendered() &&
			!p.IsRendering() { // not a child control
			p.LoadData(ctx, p)
		}
	}
	return
}

// DrawInnerHtml is called by the framework to draw the control's inner html.
func (d *DataPager) DrawInnerHtml(ctx context.Context, w io.Writer) (err error) {
	h := d.Self.(DataPagerI).PreviousButtonsHtml()
	pageStart, pageEnd := d.CalcBunch()
	for i := pageStart; i <= pageEnd; i++ {
		h += d.Self.(DataPagerI).PageButtonsHtml(i)
	}

	h += d.Self.(DataPagerI).NextButtonsHtml()
	_, err = io.WriteString(w, h)
	return
}

// PreviousButtonsHtml returns the html to draw the previous buttons. Subclasses can override this to
// change how the Previous buttons are drawn.
func (d *DataPager) PreviousButtonsHtml() string {
	var prev string
	var actionValue string

	pageNum := d.PagedControl().PageNum()

	actionValue = strconv.Itoa(pageNum - 1)

	attr := html.NewAttributes().
		Set("id", d.ID()+"_arrow_"+actionValue).
		SetClass("arrow previous")

	if pageNum <= 1 {
		attr.SetDisabled(true)
		attr.SetStyle("cursor", "not-allowed")
	}

	prev = d.ButtonProxy().ButtonHtml(d.LabelForPrevious, actionValue, attr, false)

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

	p := d.PagedControl()
	pageNum := p.PageNum()
	actionValue = strconv.Itoa(pageNum + 1)

	attr := html.NewAttributes().
		Set("id", d.ID()+"_arrow_"+actionValue).
		SetClass("arrow next")

	_, pageEnd := d.CalcBunch()
	pageCount := p.CalcPageCount()

	if pageNum >= pageCount {
		attr.SetDisabled(true)
		attr.SetStyle("cursor", "not-allowed")
	}

	next = d.ButtonProxy().ButtonHtml(d.LabelForNext, actionValue, attr, false)

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
	p := d.PagedControl()
	pageNum := p.PageNum()

	if pageNum == i {
		attr.AddClass("selected")
		attr.Set("aria-selected", "true")
		attr.Set("tabindex", "0")
	} else {
		attr.Set("aria-selected", "false")
		attr.Set("tabindex", "-1")
		// TODO: We need javascript to respond to arrow keys to set the focus on the buttons. User could then press space to click on button.
	}
	return d.ButtonProxy().ButtonHtml(actionValue, actionValue, attr, false)
}

// MarshalState is an internal function to save the state of the control
func (d *DataPager) MarshalState(m maps.Setter) {
	m.Set("pageNum", d.PagedControl().PageNum())
}

// UnmarshalState is an internal function to restore the state of the control
func (d *DataPager) UnmarshalState(m maps.Loader) {
	if v, ok := m.Load("pageNum"); ok {
		if i, ok := v.(int); ok {
			d.PagedControl().SetPageNum(i) // admittedly, multiple pagers will repeat the same call, but not likely to effect performance
		}
	}
}

func (d *DataPager) PagedControl() PagedControlI {
	return d.Page().GetControl(d.pagedControlID).(PagedControlI)
}

func (d *DataPager) Serialize(e page.Encoder) (err error) {
	if err = d.ControlBase.Serialize(e); err != nil {
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

	if err = e.Encode(d.pagedControlID); err != nil {
		return
	}

	return
}

func (d *DataPager) Deserialize(dec page.Decoder) (err error) {
	if err = d.ControlBase.Deserialize(dec); err != nil {
		panic(err)
	}

	if err = dec.Decode(&d.maxPageButtons); err != nil {
		panic(err)
	}

	if err = dec.Decode(&d.ObjectName); err != nil {
		panic(err)
	}

	if err = dec.Decode(&d.ObjectPluralName); err != nil {
		panic(err)
	}

	if err = dec.Decode(&d.LabelForNext); err != nil {
		panic(err)
	}

	if err = dec.Decode(&d.LabelForPrevious); err != nil {
		panic(err)
	}

	if err = dec.Decode(&d.pagedControlID); err != nil {
		panic(err)
	}

	return
}

func GetDataPager(c page.ControlI, id string) DataPagerI {
	return c.Page().GetControl(id).(DataPagerI)
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
}


// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c DataPagerCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	if !parent.Page().HasControl(c.PagedControl) {
		panic ("you must declare the paged control before the data pager")
	}
	p := parent.Page().GetControl(c.PagedControl).(PagedControlI)
	ctrl := NewDataPager(parent, c.ID, p)
	c.Init(ctx, ctrl)
	return ctrl
}

// Init is called by implementations of Buttons to initialize a control with the
// creator. You do not normally need to call this.
func (c DataPagerCreator) Init(ctx context.Context, ctrl DataPagerI) {
	if c.MaxPageButtons > 0 {
		ctrl.SetMaxPageButtons(c.MaxPageButtons)
	}
	ctrl.SetObjectNames(c.ObjectName, c.ObjectPluralName)
	if c.LabelForNext != "" {
		ctrl.SetLabels(c.LabelForPrevious, c.LabelForNext)
	}

	ctrl.ApplyOptions(ctx, c.ControlOptions)
}

func init() {
	page.RegisterControl(&DataPager{})
}