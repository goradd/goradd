package panels

import (
	"context"
	. "github.com/goradd/goradd/pkg/bootstrap/control"
	"github.com/goradd/goradd/pkg/bootstrap/examples"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
)

type TabsPanel struct {
	control.Panel
}

func NewTabsPanel(ctx context.Context, parent page.ControlI) {
	p := new(TabsPanel)
	p.Init(p, ctx, parent, "TabsPanel")
}

func (f *TabsPanel) Init(self any, ctx context.Context, parent page.ControlI, id string) {
	f.Panel.Init(self, parent, id)
	f.AddControls(ctx,
		TabsCreator{
			ID:       "tabs-panel",
			TabStyle: TabStyleTabs,
			Children: control.Children(
				control.PanelCreator{
					Text: "Tab 1",
					Children: control.Children(
						control.PanelCreator{
							Text: "First tab content",
						},
					),
				},
				control.PanelCreator{
					Text: "Tab 2",
				},
				control.PanelCreator{
					Text: "Tab 3",
				},
			),
		},
		TabsCreator{
			ID:       "pills-panel",
			TabStyle: TabStylePills,
			Children: control.Children(
				control.PanelCreator{
					Text: "Tab 1",
					Children: control.Children(
						control.PanelCreator{
							Text: "First tab content",
						},
					),
				},
				control.PanelCreator{
					Text: "Tab 2",
				},
				control.PanelCreator{
					Text: "Tab 3",
				},
			),
		},
		TabsCreator{
			ID:       "underline-panel",
			TabStyle: TabStyleUnderline,
			Children: control.Children(
				control.PanelCreator{
					Text: "Tab 1",
					Children: control.Children(
						control.PanelCreator{
							Text: "First tab content",
						},
					),
				},
				control.PanelCreator{
					Text: "Tab 2",
				},
				control.PanelCreator{
					Text: "Tab 3",
				},
			),
		},
	)
}

func init() {
	examples.RegisterPanel("tabs", "Tabs", NewTabsPanel, 6)
	page.RegisterControl(&TabsPanel{})
}
