package page

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/log"
	"github.com/spekary/goradd/orm/db"
	"github.com/spekary/goradd/util/types"
	"goradd/config"
	"path/filepath"
	"strings"
)

const htmlVarFormstate string = "Goradd__FormState"
const htmlVarParams string = "Goradd__Params"

type FormI interface {
	ControlI	// Note we are not inheriting from localpage here, to avoid import loop and because its not really necessary
	// Create the objects on the form without necessarily initializing them
	Init(ctx context.Context, self FormI, path string, id string)
	CreateControls(ctx context.Context)
	LoadControls(ctx context.Context)
	AddRelatedFiles()
	AddHeadTags()
	DrawHeaderTags(ctx context.Context, buf *bytes.Buffer)
	Response() *Response
	renderAjax(ctx context.Context, buf *bytes.Buffer) error
	AddStyleSheetFile(path string, attributes *html.Attributes)
	AddJavaScriptFile(path string, forceHeader bool, attributes *html.Attributes)
	DisplayAlert(ctx context.Context, msg string)

	// Lifecycle calls
	Run(ctx context.Context) error
	Exit(ctx context.Context, err error)
}

// FormBase is the basic form controller structure for the application and also serves as the drawing mechanism for the
// <form> tag in the html output. Normally, you should not descend your forms from here, but rather the version in
// your local goradd directory so that you can easily make modifications to the way forms work in your application.
type FormBase struct {
	Control
	response Response // don't serialize this

	// serialized lists of related files
	headerStyleSheets   *types.OrderedMap
	importedStyleSheets *types.OrderedMap // when refreshing, these get moved to the headerStyleSheets
	headerJavaScripts   *types.OrderedMap
	bodyJavaScripts     *types.OrderedMap
	importedJavaScripts *types.OrderedMap // when refreshing, these get moved to the bodyJavaScripts
}

func (f *FormBase) Init(ctx context.Context, self FormI, path string, id string) {
	var p = &Page{}
	p.Init(ctx, path)

	f.page = p
	if id == "" {
		panic("Forms must have an id assigned")
	}
	f.Control.id = id
	f.Control.Init(self, nil)
	f.Tag = "form"
	self.AddRelatedFiles()
	self.CreateControls(ctx)
	self.LoadControls(ctx)

	/*	TODO: Add a dialog and designer click if in design mode
		            if (defined('QCUBED_DESIGN_MODE') && QCUBED_DESIGN_MODE == 1) {
	                // Attach custom event to dialog to handle right click menu items sent by form

	                $dlg = new Q\ModelConnector\EditDlg ($objClass, 'qconnectoreditdlg');

	                $dlg->addAction(
	                    new Q\Event\On('qdesignerclick'),
	                    new Q\Action\Ajax ('ctlDesigner_Click', null, null, 'ui')
	                );
	            }

	*/
}

func (f *FormBase) this() FormI {
	return f.Self.(FormI)
}

// AddRelatedFiles adds related javascript and style sheet files. Override This to get these files from a different location,
// or to load additional files. The order is important, so if you override This, be sure these files get loaded
// before other files.
func (f *FormBase) AddRelatedFiles() {
	f.AddJavaScriptFile("http://code.jquery.com/jquery-3.3.1.min.js", false, html.NewAttributes().Set("integrity", "sha256-FgpCb/KJQlLNfOu91ta32o/NMZxltwRo8QtmkMRdAu8="))
	f.AddJavaScriptFile(config.GoraddAssets()+"/js/goradd.js", false, nil)
	f.AddStyleSheetFile(config.GoraddAssets()+"/css/goradd.css", nil)
}

// CreateControls is a stub function for you to implement in an overriding object. This is where you will create your
// controls.
func (f *FormBase) CreateControls(ctx context.Context) {
}

