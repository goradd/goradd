package control

import (
	"context"
	"github.com/goradd/goradd/pkg/page/control"
)

// The local FormBase override. All generated forms descend from this one. You can change how all the forms in your
// application work by making modifications here, and then making sure all your forms include this one.
type FormBase struct {
	control.FormBase
}

func (f *FormBase) Init(ctx context.Context, id string) {
	f.FormBase.Init(ctx, id)

	// additional initializations. For example, your custom page template.
	//f.Page().SetDrawFunction()
}

// You can put overrides that should apply to all your forms here.
func (f *FormBase) AddRelatedFiles() {
	f.FormBase.AddRelatedFiles() // add default files
	//f.AddJQueryUI()				// Load JQuery UI if needed
	//bsconfig.LoadBootstrap(f) 	// Load Bootstrap if needed

	// Load you own site-wide css and js files below
	//f.AddStyleSheetFile(filepath.Join(config2.ProjectAssets(), "css","my.css"), nil)
}

// AddHeadTags adds tags for the header of the page
func (f *FormBase) AddHeadTags() {
	f.FormBase.AddHeadTags() // call default first

	/* Uncomment this to add a favicon
	f.Page().AddHtmlHeaderTag(
		html.VoidTag{
			Tag: "link",
			Attr: html.Attributes{
				"rel":    "icon",
				"type": "image/x-icon",
				"href": "/favicon.ico",
			},
		},
	)
	*/

	/* Uncomment this for bootstrap
	f.Page().AddHtmlHeaderTag(
		html.VoidTag{
			Tag: "meta",
			Attr: html.Attributes{
				"name":    "viewport",
				"content": "width=device-width, initial-scale=1, shrink-to-fit=no",
			},
		},
	)
	 */
}

// AddJQuery adds the jquery javascript to the form
/* Uncomment this to change how JQuery gets loaded
func (f *FormBase) AddJQuery() {
	if !config.Release {
		f.AddJavaScriptFile(filepath.Join(config.GoraddAssets(), "js", "jquery3.js"), false, nil)
	} else {
		f.AddJavaScriptFile("https://code.jquery.com/jquery-3.3.1.min.js", false,
			html.NewAttributes().Set("integrity", "sha256-FgpCb/KJQlLNfOu91ta32o/NMZxltwRo8QtmkMRdAu8=").
				Set("crossorigin", "anonymous"))
	}
}
*/

// AddJQueryUI adds the JQuery UI javascript to the form. This is not loaded by default, but many add-ons
// use it, so its here for convenience.
/* Uncomment this to change how JQueryUI gets loaded
func (f *FormBase) AddJQueryUI() {
	if !config.Release {
		f.AddJavaScriptFile(filepath.Join(config.GoraddAssets(), "js", "jquery-ui.js"), false, nil)
	} else {
		f.AddJavaScriptFile("https://code.jquery.com/ui/1.12.1/jquery-ui.min.js", false,
			html.NewAttributes().Set("integrity", "sha256-VazP97ZCwtekAsvgPBSUwPFKdrwD3unUfSGVYrahUqU=").
				Set("crossorigin", "anonymous"))
	}
}
*/
