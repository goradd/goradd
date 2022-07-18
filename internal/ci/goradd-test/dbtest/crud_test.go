package dbtest

import (
	"github.com/stretchr/testify/assert"
	"goradd-project/gen/goradd/model"
	"goradd-project/gen/goradd/model/node"
	model2 "goradd-project/gen/goraddUnit/model"
	node2 "goradd-project/gen/goraddUnit/model/node"

	//model2 "goradd-project/gen/goraddUnit/model"
	//node2 "goradd-project/gen/goraddUnit/model/node"
	"testing"
)

func TestInsertNotInitialized(t *testing.T) {
	ctx := getContext()
	person := model.NewPerson()

	assert.Panics(t, func() {
		person.Save(ctx)
	}, "Not initialized, but required fields will panic on insert")

}

func TestCrudForwardOneManyNull(t *testing.T) {
	ctx := getContext()
	address := model.NewAddress()
	address.SetStreet("1 Center St")
	address.SetCity("Mendenhall")

	person := model.NewPerson()
	person.SetFirstName("John")
	person.SetLastName("Perkins")
	address.SetPerson(person)
	address.Save(ctx)

	address2 := model.LoadAddress(ctx, address.ID(), node.Address().Person())

	assert.Equal(t, "Mendenhall", address2.City())
	assert.Equal(t, "Perkins", address2.Person().LastName())

	address2.Person().SetFirstName("Derek")
	address2.Save(ctx)

	address3 := model.LoadAddress(ctx, address2.ID(), node.Address().Person())
	assert.Equal(t, "Derek", address3.Person().FirstName())

	detachedPersonID := address3.PersonID()

	person2 := model.NewPerson()
	person2.SetFirstName("Jim")
	person2.SetLastName("Gordon")
	address3.SetPerson(person2)
	address3.Save(ctx)

	address4 := model.LoadAddress(ctx, address3.ID(), node.Address().Person())
	assert.Equal(t, "Gordon", address4.Person().LastName())

	address4.Person().Delete(ctx) // should delete address too

	address5 := model.LoadAddress(ctx, address.ID(), node.Address().Person())
	assert.Nil(t, address5)

	person3 := model.LoadPerson(ctx, detachedPersonID)
	person3.Delete(ctx)
}

func TestCrudReverseOneManyCascade(t *testing.T) {
	ctx := getContext()
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

	assert.NotEqual(t, "", person.Addresses()[0].ID(), "New reference has a new id")
	assert.Equal(t, "Eisleben", person.Address(person.Addresses()[0].ID()).City(), "New references are indexed")

	person2 := model.LoadPerson(ctx, person.ID(), node.Person().Addresses())

	assert.Equal(t, "Martin", person2.FirstName())
	assert.Equal(t, "All Saints Church", person2.Addresses()[1].Street())

	addrID := person2.Addresses()[1].ID()

	person2.Addresses()[1].SetStreet("Here and There")
	person2.Save(ctx)

	person3 := model.LoadPerson(ctx, person.ID(), node.Person().Addresses())
	assert.Equal(t, "Here and There", person3.Addresses()[1].Street())

	address3 := model.NewAddress()
	address3.SetStreet("Wartburg Castle")
	address3.SetCity("Eisenach")

	person3.SetAddresses([]*model.Address{address1, address3})
	person3.Save(ctx)

	address4 := model.LoadAddress(ctx, person2.Addresses()[1].ID())
	assert.Nil(t, address4, "Old address was deleted")

	person4 := model.LoadPerson(ctx, person.ID(), node.Person().Addresses())
	assert.Equal(t, "Eisleben", person4.Addresses()[0].City(), "Previous address remained untouched.")
	assert.Equal(t, "Eisenach", person4.Addresses()[1].City(), "New address was added.")

	person4.Delete(ctx)

	address5 := model.LoadAddress(ctx, addrID)
	assert.Nil(t, address5, "Associated addresses were deleted too")
}

