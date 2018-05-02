package control

import (
	localPage "goradd/page"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/util"
	"github.com/spekary/goradd/page/action"
	"context"
	"strconv"
	"github.com/spekary/goradd/javascript"
	"bytes"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/util/types"
)

const (
	PageClick  = iota + 1000
)

// DataPagebar is a toolbar designed to aid scrolling through a large set of data. It is implemented using Aria design
// best practices. It is designed to be paired with a Table or DataRepeater to aid in navigating through the data.
// It is similar to a Paginator, but a paginator is for navigating through a whole series of pages and not just for
// data on one page.
type DataPagebar struct {
	localPage.Control

	totalItems int
	pageSize int
	pageNum int
	maxPageButtons int
	objectName string
	objectPluralName string
	labelForNext string
	labelForPrevious string

	paginatedControl page.ControlI	// TODO: Switch to a Scrollerer or ScrollI. Scrolled control should be able to draw just a portion of itself

	proxy *Proxy
}

func NewPaginator(parent page.ControlI, paginatedControl page.ControlI) *DataPagebar {
	p := DataPagebar{paginatedControl:paginatedControl}
	p.Init(parent)
	return &p
}

func (d *DataPagebar) Init(parent page.ControlI) {
	d.Control.Init(d, parent)
	d.labelForNext = d.T("Next")
	d.labelForPrevious = d.T("Previous")
	d.maxPageButtons = 10
	d.proxy = NewProxy(d)
	d.proxy.OnClick(action.Ajax(d.ID(), PageClick))
	d.Tag = "div"
	d.SetAttribute("role", "tablist")
	d.AddClass("data-scrollbar")
	d.pageNum = 1
}

func (d *DataPagebar) Action(ctx context.Context, params page.ActionParams) {
	switch params.Id {
	case PageClick:
		d.SetPageNum(javascript.NumberInt(params.Values.Control))
	}
}

func (d *DataPagebar) SetTotalItems(count int) {
	d.totalItems = count
	d.limitPageNumber()
	d.Refresh()
}

func (d *DataPagebar) TotalItems() int {
	return d.totalItems
}

func (d *DataPagebar)  SetPageSize(size int) {
	d.pageSize = size
}

func (d *DataPagebar) PageSize() int {
	return d.pageSize
}

func (d *DataPagebar) PageNum() int {
	return d.pageNum
}

func (d *DataPagebar) SetPageNum(n int) {
	if d.pageNum != n {
		d.pageNum = n
		d.Refresh()
		d.refreshPaginatedControl()
	}
}

func (d *DataPagebar) refreshPaginatedControl() {
	d.paginatedControl.Refresh()
}

// SetMaxPageButtons sets the maximum number of buttons that will be displayed in the paginator.
func (d *DataPagebar) SetMaxPageButtons(b int) {
	d.maxPageButtons = b
}

func (d *DataPagebar) SetObjectNames(singular string, plural string) {
	d.objectName = singular
	d.objectPluralName = plural
}

func (d *DataPagebar) SliceOffsets() (start, end int) {
	start = (d.pageNum - 1) * d.pageSize
	end = util.MinInt(start + d.pageSize, d.totalItems)
	return
}

// SetLabels sets the previous and next labels. Translate these first.
func (d *DataPagebar) SetLabels(previous string, next string) {
	d.labelForPrevious = previous
	d.labelForNext = next
}

func (d *DataPagebar) limitPageNumber() {
	pageCount := d.calcPageCount()

	if d.pageNum > pageCount {
		if pageCount <= 1 {
			d.pageNum = 1
		} else {
			d.pageNum = pageCount
		}
	}
}

