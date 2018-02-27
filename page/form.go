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
	f.SetAttribute("action", page.GetPageBase().path)
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
func (c *FormBase) Draw(ctx context.Context, buf *bytes.Buffer) (err error) {
	err = c.this().PreRender(ctx, buf)
	buf.WriteString(`<form ` + c.this().Attributes().String() + ">\n")
	if err = c.this().DrawTemplate(ctx, buf); err != nil {
		return // the template is required
	}
	// Render controls that are marked to auto render if the form did not render them
	for _,ctrl := range c.children {
		if ctrl.ShouldAutoRender() &&
			!ctrl.WasRendered() {

			err = ctrl.Draw(ctx, buf)

			if err != nil {
				break
			}
		}
	}
	buf.WriteString("\n<\form>\n")
	c.this().PostRender(ctx, buf)
	return
}

func (c *FormBase) PreRender(ctx context.Context, buf *bytes.Buffer) (err error) {
	if err = c.Control.PreRender(ctx, buf); err != nil {
		return
	}

	c.SetAttribute("method", "post")
	c.SetAttribute("action", c.page.Path())

	return
}

func (c *FormBase) PostRender(ctx context.Context, buf *bytes.Buffer) (err error) {

	c.drawBodyScriptFiles(ctx, buf)

	// Render control level JavaScript

	var r = &GetContext(ctx).Response

	// Go through all controls and gather up any JS or CSS to run or Form Attributes to modify
	// Controls should use the response to execute commands on controls or execute general javascript.

	c.this().getScripts(r)
	s := r.JavaScript()

	// TODO: Remove jQuery dependency. Probably should attach to a window.load so that initializing functions can get accurate measurements of themselves since domready does not guarantee css has been loaded.
	s = fmt.Sprintf(`<script type="text/javascript">jQuery(document).ready(function($j) { %s; });</script>`, s)
	buf.WriteString(s)

	err = c.Control.PostRender(ctx, buf)

	c.resetFlags()

	return
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
func (c *FormBase) AddJavaScriptFile(path string, forceHeader bool, attributes *html.Attributes) {
	if forceHeader && c.isOnPage {
		panic ("You cannot force a JavaScript file to be in the header if you insert it after the page is drawn.")
	}
	if c.isOnPage {
		if c.importedJavaScripts == nil {
			c.importedJavaScripts = types.NewOrderedMap()
		}
		c.importedJavaScripts.Set(path, attributes)
	} else if forceHeader {
		if c.headerJavaScripts == nil {
			c.headerJavaScripts = types.NewOrderedMap()
		}
		c.headerJavaScripts.Set(path, attributes)
	} else {
		if c.bodyJavaScripts == nil {
			c.bodyJavaScripts = types.NewOrderedMap()
		}
		c.bodyJavaScripts.Set(path, attributes)
	}
}

// Add a javascript file that is a concatenation of other javascript files the system uses.
// This allows you to concatenate and minimize all the javascript files you are using without worrying about
// libraries and controls that are adding the individual files through the AddJavaScriptFile function
func (c *FormBase) AddMasterJavaScriptFile(url string, attributes []string, files []string) {
	// TODO
}

// AddStyleSheetFile registers a StyleSheet file such that it will get loaded on the page.
// The file will be loaded on the page at initial draw in the header, or will be inserted into the file if the page
// is already drawn. The path is either a url, or an internal path to the location of the file
// in the development environment. Development files will automatically get copied to the local assets directory for easy
// deployment and so that the MUX can find the file and serve it (this happens at draw time).
// The attributes will be extra attributes included with the tag,
// which is useful for things like crossorigin and integrity attributes.
func (c *FormBase) AddStyleSheetFile(path string, attributes *html.Attributes) {
	if c.isOnPage {
		if c.importedStyleSheets == nil {
			c.importedStyleSheets = types.NewOrderedMap()
		}
		c.importedStyleSheets.Set(path, attributes)
	} else {
		if c.headerStyleSheets == nil {
			c.headerStyleSheets = types.NewOrderedMap()
		}
		c.headerStyleSheets.Set(path, attributes)
	}
}

// DrawHeaderTags is called by the page drawing routine to draw its header tags
// If you override this, be sure to call this version too
func (c *FormBase) DrawHeaderTags(ctx context.Context, buf *bytes.Buffer) {
	if c.headerStyleSheets != nil {
		c.headerStyleSheets.Range(func (path string, attr interface{}) bool {
			var attributes *html.Attributes = attr.(*html.Attributes)
			if attributes == nil {
				attributes = html.NewAttributes()
			}
			attributes.Set("rel", "stylesheet")
			if path[:4] == "http" {
				attributes.Set("href", path)
			} else {
				_,fileName := filepath.Split(path)
				attributes.Set("href", RegisterAssetFile("/css/" + fileName, path))
			}
			buf.WriteString(html.RenderVoidTag("link", attributes))
			return true
		})
	}
	if c.headerJavaScripts != nil {
		c.headerJavaScripts.Range(func (path string, attr interface{}) bool {
			var attributes *html.Attributes = attr.(*html.Attributes)
			if attributes == nil {
				attributes = html.NewAttributes()
			}
			if path[:4] == "http" {
				attributes.Set("src", path)
			} else {
				_,fileName := filepath.Split(path)
				attributes.Set("src", RegisterAssetFile("/js/" + fileName, path))
			}
			buf.WriteString(html.RenderTag("script", attributes, ""))
			return true
		})
	}
}

func (c *FormBase) drawBodyScriptFiles(ctx context.Context, buf *bytes.Buffer) {
	c.bodyJavaScripts.Range(func (path string, attr interface{}) bool {
		var attributes *html.Attributes = attr.(*html.Attributes)
		if attributes == nil {
			attributes = html.NewAttributes()
		}
		if path[:4] == "http" {
			attributes.Set("src", path)
		} else {
			_,fileName := filepath.Split(path)
			attributes.Set("src", RegisterAssetFile("/js/" + fileName, path))
		}
		buf.WriteString(html.RenderTag("script", attributes, ""))
		return true
	})

}