package tutorial

import (
	"context"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/event"
	"strings"

	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
)

const IndexFormPath = "/goradd/tutorial.g"
const IndexFormId = "IndexForm"

const (
	ViewSourceAction = iota + 1
)

type IndexForm struct {
	FormBase
	detailPanel      *Panel
	viewSourceButton *Button
	sourcePanel		 *SourcePanel

	currentPageRecord pageRecord
}

type createFunction func(ctx context.Context, parent page.ControlI) page.ControlI
type pageRecord struct {
	order     int
	id 	  string
	title string
	f     createFunction
	files []string
}
type pageRecordList []pageRecord

var pages = make(map[string]pageRecordList)

func (p pageRecordList) Less(i, j int) bool {
	return p[i].order < p[j].order
}

func NewIndexForm(ctx context.Context) page.FormI {
	f := &IndexForm{}
	f.Init(ctx, f, IndexFormPath, IndexFormId)
	f.AddRelatedFiles()

	f.detailPanel = NewPanel(f, "detailPanel")
	f.viewSourceButton = NewButton(f, "viewSourceButton")
	f.viewSourceButton.SetLabel("View Source")
	f.viewSourceButton.On(event.Click(), action.Ajax(f.ID(), ViewSourceAction))

	f.sourcePanel = NewSourcePanel(f)
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

		id := parts[1]

		for _,pr := range pl {
			if pr.id == id {
				pr.f(ctx, f.detailPanel)
				f.currentPageRecord = pr
				break
			}
		}
	} else {
		NewDefaultPanel(ctx, f.detailPanel)
	}
}


func (f *IndexForm) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case ViewSourceAction:
		f.sourcePanel.show(f.currentPageRecord.files)
	}
}


func init() {
	page.RegisterPage(IndexFormPath, NewIndexForm, IndexFormId)
}

func RegisterTutorialPage(category string, order int, id string, title string, f createFunction, files []string) {
	v, ok := pages[category]
	if !ok {
		pages[category] = pageRecordList{pageRecord{order, id, title, f, files}}
	} else {
		v = append(v, pageRecord{order, id, title, f, files})
		pages[category] = v
	}
}