func TestCrudForwardOneManySetNull(t *testing.T) {
	ctx := getContext()
	f := model2.NewForwardNull()
	f.SetName("testForward")

	r := model2.NewReverse()
	r.SetName("testReverse")
	f.SetReverse(r)

	f.Save(ctx)

	f2 := model2.LoadForwardNull(ctx, f.ID(), node2.ForwardNull().Reverse())

	assert.Equal(t, "testForward", f2.Name())
	assert.Equal(t, "testReverse", f2.Reverse().Name())

	f2.Reverse().SetName("testReverse2")
	f2.Save(ctx)

	f3 := model2.LoadForwardNull(ctx, f2.ID(), node2.ForwardNull().Reverse())
	assert.Equal(t, "testReverse2", f3.Reverse().Name())
	detached := f2.ReverseID()

	r2 := model2.NewReverse()
	r2.SetName("testReverse3")
	f3.SetReverse(r2)
	f3.Save(ctx)

	f4 := model2.LoadForwardNull(ctx, f3.ID(), node2.ForwardNull().Reverse())
	assert.Equal(t, "testReverse3", f4.Reverse().Name())

	f4.Reverse().Delete(ctx)

	f5 := model2.LoadForwardNull(ctx, f4.ID(), node2.ForwardNull().Reverse())
	assert.Nil(t, f5.Reverse())
	f5.Delete(ctx)

	r3 := model2.LoadReverse(ctx, detached)
	r3.Delete(ctx)
}

func TestCrudReverseOneManySetNull(t *testing.T) {
	ctx := getContext()

	r := model2.NewReverse()
	r.SetName("testReverse")

	f1 := model2.NewForwardNull()
	f1.SetName("testForward1")

	f2 := model2.NewForwardNull()
	f2.SetName("testForward2")

	r.SetForwardNulls([]*model2.ForwardNull{f1, f2})
	r.Save(ctx)

	r2 := model2.LoadReverse(ctx, r.ID(), node2.Reverse().ForwardNulls())

	assert.Equal(t, "testReverse", r2.Name())
	assert.Equal(t, "testForward2", r2.ForwardNulls()[1].Name())

	r2.ForwardNulls()[1].SetName("Other")
	r2.Save(ctx)

	r3 := model2.LoadReverse(ctx, r.ID(), node2.Reverse().ForwardNulls())
	assert.Equal(t, r3.ForwardNulls()[1].Name(), "Other")

	f3 := model2.NewForwardNull()
	f3.SetName("testForward3")

	r2.SetForwardNulls([]*model2.ForwardNull{f2, f3})
	r2.Save(ctx)

	f4 := model2.LoadForwardNull(ctx, r.ForwardNulls()[0].ID())
	assert.NotNil(t, f4, "Old forward reference was NOT deleted")
	assert.True(t, f4.ReverseIDIsNull(), "Old forward reference was set to NULL")

	r4 := model2.LoadReverse(ctx, r2.ID(), node2.Reverse().ForwardNulls())
	assert.Equal(t, "Other", r4.ForwardNulls()[0].Name())
	assert.Equal(t, "testForward3", r4.ForwardNulls()[1].Name())

	f4.Delete(ctx)

	id1 := r4.ForwardNulls()[0].ID()
	id2 := r4.ForwardNulls()[1].ID()
	r4.Delete(ctx)

	f5 := model2.LoadForwardNull(ctx, id1)
	assert.True(t, f5.ReverseIDIsNull(), "Associated forward reference was set to nil")
	f5.Delete(ctx)
	f5 = model2.LoadForwardNull(ctx, id2)
	assert.True(t, f5.ReverseIDIsNull(), "Associated forward reference was set to nil")
	f5.Delete(ctx)

}

