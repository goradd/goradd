package browsertest

import (
	"fmt"
	"github.com/goradd/maps"
)

type testRunnerFunction func(*TestForm)

var tests = new(maps.SliceMap[string, any])

// RegisterTestFunction registers the test and presents it in the list of tests that can be performed.
func RegisterTestFunction(name string, f testRunnerFunction) {
	if ok := tests.Has(name); ok {
		panic(fmt.Sprintf("Test %s has already been registered.", name))
	}
	tests.Set(name, f)
}
