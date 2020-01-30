package tutorial

import (
	"context"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/page/action"
	"path"
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
	Cat string
	Num int
}


func (f *IndexForm) Init(ctx context.Context, formID string) {
	f.FormBase.Init(ctx, formID)
	f.AddRelatedFiles()

	f.AddControls(ctx,
		PanelCreator{
			ID:"detailPanel",
		},
		ButtonCreator{
			ID: "viewSourceButton",
			Text: "View Source",
			OnClick: action.Ajax(f.ID(), ViewSourceAction),
		},
	)
}

func (f *IndexForm) AddRelatedFiles() {
	f.FormBase.AddRelatedFiles()
	f.AddStyleSheetFile(path.Join(config.GoraddAssets(), "css", "welcome.css"), nil)
}

func (f *IndexForm) LoadControls(ctx context.Context) {
	if pageID, ok := page.GetContext(ctx).FormValue("pageID"); ok {
		// pageID is a category and integer id combined
		parts := strings.Split(pageID, "-")
		if len(parts) != 2 {
			return
		}
		cat := parts[0]
		id := parts[1]

		pl, ok := pages[cat]
		if !ok {
			return
		}

		for i,pr := range pl {
			if pr.id == id {
				pr.f(ctx, GetPanel(f, "detailPanel"))
				f.Cat = cat
				f.Num = i
				break
			}
		}
	} else {
		NewDefaultPanel(ctx, GetPanel(f, "detailPanel"))
	}
}

func (f *IndexForm) ShowSourceDialog() {
	d, isNew := GetDialogPanel(f,"sourceDialog")
	if isNew {
		d.SetTitle("Source")
		d.AddCloseButton("Close", "close")
		d.SetHasCloseBox(true)
		NewSourcePanel(d, "sourcePanel")
	}
	d.Show()
}


func (f *IndexForm) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case ViewSourceAction:
		if l,ok := pages[f.Cat]; ok {
			pr := l[f.Num]
			f.ShowSourceDialog()
			GetSourcePanel(f).show(pr.files)
		}
	}
}


func init() {
	page.RegisterForm(IndexFormPath, &IndexForm{}, IndexFormId)
}

