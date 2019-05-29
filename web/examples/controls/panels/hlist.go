package panels

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/control/data"
	"github.com/goradd/goradd/pkg/url"
	"github.com/goradd/goradd/test/browsertest"
	"github.com/goradd/goradd/web/examples/controls"
)

// Example data
type person struct {
	name string
}

type project struct {
	name string
	people []person
}

var projects = []project{
	{"Acme Widget", []person{{"Isaiah"}, {"Alaina"}, {"Gabriel"}}},
	{"Ace Thingers", []person{{"Agustin"}, {"Abby"}, {"Josiah"}, {"Anthony"}}},
	{"Big Business", []person{{"Shannon"}}},
	{"Small Frys", []person{{"April"}, {"Ruben"}, {"Karriyma"}, {"McKenzie"}}},
}


type HListPanel struct {
	Panel

	OList *OrderedList
	UList *UnorderedList

	SubmitAjax      *Button
	SubmitServer    *Button
}

func NewHListPanel(ctx context.Context, parent page.ControlI)  {
	itemList := []ListValue{
		{"First", 1},
		{"Second", 2},
		{"Third", 3},
		{"Fourth", 4},
		{"Fifth", 5},
	}

	p := &HListPanel{}
	p.Panel.Init(p, parent, "HListPanel")

	p.OList = NewOrderedList(p,"orderedList")
	p.OList.SetData(itemList)

	p.UList = NewUnorderedList(p,"unorderedList")
	p.UList.SetDataProvider(p)

	p.SubmitAjax = NewButton(p, "ajaxButton")
	p.SubmitAjax.SetText("Submit Ajax")
	p.SubmitAjax.OnSubmit(action.Ajax(p.ID(), AjaxSubmit))

	p.SubmitServer = NewButton(p, "serverButton")
	p.SubmitServer.SetText("Submit Server")
	p.SubmitServer.OnSubmit(action.Server(p.ID(), ServerSubmit))

}

func (p *HListPanel) BindData(ctx context.Context, s data.DataManagerI) {
	// This is an example of how to populate a hierarchical list using a data binder.
	// One use of this is to query the database, and then walk the results.
	p.UList.Clear()
	for _,proj := range projects {
		listItem := NewListItem(proj.name)
		for _,per := range proj.people {
			listItem.AddItem(per.name)
		}
		p.UList.AddListItems(listItem)
	}
}


func init() {
	controls.RegisterPanel("hlist", "Nested Lists", NewHListPanel, 8)

	//browsertest.RegisterTestFunction("Select List Ajax Submit", testHListAjaxSubmit)
	//browsertest.RegisterTestFunction("Select List Server Submit", testHListServerSubmit)
}

// testPlain exercises the plain text box
func testHListAjaxSubmit(t *browsertest.TestForm)  {
	var myUrl = url.NewBuilder(controlsFormPath).SetValue("control", "HList").SetValue("testing", 1).String()
	f := t.LoadUrl(myUrl)

	testHListSubmit(t, f, f.Page().GetControl("ajaxButton"))

	t.Done("Complete")
}

func testHListServerSubmit(t *browsertest.TestForm)  {
	var myUrl = url.NewBuilder(controlsFormPath).SetValue("control", "HList").SetValue("testing", 1).String()
	f := t.LoadUrl(myUrl)

	testHListSubmit(t, f, f.Page().GetControl("serverButton"))

	t.Done("Complete")
}

// testHListSubmit does a variety of submits using the given button. We use this to double check the various
// results we might get after a submission, as well as nsure that the ajax and server submits produce
// the same results.
func testHListSubmit(t *browsertest.TestForm, f page.FormI, btn page.ControlI) {

	// For testing purposes, we need to use the id of the list item, rather than the value of the list item,
	// since that is what is presented in the html.
	//select1 := f.Page().GetControl("orderedList").(*OrderedList)
}
