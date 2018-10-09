package control

import (
	"context"
	"fmt"
	"github.com/spekary/gengen/maps"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	"goradd-project/config"
	"goradd-project/override/control_base"
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
	control_base.Table
	selectedID string
}

func NewSelectTable(parent page.ControlI, id string) *SelectTable {
	t := &SelectTable{}
	t.Init(t, parent, id)
	return t
}

func (t *SelectTable) Init(self page.ControlI, parent page.ControlI, id string) {
	t.Table.Init(self, parent, id)
	t.ParentForm().AddJavaScriptFile(config.GoraddAssets() + "/js/jquery.scrollIntoView.js", false, nil)
	t.ParentForm().AddJavaScriptFile(config.GoraddAssets() + "/js/select-table.js", false, nil)
}

func (t *SelectTable) this() SelectTableI {
	return t.Self.(SelectTableI)
}

func (t *SelectTable) GetRowAttributes(row int, data interface{}) (a *html.Attributes) {
	var id string

	if t.RowStyler() != nil {
		a = t.RowStyler().Attributes(row, data)
		id = a.Get("id") // styler might be giving us an id
	} else {
		a = html.NewAttributes()
	}

	// try to guess the id from the data
	if id == "" {
		switch obj := data.(type) {
		case IDer:
			id = obj.ID()
		case PrimaryKeyer:
			id = obj.PrimaryKey()
		case map[string]string:
			id,_ = obj["id"]
		case maps.StringGetter:
			id = obj.Get("id")
		}
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
	t.ParentForm().Response().ExecuteControlCommand(t.ID(), "selectTable", "option", "selectedId", id)
}

func (t *SelectTable) MarshalState(m maps.Setter) {
	m.Set("selId", t.selectedID)
}

func (t *SelectTable) UnmarshalState(m maps.Loader) {
	if v,ok := m.Load("selId"); ok {
		if id, ok := v.(string); ok {
			t.selectedID = id
		}
	}
}

func (t *SelectTable) PutCustomScript(ctx context.Context, response *page.Response) {
	options := map[string]interface{}{}
	options["selectedId"] = t.selectedID

	response.ExecuteControlCommand(t.ID(), "selectTable", page.PriorityHigh, options)
}
