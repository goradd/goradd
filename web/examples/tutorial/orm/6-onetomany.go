package orm

import (
	"context"
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/sys"
	"github.com/goradd/goradd/web/examples/gen/goradd/model"
	"github.com/goradd/goradd/web/examples/tutorial"
	"path/filepath"
)

type OneManyPanel struct {
	Panel
}

func NewOneManyPanel(ctx context.Context, parent page.ControlI) page.ControlI {
	p := &OneManyPanel{}
	p.Self = p
	p.Init(ctx, parent, "")
	return p
}

func (p *OneManyPanel) Init(ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(parent, id)
}


func init() {
	page.RegisterControl(&OneManyPanel{})

	dir := sys.SourceDirectory()
	tutorial.RegisterTutorialPage("orm", 6, "onetomany", "One-to-Many References", NewOneManyPanel,
		[]string {
			sys.SourcePath(),
			filepath.Join(dir, "6-onetomany.tpl.got"),
		})
}

func (p *OneManyPanel) addRecord(ctx context.Context) string {
	address := model.NewAddress()
	address.SetStreet("My Street")
	address.SetCity("Bucharest")

	person := model.NewPerson()
	person.SetFirstName("Richard")
	person.SetLastName("Wurmbrand")
	db.ExecuteTransaction(ctx, model.Database(), func() {
		person.Save(ctx)
		id := person.ID()
		address.SetPersonID(id)
		address.Save(ctx)
	})
	return address.ID()
}

func (p *OneManyPanel) addRecordSimpler(ctx context.Context) string {
	address := model.NewAddress()
	address.SetStreet("1 Center St")
	address.SetCity("Mendenhall")

	person := model.NewPerson()
	person.SetFirstName("John")
	person.SetLastName("Perkins")
	address.SetPerson(person)
	address.Save(ctx)
	return address.ID()
}


func (p *OneManyPanel) addMany(ctx context.Context) string {
	person := model.NewPerson()
	person.SetFirstName("Martin")
	person.SetLastName("Luther")

	address1 := model.NewAddress()
	address1.SetStreet("My Street")
	address1.SetCity("Eisleben")

	address2 := model.NewAddress()
	address2.SetStreet("All Saints Church")
	address2.SetCity("Wittenburg")

	person.SetAddresses([]*model.Address{address1, address2})
	person.Save(ctx)
	return person.ID()
}
