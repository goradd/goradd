package control

import (
	localPage "goradd/page"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/util"
)

type Paginator struct {
	localPage.Control

	totalItems int
	pageSize int
	pageNum int
	maxPageButtons int
	objectName string
	objectPluralName string
	labelForNext string
	labelForPrevious string

	paginatedControl page.ControlI	// TODO: Switch to a Paginater or PaginateI
}

func NewPaginator(parent page.ControlI, paginatedControl page.ControlI) *Paginator {
	p := Paginator{paginatedControl:paginatedControl}
	p.Init(parent)
	return &p
}

func (p *Paginator) Init(parent page.ControlI) {
	p.Control.Init(p, parent)
	p.labelForNext = p.T("Next")
	p.labelForPrevious = p.T("Previous")
	p.maxPageButtons = 10
}

func (p *Paginator) SetTotalItems(count int) {
	p.totalItems = count
	p.limitPageNumber()
	p.Refresh()
}

func (p *Paginator) TotalItems() int {
	return p.totalItems
}

func (p *Paginator)  SetPageSize(size int) {
	p.pageSize = size
}

func (p *Paginator) PageSize() int {
	return p.pageSize
}

func (p *Paginator) SetPageNum(n int) {
	p.pageNum = n
	p.Refresh()
	p.refreshPaginatedControl()
}

func (p *Paginator) refreshPaginatedControl() {
	p.paginatedControl.Refresh()
}

// SetMaxPageButtons sets the maximum number of buttons that will be displayed in the paginator.
func (p *Paginator) SetMaxPageButtons(b int) {
	p.maxPageButtons = b
}

func (p *Paginator) SetObjectNames(singular string, plural string) {
	p.objectName = singular
	p.objectPluralName = plural
}

// SetLabels sets the previous and next labels. Translate these first.
func (p *Paginator) SetLabels(previous string, next string) {
	p.labelForPrevious = previous
	p.labelForNext = next
}

func (p *Paginator) limitPageNumber() {
	pageCount := p.calcPageCount()

	if p.pageNum > pageCount {
		if pageCount <= 1 {
			p.pageNum = 1
		} else {
			p.pageNum = pageCount
		}
	}
}

func (p *Paginator) calcPageCount() int {
	return (p.totalItems - 1)/p.pageSize + 1
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
 */

func (p *Paginator) calcBunch() (pageStart, pageEnd int) {

	pageCount := p.calcPageCount();
	if pageCount <= p.maxPageButtons {
		return 1, pageCount
	} else {
		minEndOfBunch := util.MinInt(p.maxPageButtons-2, pageCount)
		maxStartOfBunch := util.MaxInt(pageCount-p.maxPageButtons+3, 1)

		leftOfBunchCount := (p.maxPageButtons - 5) / 2
		rightOfBunchCount := (p.maxPageButtons - 4) / 2

		leftBunchTrigger := leftOfBunchCount + 4
		rightBunchTrigger := maxStartOfBunch + (p.maxPageButtons-7)/2

		if p.pageNum < leftBunchTrigger {
			pageStart = 1
		} else {
			pageStart = util.MinInt(maxStartOfBunch, p.pageNum-leftOfBunchCount)
		}

		if p.pageNum > rightBunchTrigger {
			pageEnd = pageCount
		} else {
			pageEnd = util.MaxInt(minEndOfBunch, p.pageNum+rightOfBunchCount)
		}
		return
	}
}