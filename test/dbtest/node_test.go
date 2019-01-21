package dbtest

import (
	"github.com/goradd/goradd/pkg/orm/query"
	"goradd-project/gen/goradd/model"
	"goradd-project/gen/goradd/model/node"
	"testing"
)

func TestNodeEquality(t *testing.T) {

	n := node.Person()
	if !n.Equals(n) {
		t.Error("Table node not equal to self")
	}

	n = node.Project().Manager()
	if !n.Equals(n) {
		t.Error("Reference node not equal to self")
	}

	n2 := node.Person().ProjectsAsManager()
	if !n2.Equals(n2) {
		t.Error("Reverse Reference node not equal to self")
	}

	n3 := node.Person().ProjectsAsTeamMember()
	if !n3.Equals(n3) {
		t.Error("Many-Many node not equal to self")
	}

	n4 := query.NewValueNode(model.PersonTypeContractor)
	if !n4.Equals(n4) {
		t.Error("Type node not equal to self")
	}

}
