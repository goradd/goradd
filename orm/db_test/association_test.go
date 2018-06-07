package db

import (
	"context"
	. "github.com/spekary/goradd/orm/op"
	"github.com/stretchr/testify/assert"
	"goradd/gen/goradd/model"
	"goradd/gen/goradd/model/node"
	"testing"
)

func TestMany2(t *testing.T) {

	ctx := context.Background()

	// All People Who Are on a Project Managed by Karen Wolfe (Person ID #7)
	people := model.QueryPeople().
		OrderBy(node.Person().LastName(), node.Person().FirstName()).
		Where(Equal(node.Person().ProjectsAsTeamMember().Manager().LastName(), "Wolfe")).
		Distinct().
		Select(node.Person().LastName(), node.Person().FirstName()).
		Load(ctx)

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
	ctx := context.Background()

	// All people who are inactive
	people := model.QueryPeople().
		OrderBy(node.Person().LastName(), node.Person().FirstName()).
		Where(Equal(node.Person().PersonTypes().ID(), model.PERSON_TYPE_INACTIVE)).
		Distinct().
		Select(node.Person().LastName(), node.Person().FirstName()).
		Load(ctx)

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
	ctx := context.Background()

	people := model.QueryPeople().
		OrderBy(node.Person().LastName(), node.Person().FirstName(), node.Person().ProjectsAsTeamMember().Name()).
		Where(Equal(node.Person().ProjectsAsTeamMember().Manager().LastName(), "Wolfe")).
		Select(node.Person().LastName(), node.Person().FirstName(), node.Person().ProjectsAsTeamMember().Name()).
		Load(ctx)

	name := people[0].ProjectsAsTeamMember()[0].Name()

	assert.Equal(t, "ACME Payment System", name)
}

