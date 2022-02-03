package config

import (
	_ "github.com/goradd/goradd/pkg/bootstrap/assets"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
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
		form.Page().AddHtmlHeaderTag(html.VoidTag{
			Tag:  "meta",
			Attr: html.NewAttributes().
				AddAttributeValue("name", "viewport").
				AddAttributeValue("content","width=device-width, initial-scale=1, shrink-to-fit=no"),
			})
		if config.Release {
			form.AddJavaScriptFile("https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js", false,
				html.NewAttributes().Set("integrity", "sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1").Set("crossorigin", "anonymous"))
			form.AddStyleSheetFile("https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css",
				html.NewAttributes().Set("integrity", "sha384-1BmE4kWBq78iYhFldvKuhfTAU6auU8tT94WrHftjDbrCEXSU1oBoqyl2QvZ6jIW3").Set("crossorigin", "anonymous"))
		} else {
			form.AddJavaScriptFile(path.Join(config.AssetPrefix, "bootstrap", "js", "bootstrap.bundle.js"), false, nil)
			form.AddStyleSheetFile(path.Join(config.AssetPrefix, "bootstrap", "css", "bootstrap.css"), nil)
		}
	}
}
