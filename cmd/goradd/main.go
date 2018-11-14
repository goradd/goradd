// The main package runs a web server that works as an aid to installation, updating and building your application.
// Most of the code is in the buildtools directory
package goradd

import (
	"github.com/spekary/goradd/buildtools"
)


// Create other flags you might care about here

func main() {
	buildtools.Launch()
}

