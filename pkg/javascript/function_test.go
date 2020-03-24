package javascript_test

import (
	"fmt"
	. "github.com/goradd/goradd/pkg/javascript"
)

func ExampleFunctionCall() {
	c := Function("substr", "str", 1, 2)
	fmt.Println(c.JavaScript())
	// Output: str.substr(1,2)
}
