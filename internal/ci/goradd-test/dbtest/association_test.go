package dbtest

import (
	. "github.com/goradd/goradd/pkg/orm/op"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"goradd-project/gen/goradd/model"
	"goradd-project/gen/goradd/model/node"
	"testing"
)

func TestMany2(t *testing.T) {

	ctx := getContext()

	// All People Who Are on a Project Managed by Karen Wolfe (Person ID #7)
	people := model.QueryPeople(ctx).
		OrderBy(node.Person().LastName(), node.Person().FirstName()).
		Where(Equal(node.Person().Projects().Manager().LastName(), "Wolfe")).
		Distinct().
		Select(node.Person().LastName(), node.Person().FirstName()).
		Load()

	names := []string{}
	for _, p := range people {
		names = append(names, p.FirstName()+" "+p.LastName())
	}
	names2 := []string{
		"Brett Carlisle",
		"John Doe",
		"Samantha Jones",
		"Jacob Pratt",
		"Kendall Public",
		"Ben Robinson",
		"Alex Smith",
		"Wendy Smith",
		"Karen Wolfe",
	}

	assert.Equal(t, names2, names)
}

func TestManyTypes(t *testing.T) {
	ctx := getContext()

	// All people who are inactive
	people := model.QueryPeople(ctx).
		OrderBy(node.Person().LastName(), node.Person().FirstName()).
		Where(Equal(node.Person().PersonTypes().ID(), model.PersonTypeInactive)).
		Distinct().
		Select(node.Person().LastName(), node.Person().FirstName()).
		Load()

	names := []string{}
	for _, p := range people {
		names = append(names, p.FirstName()+" "+p.LastName())
	}
	names2 := []string{
		"Linda Brady",
		"John Doe",
		"Ben Robinson",
	}
	assert.Equal(t, names2, names)
}

func TestManySelect(t *testing.T) {
	ctx := getContext()

	people := model.QueryPeople(ctx).
		OrderBy(node.Person().LastName(), node.Person().FirstName(), node.Person().Projects().Name()).
		Where(Equal(node.Person().Projects().Manager().LastName(), "Wolfe")).
		Select(node.Person().LastName(), node.Person().FirstName(), node.Person().Projects().Name()).
		Load()

	person := people[0]
	projects := person.Projects()
	name := projects[0].Name()

	assert.Equal(t, "ACME Payment System", name)
}

func Test2Nodes(t *testing.T) {
	ctx := getContext()
	milestones := model.QueryMilestones(ctx).
		Join(node.Milestone().Project().Manager()).
		Where(Equal(node.Milestone().ID(), 1)). // Filter out people who are not managers
		Load()

	assert.True(t, milestones[0].NameIsValid(), "Milestone 1 has a name")
	assert.Equal(t, "Milestone A", milestones[0].Name(), "Milestone 1 has name of Milestone A")
	assert.True(t, milestones[0].Project().NameIsValid(), "Project 1 should have a name")
	assert.True(t, milestones[0].Project().Manager().FirstNameIsValid(), "Person 7 has a name")
	assert.Equal(t, "Karen", milestones[0].Project().Manager().FirstName(), "Person 7 has first name of Karen")
}

func TestForwardMany(t *testing.T) {
	ctx := getContext()
	milestones := model.QueryMilestones(ctx).
		Join(node.Milestone().Project().TeamMembers()).
		OrderBy(node.Milestone().Project().TeamMembers().LastName(), node.Milestone().Project().TeamMembers().FirstName()).
		Where(Equal(node.Milestone().ID(), 1)). // Filter out people who are not managers
		Load()

	names := []string{}
	for _, p := range milestones[0].Project().TeamMembers() {
		names = append(names, p.FirstName()+" "+p.LastName())
	}
	names2 := []string{
		"Samantha Jones",
		"Kendall Public",
		"Alex Smith",
		"Wendy Smith",
		"Karen Wolfe",
	}
	assert.Equal(t, names2, names)

}

func TestManyForward(t *testing.T) {
	ctx := getContext()
	people := model.QueryPeople(ctx).
		OrderBy(node.Person().ID(), node.Person().Projects().Name()).
		Select(node.Person().Projects().Manager().FirstName(), node.Person().Projects().Manager().LastName()).
		Load()

	names := []string{}
	var p *model.Project
	for _, p = range people[0].Projects() {
		names = append(names, p.Manager().FirstName()+" "+p.Manager().LastName())
	}
	names2 := []string{
		"Karen Wolfe",
		"John Doe",
	}
	assert.Equal(t, names2, names)

}