func TestCrudForwardOneManyRestrict(t *testing.T) {
	ctx := getContext()
	f := model2.NewForwardRestrict()
	f.SetName("testForward")

	r := model2.NewReverse()
	r.SetName("testReverse")
	f.SetReverse(r)

	f.Save(ctx)

	f2 := model2.LoadForwardRestrict(ctx, f.ID(), node2.ForwardRestrict().Reverse())

	assert.Equal(t, "testForward", f2.Name())
	assert.Equal(t, "testReverse", f2.Reverse().Name())

	f2.Reverse().SetName("testReverse2")
	f2.Save(ctx)

	f3 := model2.LoadForwardRestrict(ctx, f2.ID(), node2.ForwardRestrict().Reverse())
	assert.Equal(t, "testReverse2", f3.Reverse().Name())
	detached := f2.ReverseID()

	r2 := model2.NewReverse()
	r2.SetName("testReverse3")
	f3.SetReverse(r2)
	f3.Save(ctx)

	f4 := model2.LoadForwardRestrict(ctx, f3.ID(), node2.ForwardRestrict().Reverse())
	assert.Equal(t, "testReverse3", f4.Reverse().Name())

	assert.Panics(t, func() { f4.Reverse().Delete(ctx) })
	id1 := f4.ID()
	id2 := f4.Reverse().ID()

	f5 := model2.LoadForwardRestrict(ctx, id1)
	f5.Delete(ctx)

	r3 := model2.LoadReverse(ctx, id2)
	r3.Delete(ctx)

	r4 := model2.LoadReverse(ctx, detached)
	r4.Delete(ctx)

}

func TestCrudReverseOneManyRestrict(t *testing.T) {
	ctx := getContext()

	r := model2.NewReverse()
	r.SetName("testReverse")

	f1 := model2.NewForwardRestrict()
	f1.SetName("testForward1")

	f2 := model2.NewForwardRestrict()
	f2.SetName("testForward2")

	r.SetForwardRestricts([]*model2.ForwardRestrict{f1, f2})
	r.Save(ctx)

	r2 := model2.LoadReverse(ctx, r.ID(), node2.Reverse().ForwardRestricts())

	assert.Equal(t, "testReverse", r2.Name())
	assert.Equal(t, "testForward2", r2.ForwardRestricts()[1].Name())

	r2.ForwardRestricts()[1].SetName("Other")
	r2.Save(ctx)

	r3 := model2.LoadReverse(ctx, r.ID(), node2.Reverse().ForwardRestricts())
	assert.Equal(t, r3.ForwardRestricts()[1].Name(), "Other")

	f3 := model2.NewForwardRestrict()
	f3.SetName("testForward3")

	r2.SetForwardRestricts([]*model2.ForwardRestrict{f2, f3})
	r2.Save(ctx)

	f4 := model2.LoadForwardRestrict(ctx, r.ForwardRestricts()[0].ID())
	// Note that whether we delete here or set to NULL will depend on whether the forward looking field is nullable.
	// SQL does not really have a way of describing this operation, so we use the nullable to tell us what to do.
	// TODO: Create a unit test for the nullable situation
	assert.Nil(t, f4, "Old forward reference was deleted")

	r4 := model2.LoadReverse(ctx, r2.ID(), node2.Reverse().ForwardRestricts())
	assert.Equal(t, "Other", r4.ForwardRestricts()[0].Name())
	assert.Equal(t, "testForward3", r4.ForwardRestricts()[1].Name())

	assert.Panics(t, func() {
		r4.Delete(ctx)
	})

	r4.ForwardRestricts()[0].Delete(ctx)
	r4.ForwardRestricts()[1].Delete(ctx)
	r4.Delete(ctx)
}

