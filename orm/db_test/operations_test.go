package db

import (
	"testing"
	"grlocal/model"
	"context"
	. 	"github.com/spekary/goradd/orm/op"
	"grlocal/model/node"
	"github.com/stretchr/testify/assert"
	"github.com/spekary/goradd/orm/db"
	"github.com/spekary/goradd/datetime"

)


func TestEqualBasic(t *testing.T) {
	ctx := context.Background()
	projects := model.QueryProjects().
		Where(Equal(node.Project().Num(), 2)).
		OrderBy(node.Project().Num()).
		Load(ctx)

	assert.EqualValues(t,2, projects[0].Num(),"Did not find correct project.")

}

func TestLogical(t *testing.T) {
	type testCase struct {
		testNode db.NodeI
		objectNum int
		expectedId interface{}
		count int
		desc string
	}
	var tests = []testCase {
		{GreaterThan(node.Project().Num(), 3), 0, 4, 1,"Greater than uint test"},
		{GreaterThan(node.Project().StartDate(), datetime.NewDateTime("1/1/2006", datetime.UsDate)), 0, 2, 2, "Greater than datetime test"},
		{GreaterThan(node.Project().Spent(), 10000), 1, 2, 2,"Greater than float test"},
		{LessThan(node.Project().Num(), 3), 1, 2, 2, "Less than uint test"},
		{LessThan(node.Project().EndDate(), datetime.NewDateTime("1/1/2006", datetime.UsDate)), 1, 4, 2,"Less than date test"},
		{IsNull(node.Project().EndDate()), 0, 2, 1,"Is Null test"},
		{IsNotNull(node.Project().EndDate()), 0, 1, 3,"Is Not Null test"},
		{GreaterOrEqual(node.Project().ProjectStatusTypeID(), 2), 1, 4, 2,"Greater or Equal test"},
		{LessOrEqual(node.Project().StartDate(), datetime.NewDateTime("2/15/2006", datetime.UsDate)), 2, 4, 3,"Less or equal date test"},
		{Or(Equal(node.Project().Num(), 1), Equal(node.Project().Num(), 4)), 1, 4, 2,"Or test"},
		{Xor(Equal(node.Project().Num(), 3), Equal(node.Project().ProjectStatusTypeID(), 1)), 0, 2, 1,"Xor test"},
		{Not(Xor(Equal(node.Project().Num(), 3), Equal(node.Project().ProjectStatusTypeID(), 1))), 0, 1, 3,"Not test"},
		{Like(node.Project().Name(), "%ACME%"), 1, 4, 2, "Like test"},
		{In(node.Project().Num(), 2,3,4), 1, 3, 3, "In test"},

	}

	ctx := context.Background()

	var projects []*model.Project
	for _,c := range tests {
		projects = model.QueryProjects().
			Where(c.testNode).
			OrderBy(node.Project().Num()).
			Load(ctx)

		assert.EqualValues(t,c.expectedId, projects[c.objectNum].Num(),c.desc)
		assert.EqualValues(t,c.count, len(projects),c.desc + " - count")
	}
}


func TestCount2(t *testing.T) {
	ctx := context.Background()
	count := model.QueryPeople().
		Count(ctx, true, node.Person().LastName())

	assert.EqualValues(t,10, count)

}

func TestCalculations(t *testing.T) {
	type testCase struct {
		testNode db.NodeI
		objectNum int
		expectedValue interface{}
		desc string
	}
	var tests = []testCase {
		{Add(node.Project().Spent(), node.Project().Budget()), 0, "19811.00","Add test"},
		{Subtract(node.Project().Spent(), 2000), 0, "8250.75","Subtract test"},
		{Multiply(node.Project().Num(), 3), 3, "12","Multiply test"},
		{Mod(node.Project().Num(), 2), 2, "1","Mod test"},
		{Round(Divide(node.Project().Num(), 2)), 3, "2","Mod test"},
	}

	ctx := context.Background()

	var projects []*model.Project
	for _,c := range tests {
		projects = model.QueryProjects().
			Alias("Value", c.testNode).
			OrderBy(node.Project().Num()).
			Load(ctx)

		assert.EqualValues(t,c.expectedValue, projects[c.objectNum].GetAlias("Value").String(),c.desc)
	}
}

func TestAggregates(t *testing.T) {
	ctx := context.Background()
	projects := model.QueryProjects().
		Alias("sum", Sum(node.Project().Spent())).
		OrderBy(node.Project().ProjectStatusTypeID()).
		GroupBy(node.Project().ProjectStatusTypeID()).
		Load(ctx)

	assert.EqualValues(t,77400.5, projects[0].GetAlias("sum").Float())

	projects2 := model.QueryProjects().
		Alias("min", Min(node.Project().Spent())).
		OrderBy(node.Project().ProjectStatusTypeID()).
		//GroupBy(node.Project().ProjectStatusTypeID()).
		Load(ctx)

	assert.EqualValues(t,4200.50, projects2[0].GetAlias("min").Float())
}