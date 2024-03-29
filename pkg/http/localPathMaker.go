package http

import (
	"github.com/goradd/goradd/pkg/config"
	"path"
)

// LocalPathMaker converts an HTTP path rooted to the application, to a path accessible by the server.
type LocalPathMaker func(string) string

var localPathMaker LocalPathMaker = defaultLocalPathMaker

func defaultLocalPathMaker(p string) string {
	var hasSlash bool
	if p == "" {
		panic(`cannot make a local path to an empty path. If you are trying to refer to the root, use '/'.`)
	}
	if p[len(p)-1] == '/' {
		hasSlash = true
	}
	// We want to prevent local paths (here/there) and external paths (http://here/there) from getting additional paths
	if p[0] == '/' && config.ProxyPath != "" {
		p = path.Join(config.ProxyPath, p) // will strip trailing slashes
		if hasSlash {
			p = p + "/"
		}
	}
	return p

}

// MakeLocalPath turns a path that points to a resource on this computer into a path that will reach
// that resource from a browser. It takes into account a variety of settings that may affect the
// path and that will depend on how the app is deployed.
// You can inject your own local path maker using SetLocalPathMaker
func MakeLocalPath(p string) string {
	return localPathMaker(p)
}

// SetLocalPathMaker sets the local path maker to the given one.
//
// The default local path maker will prepend config.ProxyPath to all local paths.
func SetLocalPathMaker(f LocalPathMaker) {
	localPathMaker = f
}
