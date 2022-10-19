// Package controls contains example pages that demonstrate how to set up and
// use various Goradd controls. Most Goradd controls mirror html controls.
package controls

import (
	"context"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/http"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
	. "github.com/goradd/goradd/pkg/page/control/list"
	"github.com/goradd/goradd/pkg/url"
	"path"
)

const ControlsFormPath = "/goradd/examples/controls.g"
const ControlsFormId = "ControlsForm"

const ControlsFormListID = "listPanel"
const ControlsFormDetailID = "detailPanel"

const (
	TestButtonAction = iota + 1
)

type ControlsForm struct {
	FormBase
}

type createFunction func(ctx context.Context, parent page.ControlI)
type controlEntry struct {
	key  string
	name string
	f    createFunction
}

var controls []controlEntry

func (f *ControlsForm) Init(ctx context.Context, formID string) {
	f.FormBase.Init(ctx, formID)

	f.AddControls(ctx,
		UnorderedListCreator{
			ID:           ControlsFormListID,
			DataProvider: f,
		},
		PanelCreator{
			ID: ControlsFormDetailID,
		},
	)
}

func (f *ControlsForm) AddRelatedFiles() {
	f.FormBase.AddRelatedFiles()
	f.AddStyleSheetFile(path.Join(config.AssetPrefix, "goradd", "welcome", "css", "welcome.css"), nil)
}

func (f *ControlsForm) LoadControls(ctx context.Context) {
	var createF createFunction
	if _, ok := page.GetContext(ctx).FormValue("testing"); ok {
		f.SetAttribute("novalidate", true) // bypass html validation for testing
	}

	if id, ok := page.GetContext(ctx).FormValue("control"); ok {
		for _, c := range controls {
			if c.key == id {
				createF = c.f
			}
		}
	}

	if createF == nil {
		createF = controls[0].f
	}

	createF(ctx, GetPanel(f, ControlsFormDetailID))
}

func (f *ControlsForm) BindData(ctx context.Context, s DataManagerI) {
	pageContext := page.GetContext(ctx)
	list := GetUnorderedList(f, ControlsFormListID)
	list.Clear()
	for _, c := range controls {
		item := list.Add(c.name, c.key)
		a := url.
			NewBuilderFromUrl(pageContext.URL).
			SetValue("control", c.key).
			String()
		item.SetAnchor(http.MakeLocalPath(a))
	}
}

func RegisterPanel(key string,
	name string,
	f createFunction) {

	for _, c := range controls {
		if c.key == key {
			panic("panel " + key + " is already registered")
		}
	}
	controls = append(controls, controlEntry{key, name, f})
}

func init() {
	page.RegisterForm(ControlsFormPath, &ControlsForm{}, ControlsFormId)
}
