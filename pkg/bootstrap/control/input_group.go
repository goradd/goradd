package control

import (
	"context"
	"github.com/goradd/goradd/pkg/bootstrap/config"
	"github.com/goradd/goradd/pkg/page"
	grctl "github.com/goradd/goradd/pkg/page/control"
)

type InputGroupI interface {
	grctl.PanelI
}

type InputGroup struct {
	grctl.Panel
}


// NewInputGroup creates a new input group
func NewInputGroup(parent page.ControlI, id string) *InputGroup {
	b := new(InputGroup)
	b.Self = b
	b.Init(parent, id)
	return b
}

func (g *InputGroup) Init(parent page.ControlI, id string) {
	g.Panel.Init(parent, id)
	config.LoadBootstrap(g.ParentForm())
}

func (g *InputGroup) this() InputGroupI {
	return g.Self.(InputGroupI)
}


type InputGroupCreator struct {
	// ID is the control id
	ID string
	Prepend []page.Creator
	Child page.Creator
	Append []page.Creator
	page.ControlOptions
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c InputGroupCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewInputGroup(parent, c.ID)
	ctrl.AddClass("input-group")
	var children []page.Creator

	if c.Prepend != nil {
		children = append(children,
			grctl.PanelCreator{
				Children:       c.Prepend,
				ControlOptions: page.ControlOptions{
					Class: "input-group-prepend",
				},
			},
		)
	}
	children = append(children, c.Child)
	if c.Append != nil {
		children = append(children,
			grctl.PanelCreator{
				Children:       c.Append,
				ControlOptions: page.ControlOptions{
					Class: "input-group-append",
				},
			},
		)
	}
	ctrl.ApplyOptions(ctx, c.ControlOptions)
	ctrl.AddControls(ctx, children...)

	return ctrl
}


func init() {
	page.RegisterControl(new (InputGroup))
}