package page

import (
	"context"
	"bytes"
	"goradd/config"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/util/types"
	"path/filepath"
	"fmt"
)


type FormI interface {
	ControlI
	// Create the objects on the form without necessarily initializing them
	CreateControls(ctx context.Context)
	InitializeControls(ctx context.Context)
	AddRelatedFiles()
	DrawHeaderTags(ctx context.Context, buf *bytes.Buffer)
}

type FormBase struct {
	Control
	response Response

	// serialized lists of related files
	headerStyleSheets *types.OrderedMap
	importedStyleSheets *types.OrderedMap // when refreshing, these get moved to the headerStyleSheets
	headerJavaScripts *types.OrderedMap
	bodyJavaScripts *types.OrderedMap
	importedJavaScripts *types.OrderedMap // when refreshing, these get moved to the bodyJavaScripts
}

func (f *FormBase) Init(self FormI, page PageI, id string) {
	f.page = page
	f.Control.Init(self, nil, id)
	f.Tag = "form"
	self.AddRelatedFiles()
}

// AddRelatedFiles adds related javascript and style sheet files. Override this to get these files from a different location,
// or to load additional files. The order is important, so if you override this, be sure these files get loaded
// before other files.
func (f *FormBase) AddRelatedFiles() {
	f.AddJavaScriptFile("http://code.jquery.com/jquery-3.3.1.min.js", false, html.NewAttributes().Set("integrity", "sha256-FgpCb/KJQlLNfOu91ta32o/NMZxltwRo8QtmkMRdAu8="))
	f.AddJavaScriptFile(config.GoraddAssets() + "/js/goradd.js", false, nil)
	f.AddStyleSheetFile(config.GoraddAssets() + "/css/goradd.css", nil)
}

func (f *FormBase) CreateControls(ctx context.Context) {
}

func (f *FormBase) InitializeControls(ctx context.Context) {
}


// Draw renders the form. Even though forms are technically controls, we use a custom drawing
// routine for performance reasons and for control.
func (f *FormBase) Draw(ctx context.Context, buf *bytes.Buffer) (err error) {
	err = f.this().PreRender(ctx, buf)
	buf.WriteString(`<form ` + f.this().DrawingAttributes().String() + ">\n")
	if err = f.this().DrawTemplate(ctx, buf); err != nil {
		return // the template is required
	}
	// Render controls that are marked to auto render if the form did not render them
	for _,ctrl := range f.children {
		if ctrl.ShouldAutoRender() &&
			!ctrl.WasRendered() {

			err = ctrl.Draw(ctx, buf)

			if err != nil {
				break
			}
		}
	}

	// Go through all controls and gather up any JS or CSS to run or Form Attributes to modify
	// Controls should use the response to execute commands on controls or execute general javascript.
	// This must be done before we save the form state
	f.this().getScripts(&f.response)
	f.resetFlags()
	formstate := f.saveState() // From this point on we should not change any controls, just draw

	// Render hidden controls

	// Place holder for postBack and postAjax functions to place their data
	buf.WriteString(`<input type="hidden" name="Goradd__Params" id="Goradd__Params" value="" />` + "\n")

	// Serialize and write out the formstate
	buf.WriteString(fmt.Sprintf(`<input type="hidden" name="Goradd__FormState" id="Goradd__Formstate" value="%s" />`, formstate))

	f.drawBodyScriptFiles(ctx, buf)	// Fixing a bug?

	buf.WriteString("\n</form>\n")

	// Draw things that come after the form tag


	// Write out the control scripts gathered above
	s := `goradd.initForm();` + "\n";
	s += f.response.JavaScript()
	f.response = NewResponse()	// Reset
	s = fmt.Sprintf(`<script>jQuery(document).ready(function($j) { %s; });</script>`, s)
	buf.WriteString(s)

	f.this().PostRender(ctx, buf)

	return
}

func (f *FormBase) DrawingAttributes() *html.Attributes {
	a := f.Control.DrawingAttributes()
	a.Set("action", f.Page().GetPageBase().path)
	a.Set("data-goradd", "form")
	return a
}

func (f *FormBase) PreRender(ctx context.Context, buf *bytes.Buffer) (err error) {
	if err = f.Control.PreRender(ctx, buf); err != nil {
		return
	}

	f.SetAttribute("method", "post")
	f.SetAttribute("action", f.page.Path())

	return
}

// saveState saves the state of the form in the page cache.
// This version keeps the page in memory. Future versions may serialize formstates to store them on disk.
func (f *FormBase) saveState() string {
	var s = f.page.StateId()
	GetPageManager().cache.Set(s, f.page) // the page should already exist in the cache. This just tells the cache that we used it, so make it current.
	return f.page.StateId()
}

