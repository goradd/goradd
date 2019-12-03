package dbtest

import (
	"encoding/json"
	"goradd-project/gen/goradd/model"
	"testing"

	. "github.com/goradd/goradd/pkg/orm/op"
	"goradd-project/gen/goradd/model/node"
	"strconv"

	"github.com/goradd/goradd/pkg/datetime"
	"github.com/goradd/goradd/pkg/orm/query"
	"github.com/stretchr/testify/assert"
)


func TestBasic(t *testing.T) {

	ctx := getContext()

	people := model.QueryPeople(ctx).
		OrderBy(node.Person().ID()).
		Load(ctx)
	if len(people) != 12 {
		t.Error("12 people not found")
	}
	if people[0].FirstName() != "John" {
		t.Error("First person is not John")
	}
}

func TestSort(t *testing.T) {
	ctx := getContext()
	people := model.QueryPeople(ctx).
		OrderBy(node.Person().LastName()).
		Load(ctx)

	if people[0].LastName() != "Brady" {
		t.Error("Person found not Brady, found " + people[0].LastName())
	}

	people = model.QueryPeople(ctx).
		OrderBy(node.Person().FirstName()).
		Load(ctx)
	if people[0].FirstName() != "Alex" {
		t.Error("Person found not Alex, found " + people[0].FirstName())
	}


	// Testing for regression bug with multiple sorts
	people = model.QueryPeople(ctx).
		OrderBy(node.Person().LastName().Descending(), node.Person().FirstName().Ascending()).
		Load(ctx)
	if people[0].FirstName() != "Karen" || people[0].LastName() != "Wolfe" {
		t.Error("Person found not Karen Wolfe, found " + people[0].FirstName() + " " + people[0].LastName())
	}
}

func TestWhere(t *testing.T) {
	ctx := getContext()
	_ = query.Value("Smith").(query.NodeI)
	people := model.QueryPeople(ctx).
		Where(Equal(node.Person().LastName(), "Smith")).
		OrderBy(node.Person().FirstName().Descending(), node.Person().LastName()).
		Load(ctx)

	if people[0].FirstName() != "Wendy" {
		t.Error("Person found not Wendy, found " + people[0].FirstName())
	}
}

func TestReference(t *testing.T) {
	ctx := getContext()
	projects := model.QueryProjects(ctx).
		Join(node.Project().Manager()).
		OrderBy(node.Project().ID()).
		Load(ctx)

	if projects[0].Manager().FirstName() != "Karen" {
		t.Error("Person found not Karen, found " + projects[0].Manager().FirstName())
	}

}

func TestManyMany(t *testing.T) {
	ctx := getContext()
	projects := model.QueryProjects(ctx).
		Join(node.Project().TeamMembers()).
		OrderBy(node.Project().ID()).
		Load(ctx)

	if len(projects[0].TeamMembers()) != 5 {
		t.Error("Did not find 5 team members in project 1. Found: " + strconv.Itoa(len(projects[0].TeamMembers())))
	}

}

func TestReverseReference(t *testing.T) {
	ctx := getContext()
	people := model.QueryPeople(ctx).
		Join(node.Person().ProjectsAsManager()).
		OrderBy(node.Person().ID()).
		Load(ctx)

	if people[0].FirstName() != "John" {
		t.Error("Did not find person 0.")
	}

	if len(people[6].ProjectsAsManager()) != 2 {
		t.Error("Did not find 2 ProjectsAsManagers.")
	}

}

func TestBasicType(t *testing.T) {
	ctx := getContext()
	projects := model.QueryProjects(ctx).
		OrderBy(node.Project().ID()).
		Load(ctx)

	if projects[0].ProjectStatusType() != model.ProjectStatusTypeCompleted {
		t.Error("Did not find correct project type.")
	}
}

func TestManyType(t *testing.T) {
	ctx := getContext()
	people := model.QueryPeople(ctx).
		OrderBy(node.Person().ID(), node.Person().PersonTypes().ID().Descending()).
		Join(node.Person().PersonTypes()).
		Load(ctx)

	if len(people[0].PersonTypes()) != 2 {
		t.Error("Did not expand to 2 person types.")
	}

	if people[0].PersonTypes()[0] != model.PersonTypeInactive {
		t.Error("Did not find correct person type.")
	}

}

