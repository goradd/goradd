package panels

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
	. "github.com/goradd/goradd/pkg/page/control/list"
	"github.com/goradd/goradd/pkg/url"
	"github.com/goradd/goradd/test/browsertest"
)

// Example data
type person struct {
	name string
}

type hlistProject struct {
	name   string
	people []person
}

var hlistProjects = []hlistProject{
	{"Acme Widget", []person{{"Isaiah"}, {"Alaina"}, {"Gabriel"}}},
	{"Ace Thingers", []person{{"Agustin"}, {"Abby"}, {"Josiah"}, {"Anthony"}}},
	{"Big Business", []person{{"Shannon"}}},
	{"Small Frys", []person{{"April"}, {"Ruben"}, {"Karriyma"}, {"McKenzie"}}},
}

type HListPanel struct {
	Panel
}

func NewHListPanel(ctx context.Context, parent page.ControlI) {
	p := new(HListPanel)
	p.Init(p, ctx, parent, "HListPanel")
}

func (p *HListPanel) Init(self any, ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(self, parent, id)
	itemList := []ListValue{
		{"First", 1},
		{"Second", 2},
		{"Third", 3},
		{"Fourth", 4},
		{"Fifth", 5},
	}
	p.AddControls(ctx,
		OrderedListCreator{
			ID:    "orderedList",
			Items: itemList,
		},
		UnorderedListCreator{
			ID:             "unorderedList",
			DataProviderID: "HListPanel",
		},
	)
}

func (p *HListPanel) BindData(ctx context.Context, s DataManagerI) {
	// This is an example of how to populate a hierarchical list using a data binder.
	// One use of this is to query the database, and then walk the results.
	ulist := s.(*UnorderedList)
	ulist.Clear()
	for _, proj := range hlistProjects {
		listItem := NewItem(proj.name)
		for _, per := range proj.people {
			listItem.Add(per.name)
		}
		ulist.AddItems(listItem)
	}
}

func init() {
	//browsertest.RegisterTestFunction("Select List Ajax Submit", testHListAjaxSubmit)
	//browsertest.RegisterTestFunction("Select List Server Submit", testHListServerSubmit)
}

// testPlain exercises the plain text box
func testHListAjaxSubmit(t *browsertest.TestForm) {
	var myUrl = url.NewBuilder(controlsFormPath).SetValue("control", "HList").SetValue("testing", 1).String()
	t.LoadUrl(myUrl)

	testHListSubmit(t, "ajaxButton")

	t.Done("Complete")
}

func testHListServerSubmit(t *browsertest.TestForm) {
	var myUrl = url.NewBuilder(controlsFormPath).SetValue("control", "HList").SetValue("testing", 1).String()
	t.LoadUrl(myUrl)

	testHListSubmit(t, "serverButton")

	t.Done("Complete")
}

// testHListSubmit does a variety of submits using the given button. We use this to double check the various
// results we might get after a submission, as well as nsure that the ajax and server submits produce
// the same results.
func testHListSubmit(t *browsertest.TestForm, btnID string) {

	// For testing purposes, we need to use the id of the list item, rather than the value of the list item,
	// since that is what is presented in the html.
	//select1 := f.Page().GetControl("orderedList").(*OrderedList)
}

func init() {
	page.RegisterControl(&HListPanel{})
}
