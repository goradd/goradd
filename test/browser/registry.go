package browser

type testRunnerFunction func(*TestForm)

var tests = make(map[string]testRunnerFunction)

// RegisterTestFunction registers the test and presents it in the list of tests that can be performed.
func RegisterTestFunction (name string, f testRunnerFunction) {
	tests[name] = f
}

