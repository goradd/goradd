package table

import (
	"context"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/javascript"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/page/action"
	"github.com/spekary/goradd/page/control/control_base"
	"github.com/spekary/goradd/page/event"
	"strings"
)

const (
	AllClickAction = iota + 1000
)

type CheckboxColumnI interface {
	ColumnI
	CheckboxAttributes(data interface{}) *html.Attributes
}


// CheckboxColumn is a table that contains a checkbox. You must provide it a CheckboxProvider to connect ids and default data
// to the checkbox. Use Changes() to get the list of checkbox ids that have changed since the list was initially drawn.
type CheckboxColumn struct {
	ColumnBase
	showCheckAll bool
	checkboxer   CheckboxProvider
	changes      map[string]bool // records changes
}

// NewChecboxColumn creates a new table table that contains a checkbox. You must provide the id of the parent table,
// add a CheckboxProvider which will connect checkbox states to data states
//
// The table will keep track of what checkboxes have been clicked and the new values. Call Changes() to get those
// changes. Or, if you are recording your changes in real time, attach a CheckboxColumnClick event to the table.
func NewCheckboxColumn(p CheckboxProvider) *CheckboxColumn {
	if p == nil {
		panic("A checkbox attribute provider is required.")
	}

	i := CheckboxColumn{checkboxer: p}
	i.Init()
	return &i
}

func (c *CheckboxColumn) Init() {
	c.ColumnBase.Init(c)
	c.isHtml = true
	c.changes = map[string]bool{}
}

func (c *CheckboxColumn) this() CheckboxColumnI {
	return c.Self.(CheckboxColumnI)
}

func (c *CheckboxColumn) SetShowCheckAll(s bool) *CheckboxColumn {
	c.showCheckAll = s
	return c
}

func (c *CheckboxColumn) HeaderCellHtml(ctx context.Context, row int, col int) (h string) {
	if c.showCheckAll {
		a := c.this().CheckboxAttributes(nil)
		a.Set("type", "checkbox")
		h = html.RenderVoidTag("input", a)
	}
	if c.IsSortable() {
		h += c.RenderSortButton("")
	}
	return
}

// CheckboxAttributes returns the attributes for the input tag that will display the checkbox.
// If data is nil, it indicates a checkAll box.
func (c *CheckboxColumn) CheckboxAttributes(data interface{}) *html.Attributes {
	p := c.checkboxer
	a := p.Attributes(data)
	if a == nil {
		a = html.NewAttributes()
	}
	var id string
	if data == nil {
		pubid := c.ID() + "_all"
		a.Set("id", pubid)
		a.SetDataAttribute("grAll", "1")
	} else if id = p.ID(data); id != "" {
		// TODO: optionally encrypt the id in case its a database id. Difficult since database ids might themselves be large hashes (aka Google data store)
		pubid := c.ID() + "_" + id
		a.Set("id", pubid)
		a.SetDataAttribute("grCheckcol", "1")
	} else {
		panic("A checkbox id is required.")
	}

	if newVal, ok := c.changes[id]; ok { // If we have recorded a change, use this value on  refresh.
		if newVal {
			a.Set("checked", "")
		}
	} else if p.IsChecked(data) { // otherwise, use the data value
		a.Set("checked", "")
	}
	if !a.Has("value") { // value is required by html
		a.Set("value", "1")
	}
	a.SetDataAttribute("grTrackchanges", "1")
	return a
}

func (c *CheckboxColumn) CellText(ctx context.Context, row int, col int, data interface{}) string {
	a := c.this().CheckboxAttributes(data)
	a.Set("type", "checkbox")
	return html.RenderVoidTag("input", a)
}

// Changes returns a map of ids corresponding to checkboxes that have changed. Both true and false values indicate the
// current state of that particular checkbox. Note that if a user checks a box, then checks it again, even though it
// is back to its original value, it will still show up in the changes list.
func (c *CheckboxColumn) Changes() map[string]bool {
	return c.changes
}

// UpdateFormValues will look for changes to our checkboxes and record those changes.
func (c *CheckboxColumn) UpdateFormValues(ctx *page.Context) {
	for k, v := range ctx.CheckableValues() {
		index := strings.LastIndexAny(k, "_")
		if index > 0 {
			column_id := k[:index]
			check_id := k[index+1:]

			if column_id == c.ID() {
				c.changes[check_id] = control_base.ConvertToBool(v)
			}
		}
	}
}

// AddActions adds actions to the table that the table can respond to.
func (c *CheckboxColumn) AddActions(table page.ControlI) {
	table.On(event.CheckboxColumnClick().Selector(`input[data-gr-all]`), action.Ajax(c.ID(), ColumnAction).ActionValue(AllClickAction), action.PrivateAction{})
}

func (c *CheckboxColumn) Action(ctx context.Context, params page.ActionParams) {
	switch javascript.NumberInt(params.Values.Action) {
	case AllClickAction:
		p := params.Values.Event.(map[string]interface{})
		c.allClick(p["id"].(string), p["checked"].(bool), javascript.NumberInt(p["row"]), javascript.NumberInt(p["col"]))
	}
}

// The check all checkbox has been checked.
func (c *CheckboxColumn) allClick(id string, checked bool, row int, col int) {
	all := c.checkboxer.All()

	if all != nil {
		for k, v := range all {
			if v == checked {
				c.changes[k] = checked
			} else {
				delete(c.changes, k)
			}
		}
		// Fire javascript to check all visible
		//js := fmt.Sprintf(`$j('input[data-gr-checkcol]').prop('checked', %t)`, checked)
		//c.parentTable.FormBase().Response().ExecuteJavaScript(js, page.PriorityStandard)
		c.parentTable.GetForm().Response().ExecuteSelectorFunction(`input[data-gr-checkcol]`, `prop`, page.PriorityStandard, `checked`, checked)

	} else {
		// Fire javascript to check all visible and trigger a change
		if checked {
			c.parentTable.GetForm().Response().ExecuteSelectorFunction(`input[data-gr-checkcol]:not(:checked)`, `trigger`, page.PriorityStandard, `click`)

		} else {
			c.parentTable.GetForm().Response().ExecuteSelectorFunction(`input[data-gr-checkcol]:checked`, `trigger`, page.PriorityStandard, `click`)
		}
	}

}

type CheckboxProvider interface {
	// Id should return a unique id corresponding to the data. It is used to track the checked state of the checkbox.
	ID(data interface{}) string
	// IsChecked should return true if the checkbox corresponding to the row data should initially be checked. After the
	// initial draw, the table will keep track of the state of the checkbox, meaning you do not need to live update your data.
	// If you are using the table just as a selction of items to act on, just return false here.
	IsChecked(data interface{}) bool
	// Attributes returns the attributes that will be applied to the checkbox corresponding to the data row.
	// Use this primarily for providing custom attributes. Return nil if you have no custom attributes.
	Attributes(data interface{}) *html.Attributes
	// If you enable the checkAll box, you can use this to return a map of all the ids and their inital values here. This is
	// mostly helpful if your table is not showing all the rows at once (i.e. you are using a paginator or scroller and
	// only showing a subset of data at one time). If your table is show a checkAll box, and you return nil here, the
	// checkAll will only perform a javascript checkAll, and thus only check the visible items.
	All() map[string]bool
}
