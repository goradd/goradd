package widget

import (
	"context"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/event"
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
	ItemTable  *SelectTable
	NewButton  *Button
	EditButton *Button
}

func NewItemListPanel(parent page.ControlI, id string) *ItemListPanel {
	p := &ItemListPanel{}
	p.Init(p, parent, id)
	p.ParentForm().AddStyleSheetFile(config.GoraddAssets()+"/css/item-list-panel.css", nil)

	return p
}

func (p *ItemListPanel) Init(self ItemListPanelI, parent page.ControlI, id string) {
	p.Panel.Init(self, parent, id)
	p.FilterPanel = NewPanel(p, p.ID()+"-filter")
	p.ScrollPanel = NewPanel(p, p.ID()+"-scroller")
	p.ButtonPanel = NewPanel(p, p.ID()+"-btnpnl")

	p.FilterPanel.AddClass("filter")
	p.ScrollPanel.AddClass("scroller")
	p.ButtonPanel.AddClass("buttons")

	p.FilterText = NewTextbox(p.FilterPanel, p.ID()+"-filtertxt")
	p.FilterText.SetPlaceholder(p.ParentForm().ΩT("Search"))
	p.FilterText.SetType(TextboxTypeSearch)

	p.ItemTable = NewSelectTable(p.ScrollPanel, p.ID()+"-table")

	p.NewButton = NewButton(p.ButtonPanel, p.ID()+"-newbtn")
	p.NewButton.SetText(p.ParentForm().ΩT("New"))

	p.EditButton = NewButton(p.ButtonPanel, p.ID()+"-editbtn")
	p.EditButton.SetText(p.ParentForm().ΩT("Edit"))

	p.FilterText.On(event.Input().Delay(300), action.Ajax(p.ID(), filterChanged))
	p.FilterText.On(event.EnterKey().Terminating(), action.Ajax(p.ID(), filterChanged))

}

func (p *ItemListPanel) Load(ctx context.Context) {
	p.ItemTable.SaveState(ctx, true)
}

func (p *ItemListPanel) this() ItemListPanelI {
	return p.Self.(ItemListPanelI)
}

func (c *ItemListPanel) ΩDrawingAttributes() html.Attributes {
	a := c.Panel.ΩDrawingAttributes()
	a.SetDataAttribute("grctl", "itemlistpnl")
	return a
}

func (f *ItemListPanel) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case filterChanged:
		f.ItemTable.Refresh() // TODO: Change this to some kind of data only refresh so that when control is redrawn the scroll position is maintained
	}
}