func TestConditionalJoin(t *testing.T) {
	ctx := getContext()

	projects := model.QueryProjects(ctx).
		OrderBy(node.Project().Name()).
		Join(node.Project().Manager(), Equal(node.Project().Manager().LastName(), "Wolfe")).
		Join(node.Project().TeamMembers(), Equal(node.Project().TeamMembers().LastName(), "Smith")).
		Load()

	// Reverse references
	people := model.QueryPeople(ctx).
		Join(node.Person().Addresses(), Equal(node.Person().Addresses().City(), "New York")).
		Join(node.Person().ProjectsAsManager(), Equal(node.Person().ProjectsAsManager().StatusID(), model.ProjectStatusOpen)).
		Join(node.Person().ProjectsAsManager().Milestones()).
		Join(node.Person().Login(), Like(node.Person().Login().Username(), "b%")).
		OrderBy(node.Person().LastName(), node.Person().FirstName(), node.Person().ProjectsAsManager().Name()).
		Load()

	assert.Equal(t, "John", people[2].FirstName(), "John Doe is the 3rd Person.")
	assert.Len(t, people[2].ProjectsAsManager(), 1, "John Doe manages 1 Project.")
	assert.Len(t, people[2].ProjectsAsManager()[0].Milestones(), 1, "John Doe has 1 Milestone")

	// Groups that are not expanded by the conditional join are still created as empty arrays. NoSql databases will need to do this too.
	// This makes it a little easier to write code that uses it, becuase you don't have to test for nil
	assert.Len(t, people[0].ProjectsAsManager(), 0)

	// Check parallel reverse reference with condition
	assert.Len(t, people[7].Addresses(), 2, "Ben Robinson has 2 Addresses")
	assert.Len(t, people[2].Addresses(), 0, "John Doe has no Addresses")

	// Reverse reference unique
	assert.Equal(t, "brobinson", people[7].Login().Username(), "Ben Robinson's Login was selected")
	assert.Nil(t, people[2].Login(), "John Doe's Login was not selected")

	// Forward reference
	assert.Nil(t, projects[2].Manager(), "")
	assert.Equal(t, projects[0].Manager().FirstName(), "Karen")

	// Many-many
	assert.Len(t, projects[3].TeamMembers(), 2, "Project 4 has 2 team members with last name Smith")
	assert.Equal(t, "Smith", projects[3].TeamMembers()[0].LastName(), "The first team member from project 4 has a last name of smith")
}

func TestConditionalExpand(t *testing.T) {
	ctx := getContext()

	// Reverse references
	people := model.QueryPeople(ctx).
		Join(node.Person().Addresses(), Equal(node.Person().Addresses().City(), "Mountain View")).
		Join(node.Person().ProjectsAsManager(), Like(node.Person().ProjectsAsManager().Name(), "%Website%")).
		Join(node.Person().ProjectsAsManager().Milestones()).
		Expand(node.Person().ProjectsAsManager()).
		OrderBy(node.Person().LastName(), node.Person().FirstName(), node.Person().ProjectsAsManager().Name()).
		Load()

	assert.Equal(t, "Karen", people[11].FirstName(), "Karen is the 12th Person.")
	assert.Len(t, people[11].ProjectsAsManager(), 1, "Karen Wolfe selected 1 project.")
	assert.Len(t, people[11].Addresses(), 1, "Karen Wolfe has 1 Address.")

	assert.Len(t, people[2].ProjectsAsManager(), 0, "John Doe selected no projects.")
	assert.Len(t, people[2].Addresses(), 0, "John Doe has no Addresses")

}

func TestSelectByID(t *testing.T) {
	ctx := getContext()

	projects := model.QueryProjects(ctx).
		OrderBy(node.Project().Name().Descending()).
		Load()

	require.Len(t, projects, 4)
	id := projects[3].ID()

	// Reverse references
	people := model.QueryPeople(ctx).
		Join(node.Person().ProjectsAsManager()).
		Where(Equal(node.Person().LastName(), "Wolfe")).
		Load()

	p := people[0]
	require.NotNil(t, p)
	m := p.ProjectAsManager(id)
	require.NotNil(t, m, "Could not fine project as manager: "+id)
	assert.Equal(t, m.Name(), "ACME Payment System")
}

func Test2ndLoad(t *testing.T) {
	ctx := getContext()
	projects := model.QueryProjects(ctx).
		OrderBy(node.Project().Manager().FirstName()).
		Load()

	mgr := projects[0].LoadManager(ctx)
	assert.Equal(t, "Doe", mgr.LastName())

}

func TestSetPrimaryKeys(t *testing.T) {
	ctx := getContext()
	person := model.LoadPerson(ctx, "1", node.Person().Projects())
	assert.Len(t, person.Projects(), 2)
	person.SetProjectPrimaryKeys([]string{"1", "2", "3"})
	person.Save(ctx)

	person2 := model.LoadPerson(ctx, "1", node.Person().Projects())
	assert.Len(t, person2.Projects(), 3)

	person2.SetProjectPrimaryKeys([]string{"3", "4"})
	person2.Save(ctx)
}
