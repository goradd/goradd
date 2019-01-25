package javascript_test

import (
	"fmt"
	. "github.com/goradd/goradd/pkg/javascript"
)

func ExampleClosure() {
	c := Closure ("return a == b;", "a", "b")
	fmt.Println(c.JavaScript())
	// Output: function(a, b) {return a == b;}
}


func ExampleClosureCall() {
	c := ClosureCall ("return this == b;", "a", "b")
	fmt.Println(c.JavaScript())
	// Output: (function(b) {return this == b;}).call(a)
}