func TestReverse2(t *testing.T) {
	ctx := context.Background()
	people := model.QueryPeople().
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

func Test2Nodes(t *testing.T) {
	ctx := context.Background()
	milestones := model.QueryMilestones().
		Join(node.Milestone().Project().Manager()).
		Where(Equal(node.Milestone().ID(), 1)). // Filter out people who are not managers
		Load(ctx)

	assert.True(t, milestones[0].NameIsValid(), "Milestone 1 has a name")
	assert.Equal(t, "Milestone A", milestones[0].Name(), "Milestone 1 has name of Milestone A")
	assert.True(t, milestones[0].Project().NameIsValid(), "Project 1 has a name")
	assert.Equal(t, "ACME Website Redesign", milestones[0].Project().Name(), "Project 1 has name of ACME Website Redesign")
	assert.True(t, milestones[0].Project().Manager().FirstNameIsValid(), "Person 7 has a name")
	assert.Equal(t, "Karen", milestones[0].Project().Manager().FirstName(), "Person 7 has first name of Karen")
}

func TestForwardMany(t *testing.T) {
	ctx := context.Background()
	milestones := model.QueryMilestones().
		Join(node.Milestone().Project().TeamMembers()).
		OrderBy(node.Milestone().Project().TeamMembers().LastName(), node.Milestone().Project().TeamMembers().FirstName()).
		Where(Equal(node.Milestone().ID(), 1)). // Filter out people who are not managers
		Load(ctx)

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

// Complex test finding all the team members of all the projects a person is managing, ordering by last name
func TestReverseMany(t *testing.T) {
	ctx := context.Background()
	people := model.QueryPeople().
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

	// Test an intermediate expansion
	people = model.QueryPeople().
		Join(node.Person().ProjectsAsManager().TeamMembers()).
		OrderBy(node.Person().ID(), node.Person().ProjectsAsManager().TeamMembers().LastName(), node.Person().ProjectsAsManager().TeamMembers().FirstName()).
		Expand(node.Person().ProjectsAsManager()).
		Load(ctx)

	names = []string{}
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

func TestManyForward(t *testing.T) {
	ctx := context.Background()
	people := model.QueryPeople().
		OrderBy(node.Person().ID(), node.Person().ProjectsAsTeamMember().Name()).
		Select(node.Person().ProjectsAsTeamMember().Manager().FirstName(), node.Person().ProjectsAsTeamMember().Manager().LastName()).
		Load(ctx)

	names := []string{}
	var p *model.Project
	for _, p = range people[0].ProjectsAsTeamMember() {
		names = append(names, p.Manager().FirstName()+" "+p.Manager().LastName())
	}
	names2 := []string{
		"Karen Wolfe",
		"John Doe",
	}
	assert.Equal(t, names2, names)

}

func TestUniqueReverse(t *testing.T) {
	ctx := context.Background()
	person := model.QueryPeople().
		Where(Equal(node.Person().LastName(), "Doe")).
		Get(ctx)
	assert.Nil(t, person.Login())

	person = model.QueryPeople().
		Where(Equal(node.Person().LastName(), "Doe")).
		Join(node.Person().Login()).
		Load(ctx)[0]
	assert.Equal(t, "jdoe", person.Login().Username())
}

func TestReverseReferences(t *testing.T) {
	// Test early binding
	ctx := context.Background()
	people := model.QueryPeople().
		Join(node.Person().ProjectsAsManager()).
		OrderBy(node.Person().LastName(), node.Person().FirstName(), node.Person().ProjectsAsManager().Name()).
		Load(ctx)
	person := people[2]
	person.SetFirstName("test")
	projects := person.ProjectsAsManager()
	person2 := projects[0].Manager()
	assert.Equal(t, "test", person2.FirstName())

	// Test forward reference looking back
	project := model.QueryProjects().
		Join(node.Project().Manager().ProjectsAsManager()).
		Load(ctx)[0]

	project.SetName("test")
	person = project.Manager()
	project2 := person.ProjectsAsManager()[0]
	assert.Equal(t, "test", project2.Name())

	// Test unique reverse reference
	person = model.QueryPeople().
		Join(node.Person().Login()).
		Load(ctx)[0]

	person.SetFirstName("test")
	login := person.Login()
	person2 = login.Person()
	assert.Equal(t, "test", person2.FirstName())

	// Test single expansion
	person = model.QueryPeople().
		Join(node.Person().ProjectsAsManager().Manager()).
		Expand(node.Person().ProjectsAsManager()).
		Load(ctx)[0]
	person.SetFirstName("test")
	projects = person.ProjectsAsManager()
	person2 = projects[0].Manager()
	assert.Equal(t, "test", person2.FirstName())

	// Test ManyMany
	project = model.QueryProjects().
		Join(node.Project().TeamMembers().ProjectsAsTeamMember()).
		Load(ctx)[0]

	project.SetName("test")
	people = project.TeamMembers()
	project2 = people[0].ProjectAsTeamMember()
	assert.Equal(t, "test", project2.Name())
}

func TestConditionalJoin(t *testing.T) {
	ctx := context.Background()

	projects := model.QueryProjects().
		OrderBy(node.Project().Name()).
		Join(node.Project().Manager(), Equal(node.Project().Manager().LastName(), "Wolfe")).
		Join(node.Project().TeamMembers(), Equal(node.Project().TeamMembers().LastName(), "Smith")).
		Load(ctx)

	// Reverse references
	people := model.QueryPeople().
		Join(node.Person().Addresses(), Equal(node.Person().Addresses().City(), "New York")).
		Join(node.Person().ProjectsAsManager(), Equal(node.Person().ProjectsAsManager().ProjectStatusTypeID(), model.PROJECT_STATUS_TYPE_OPEN)).
		Join(node.Person().ProjectsAsManager().Milestones()).
		Join(node.Person().Login(), Like(node.Person().Login().Username(), "b%")).
		OrderBy(node.Person().LastName(), node.Person().FirstName(), node.Person().ProjectsAsManager().Name()).
		Load(ctx)

	assert.Equal(t, "John", people[2].FirstName(), "John Does is the 3rd Person.")
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
	ctx := context.Background()

	// Reverse references
	people := model.QueryPeople().
		Join(node.Person().Addresses(), Equal(node.Person().Addresses().City(), "Mountain View")).
		Join(node.Person().ProjectsAsManager(), Like(node.Person().ProjectsAsManager().Name(), "%Website%")).
		Join(node.Person().ProjectsAsManager().Milestones()).
		Expand(node.Person().ProjectsAsManager()).
		OrderBy(node.Person().LastName(), node.Person().FirstName(), node.Person().ProjectsAsManager().Name()).
		Load(ctx)

	assert.Equal(t, "Karen", people[11].FirstName(), "Karen is the 12th Person.")
	assert.Len(t, people[11].ProjectsAsManager(), 1, "Karen Wolfe selected 1 project.")
	assert.Len(t, people[11].Addresses(), 1, "Karen Wolfe has 1 Address.")

	assert.Len(t, people[2].ProjectsAsManager(), 0, "John Doe selected no projects.")
	assert.Len(t, people[2].Addresses(), 0, "John Doe has no Addresses")

}

func TestSelectByID(t *testing.T) {
	ctx := context.Background()

	projects := model.QueryProjects().
		OrderBy(node.Project().Name().Descending()).
		Load(ctx)

	id := projects[3].ID()

	// Reverse references
	people := model.QueryPeople().
		Join(node.Person().ProjectsAsManager()).
		Where(Equal(node.Person().LastName(), "Wolfe")).
		Load(ctx)

	assert.Equal(t, people[0].ProjectAsManager(id).Name(), "ACME Payment System")
}
