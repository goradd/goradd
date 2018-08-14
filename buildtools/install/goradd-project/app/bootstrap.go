package app

import (
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	"goradd-project/config"
)

// Configuration options for Bootstrap

// LoadBootstrap loads the various bootstrap files required by bootstrap. It is called automatically
// by the bootstrap components, but this gives you an opportunity to customize where the client
// gets the files.
func LoadBootstrap(form page.FormI) {
	switch config.Mode {
	case config.AppModeDevelopment:
		// Get files locally in case you are off the grid temporarily
		form.AddJavaScriptFile(config.GoraddDir+"/bootstrap/assets/js/bootstrap.bundle.js", false, nil)
		form.AddStyleSheetFile(config.GoraddDir+"/bootstrap/assets/css/bootstrap.min.css", nil)
	case config.AppModeRelease:
		form.AddJavaScriptFile("https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.3/umd/popper.min.js", false,
			html.NewAttributes().Set("integrity", "sha384-ZMP7rVo3mIykV+2+9J3UJ46jBk0WLaUAdn689aCwoqbBJiSnjAK/l8WvCWPIPm49").Set("crossorigin", "anonymous"))
		form.AddJavaScriptFile("https://stackpath.bootstrapcdn.com/bootstrap/4.1.1/js/bootstrap.min.js", false,
			html.NewAttributes().Set("integrity", "sha384-smHYKdLADwkXOn1EmN1qk/HfnUcbVRZyYmZ4qpPea6sjB/pTJ0euyQp0Mk8ck+5T").Set("crossorigin", "anonymous"))
		form.AddStyleSheetFile("https://stackpath.bootstrapcdn.com/bootstrap/4.1.1/css/bootstrap.min.css",
			html.NewAttributes().Set("integrity", "sha384-WskhaSGFgHYWDcbwN70/dfYBj47jz9qbsMId/iRN3ewGhXQFZCSftd1LZCfmhktB").Set("crossorigin", "anonymous"))

	}
}

func init() {
	page.RegisterCssFile("bootstrap.min.css", config.GoraddDir+"/bootstrap/assets/css/bootstrap.min.css")
	page.RegisterJsFile("bootstrap.bundle.js", config.GoraddDir+"/bootstrap/assets/js/bootstrap.bundle.js")
}