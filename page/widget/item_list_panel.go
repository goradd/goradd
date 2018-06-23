package widget

import (
	. "github.com/spekary/goradd/page/control"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/html"
	"goradd/config"
	"github.com/spekary/goradd/page/control/table"
	"context"
	"github.com/spekary/goradd/page/event"
	"github.com/spekary/goradd/page/action"
)

const (
	filterChanged = iota + 4000
)

type ItemListPanelI interface {
	PanelI
}

type ItemListPanel struct {
	Panel
	FilterPanel *Panel
	ScrollPanel *Panel
	ButtonPanel *Panel

	FilterText *Textbox
	ItemTable *table.SelectTable
	NewButton *Button
	EditButton *Button
}

func NewItemListPanel(parent page.ControlI, id string) *ItemListPanel {
	p := &ItemListPanel{}
	p.Init(p, parent, id)
	p.Form().AddStyleSheetFile(config.GoraddAssets() + "/css/item-list-panel.css", nil)

	return p
}

func (p *ItemListPanel) Init(self ItemListPanelI, parent page.ControlI, id string) {
	p.Panel.Init(self, parent, id)
	p.FilterPanel = NewPanel(p, p.ID() + "_filter")
	p.ScrollPanel = NewPanel(p, p.ID() + "_scroller")
	p.ButtonPanel = NewPanel(p, p.ID() + "_btnpnl")

	p.FilterPanel.AddClass("filter")
	p.ScrollPanel.AddClass("scroller")
	p.ButtonPanel.AddClass("buttons")

	p.FilterText = NewTextbox(p.FilterPanel, p.ID() + "_filtertxt")
	p.FilterText.SetPlaceholder(p.Form().T("Search"))
	p.FilterText.SetType(TextboxTypeSearch)

	p.ItemTable = table.NewSelectTable(p.ScrollPanel, p.ID() + "_table")

	p.NewButton = NewButton(p.ButtonPanel, p.ID() + "_newbtn")
	p.NewButton.SetText(p.Form().T("New"))

	p.EditButton = NewButton(p.ButtonPanel, p.ID() + "_editbtn")
	p.EditButton.SetText(p.Form().T("Edit"))

	p.FilterText.On(event.Input().Delay(300), action.Ajax(p.ID(), filterChanged))
	p.FilterText.On(event.EnterKey().Terminating(), action.Ajax(p.ID(), filterChanged))

}

func (p *ItemListPanel) Load(ctx context.Context) {
	p.ItemTable.SaveState(ctx, true)
}

func (p *ItemListPanel) this() ItemListPanelI {
	return p.Self.(ItemListPanelI)
}

func (c *ItemListPanel) DrawingAttributes() *html.Attributes {
	a := c.Panel.DrawingAttributes()
	a.SetDataAttribute("grctl", "itemlistpnl")
	return a
}

func (f *ItemListPanel) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case filterChanged:
		f.ItemTable.Refresh() // TODO: Change this to some kind of data only refresh so that when control is redrawn the scroll position is maintained
	}
}