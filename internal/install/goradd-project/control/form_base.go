package control

import (
	"context"
	// bsconfig "github.com/goradd/goradd/pkg/bootstrap/config"
	"github.com/goradd/goradd/pkg/page/control"
)

// FormBase is the local form override. All generated forms descend from this one. You can change how all the forms in your
// application work by making modifications here, and then making sure all your forms include this one.
//
// For example:
//
//	  type MyForm struct {
//			control.FormBase
//			myValue int
//	  }
type FormBase struct {
	control.FormBase
}

// Init initializes the form values. You can use this particular version to define and include a custom page drawing
// template, among other things.
func (f *FormBase) Init(self any, ctx context.Context, id string) {
	f.FormBase.Init(self, ctx, id)

	// additional initializations. For example, your custom page template.
	//f.Page().SetDrawFunction()
}

// AddRelatedFiles is the place to add css, javascript and other files that should be loaded for all forms.
//
// To add css files, call [control.FormBase.AddStyleSheetFile].
//
// To add javascript files, call [control.FormBase.AddJavaScriptFile]
//
// If you are using Bootstrap, you should also call [github.com/goradd/goradd/pkg/bootstrap/config.LoadBootstrap] to initialize it here.
func (f *FormBase) AddRelatedFiles() {
	f.FormBase.AddRelatedFiles() // add default files

	// Load Bootstrap if needed. Its normally added by all bootstrap controls, but if you are
	// creating a form that just uses straight bootstrap css, without using the goradd bootstrap controls
	// then this will make it so that those pages will work too.
	// bsconfig.LoadBootstrap(f)

	// Load your own site-wide css and js files below
	//f.AddStyleSheetFile(filepath.Join(config2.ProjectAssets(), "css","my.css"), nil)
}

// AddHeadTags is the place to add head tags for the header of the page by calling [control.Page.AddHtmlHeaderTag]
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
}
