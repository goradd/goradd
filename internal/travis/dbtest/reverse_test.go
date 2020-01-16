package dbtest

import (
	. "github.com/goradd/goradd/pkg/orm/op"
	"github.com/stretchr/testify/assert"
	"goradd-project/gen/goradd/model"
	"goradd-project/gen/goradd/model/node"
	"testing"
)

func TestReverse2(t *testing.T) {
	ctx := getContext()
	people := model.QueryPeople(ctx).
		Join(node.Person().ProjectsAsManager()).
		OrderBy(node.Person().ID(), node.Person().ProjectsAsManager().Name()).
		Where(IsNotNull(node.Person().ProjectsAsManager().ID())). // Filter out people who are not managers
		Select(node.Person().ProjectsAsManager().Name()).
		Load(ctx)

	if len(people[2].ProjectsAsManager()) != 2 {
		t.Error("Did not find 2 ProjectsAsManagers.")
	}

	assert.Len(t, people[2].ProjectsAsManager(), 2)
	assert.Equal(t, people[2].ProjectsAsManager()[0].Name(), "ACME Payment System")
	assert.True(t, people[2].ProjectsAsManager()[0].IDIsValid())
	assert.False(t, people[2].ProjectsAsManager()[0].NumIsValid())
}

// Complex test finding all the team members of all the projects a person is managing, ordering by last name
func TestReverseMany(t *testing.T) {
	ctx := getContext()
	people := model.QueryPeople(ctx).
		OrderBy(node.Person().ID(), node.Person().ProjectsAsManager().TeamMembers().LastName(), node.Person().ProjectsAsManager().TeamMembers().FirstName()).
		Select(node.Person().ProjectsAsManager().TeamMembers().FirstName(), node.Person().ProjectsAsManager().TeamMembers().LastName()).
		Load(ctx)

	names := []string{}
	for _, p := range people[0].ProjectsAsManager()[0].TeamMembers() {
		names = append(names, p.FirstName()+" "+p.LastName())
	}
	names2 := []string{
		"John Doe",
		"Mike Ho",
		"Samantha Jones",
		"Jennifer Smith",
		"Wendy Smith",
	}
	assert.Equal(t, names2, names)

	names = []string{}
	for _, pr := range people[6].ProjectsAsManager() {
		for _, p := range pr.TeamMembers() {
			names = append(names, p.FirstName()+" "+p.LastName())
		}
	}
	assert.Len(t, names, 12) // Includes duplicates. If we ever get Distinct to manually remove duplicates, we should fix this.
}

func TestReverseManyExpansion(t *testing.T) {
	ctx := getContext()
	// Test an intermediate expansion
	people := model.QueryPeople(ctx).
		Join(node.Person().ProjectsAsManager().TeamMembers()).
		OrderBy(node.Person().ID(), node.Person().ProjectsAsManager().TeamMembers().LastName(), node.Person().ProjectsAsManager().TeamMembers().FirstName()).
		Expand(node.Person().ProjectsAsManager()).
		Load(ctx)

	names2 := []string{
		"John Doe",
		"Mike Ho",
		"Samantha Jones",
		"Jennifer Smith",
		"Wendy Smith",
	}
	names := []string{}
	for _, p := range people[0].ProjectsAsManager()[0].TeamMembers() {
		names = append(names, p.FirstName()+" "+p.LastName())
	}
	assert.Equal(t, names2, names)

	names = []string{}
	for _, pr := range people[6].ProjectsAsManager() {
		for _, p := range pr.TeamMembers() {
			names = append(names, p.FirstName()+" "+p.LastName())
		}
	}

	// Should only select first group
	names4 := []string{
		"Brett Carlisle",
		"John Doe",
		"Samantha Jones",
		"Jacob Pratt",
		"Kendall Public",
		"Ben Robinson",
		"Alex Smith",
	}
	assert.Equal(t, names4, names)

}

func TestUniqueReverse(t *testing.T) {
	ctx := getContext()
	person := model.QueryPeople(ctx).
		Where(Equal(node.Person().LastName(), "Doe")).
		Get(ctx)
	assert.Nil(t, person.Login())

	person = model.QueryPeople(ctx).
		Where(Equal(node.Person().LastName(), "Doe")).
		Join(node.Person().Login()).
		Load(ctx)[0]
	assert.Equal(t, "jdoe", person.Login().Username())
}


// TestReverseReferenceManySave is testing save and delete for a reverse reference that cannot be null.
func TestReverseReferenceManySave(t *testing.T) {
	ctx := getContext()
	// Test insert
	person := model.NewPerson()
	person.SetFirstName("Sam")
	person.SetLastName("I Am")

	addr1 := model.NewAddress()
	addr1.SetCity("Here")
	addr1.SetStreet("There")

	addr2 := model.NewAddress()
	addr2.SetCity("Near")
	addr2.SetStreet("Far")

	person.SetAddresses([]*model.Address{
		addr1, addr2,
	})

	person.Save(ctx)
	id := person.ID()

	addr1Id := addr1.ID()
	assert.NotEmpty(t, addr1Id)

	addr3 := person.Address(addr1Id)
	assert.Equal(t, "There", addr3.Street(), "Successfully attached the new addresses onto the person object.")

	person2 := model.LoadPerson(ctx, id, node.Person().Addresses())

	assert.Equal(t, "Sam", person2.FirstName(), "Retrieved the correct person")
	assert.Equal(t, 2, len(person2.Addresses()), "Retrieved the addresses attached to the person")

	person2.Delete(ctx)

	person3 := model.LoadPerson(ctx, id, node.Person().Addresses())
	assert.Nil(t, person3, "Successfully deleted the new person")

	addr4 := model.LoadAddress(ctx, addr1Id)
	assert.Nil(t, addr4, "Successfully deleted the address attached to the person")

}

// Testing a reverse reference with a unique index, which will cause a one-to-one relationship.
// This tests save and delete
func TestReverseReferenceUniqueSave(t *testing.T) {
	ctx := getContext()

	person := model.NewPerson()
	person.SetFirstName("Sam")
	person.SetLastName("I Am")

	e1 := model.NewEmployeeInfo()
	e1.SetEmployeeNumber(12345)
	person.SetEmployeeInfo(e1)

	person.Save(ctx)
	id := person.ID()

	e1Id := e1.ID()
	assert.NotEmpty(t, e1Id)

	e2 := person.EmployeeInfo()
	assert.Equal(t, e1Id, e2.ID(), "Successfully attached the new employee info object onto the person object.")

	person2 := model.LoadPerson(ctx, id, node.Person().EmployeeInfo())

	assert.Equal(t, "Sam", person2.FirstName(), "Retrieved the correct person")
	assert.Equal(t, e1Id, person2.EmployeeInfo().ID(), "Retrieved the employee info attached to the person")

	person2.Delete(ctx)

	person3 := model.LoadPerson(ctx, id, node.Person().EmployeeInfo())
	assert.Nil(t, person3, "Successfully deleted the new person")

	e4 := model.LoadEmployeeInfo(ctx, e1Id)
	assert.Nil(t, e4, "Successfully deleted the employee info attached to the person")

}