// LoadControls is a stub function for you to implement in an overriding object. This is where you would
// initialize your controls to initial values if not the default. Note that you should also call SetSaveState on
// controls here, but only after initializing the control
func (f *FormBase) LoadControls(ctx context.Context) {
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
	if err = f.RenderAutoControls(ctx, buf); err != nil {
		panic(err)
	}

	f.resetDrawingFlags()
	formstate := f.saveState() // From This point on we should not change any controls, just draw

	// Render hidden controls

	// Place holder for postBack and postAjax functions to place their data
	buf.WriteString(`<input type="hidden" name="` + htmlVarParams + `" id="` + htmlVarParams + `" value="" />` + "\n")

	// Serialize and write out the formstate
	buf.WriteString(fmt.Sprintf(`<input type="hidden" name="`+htmlVarFormstate+`" id="`+htmlVarFormstate+`" value="%s" />`, formstate))

	f.drawBodyScriptFiles(ctx, buf) // Fixing a bug?

	buf.WriteString("\n</form>\n")

	// Draw things that come after the form tag

	// Write out the control scripts gathered above
	s := `goradd.initForm();` + "\n"
	s += f.response.JavaScript()
	f.response = NewResponse() // Reset
	s = fmt.Sprintf(`<script>jQuery(document).ready(function($j) { %s; });</script>`, s)
	buf.WriteString(s)

	f.this().PostRender(ctx, buf)

	f.outputSqlProfile(ctx, buf)
	return
}

// outputSqlProfile looks for sql profiling information and sends it to the browser if found
func (f *FormBase) outputSqlProfile(ctx context.Context, buf *bytes.Buffer) {
	if profiles := db.GetProfiles(ctx); profiles != nil {
		var head = `<h4 onclick="$j('#sqlprofilelist').toggle();" style="position:fixed; bottom:0">SQL Profile <i class="fas fa-arrow-circle-down" ></i></h4>`
		var s string
		for _, profile := range profiles {
			dif := profile.EndTime.Sub(profile.BeginTime)
			sql := strings.Replace(profile.Sql, "\n", "<br />", -1)
			s += fmt.Sprintf(`<p class="profile"><div>Time: %s Begin: %s End: %s</div><div>%s</div></p>`,
				dif.String(), profile.BeginTime.Format("3:04:05.000"), profile.EndTime.Format("3:04:05.000"), sql)
		}
		s = html.RenderTag("div", html.NewAttributes().SetID("sqlprofilelist").SetDisplay("none"), s)
		buf.WriteString(html.RenderTag("div", html.NewAttributes().SetID("sqlprofile"), head + s))
	}
}

// Assembles the ajax response for the entire form and draws it to the return buffer
func (f *FormBase) renderAjax(ctx context.Context, buf *bytes.Buffer) (err error) {
	var buf2 []byte
	var pagestate string

	if !f.response.hasExclusiveCommand() { // skip drawing if we are in a high priority situation
		// gather modified controls
		f.DrawAjax(ctx, &f.response)
	}

	pagestate = f.saveState()
	var grctx = GetContext(ctx)
	if pagestate != grctx.pageStateId {
		panic("page state changed")
	}
	//f.response.SetControlValue(htmlVarFormstate, formstate)
	// TODO: render imported style sheets and java scripts
	f.resetDrawingFlags()
	buf2, err = json.Marshal(&f.response)
	f.response = NewResponse() // Reset
	buf.Write(buf2)
	log.FrameworkDebug("renderAjax - ", string(buf2))
	return
}

func (f *FormBase) DrawingAttributes() *html.Attributes {
	a := f.Control.DrawingAttributes()
	a.Set("method", "post")
	a.Set("action", f.Page().GetPageBase().path)
	a.SetDataAttribute("grctl", "form")
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
	var s = f.page.StateID()
	pageCache.Set(s, f.page) // the page should already exist in the cache. This just tells the cache that we used it, so make it current.
	return f.page.StateID()
}

