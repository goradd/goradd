package control

import (
	"bytes"
	"context"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/javascript"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/page/action"
	"github.com/spekary/goradd/util"
	"github.com/spekary/goradd/util/types"
	localPage "goradd-project/override/page"
	"strconv"
	"goradd-project/config"
	"github.com/spekary/goradd/page/control/data"
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
	getDataPagers() []DataPagerI
}

// PaginatedControl is a mixin that makes a control controllable by a data pager
type PaginatedControl struct {
	totalItems       int
	pageSize         int
	pageNum          int
	dataPagers 		 []DataPagerI
}

func (c *PaginatedControl) SetTotalItems(count uint) {
	c.totalItems = int(count)
	c.limitPageNumber()
}

func (c *PaginatedControl) TotalItems() int {
	return c.totalItems
}

func (c *PaginatedControl) SetPageSize(size int) {
	c.pageSize = size
}

func (c *PaginatedControl) PageSize() int {
	return c.pageSize
}

func (c *PaginatedControl) PageNum() int {
	return c.pageNum
}

func (c *PaginatedControl) SetPageNum(n int) {
	if c.pageNum != n {
		c.pageNum = n
	}
}

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

func (c *PaginatedControl) CalcPageCount() int {
	if c.pageSize == 0 || c.totalItems == 0 {
		return 0
	}
	return (c.totalItems-1)/c.pageSize + 1
}

func (c *PaginatedControl) getDataPagers() []DataPagerI{
	return c.dataPagers
}



// DataPagerI is the data pager interface that allows this object to call into subclasses.
type DataPagerI interface {
	page.ControlI
	PreviousButtonsHtml() string
	NextButtonsHtml() string
	PageButtonsHtml(i int) string
}

// DataPager is a toolbar designed to aid scrolling through a large set of data. It is implemented using Aria design
// best practices. It is designed to be paired with a Table or DataRepeater to aid in navigating through the data.
// It is similar to a Paginator, but a paginator is for navigating through a whole series of pages and not just for
// data on one override.
type DataPager struct {
	localPage.Control

	maxPageButtons   int
	ObjectName       string
	ObjectPluralName string
	LabelForNext     string
	LabelForPrevious string

	paginatedControl PaginatedControlI
	Proxy *Proxy
}

func NewDataPager(parent page.ControlI, id string, paginatedControl PaginatedControlI) *DataPager {
	d := DataPager{}
	d.Init(&d, parent, id, paginatedControl)
	return &d
}

func (d *DataPager) Init(self page.ControlI, parent page.ControlI, id string, paginatedControl PaginatedControlI) {
	d.Control.Init(self, parent, id)
	d.Tag = "div"
	d.LabelForNext = d.T("Next")
	d.LabelForPrevious = d.T("Previous")
	d.maxPageButtons = config.MaxPageButtons
	paginatedControl.AddDataPager(self.(DataPagerI))
	d.paginatedControl = paginatedControl
	d.Proxy = NewProxy(d)
	d.Proxy.OnClick(action.Ajax(d.ID(), PageClick))
	d.SetAttribute("role", "tablist")
	d.paginatedControl.SetPageNum(1)
}

func (d *DataPager) DrawingAttributes() *html.Attributes {
	a := d.Control.DrawingAttributes()
	a.SetDataAttribute("grctl", "datapager")
	return a
}

