package column

import (
	"context"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	"strings"
	"github.com/spekary/goradd/page/control/control_base"
)

// Checkbox is a column that contains a checkbox. You must provide it a CheckboxProvider to connect ids and default data
// to the checkbox. Use Changes() to get the list of checkbox ids that have changed since the list was initially drawn.
type Checkbox struct {
	ColumnBase
	showCheckAll bool
	parentTableId string
	checkboxer CheckboxProvider
	changes map[string]bool			// records changes
}

type CheckboxI interface {
	ColumnI
	CheckboxAttributes(data interface{}) *html.Attributes

}

// NewChecboxColumn creates a new table column that contains a checkbox. You must provide the id of the parent table,
// add a CheckboxProvider which will connect checkbox states to data states
//
// The column will keep track of what checkboxes have been clicked and the new values. Call Changes() to get those
// changes. Or, if you are recording your changes in real time, attach a CheckboxColumnClick event to the table.
func NewCheckboxColumn(parentTableId string, p CheckboxProvider) *Checkbox {
	if p == nil {
		panic("A checkbox attribute provider is required.")
	}
	if parentTableId == "" {
		panic("The parent table id is required.")
	}

	i := Checkbox{parentTableId:parentTableId, checkboxer:p}
	i.Init()
	return &i
}

func (c *Checkbox) Init() {
	c.ColumnBase.Init(c)
	c.dontEscape = true
	c.changes = map[string]bool {}
}

func (c *Checkbox) This() CheckboxI {
	return c.Self.(CheckboxI)
}

func (c *Checkbox) SetShowCheckAll(s bool) {
	c.showCheckAll = s
}

func (c *Checkbox) HeaderCellText(ctx context.Context, row int, col int) string {
	if c.showCheckAll {
		a := c.This().CheckboxAttributes(nil)
		a.Set("type", "checkbox")
		return html.RenderVoidTag("input", a)
	}
	return ""
}

func (c *Checkbox) CheckboxAttributes(data interface{}) *html.Attributes {
	p := c.checkboxer
	a := p.Attributes(data)
	if a == nil {
		a = html.NewAttributes()
	}
	var id string
	if id = p.Id(data); id != "" {
		// TODO: optionally encrypt the id in case its a database id. Difficult since database ids might themselves be large hashes (aka Google data store)
		pubid := c.Id() + "_" + id
		a.Set("id", pubid)
	} else {
		panic("A checkbox id is required.")
	}

	if newVal, ok := c.changes[id]; ok {	// If we have recorded a change, use this value on  refresh.
		if newVal {
			a.Set("checked", "")
		}
	} else if p.IsChecked(data) {	// otherwise, use the data value
		a.Set("checked", "")
	}
	if !a.Has("value") {		// value is required by html
		a.Set("value", "1")
	}
	a.SetDataAttribute("grTrackchanges", "1")
	a.SetDataAttribute("grCheckcol", "1")
	return a
}

func (c *Checkbox) CellText(ctx context.Context, row int, col int, data interface{}) string {
	a := c.This().CheckboxAttributes(data)
	a.Set("type", "checkbox")
	return html.RenderVoidTag("input", a)
}

// Changes returns a map of ids corresponding to checkboxes that have changed. Both true and false values indicate the
// current state of that particular checkbox. Note that if a user checks a box, then checks it again, even though it
// is back to its original value, it will still show up in the changes list.
func (c *Checkbox) Changes() map[string]bool {
	return c.changes
}

// UpdateFormValues will look for changes to our checkboxes and record those changes.
func (c *Checkbox) UpdateFormValues(ctx *page.Context) {
	for k,v := range ctx.CheckableValues() {
		index := strings.LastIndexAny(k, "_")
		if index > 0 {
			column_id := k[:index]
			check_id := k[index+1:]

			if column_id == c.Id() {
				c.changes[check_id] = control_base.ConvertToBool(v)
			}
		}
	}
}


type CheckboxProvider interface {
	// Id should return a unique id corresponding to the data. It is used to track the checked state of the checkbox.
	Id(data interface{}) string
	// IsChecked should return true if the checkbox corresponding to the row data should initially be checked. After the
	// initial draw, the column will keep track of the state of the checkbox, meaning you do not need to live update your data.
	// If you are using the column to select items to act on, just return false here.
	IsChecked(data interface{}) bool
	// Attributes returns the attributes that will be applied to the checkbox corresponding to the data row.
	// Use this primarily for providing custom attributes. It is OK to return nil.
	Attributes(data interface{})  *html.Attributes
}