func TestCrudForwardOneOneCascade(t *testing.T) {
	ctx := getContext()
	f := model2.NewForwardCascadeUnique()
	f.SetName("testForward")

	r := model2.NewReverse()
	r.SetName("testReverse")
	f.SetReverse(r)

	f.Save(ctx)

	f2 := model2.LoadForwardCascadeUnique(ctx, f.ID(), node2.ForwardCascadeUnique().Reverse())

	assert.Equal(t, "testForward", f2.Name())
	assert.Equal(t, "testReverse", f2.Reverse().Name(), "Forward looking unique linked record was saved.")

	f2.Reverse().SetName("testReverse2")
	f2.Save(ctx)

	f3 := model2.LoadForwardCascadeUnique(ctx, f2.ID(), node2.ForwardCascadeUnique().Reverse())
	assert.Equal(t, "testReverse2", f3.Reverse().Name(), "Field in forward looking unique linked record was saved.")

	r2 := model2.NewReverse()
	r2.SetName("testReverse3")
	f3.SetReverse(r2)
	f3.Save(ctx)

	f4 := model2.LoadForwardCascadeUnique(ctx, f3.ID(), node2.ForwardCascadeUnique().Reverse())
	assert.Equal(t, "testReverse3", f4.Reverse().Name(), "Forward looking unique linked record was replaced.")
	f2.Reverse().Delete(ctx)
	f4.Reverse().Delete(ctx)

	f5 := model2.LoadForwardCascadeUnique(ctx, f4.ID(), node2.ForwardCascadeUnique().Reverse())
	assert.Nil(t, f5, "Forward looking unique linked cascade record was deleted when linked record was deleted.")
}

func TestCrudReverseOneOneCascade(t *testing.T) {
	ctx := getContext()

	r := model2.NewReverse()
	r.SetName("testReverse")

	f1 := model2.NewForwardCascadeUnique()
	f1.SetName("testForward1")

	r.SetForwardCascadeUnique(f1)
	r.Save(ctx)

	r2 := model2.LoadReverse(ctx, r.ID(), node2.Reverse().ForwardCascadeUnique())

	assert.Equal(t, "testReverse", r2.Name())
	assert.Equal(t, "testForward1", r2.ForwardCascadeUnique().Name())

	r2.ForwardCascadeUnique().SetName("Other")
	r2.Save(ctx)

	r3 := model2.LoadReverse(ctx, r.ID(), node2.Reverse().ForwardCascadeUnique())
	assert.Equal(t, r3.ForwardCascadeUnique().Name(), "Other")

	f3 := model2.NewForwardCascadeUnique()
	f3.SetName("testForward3")

	r2.SetForwardCascadeUnique(f3)
	r2.Save(ctx)

	f4 := model2.LoadForwardCascadeUnique(ctx, r.ForwardCascadeUnique().ID())
	assert.Nil(t, f4, "Old forward reference was deleted")

	r4 := model2.LoadReverse(ctx, r2.ID(), node2.Reverse().ForwardCascadeUnique())
	assert.Equal(t, "testForward3", r4.ForwardCascadeUnique().Name())

	r4.Delete(ctx)
}

func TestCrudForwardOneOneSetNull(t *testing.T) {
	ctx := getContext()
	f := model2.NewForwardNullUnique()
	f.SetName("testForward")

	r := model2.NewReverse()
	r.SetName("testReverse")
	f.SetReverse(r)

	f.Save(ctx)

	f2 := model2.LoadForwardNullUnique(ctx, f.ID(), node2.ForwardNullUnique().Reverse())

	assert.Equal(t, "testForward", f2.Name())
	assert.Equal(t, "testReverse", f2.Reverse().Name(), "Forward looking unique linked record was saved.")

	f2.Reverse().SetName("testReverse2")
	f2.Save(ctx)

	f3 := model2.LoadForwardNullUnique(ctx, f2.ID(), node2.ForwardNullUnique().Reverse())
	assert.Equal(t, "testReverse2", f3.Reverse().Name(), "Field in forward looking unique linked record was saved.")

	r2 := model2.NewReverse()
	r2.SetName("testReverse3")
	f3.SetReverse(r2)
	f3.Save(ctx)

	f4 := model2.LoadForwardNullUnique(ctx, f3.ID(), node2.ForwardNullUnique().Reverse())
	assert.Equal(t, "testReverse3", f4.Reverse().Name(), "Forward looking unique linked record was replaced.")
	f2.Reverse().Delete(ctx)
	f4.Reverse().Delete(ctx)

	f5 := model2.LoadForwardNullUnique(ctx, f4.ID(), node2.ForwardNullUnique().Reverse())
	assert.True(t, f5.ReverseIDIsNull(), "Forward looking unique linked setnull record was set to NULL after delete.")
	f2.Delete(ctx)
}

