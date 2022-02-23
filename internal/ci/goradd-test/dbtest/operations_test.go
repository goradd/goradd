package dbtest

import (
	. "github.com/goradd/goradd/pkg/orm/op"
	"github.com/goradd/goradd/pkg/orm/query"
	"github.com/goradd/goradd/pkg/time"
	"github.com/stretchr/testify/assert"
	"goradd-project/gen/goradd/model"
	"goradd-project/gen/goradd/model/node"
	"testing"
)

func TestEqualBasic(t *testing.T) {
	ctx := getContext()
	projects := model.QueryProjects(ctx).
		Where(Equal(node.Project().Num(), 2)).
		OrderBy(node.Project().Num()).
		Load()

	assert.EqualValues(t, 2, projects[0].Num(), "Did not find correct project.")

}

func TestMultiWhere(t *testing.T) {
	ctx := getContext()
	projects := model.QueryPeople(ctx).
		Where(Equal(node.Person().LastName(), "Smith")).
		Where(Equal(node.Person().FirstName(), "Alex")).
		Load()

	assert.Len(t, projects, 1)
}


func TestLogical(t *testing.T) {
	type testCase struct {
		testNode   query.NodeI
		objectNum  int
		expectedId interface{}
		count      int
		desc       string
	}
	tests := []testCase{
		{GreaterThan(node.Project().Num(), 3), 0, 4, 1, "Greater than uint test"},
		{GreaterThan(node.Project().StartDate(), time.NewDate(2006,1,1)), 0, 2, 2, "Greater than datetime test"},
		{GreaterThan(node.Project().Spent(), 10000), 1, 2, 2, "Greater than float test"},
		{LessThan(node.Project().Num(), 3), 1, 2, 2, "Less than uint test"},
		{LessThan(node.Project().EndDate(), time.NewDate(2006,1,1)), 1, 4, 2, "Less than date test"},
		{IsNull(node.Project().EndDate()), 0, 2, 1, "Is Null test"},
		{IsNotNull(node.Project().EndDate()), 0, 1, 3, "Is Not Null test"},
		{GreaterOrEqual(node.Project().ProjectStatusTypeID(), 2), 1, 4, 2, "Greater or Equal test"},
		{LessOrEqual(node.Project().StartDate(), time.NewDate(2006,2,15)), 2, 4, 3, "Less or equal date test"},
		{Or(Equal(node.Project().Num(), 1), Equal(node.Project().Num(), 4)), 1, 4, 2, "Or test"},
		{Xor(Equal(node.Project().Num(), 3), Equal(node.Project().ProjectStatusTypeID(), 1)), 0, 2, 1, "Xor test"},
		{Not(Xor(Equal(node.Project().Num(), 3), Equal(node.Project().ProjectStatusTypeID(), 1))), 0, 1, 3, "Not test"},
		{Like(node.Project().Name(), "%ACME%"), 1, 4, 2, "Like test"},
		{In(node.Project().Num(), 2, 3, 4), 1, 3, 3, "In test"},
	}

	ctx := getContext()

	for i, c := range tests {
		t.Run(c.desc, func(t *testing.T) {
			var projects []*model.Project

			projects = model.QueryProjects(ctx).
				Where(c.testNode).
				OrderBy(node.Project().Num()).
				Load()

			if len(projects) <= c.objectNum {
				t.Errorf("Test case produced out of range error. Test case #: %d", i)
			} else {
				assert.EqualValues(t, c.expectedId, projects[c.objectNum].Num(), c.desc)
				assert.EqualValues(t, c.count, len(projects), c.desc+" - count")
			}
		})
	}
}

func TestCount2(t *testing.T) {
	ctx := getContext()
	count := model.QueryPeople(ctx).
		Count(true, node.Person().LastName())

	assert.EqualValues(t, 10, count)

}

func TestCalculations(t *testing.T) {
	type testCase struct {
		testNode      query.NodeI
		objectNum     int
		expectedValue interface{}
		desc          string
	}
	var intTests = []testCase{
		{Multiply(node.Project().Num(), 3), 3, 12, "Multiply test"},
		{Mod(node.Project().Num(), 2), 2, 1, "Mod test"},
		{Round(Divide(node.Project().Num(), 2)), 3, 2, "Mod test"},
	}

	var floatTests = []testCase{
		{Add(node.Project().Spent(), node.Project().Budget()), 0, 19811.00, "Add test"},
		{Subtract(node.Project().Spent(), 2000), 0, 8250.75, "Subtract test"},
	}

	ctx := getContext()

	var projects []*model.Project
	for _, c := range intTests {
		projects = model.QueryProjects(ctx).
			Alias("Value", c.testNode).
			OrderBy(node.Project().Num()).
			Load()

		assert.EqualValues(t, c.expectedValue, projects[c.objectNum].GetAlias("Value").Int(), c.desc)
	}

	for _, c := range floatTests {
		projects = model.QueryProjects(ctx).
			Alias("Value", c.testNode).
			OrderBy(node.Project().Num()).
			Load()

		assert.EqualValues(t, c.expectedValue, projects[c.objectNum].GetAlias("Value").Float(), c.desc)
	}

}

func TestAggregates(t *testing.T) {
	ctx := getContext()
	projects := model.QueryProjects(ctx).
		Alias("sum", Sum(node.Project().Spent())).
		OrderBy(node.Project().ProjectStatusTypeID()).
		GroupBy(node.Project().ProjectStatusTypeID()).
		Load()

	assert.EqualValues(t, 77400.5, projects[0].GetAlias("sum").Float())

	projects2 := model.QueryProjects(ctx).
		Alias("min", Min(node.Project().Spent())).
		OrderBy(node.Project().ProjectStatusTypeID()).
		GroupBy(node.Project().ProjectStatusTypeID()).
		Load()

	assert.EqualValues(t, 4200.50, projects2[0].GetAlias("min").Float())
}

func TestAliases(t *testing.T) {
	ctx := getContext()
	nVoyel := node.Person().ProjectsAsManager().Milestones()
	nVoyel.SetAlias("voyel")
	nConson := node.Person().ProjectsAsManager().Milestones()
	nConson.SetAlias("conson")

	people := model.QueryPeople(ctx).
		OrderBy(node.Person().LastName(), node.Person().FirstName()).
		Where(IsNotNull(nConson)).
		Join(nVoyel, In(nVoyel.Name(), "Milestone A", "Milestone E", "Milestone I")).
		Join(nConson, NotIn(nConson.Name(), "Milestone A", "Milestone E", "Milestone I")).
		GroupBy(node.Person().ID(),node.Person().FirstName(), node.Person().LastName()).
		Alias("min_voyel", Min(nVoyel.Name())).
		Alias("min_conson", Min(nConson.Name())).
		Load()

	assert.EqualValues(t, 3, len(people))
	assert.Equal(t, "Doe", people[0].LastName())
	assert.Equal(t, "Ho", people[1].LastName())
	assert.Equal(t, "Wolfe", people[2].LastName())

	assert.True(t, people[0].GetAlias("min_voyel").IsNil())
	assert.Equal(t, "Milestone F", people[0].GetAlias("min_conson").String())

	assert.Equal(t, "Milestone E", people[1].GetAlias("min_voyel").String())
	assert.Equal(t, "Milestone D", people[1].GetAlias("min_conson").String())

	assert.Equal(t, "Milestone A", people[2].GetAlias("min_voyel").String())
	assert.Equal(t, "Milestone B", people[2].GetAlias("min_conson").String())
}
