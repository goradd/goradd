package javascript_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/goradd/goradd/pkg/javascript"
	"github.com/stretchr/testify/assert"
)

func ExampleClosure() {
	c := NewClosure("return a == b;", "a", "b")
	fmt.Println(c.JavaScript())
	// Output: function(a, b) {return a == b;}
}

func ExampleClosureCall() {
	c := NewClosureCall("return this == b;", "a", "b")
	fmt.Println(c.JavaScript())
	// Output: (function(b) {return this == b;}).call(a)
}

func Test_closure_MarshalJSON(t *testing.T) {
	type T struct {
		Func         string   `json:"func"`
		GoraddObject string   `json:"goraddObject"`
		Params       []string `json:"params"`
		Context      string   `json:"call,omitempty"`
	}

	c := NewClosure("return (b + 1);", "b")
	c2 := NewClosureCall("return (this.b + a);", "this", "a")
	a := []interface{}{c, c2}
	b, err := json.Marshal(a)
	assert.NoError(t, err)

	var j []T
	err = json.Unmarshal(b, &j)
	assert.NoError(t, err)
	assert.Equal(t, "return (b + 1);", j[0].Func)
	assert.Equal(t, "closure", j[0].GoraddObject)
	assert.Empty(t, j[0].Context)
	assert.Equal(t, "b", j[0].Params[0])

	assert.Equal(t, "return (this.b + a);", j[1].Func)
	assert.Equal(t, "closure", j[1].GoraddObject)
	assert.Equal(t, "this", j[1].Context)
	assert.Equal(t, "a", j[1].Params[0])
}