func (d *DataPager) Action(ctx context.Context, params page.ActionParams) {
	switch params.ID {
	case PageClick:
		d.paginatedControl.SetPageNum(javascript.NumberInt(params.Values.Control))
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

func (d *DataPager) SetObjectNames(singular string, plural string) {
	d.ObjectName = singular
	d.ObjectPluralName = plural
}

// SliceOffsets returns the start and end values to use to specify a portion of a slice corresponding to the
// data the pager refers to
func (d *DataPager) SliceOffsets() (start, end int) {
	start = (d.paginatedControl.PageNum() - 1) * d.paginatedControl.PageSize()
	end = util.MinInt(start+ d.paginatedControl.PageSize(),  d.paginatedControl.TotalItems())
	return
}

// SqlLimits returns the limits you would use in a sql database limit clause
func (d *DataPager) SqlLimits() (maxRowCount, offset int) {
	offset = (d.paginatedControl.PageNum() - 1) * d.paginatedControl.PageSize()
	maxRowCount = d.paginatedControl.PageSize()
	return
}

// SetLabels sets the previous and next labels. Translate these first.
func (d *DataPager) SetLabels(previous string, next string) {
	d.LabelForPrevious = previous
	d.LabelForNext = next
}



/**
 * "Bunch" is defined as the collection of numbers that lies in between the pair of Ellipsis ("...")
 *
 * LAYOUT
 *
 * For IndexCount of 10
 * 2   213   2 (two items to the left of the bunch, and then 2 indexes, selected index, 3 indexes, and then two items to the right of the bunch)
 * e.g. 1 ... 5 6 *7* 8 9 10 ... 100
 *
 * For IndexCount of 11
 * 2   313   2
 *
 * For IndexCount of 12
 * 2   314   2
 *
 * For IndexCount of 13
 * 2   414   2
 *
 * For IndexCount of 14
 * 2   415   2
 *
 *
 *
 * START/END PAGE NUMBERS FOR THE BUNCH
 *
 * For IndexCount of 10
 * 1 2 3 4 5 6 7 8 .. 100
 * 1 .. 4 5 *6* 7 8 9 .. 100
 * 1 .. 92 93 *94* 95 96 97 .. 100
 * 1 .. 93 94 95 96 97 98 99 100
 *
 * For IndexCount of 11
 * 1 2 3 4 5 6 7 8 9 .. 100
 * 1 .. 4 5 6 *7* 8 9 10 .. 100
 * 1 .. 91 92 93 *94* 95 96 97 .. 100
 * 1 .. 92 93 94 95 96 97 98 99 100
 *
 * For IndexCount of 12
 * 1 2 3 4 5 6 7 8 9 10 .. 100
 * 1 .. 4 5 6 *7* 8 9 10 11 .. 100
 * 1 .. 90 91 92 *93* 94 95 96 97 .. 100
 * 1 .. 91 92 93 94 95 96 97 98 99 100
 *
 * For IndexCount of 13
 * 1 2 3 4 5 6 7 8 9 11 .. 100
 * 1 .. 4 5 6 7 *8* 9 10 11 12 .. 100
 * 1 .. 89 90 91 92 *93* 94 95 96 97 .. 100
 * 1 .. 90 91 92 93 94 95 96 97 98 99 100

Note: there are likely better ways to do this. Some innovative ones are to have groups of 10s, and then 100s etc.
Or, use the ellipsis as a dropdown menu for more selections
*/

func (d *DataPager) CalcBunch() (pageStart, pageEnd int) {

	pageCount := d.paginatedControl.CalcPageCount()
	pageNum := d.paginatedControl.PageNum()

	if pageCount <= d.maxPageButtons {
		return 1, pageCount
	} else {
		minEndOfBunch := util.MinInt(d.maxPageButtons-2, pageCount)
		maxStartOfBunch := util.MaxInt(pageCount-d.maxPageButtons+3, 1)

		leftOfBunchCount := (d.maxPageButtons - 5) / 2
		rightOfBunchCount := (d.maxPageButtons - 4) / 2

		leftBunchTrigger := leftOfBunchCount + 4
		rightBunchTrigger := maxStartOfBunch + (d.maxPageButtons-7)/2

		if pageNum < leftBunchTrigger {
			pageStart = 1
		} else {
			pageStart = util.MinInt(maxStartOfBunch, pageNum-leftOfBunchCount)
		}

		if pageNum > rightBunchTrigger {
			pageEnd = pageCount
		} else {
			pageEnd = util.MaxInt(minEndOfBunch, pageNum+rightOfBunchCount)
		}
		return
	}
}

func (d *DataPager) PreRender(ctx context.Context, buf *bytes.Buffer) (err error) {
	err = d.Control.PreRender(ctx, buf)

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

func (d *DataPager) DrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	h := d.Self.(DataPagerI).PreviousButtonsHtml()
	pageStart, pageEnd := d.CalcBunch()
	for i := pageStart; i <= pageEnd; i++ {
		h += d.Self.(DataPagerI).PageButtonsHtml(i)
	}

	h += d.Self.(DataPagerI).NextButtonsHtml()
	_, err = buf.WriteString(h)
	return
}

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

	if pageNum >= pageCount-1 {
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

func (d *DataPager) PageButtonsHtml(i int) string {
	actionValue := strconv.Itoa(i)
	attr := html.NewAttributes().Set("id", d.ID()+"_page_"+actionValue).Set("role", "tab").AddClass("override")
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

// MarshalState is an internal function to save the state of the control
func (d *DataPager) MarshalState(m types.MapI) {
	m.Set("pageNum", d.paginatedControl.PageNum())
}

// UnmarshalState is an internal function to restore the state of the control
func (d *DataPager) UnmarshalState(m types.MapI) {
	if m.Has("pageNum") {
		i, _ := m.GetInt("pageNum")
		d.paginatedControl.SetPageNum (i) // admittedly, multiple pagers will repeat the same call, but not likely to effect performance
	}
}

func (d *DataPager) PaginatedControl() PaginatedControlI {
	return d.paginatedControl
}