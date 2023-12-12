package config

import (
	_ "github.com/goradd/goradd/pkg/bootstrap/assets"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/html5tag"
	"path"
)

// Loader is the injected loader. Set it during your application's initialization
// if you want to load bootstrap differently than below.
var Loader func(page.FormI)

// Configuration options for Bootstrap

// LoadBootstrap loads the various asset files required by bootstrap. It is called automatically
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
		form.Page().AddHtmlHeaderTag(html5tag.VoidTag{
			Tag: "meta",
			Attr: html5tag.NewAttributes().
				AddValues("name", "viewport").
				AddValues("content", "width=device-width, initial-scale=1, shrink-to-fit=no"),
		})
		if config.Release {
			form.AddJavaScriptFile("https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/js/bootstrap.bundle.min.js", false,
				html5tag.NewAttributes().Set("integrity", "sha384-C6RzsynM9kWDrMNeT87bh95OGNyZPhcTNXj1NW7RuBCsyN/o0jlpcV8Qyq46cDfL").Set("crossorigin", "anonymous"))
			form.AddStyleSheetFile("https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css",
				html5tag.NewAttributes().Set("integrity", "sha384-T3c6CoIi6uLrA9TneNEoa7RxnatzjcDSCmG1MXxSR1GAsXEV/Dwwykc2MPK8M2HN").Set("crossorigin", "anonymous"))
		} else {
			form.AddJavaScriptFile(path.Join(config.AssetPrefix, "bootstrap", "js", "bootstrap.bundle.js"), false, nil)
			form.AddStyleSheetFile(path.Join(config.AssetPrefix, "bootstrap", "css", "bootstrap.css"), nil)
		}
	}
}
