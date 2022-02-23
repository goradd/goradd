package javascript_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/goradd/goradd/pkg/javascript"
	"github.com/stretchr/testify/assert"
)

func ExampleFunctionCall() {
	c := NewFunctionCall("substr", "str", 1, 2)
	fmt.Println(c.JavaScript())
	// Output: str.substr(1,2)
}

func TestFunctionCall_MarshalJSON(t *testing.T) {
	c := NewFunctionCall("myFunc", "this", "b")
	a := []interface{}{c}
	b, err := json.Marshal(a)
	assert.NoError(t, err)
	a = nil
	err = json.Unmarshal(b, &a)
	assert.NoError(t, err)
	assert.Equal(t, "this", a[0].(map[string]interface{})["context"].(string))
	assert.Equal(t, "myFunc", a[0].(map[string]interface{})["func"].(string))
	assert.Equal(t, "function", a[0].(map[string]interface{})["goraddObject"].(string))
	assert.Equal(t, "b", a[0].(map[string]interface{})["params"].([]interface{})[0].(string))
}