func TestManyManySingles(t *testing.T) {
	ctx := getContext()
	projects := model.QueryProjects(ctx).
		Expand(node.Project().TeamMembers()).
		OrderBy(node.Project().ID(), node.Project().TeamMembers().FirstName()).
		Load(ctx)

	if projects[4].Name() != "ACME Website Redesign" { // should have 5 lines here that are all project 1
		t.Error("Did not find expanded project ACME Website Redesign.")
	}

	if projects[5].Name() != "State College HR System" { // should have 5 lines here that are all project 1
		t.Error("Did not find expanded project State College HR System.")
	}

	if projects[3].TeamMember().FirstName() != "Samantha" {
		t.Error("Did not find Samantha. Found: " + projects[3].TeamMember().FirstName())
	}

}

func TestReverseReferenceSingles(t *testing.T) {
	ctx := getContext()
	people := model.QueryPeople(ctx).
		Expand(node.Person().ProjectsAsManager()).
		OrderBy(node.Person().ID()).
		Load(ctx)

	if people[7].FirstName() != "Karen" {
		t.Error("Did not find expanded person Karen Wolfe.")
	}

	if people[7].ProjectsAsManager()[0].ManagerID() != people[7].ID() {
		t.Error("Did not find Project 2.")
	}

}

func TestManyTypeSingles(t *testing.T) {
	ctx := getContext()
	people := model.QueryPeople(ctx).
		OrderBy(node.Person().ID(), node.Person().PersonTypes().ID().Descending()).
		Expand(node.Person().PersonTypes()).
		Load(ctx)

	if people[1].PersonType() != model.PersonTypeManager {
		t.Error("Did not find correct person type.")
	}

}

func TestAlias(t *testing.T) {
	ctx := getContext()
	projects := model.QueryProjects(ctx).
		Where(Equal(node.Project().ID(), 1)).
		Alias("Difference", Subtract(node.Project().Budget(), node.Project().Spent())).
		Load(ctx)

	v := projects[0].GetAlias("Difference").Float()
	assert.EqualValues(t, -690.5, v)
}

func TestAlias2(t *testing.T) {
	ctx := getContext()
	projects := model.QueryProjects(ctx).
		Alias("a", node.Project().Num()).
		Alias("b", node.Project().Name()).
		Alias("c", node.Project().Spent()).
		Alias("d", node.Project().StartDate()).
		Alias("e", IsNull(node.Project().EndDate())).
		OrderBy(node.Project().ID()).
		Load(ctx)

	project := projects[0]
	assert.Equal(t, 1, project.GetAlias("a").Int())
	assert.Equal(t, "ACME Website Redesign", project.GetAlias("b").String())
	assert.Equal(t, 10250.75, project.GetAlias("c").Float())
	d,_ := datetime.FromSqlDateTime("2004-03-01")
	assert.EqualValues(t, d, project.GetAlias("d").DateTime())
	assert.Equal(t, false, project.GetAlias("e").Bool())

}

func TestCount(t *testing.T) {
	ctx := getContext()
	count := model.QueryProjects(ctx).
		Count(ctx, false)

	assert.EqualValues(t, 4, count)
}

func TestGroupBy(t *testing.T) {
	ctx := getContext()
	projects := model.QueryProjects(ctx).
		Alias("teamMemberCount", Count(node.Project().TeamMembers())).
		GroupBy(node.Project()).
		Load(ctx)

	assert.EqualValues(t, 5, projects[0].GetAlias("teamMemberCount").Int())

}

func TestSelect(t *testing.T) {
	ctx := getContext()
	projects := model.QueryProjects(ctx).
		Select(node.Project().Name()).
		Load(ctx)

	project := projects[0]
	assert.True(t, project.NameIsValid())
	assert.False(t, project.ManagerIDIsValid())
	assert.True(t, project.IDIsValid())
}

func TestLimit(t *testing.T) {
	ctx := getContext()
	people := model.QueryPeople(ctx).
		OrderBy(node.Person().ID()).
		Limit(2, 3).
		Load(ctx)

	assert.EqualValues(t, "Mike", people[0].FirstName())
	assert.Len(t, people, 2)
}

func TestSaveAndDelete(t *testing.T) {
	ctx := getContext()

	person := model.NewPerson()
	person.SetFirstName("Test1")
	person.SetLastName("Last1")
	person.Save(ctx)

	people := model.QueryPeople(ctx).
		Where(
			And(
				Equal(
					node.Person().FirstName(), "Test1"),
				Equal(
					node.Person().LastName(), "Last1"))).
		Load(ctx)

	assert.EqualValues(t, person.ID(), people[0].ID())

	people[0].Delete(ctx)

	people = model.QueryPeople(ctx).
		Where(
			And(
				Equal(
					node.Person().FirstName(), "Test1"),
				Equal(
					node.Person().LastName(), "Last1"))).
		Load(ctx)

	assert.Len(t, people, 0, "Deleted the person")
}

