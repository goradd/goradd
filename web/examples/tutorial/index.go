package tutorial

import (
	"context"
	"strconv"
	"strings"

	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
)

const IndexFormPath = "/goradd/tutorial.g"
const IndexFormId = "IndexForm"

const (
	TestButtonAction = iota + 1
)

type IndexForm struct {
	FormBase
	detail *Panel
}

type createFunction func(ctx context.Context, parent page.ControlI) page.ControlI
type pageRecord struct {
	i     int
	title string
	f     createFunction
}
type pageRecordList []pageRecord

var pages = make(map[string]pageRecordList)

func (p pageRecordList) Less(i, j int) bool {
	return p[i].i < p[j].i
}

func NewIndexForm(ctx context.Context) page.FormI {
	f := &IndexForm{}
	f.Init(ctx, f, IndexFormPath, IndexFormId)
	f.AddRelatedFiles()

	f.detail = NewPanel(f, "detailPanel")

	return f
}

func (f *IndexForm) LoadControls(ctx context.Context) {
	if pageID, ok := page.GetContext(ctx).FormValue("pageID"); ok {
		// pageID is a category and integer id combined
		parts := strings.Split(pageID, "-")
		if len(parts) != 2 {
			return
		}

		pl, ok := pages[parts[0]]
		if !ok {
			return
		}

		i, err := strconv.Atoi(parts[1])
		if err != nil {
			return
		}

		if len(pl) <= i {
			return
		}

		pr := pl[i]
		pr.f(ctx, f.detail)

	} else {
		NewDefaultPanel(ctx, f.detail)
	}
}

func RegisterTutorialPage(category string, i int, title string, f createFunction) {
	v, ok := pages[category]
	if !ok {
		pages[category] = pageRecordList{pageRecord{i, title, f}}
	} else {
		v = append(v, pageRecord{i, title, f})
		pages[category] = v
	}
}

func init() {
	page.RegisterPage(IndexFormPath, NewIndexForm, IndexFormId)
}
