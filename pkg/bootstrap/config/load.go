package config

import (
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"path/filepath"
	"runtime"
)

// Loader is the injected loader. Set it during your application's initialization
// if you want to load bootstrap differently than below.
var Loader func(page.FormI)

// Configuration options for Bootstrap

// LoadBootstrap loads the various bootstrap files required by bootstrap. It is called automatically
// by the bootstrap components, but this gives you an opportunity to customize where the client
// gets the files.
func LoadBootstrap(form page.FormI) {
	if form.Page().HasMetaTag("viewport") {
		// already loaded
		return
	}
	if Loader != nil {
		Loader(form)
	} else {
		form.Page().AddHtmlHeaderTag(html.VoidTag{
			Tag:  "meta",
			Attr: html.NewAttributes().
				AddAttributeValue("name", "viewport").
				AddAttributeValue("content","width=device-width, initial-scale=1, shrink-to-fit=no"),
			})
		form.AddJQuery()
		if config.Release {
			form.AddJavaScriptFile("https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js", false,
				html.NewAttributes().Set("integrity", "sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1").Set("crossorigin", "anonymous"))
			form.AddJavaScriptFile("https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js", false,
				html.NewAttributes().Set("integrity", "sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM").Set("crossorigin", "anonymous"))
			form.AddStyleSheetFile("https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css",
				html.NewAttributes().Set("integrity", "sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T").Set("crossorigin", "anonymous"))
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
