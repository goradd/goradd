package examples

import (
	"context"
	bootstrap "github.com/goradd/goradd/pkg/bootstrap/control"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
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
}

func (f *ControlsForm) Init(ctx context.Context, formID string) {
	f.FormBase.Init(ctx, formID)
	f.AddRelatedFiles()
	f.AddControls(ctx,
		bootstrap.NavbarCreator{
			ID: "nav",
			Children: Children(
				bootstrap.NavGroupCreator{
					ID: "navList",
				},
			),
		},
		PanelCreator{
			ID: "detailPanel",
			ControlOptions: page.ControlOptions{
				Class: "container",
			},
		},
	)
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

func (f *ControlsForm) BindData(ctx context.Context, s DataManagerI) {
	navGroup := bootstrap.GetNavGroup(f, "navList")
	navGroup.RemoveChildren()
	sort.Slice(controls, func(i, j int) bool {
		return controls[i].order < controls[j].order
	})
	pageContext := page.GetContext(ctx)
	for _, c := range controls {
		item := bootstrap.NewNavLink(navGroup, c.key)
		a := url.
			NewBuilderFromUrl(pageContext.URL).
			SetValue("control", c.key).
			String()
		item.SetLocation(a)
	}
}

func init() {
	page.RegisterForm(ControlsFormPath, &ControlsForm{}, ControlsFormId)
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