func TestUpdate(t *testing.T) {
	ctx := getContext()

	person := model.NewPerson()
	person.SetFirstName("Test1")
	person.SetLastName("Last1")
	person.Save(ctx)

	person.SetFirstName("Test2")
	person.Update(ctx)

	people := model.QueryPeople(ctx).
		Where(
			And(
				Equal(
					node.Person().FirstName(), "Test2"),
			)).
		Load(ctx)

	assert.EqualValues(t, person.ID(), people[0].ID())

	person.Delete(ctx)

	people = model.QueryPeople(ctx).
		Where(
			And(
				Equal(
					node.Person().FirstName(), "Test2"),
			)).
		Load(ctx)

	assert.Len(t, people, 0, "Deleted the person")

}

func TestSingleEmpty(t *testing.T) {
	ctx := getContext()

	people := model.QueryPeople(ctx).
		Where(Equal(node.Person().ID(), 12345)).
		Load(ctx)

	assert.Len(t, people, 0)

}

func TestLazyLoad(t *testing.T) {
	ctx := getContext()

	projects := model.QueryProjects(ctx).
		Where(Equal(node.Project().ID(), 1)).
		Load(ctx)

	var mId string = projects[0].ID() // foreign keys are treated as strings for cross-database compatibility
	assert.Equal(t, "1", mId)

	manager := projects[0].LoadManager(ctx)
	assert.Equal(t, "7", manager.ID())
}

func TestDeleteQuery(t *testing.T) {
	ctx := getContext()

	person := model.NewPerson()
	person.SetFirstName("Test1")
	person.SetLastName("Last1")
	person.Save(ctx)

	model.QueryPeople(ctx).
		Where(
			And(
				Equal(
					node.Person().FirstName(), "Test1"),
				Equal(
					node.Person().LastName(), "Last1"))).
		Delete(ctx)

	people := model.QueryPeople(ctx).
		Where(
			And(
				Equal(
					node.Person().FirstName(), "Test1"),
				Equal(
					node.Person().LastName(), "Last1"))).
		Load(ctx)

	assert.Len(t, people, 0, "Deleted the person")
}

func TestHaving(t *testing.T) {
	// This particular test shows a quirk of SQL that requires:
	// 1) If you have an aggregate clause (like COUNT), you MUST have a GROUPBY clause, and
	// 2) If you have a GROUPBY, you MUST SELECT and only select the things you are grouping by.
	//
	// Sooo, when we see a GroupBy, we automatically also select the same nodes.
	ctx := getContext()
	projects := model.QueryProjects(ctx).
		GroupBy(node.Project().ID(), node.Project().Name()).
		OrderBy(node.Project().ID()).
		Alias("team_member_count", Count(node.Project().TeamMembers())).
		Having(GreaterThan(Count(query.Alias("team_member_count")), 5)).
		Load(ctx)

	assert.Len(t, projects, 2)
	assert.Equal(t, "State College HR System", projects[0].Name())
	assert.Equal(t, 6, projects[0].GetAlias("team_member_count").Int())
}

func TestFailedJoins(t *testing.T) {
	ctx := getContext()
	assert.Panics(t, func(){model.QueryProjects(ctx).Join(node.Person())})
	assert.Panics(t, func(){model.QueryProjects(ctx).Join(node.Project().ManagerID())})
}

func TestFailedExpand(t *testing.T) {
	ctx := getContext()
	assert.Panics(t, func(){model.QueryProjects(ctx).Expand(node.Person())})
	assert.Panics(t, func(){model.QueryProjects(ctx).Expand(node.Project().Manager())})
}

func TestFailedGroupBy(t *testing.T) {
	ctx := getContext()
	assert.Panics(t, func(){model.
		QueryProjects(ctx).
		GroupBy(node.Project().Name()).
			Select(node.Project().Name())})
}

func TestJson(t *testing.T) {
	ctx := getContext()
	p := model.LoadProject(ctx, "1",
		node.Project().Name(),
		node.Project().ProjectStatusType(),
		node.Project().Manager().FirstName())
	j,err := json.Marshal(p)
	assert.NoError(t, err)
	m := make(map[string]interface{})
	err = json.Unmarshal(j, &m)
	assert.NoError(t, err)
	assert.Equal(t, "ACME Website Redesign", m["name"])
	assert.Equal(t, "Completed", m["projectStatusType"])
	assert.Equal(t, "Karen", m["manager"].(map[string]interface{})["firstName"])
}