func TestCrudReverseOneOneSetNull(t *testing.T) {
	ctx := getContext()

	r := model2.NewReverse()
	r.SetName("testReverse")

	f1 := model2.NewForwardNullUnique()
	f1.SetName("testForward1")

	r.SetForwardNullUnique(f1)
	r.Save(ctx)

	r2 := model2.LoadReverse(ctx, r.ID(), node2.Reverse().ForwardNullUnique())

	assert.Equal(t, "testReverse", r2.Name())
	assert.Equal(t, "testForward1", r2.ForwardNullUnique().Name())

	r2.ForwardNullUnique().SetName("Other")
	r2.Save(ctx)

	r3 := model2.LoadReverse(ctx, r.ID(), node2.Reverse().ForwardNullUnique())
	assert.Equal(t, r3.ForwardNullUnique().Name(), "Other")

	f3 := model2.NewForwardNullUnique()
	f3.SetName("testForward3")

	r2.SetForwardNullUnique(f3)
	r2.Save(ctx)

	f4 := model2.LoadForwardNullUnique(ctx, r.ForwardNullUnique().ID())
	assert.True(t, f4.ReverseIDIsNull(), "Old forward reference was set to NULL")

	r4 := model2.LoadReverse(ctx, r2.ID(), node2.Reverse().ForwardNullUnique())
	assert.Equal(t, "testForward3", r4.ForwardNullUnique().Name())

	r4.ForwardNullUnique().Delete(ctx)
	r4.Delete(ctx)
	f4.Delete(ctx)
}

func TestCrudForwardOneOneRestrict(t *testing.T) {
	ctx := getContext()
	f := model2.NewForwardRestrictUnique()
	f.SetName("testForward")

	r := model2.NewReverse()
	r.SetName("testReverse")
	f.SetReverse(r)

	f.Save(ctx)

	f2 := model2.LoadForwardRestrictUnique(ctx, f.ID(), node2.ForwardRestrictUnique().Reverse())

	assert.Equal(t, "testForward", f2.Name())
	assert.Equal(t, "testReverse", f2.Reverse().Name(), "Forward looking unique linked record was saved.")

	f2.Reverse().SetName("testReverse2")
	f2.Save(ctx)

	f3 := model2.LoadForwardRestrictUnique(ctx, f2.ID(), node2.ForwardRestrictUnique().Reverse())
	assert.Equal(t, "testReverse2", f3.Reverse().Name(), "Field in forward looking unique linked record was saved.")

	r2 := model2.NewReverse()
	r2.SetName("testReverse3")
	f3.SetReverse(r2)
	f3.Save(ctx)

	f4 := model2.LoadForwardRestrictUnique(ctx, f3.ID(), node2.ForwardRestrictUnique().Reverse())
	assert.Equal(t, "testReverse3", f4.Reverse().Name(), "Forward looking unique linked record was replaced.")
	f2.Reverse().Delete(ctx)

	assert.Panics(t, func() {
		f4.Reverse().Delete(ctx)
	})
	f4.Delete(ctx)
	f4.Reverse().Delete(ctx)

}

