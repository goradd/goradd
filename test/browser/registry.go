package browser

import "fmt"

type testRunnerFunction func(*TestForm)

var tests = make(map[string]testRunnerFunction)

// RegisterTestFunction registers the test and presents it in the list of tests that can be performed.
func RegisterTestFunction (name string, f testRunnerFunction) {
	if _,ok := tests[name]; ok {
		panic(fmt.Sprintf("Test %s has already been registered.", name))
	}
	tests[name] = f
}

