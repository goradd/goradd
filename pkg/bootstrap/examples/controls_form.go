package examples

import (
	"context"
	bootstrap "github.com/goradd/goradd/pkg/bootstrap/control"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/control/data"
	"github.com/goradd/goradd/pkg/url"
	"sort"
)

const ControlsFormPath = "/goradd/examples/bootstrap.g"
const ControlsFormId = "BootstrapControlsForm"

const (
	TestButtonAction = iota + 1
)

type ControlsForm struct {
	FormBase
	list   *bootstrap.NavbarList
}

func NewControlsForm(ctx context.Context) page.FormI {
	f := &ControlsForm{}
	f.Init(ctx, f, ControlsFormPath, ControlsFormId)
	f.AddRelatedFiles()

	nav := bootstrap.NewNavbar(f, "nav")
	f.list = bootstrap.NewNavbarList(nav, "navList")

	f.list.SetDataProvider(f)
	NewPanel(f, "detailPanel").AddClass("container")

	return f
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

	createF(ctx, GetPanel(f, "detailPanel"))
}

func (f *ControlsForm) BindData(ctx context.Context, s data.DataManagerI) {
	sort.Slice(controls, func(i, j int) bool {
		return controls[i].order < controls[j].order
	})
	pageContext := page.GetContext(ctx)
	for _, c := range controls {
		item := f.list.AddItem(c.name, c.key)
		a := url.
			NewBuilderFromUrl(*pageContext.URL).
			SetValue("control", c.key).
			String()
		item.SetAnchor(a)
	}
}


func init() {
	page.RegisterPage(ControlsFormPath, NewControlsForm, ControlsFormId)
}



type createFunction func(ctx context.Context, parent page.ControlI)

type controlEntry struct {
	key   string
	name  string
	f     createFunction
	order int
}

var controls []controlEntry


func RegisterPanel(key string,
	name string,
	f createFunction,
	order int) {

	for _, c := range controls {
		if c.key == key {
			panic("panel " + key + " is already registered")
		}
	}
	controls = append(controls, controlEntry{key, name, f, order})
}
