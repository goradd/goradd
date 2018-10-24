package page

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spekary/gengen/maps"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/log"
	"github.com/spekary/goradd/orm/db"
	"github.com/spekary/goradd/session/location"
	"goradd-project/config"
	"strings"
)

const htmlVarFormstate  = "Goradd__FormState"
const htmlVarParams  = "Goradd__Params"

type FormI interface {
	ControlI	// Note we are not inheriting from localpage here, to avoid import loop and because its not really necessary
	// Create the objects on the form without necessarily initializing them
	Init(ctx context.Context, self FormI, path string, id string)
	// CreateControls(ctx context.Context)
	LoadControls(ctx context.Context)
	// AddRelatedFiles()
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

// FormBase is a base class for the Form class that is in the control package.
// It is the basic form controller structure for the application and also serves as the drawing mechanism for the
// <form> tag in the html output. Normally, you should not descend your forms from here, but rather from the
// control.Form class. You can modify the basic form class by making modifications to the goradd/page/formbase.go file.
type FormBase struct {
	Control
	response Response // don't serialize this

	// serialized lists of related files
	headerStyleSheets   *maps.SliceMap
	importedStyleSheets *maps.SliceMap // when refreshing, these get moved to the headerStyleSheets
	headerJavaScripts   *maps.SliceMap
	bodyJavaScripts     *maps.SliceMap
	importedJavaScripts *maps.SliceMap // when refreshing, these get moved to the bodyJavaScripts
}

func (f *FormBase) Init(ctx context.Context, self FormI, path string, id string) {
	var p = &Page{}
	p.Init(ctx, path)

	f.page = p
	if id == "" {
		panic("Forms must have an id assigned")
	}
	f.Control.id = id
	f.Control.Init(self, nil, id)
	f.Tag = "form"
}

func (f *FormBase) this() FormI {
	return f.Self.(FormI)
}

// AddRelatedFiles adds related javascript and style sheet files. This is the default to get the minimum goradd installation working.,
// The order is important, so if you override this, be sure these files get loaded
// before other files.
func (f *FormBase) AddRelatedFiles() {
	path, attr := config.JQueryPath()
	f.AddJavaScriptFile(path, false, html.NewAttributesFromMap(attr))
	f.AddJavaScriptFile(config.GoraddAssets()+"/js/ajaxq/ajaxq.js", false, nil) // goradd.js needs this
	f.AddJavaScriptFile(config.GoraddAssets()+"/js/goradd.js", false, nil)
	f.AddJavaScriptFile(config.GoraddAssets()+"/js/goradd-ws.js", false, nil)
	f.AddStyleSheetFile(config.GoraddAssets()+"/css/goradd.css", nil)
	f.AddStyleSheetFile("https://use.fontawesome.com/releases/v5.0.13/css/all.css",
		html.NewAttributes().Set("integrity", "sha384-DNOHZ68U8hZfKXOrtjWvjxusGo9WQnrNx2sqG0tfsghAvtVlRW3tvkXWZh58N9jp").Set("crossorigin", "anonymous"))
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
	s += fmt.Sprintf("goradd.initMessagingClient(%d, %d);\n", config.GoraddWebSocketPort, config.GoraddWebSocketTLSPort)
	s += f.response.JavaScript()
	f.response = NewResponse() // Reset
	s = fmt.Sprintf(`<script>jQuery(document).ready(function($j) { %s; });</script>`, s)
	buf.WriteString(s)

	f.this().PostRender(ctx, buf)
	return
}

// outputSqlProfile looks for sql profiling information and sends it to the browser if found
func (f *FormBase) getDbProfile(ctx context.Context) (s string)  {
	if profiles := db.GetProfiles(ctx); profiles != nil {
		for _, profile := range profiles {
			dif := profile.EndTime.Sub(profile.BeginTime)
			sql := strings.Replace(profile.Sql, "\n", "<br />", -1)
			s += fmt.Sprintf(`<p class="profile"><div>Time: %s Begin: %s End: %s</div><div>%s</div></p>`,
				dif.String(), profile.BeginTime.Format("3:04:05.000"), profile.EndTime.Format("3:04:05.000"), sql)
		}
	}
	return
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
	a.Set("action", f.Page().path)
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
// in the development environment.
func (f *FormBase) AddJavaScriptFile(path string, forceHeader bool, attributes *html.Attributes) {
	if forceHeader && f.isOnPage {
		panic("You cannot force a JavaScript file to be in the header if you insert it after the page is drawn.")
	}

	if path[:4] != "http" {
		url := GetAssetUrl(path)

		if url == "" {
			panic(path + " is not in a registered asset directory")
		}
		path = url
	}

	if f.isOnPage {
		if f.importedJavaScripts == nil {
			f.importedJavaScripts = maps.NewSliceMap()
		}
		f.importedJavaScripts.Set(path, attributes)
	} else if forceHeader {
		if f.headerJavaScripts == nil {
			f.headerJavaScripts = maps.NewSliceMap()
		}
		f.headerJavaScripts.Set(path, attributes)
	} else {
		if f.bodyJavaScripts == nil {
			f.bodyJavaScripts = maps.NewSliceMap()
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
// in the development environment. AppModeDevelopment files will automatically get copied to the local assets directory for easy
// deployment and so that the MUX can find the file and serve it (This happens at draw time).
// The attributes will be extra attributes included with the tag,
// which is useful for things like crossorigin and integrity attributes.
func (f *FormBase) AddStyleSheetFile(path string, attributes *html.Attributes) {
	if path[:4] != "http" {
		url := GetAssetUrl(path)

		if url == "" {
			panic(path + " is not in a registered asset directory")
		}
		path = url
	}

	if f.isOnPage {
		if f.importedStyleSheets == nil {
			f.importedStyleSheets = maps.NewSliceMap()
		}
		f.importedStyleSheets.Set(path, attributes)
	} else {
		if f.headerStyleSheets == nil {
			f.headerStyleSheets = maps.NewSliceMap()
		}
		f.headerStyleSheets.Set(path, attributes)
	}
}

// DrawHeaderTags is called by the page drawing routine to draw its header tags
// If you override this, be sure to call this version too
func (f *FormBase) DrawHeaderTags(ctx context.Context, buf *bytes.Buffer) {
	if f.importedStyleSheets != nil {
		if f.headerStyleSheets == nil {
			f.headerStyleSheets = maps.NewSliceMap()
		}
		f.headerStyleSheets.Merge(f.importedStyleSheets)
		f.importedStyleSheets = nil
	}

	if f.headerStyleSheets != nil {
		f.headerStyleSheets.Range(func(path string, attr interface{}) bool {
			var attributes = attr.(*html.Attributes)
			if attributes == nil {
				attributes = html.NewAttributes()
			}
			attributes.Set("rel", "stylesheet")
			attributes.Set("href", path)
			buf.WriteString(html.RenderVoidTag("link", attributes))
			return true
		})
	}

	if f.importedJavaScripts != nil {
		if f.headerJavaScripts == nil {
			f.headerJavaScripts = maps.NewSliceMap()
		}
		f.headerJavaScripts.Merge(f.importedJavaScripts)
		f.importedJavaScripts = nil
	}

	if f.headerJavaScripts != nil {
		f.headerJavaScripts.Range(func(path string, attr interface{}) bool {
			var attributes = attr.(*html.Attributes)
			if attributes == nil {
				attributes = html.NewAttributes()
			}
			attributes.Set("src", path)
			buf.WriteString(html.RenderTag("script", attributes, ""))
			return true
		})
	}
}

func (f *FormBase) drawBodyScriptFiles(ctx context.Context, buf *bytes.Buffer) {
	f.bodyJavaScripts.Range(func(path string, attr interface{}) bool {
		var attributes = attr.(*html.Attributes)
		if attributes == nil {
			attributes = html.NewAttributes()
		}
		attributes.Set("src", path)
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

// AddHeadTags is a lifecycle call that happens when a new form is created. This is where you should call
// AddHtmlHeaderTag or SetTitle on the page to set tags that appear in the <head> tag of the page.
// Head tags cannot be changed after the page is created.
func (f *FormBase) AddHeadTags() {

}

// LoadControls is a lifecycle call that happens after a form is first created. It is the place to initialize the value
// of the controls in the form based on variables sent to the form.
func (f *FormBase) LoadControls(ctx context.Context) {
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

// PushLocation pushes the URL that got us to the current page on to the location stack.
func (f *FormBase) PushLocation(ctx context.Context) {
	grctx := GetContext(ctx)
	location.Push(ctx, grctx.URL.RequestURI())
}

// PopLocation pops the most recent location off of the location stack and goes to that location.
func (f *FormBase) PopLocation(ctx context.Context) {
	if loc := location.Pop(ctx); loc != "" {
		f.ChangeLocation(loc)
	}
}