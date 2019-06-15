package config

import (
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"path/filepath"
	"runtime"
)

// Loader is the injected loader. Set it during initialization if you want to load bootstrap differently than below.
var Loader func(page.FormI)

// Configuration options for Bootstrap

// LoadBootstrap loads the various bootstrap files required by bootstrap. It is called automatically
// by the bootstrap components, but this gives you an opportunity to customize where the client
// gets the files.
func LoadBootstrap(form page.FormI) {
	if Loader != nil {
		Loader(form)
	} else {
		if config.Release {
			form.AddJavaScriptFile("https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.3/umd/popper.min.js", false,
				html.NewAttributes().Set("integrity", "sha384-ZMP7rVo3mIykV+2+9J3UJ46jBk0WLaUAdn689aCwoqbBJiSnjAK/l8WvCWPIPm49").Set("crossorigin", "anonymous"))
			form.AddJavaScriptFile("https://stackpath.bootstrapcdn.com/bootstrap/4.1.1/js/bootstrap.min.js", false,
				html.NewAttributes().Set("integrity", "sha384-smHYKdLADwkXOn1EmN1qk/HfnUcbVRZyYmZ4qpPea6sjB/pTJ0euyQp0Mk8ck+5T").Set("crossorigin", "anonymous"))
			form.AddStyleSheetFile("https://stackpath.bootstrapcdn.com/bootstrap/4.1.1/css/bootstrap.min.css",
				html.NewAttributes().Set("integrity", "sha384-WskhaSGFgHYWDcbwN70/dfYBj47jz9qbsMId/iRN3ewGhXQFZCSftd1LZCfmhktB").Set("crossorigin", "anonymous"))

		} else {
			form.AddJavaScriptFile(filepath.Join(BootstrapAssets(), "js", "bootstrap.bundle.js"), false, nil)
			form.AddStyleSheetFile(filepath.Join(BootstrapAssets(), "css", "bootstrap.min.css"), nil)
		}
	}

}

func BootstrapAssets() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filepath.Dir(filename)), "assets")
}

func init() {
	page.RegisterAssetDirectory(BootstrapAssets(), config.AssetPrefix+"bootstrap")
}