// AddJavaScriptFile registers a JavaScript file such that it will get loaded on the page.
// If forceHeader is true, the file will be listed in the header, which you should only do if the file has some
// preliminary javascript that needs to be executed before the dom loads. Otherwise, the file will be loaded after
// the dom is loaded, allowing the browser to show the page and then load the javascript in the background, giving the
// appearance of a more responsive website. If you add the file during an ajax operation, the file will be loaded
// dynamically by the goradd javascript. Controls generally should call This during the initial creation of the control if the control
// requires additional javascript to function.
//
// The path is either a url, or an internal path to the location of the file
// in the development environment. Development files will automatically get copied to the local assets directory for easy
// deployment and so that the MUX can find the file and serve it (This happens at draw time).
// The attributes are extra attributes included with the tag,
// which is useful for things like crossorigin and integrity attributes.
func (f *FormBase) AddJavaScriptFile(path string, forceHeader bool, attributes *html.Attributes) {
	if forceHeader && f.isOnPage {
		panic("You cannot force a JavaScript file to be in the header if you insert it after the page is drawn.")
	}

	// TODO: decompose path here, rather than at draw time to save some processing time.

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
// deployment and so that the MUX can find the file and serve it (This happens at draw time).
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
	if f.importedStyleSheets != nil {
		if f.headerStyleSheets == nil {
			f.headerStyleSheets = types.NewOrderedMap()
		}
		f.headerStyleSheets.Merge(f.importedStyleSheets)
		f.importedStyleSheets = nil
	}

	if f.headerStyleSheets != nil {
		f.headerStyleSheets.Range(func(path string, attr interface{}) bool {
			var attributes *html.Attributes = attr.(*html.Attributes)
			if attributes == nil {
				attributes = html.NewAttributes()
			}
			attributes.Set("rel", "stylesheet")
			if path[:4] == "http" {
				attributes.Set("href", path)
			} else {
				_, fileName := filepath.Split(path)
				attributes.Set("href", RegisterCssFile(fileName, path))
			}
			buf.WriteString(html.RenderVoidTag("link", attributes))
			return true
		})
	}

	if f.importedJavaScripts != nil {
		if f.headerJavaScripts == nil {
			f.headerJavaScripts = types.NewOrderedMap()
		}
		f.headerJavaScripts.Merge(f.importedJavaScripts)
		f.importedJavaScripts = nil
	}

	if f.headerJavaScripts != nil {
		f.headerJavaScripts.Range(func(path string, attr interface{}) bool {
			var attributes *html.Attributes = attr.(*html.Attributes)
			if attributes == nil {
				attributes = html.NewAttributes()
			}
			if path[:4] == "http" {
				attributes.Set("src", path)
			} else {
				_, fileName := filepath.Split(path)
				attributes.Set("src", RegisterJsFile(fileName, path))
			}
			buf.WriteString(html.RenderTag("script", attributes, ""))
			return true
		})
	}
}

func (f *FormBase) drawBodyScriptFiles(ctx context.Context, buf *bytes.Buffer) {
	f.bodyJavaScripts.Range(func(path string, attr interface{}) bool {
		var attributes *html.Attributes = attr.(*html.Attributes)
		if attributes == nil {
			attributes = html.NewAttributes()
		}
		if path[:4] == "http" {
			attributes.Set("src", path)
		} else {
			_, fileName := filepath.Split(path)
			attributes.Set("src", RegisterJsFile(fileName, path))
		}
		buf.WriteString(html.RenderTag("script", attributes, "") + "\n")
		return true
	})

}

func (f *FormBase) DisplayAlert(ctx context.Context, msg string) {
	f.response.displayAlert(msg)
}

func (f *FormBase) ChangeLocation(url string) {
	f.response.SetLocation(url)
}

// Response returns the form's response object that you can use to queue up javascript commands to the browser to be sent on
// the next ajax or server request
func (f *FormBase) Response() *Response {
	return &f.response
}

// AddHeadTags is a lifecycle call that happens when a new page is created. This is where you should call
// AddHtmlHeaderTag or SetTitle on the page to set tags that appear in the <head> tag of the page.
// Head tags cannot be changed after the page is created.
func (f *FormBase) AddHeadTags() {

}

// Run is a lifecycle function that gets called whenever a page is run, either by a whole page load, or an ajax call.
// Its a good place to validate that the current user should have access to the information on the page.
// Returning an error will result in the error message being displayed.
func (f *FormBase) Run(ctx context.Context) error {
	return nil
}

// Exit is a lifecycle function that gets called after the form is processed, just before control is returned to the client.
// err will be set if an error response was detected.
func (f *FormBase) Exit(ctx context.Context, err error) {
	return
}

func (f *FormBase) Refresh() {
	panic("Do not refresh the form. It cannot be drawn in ajax.")
}
