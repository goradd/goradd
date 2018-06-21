package table

import (
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/html"
	"fmt"
	"github.com/spekary/goradd/util/types"
	"goradd/config"
	"github.com/spekary/goradd/page/control"
	"context"
)

// PrimaryKeyer is an interface that is often implemented by model objects.
type PrimaryKeyer interface {
	PrimaryKey() string
}

type SelectTableI interface {
	TableI
}

// SelectTable is a table that is row selectable. To detect a row selection, trigger on event.RowSelected
type SelectTable struct {
	Table
	selectedID string
}

func NewSelectTable(parent page.ControlI) *SelectTable {
	t := &SelectTable{}
	t.Init(t, parent)
	return t
}

func (t *SelectTable) Init(self page.ControlI, parent page.ControlI) {
	t.Table.Init(self, parent)
	t.Form().AddJavaScriptFile(config.GoraddAssets() + "/js/jquery.scrollIntoView.js", false, nil)
	t.Form().AddJavaScriptFile(config.GoraddAssets() + "/js/select-table.js", false, nil)
}

func (t *SelectTable) this() SelectTableI {
	return t.Self.(SelectTableI)
}

func (t *SelectTable) GetRowAttributes(row int, data interface{}) (a *html.Attributes) {
	if t.rowStyler != nil {
		a = t.rowStyler.Attributes(row, data)
	} else {
		a = html.NewAttributes()
	}

	var id string

	// try to guess the id from the data
	switch obj := data.(type) {
	case control.IDer:
		id = obj.ID()
	case PrimaryKeyer:
		id = obj.PrimaryKey()
	case map[string]string:
		id,_ = obj["id"]
	case StringGetter:
		id = obj.Get("id")
	}
	if id != "" {
		// TODO: If configured, encrypt the id so its not publicly showing database ids
		a.SetDataAttribute("id", id)
		a.AddClass("sel")
		if id == t.selectedID {
			a.AddClass("selected")
			a.Set("aria-selected", "")
		}
	} else {
		a.AddClass("nosel")
	}
	if row % 2 == 1 {
		a.AddClass("odd")
	} else {
		a.AddClass("even")
	}

	return a
}

func (t *SelectTable) DrawingAttributes() *html.Attributes {
	a := t.Table.DrawingAttributes()
	a.SetDataAttribute("grctl", "selecttable")
	a.Set("role", "grid")
	a.Set("aria-readonly", "true")
	return a
}


func (t *SelectTable) UpdateFormValues(ctx *page.Context) {
	if data := ctx.CustomControlValue(t.ID(), "selectedId"); data != nil {
		t.selectedID = fmt.Sprintf("%v", data)
	}
}

func (t *SelectTable) SelectedID() string {
	return t.selectedID
}

func (t *SelectTable) SetSelectedID(id string) {
	t.selectedID = id
	t.Form().Response().ExecuteControlCommand(t.ID(), "selectTable", "option", "selectedId", id)
}

func (t *SelectTable) MarshalState(m types.MapI) {
	m.Set("selId", t.selectedID)
}

func (t *SelectTable) UnmarshalState(m types.MapI) {
	if m.Has("selId") {
		id, _ := m.GetString("selId")
		t.SetSelectedID(id)
	}
}

func (t *SelectTable) PutCustomScript(ctx context.Context, response *page.Response) {
	options := map[string]interface{}{}
	options["selectedId"] = t.selectedID

	response.ExecuteControlCommand(t.ID(), "selectTable", page.PriorityHigh, options)
}
