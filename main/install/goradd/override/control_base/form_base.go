package control_base

import (
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	"goradd-project/config"
)

// The local FormBase override. All framework forms descend from this one. You can change how all the forms in your
// application work by making modifications here. This struct is overridden by the one in the control package, and
// so you should descend your forms from that one.
type FormBase struct {
	page.FormBase
}

// You can put overrides that should apply to all your forms here.

func (f *FormBase) AddRelatedFiles() {
	f.AddJavaScriptFile(config.GoraddAssets()+"/js/jquery3.js", false, nil)
	f.AddJavaScriptFile(config.GoraddAssets()+"/js/ajaxq/ajaxq.js", false, nil) // goradd.js needs this
	f.AddJavaScriptFile(config.GoraddAssets()+"/js/goradd.js", false, nil)
	f.AddStyleSheetFile(config.GoraddAssets()+"/css/goradd.css", nil)
	f.AddStyleSheetFile("https://use.fontawesome.com/releases/v5.0.13/css/all.css",
		html.NewAttributes().Set("integrity", "sha384-DNOHZ68U8hZfKXOrtjWvjxusGo9WQnrNx2sqG0tfsghAvtVlRW3tvkXWZh58N9jp").Set("crossorigin", "anonymous"))

	f.AddJavaScriptFile(config.GoraddAssets()+"/js/goradd.js", false, nil)
}

// Pre-register files in case browser tries to load them before they are drawn.
// Also the place to register asset files needed by all forms
func init() {
	page.RegisterCssFile("goradd.css", config.GoraddAssets()+"/css/goradd.css")
	page.RegisterJsFile("jquery3.js", config.GoraddAssets()+"/js/jquery3.js")
	page.RegisterJsFile("ajaxq.js", config.GoraddAssets()+"/js/ajaxq/ajaxq.js")
	page.RegisterJsFile("jquery-ui.js", config.GoraddAssets()+"/js/jquery-ui.js")
	page.RegisterJsFile("goradd.js", config.GoraddAssets()+"/js/goradd.js")
}