func TestCrudReverseOneOneRestrict(t *testing.T) {
	ctx := getContext()

	r := model2.NewReverse()
	r.SetName("testReverse")

	f1 := model2.NewForwardRestrictUnique()
	f1.SetName("testForward1")

	r.SetForwardRestrictUnique(f1)
	r.Save(ctx)

	r2 := model2.LoadReverse(ctx, r.ID(), node2.Reverse().ForwardRestrictUnique())

	assert.Equal(t, "testReverse", r2.Name())
	assert.Equal(t, "testForward1", r2.ForwardRestrictUnique().Name())

	r2.ForwardRestrictUnique().SetName("Other")
	r2.Save(ctx)

	r3 := model2.LoadReverse(ctx, r.ID(), node2.Reverse().ForwardRestrictUnique())
	assert.Equal(t, r3.ForwardRestrictUnique().Name(), "Other")

	f3 := model2.NewForwardRestrictUnique()
	f3.SetName("testForward3")

	r2.SetForwardRestrictUnique(f3)
	r2.Save(ctx)

	f4 := model2.LoadForwardRestrictUnique(ctx, r.ForwardRestrictUnique().ID())
	// This restrict situation should set null
	assert.True(t, f4.ReverseIDIsNull(), "Old forward reference was set to NULL")

	r4 := model2.LoadReverse(ctx, r2.ID(), node2.Reverse().ForwardRestrictUnique())
	assert.Equal(t, "testForward3", r4.ForwardRestrictUnique().Name())

	r4.ForwardRestrictUnique().Delete(ctx)
	r4.Delete(ctx)
	f4.Delete(ctx)
}

func TestCrudManyMany(t *testing.T) {
	ctx := getContext()
	project := model.NewProject()
	project.SetName("NewProject")
	project.SetNum(100)
	project.SetStatusType(model.ProjectStatusTypeOpen)

	p1 := model.NewPerson()
	p1.SetFirstName("Me")
	p1.SetLastName("You")
	p2 := model.NewPerson()
	p2.SetFirstName("Him")
	p2.SetLastName("Her")

	// Additional test that will insert an item linked to another item that already existed
	p3 := model.LoadPerson(ctx, "1")
	project.SetManager(p3)

	project.SetTeamMembers([]*model.Person{p1, p2})
	project.Save(ctx)

	assert.Equal(t, p2.ID(), project.TeamMember(p2.ID()).ID())

	project2 := model.LoadProject(ctx, project.ID(), node.Project().TeamMembers(), node.Project().Manager())
	assert.Len(t, project2.TeamMembers(), 2)
	assert.Equal(t, "Him", project2.TeamMember(p2.ID()).FirstName())
	assert.Equal(t, "Doe", project2.Manager().LastName())

	p1.Delete(ctx)
	p2.Delete(ctx)
	project2.Delete(ctx)

	p4 := model.LoadPerson(ctx, "1")
	assert.NotNil(t, p4)
}

// TestUniqueHas will test both the Has... function that is generated by unique indexed values,
// and the ability to work with tables whose primary key is not auto generated, but rather assigned
// at save time.
func TestUniqueHas(t *testing.T) {
	ctx := getContext()
	assert.False(t, model2.HasDoubleIndexByFieldIntFieldString(ctx, 2, "we"))

	d := model2.NewDoubleIndex()
	d.SetID(1)
	d.SetFieldInt(1)
	d.SetFieldString("we")
	d.Save(ctx) // test insert

	d.SetFieldInt(2)
	d.Save(ctx) // test update

	assert.True(t, model2.HasDoubleIndexByFieldIntFieldString(ctx, 2, "we"))

	d.Delete(ctx)
}

func TestCrudIntKey(t *testing.T) {
	ctx := getContext()
	g := model.NewGift()

	g.SetNumber(4)
	g.SetName("Calling birds")
	g.Save(ctx)

	g2 := model.LoadGift(ctx, 4)
	assert.Equal(t, "Calling birds", g2.Name())
	g2.SetNumber(5)
	g2.SetName("Gold rings")
	g2.Save(ctx)

	g3 := model.LoadGift(ctx, 4)
	assert.Nil(t, g3)

	g3 = model.LoadGift(ctx, 5)
	assert.Equal(t, "Gold rings", g3.Name())

	g3.Delete(ctx)
	g3 = model.LoadGift(ctx, 5)
	assert.Nil(t, g3)
}