func (d *DataPagebar) calcPageCount() int {
	if d.pageSize == 0 || d.totalItems == 0 {
		return 0
	}
	return (d.totalItems - 1)/d.pageSize + 1
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

func (d *DataPagebar) calcBunch() (pageStart, pageEnd int) {

	pageCount := d.calcPageCount()
	if pageCount <= d.maxPageButtons {
		return 1, pageCount
	} else {
		minEndOfBunch := util.MinInt(d.maxPageButtons-2, pageCount)
		maxStartOfBunch := util.MaxInt(pageCount-d.maxPageButtons+3, 1)

		leftOfBunchCount := (d.maxPageButtons - 5) / 2
		rightOfBunchCount := (d.maxPageButtons - 4) / 2

		leftBunchTrigger := leftOfBunchCount + 4
		rightBunchTrigger := maxStartOfBunch + (d.maxPageButtons-7)/2

		if d.pageNum < leftBunchTrigger {
			pageStart = 1
		} else {
			pageStart = util.MinInt(maxStartOfBunch, d.pageNum-leftOfBunchCount)
		}

		if d.pageNum > rightBunchTrigger {
			pageEnd = pageCount
		} else {
			pageEnd = util.MaxInt(minEndOfBunch, d.pageNum+rightOfBunchCount)
		}
		return
	}
}

func (d *DataPagebar) DrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	h := d.previousButtonsHtml()
	pageStart, pageEnd := d.calcBunch()
	for i := pageStart; i <= pageEnd; i++ {
		h += d.pageButtonsHtml(i)
	}

	h += d.nextButtonsHtml()
	_,err = buf.WriteString(h)
	return
}


func (d *DataPagebar) previousButtonsHtml() string {
	var prev string
	var actionValue string
	actionValue = strconv.Itoa(d.pageNum - 1)

	attr := html.NewAttributes().
		Set("id", d.ID() + "_arrow_" + actionValue).
		SetClass("arrow previous")

	if d.pageNum <= 1 {
		attr.SetDisabled(true)
		attr.SetStyle("cursor", "not-allowed")
	}

	prev = d.proxy.ButtonHtml(d.labelForPrevious, actionValue, attr, false)

	h := prev
	pageStart, _ := d.calcBunch()
	if pageStart != 1 {
		h += d.pageButtonsHtml(1)
		h += `<span class="ellipsis">&hellip;</span>`
	}
	return h
}

func (d *DataPagebar) nextButtonsHtml() string {
	var next string
	var actionValue string

	attr := html.NewAttributes().
		Set("id", d.ID() + "_arrow_" + actionValue).
		SetClass("arrow next")

	actionValue = strconv.Itoa(d.pageNum + 1)

	_, pageEnd := d.calcBunch()
	pageCount := d.calcPageCount()

	if d.pageNum >= pageCount {
		attr.SetDisabled(true)
	}

	next = d.proxy.ButtonHtml(d.labelForNext, actionValue, attr, false)

	h := next
	if pageEnd != pageCount {
		h += d.pageButtonsHtml(pageCount) + h
		h = `<span class="ellipsis">&hellip;</span>` + h
	}
	return h
}

func (d *DataPagebar) pageButtonsHtml(i int) string {
	actionValue := strconv.Itoa(i)
	attr := html.NewAttributes().Set("id", d.ID() + "_page_" + actionValue).Set("role","tab").AddClass("page")
	if d.pageNum == i {
		attr.AddClass("selected")
		attr.Set("aria-selected", "true")
		attr.Set("tabindex", "0")
	} else {
		attr.Set("aria-selected", "false")
		attr.Set("tabindex", "-1")
		// TODO: We need javascript to respond to arrow keys to set the focus on the buttons. User could then press space to click on button.
	}
	return d.proxy.ButtonHtml(actionValue, actionValue, attr, false)
}

// MarshalState is an internal function to save the state of the control
func (d *DataPagebar) MarshalState(m types.MapI) {
	m.Set("pageNum", d.pageNum)
}

// UnmarshalState is an internal function to restore the state of the control
func (d *DataPagebar) UnmarshalState(m types.MapI) {
	if m.Has("pageNum") {
		i,_ := m.GetInt("pageNum")
		d.pageNum = i
	}
}