// AddJavaScriptFile registers a JavaScript file such that it will get loaded on the page.
// If forceHeader is true, the file will be listed in the header, which you should only do if the file has some
// preliminary javascript that needs to be executed before the dom loads. Otherwise, the file will be loaded after
// the dom is loaded, allowing the browser to show the page and then load the javascript in the background, giving the
// appearance of a more responsive website. If you add the file during an ajax operation, the file will be loaded
// dynamically by the goradd javascript. Controls generally should call this during the initial creation of the control if the control
// requires additional javascript to function.
//
// The path is either a url, or an internal path to the location of the file
// in the development environment. Development files will automatically get copied to the local assets directory for easy
// deployment and so that the MUX can find the file and serve it (this happens at draw time).
// The attributes are extra attributes included with the tag,
// which is useful for things like crossorigin and integrity attributes.
func (f *FormBase) AddJavaScriptFile(path string, forceHeader bool, attributes *html.Attributes) {
	if forceHeader && f.isOnPage {
		panic ("You cannot force a JavaScript file to be in the header if you insert it after the page is drawn.")
	}
	if f.isOnPage {
		if f.importedJavaScripts == nil {
			f.importedJavaScripts = types.NewOrderedMap()
		}
		f.importedJavaScripts.Set(path, attributes)
	} else if forceHeader {
		if f.headerJavaScripts == nil {
			f.headerJavaScripts = types.NewOrderedMap()
		}
		f.headerJavaScripts.Set(path, attributes)
	} else {
		if f.bodyJavaScripts == nil {
			f.bodyJavaScripts = types.NewOrderedMap()
		}
		f.bodyJavaScripts.Set(path, attributes)
	}
}

// Add a javascript file that is a concatenation of other javascript files the system uses.
// This allows you to concatenate and minimize all the javascript files you are using without worrying about
// libraries and controls that are adding the individual files through the AddJavaScriptFile function
func (f *FormBase) AddMasterJavaScriptFile(url string, attributes []string, files []string) {
	// TODO
}

// AddStyleSheetFile registers a StyleSheet file such that it will get loaded on the page.
// The file will be loaded on the page at initial draw in the header, or will be inserted into the file if the page
// is already drawn. The path is either a url, or an internal path to the location of the file
// in the development environment. Development files will automatically get copied to the local assets directory for easy
// deployment and so that the MUX can find the file and serve it (this happens at draw time).
// The attributes will be extra attributes included with the tag,
// which is useful for things like crossorigin and integrity attributes.
func (f *FormBase) AddStyleSheetFile(path string, attributes *html.Attributes) {
	if f.isOnPage {
		if f.importedStyleSheets == nil {
			f.importedStyleSheets = types.NewOrderedMap()
		}
		f.importedStyleSheets.Set(path, attributes)
	} else {
		if f.headerStyleSheets == nil {
			f.headerStyleSheets = types.NewOrderedMap()
		}
		f.headerStyleSheets.Set(path, attributes)
	}
}

// DrawHeaderTags is called by the page drawing routine to draw its header tags
// If you override this, be sure to call this version too
func (f *FormBase) DrawHeaderTags(ctx context.Context, buf *bytes.Buffer) {
	if f.headerStyleSheets != nil {
		f.headerStyleSheets.Range(func (path string, attr interface{}) bool {
			var attributes *html.Attributes = attr.(*html.Attributes)
			if attributes == nil {
				attributes = html.NewAttributes()
			}
			attributes.Set("rel", "stylesheet")
			if path[:4] == "http" {
				attributes.Set("href", path)
			} else {
				_,fileName := filepath.Split(path)
				attributes.Set("href", RegisterCssFile(fileName, path))
			}
			buf.WriteString(html.RenderVoidTag("link", attributes))
			return true
		})
	}
	if f.headerJavaScripts != nil {
		f.headerJavaScripts.Range(func (path string, attr interface{}) bool {
			var attributes *html.Attributes = attr.(*html.Attributes)
			if attributes == nil {
				attributes = html.NewAttributes()
			}
			if path[:4] == "http" {
				attributes.Set("src", path)
			} else {
				_,fileName := filepath.Split(path)
				attributes.Set("src", RegisterJsFile(fileName, path))
			}
			buf.WriteString(html.RenderTag("script", attributes, ""))
			return true
		})
	}
}

func (f *FormBase) drawBodyScriptFiles(ctx context.Context, buf *bytes.Buffer) {
	f.bodyJavaScripts.Range(func (path string, attr interface{}) bool {
		var attributes *html.Attributes = attr.(*html.Attributes)
		if attributes == nil {
			attributes = html.NewAttributes()
		}
		if path[:4] == "http" {
			attributes.Set("src", path)
		} else {
			_,fileName := filepath.Split(path)
			attributes.Set("src", RegisterJsFile(fileName, path))
		}
		buf.WriteString(html.RenderTag("script", attributes, "") + "\n")
		return true
	})

}

func (f *FormBase) DisplayAlert(ctx context.Context, msg string) {
	f.response.displayAlert(msg)
}