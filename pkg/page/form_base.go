package page

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"github.com/goradd/gengen/pkg/maps"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/crypt"
	"github.com/goradd/goradd/pkg/goradd"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/session"
	"github.com/goradd/goradd/pkg/session/location"
	"path/filepath"
	"strings"
)

type FormI interface {
	ControlI
	// Init initializes the base structures of the form. Do this before adding controls to the form.
	Init(ctx context.Context, self FormI, path string, id string)
	PageDrawingFunction() PageDrawFunc

	LoadControls(ctx context.Context)
	AddHeadTags()
	DrawHeaderTags(ctx context.Context, buf *bytes.Buffer)
	Response() *Response
	renderAjax(ctx context.Context, buf *bytes.Buffer) error
	AddStyleSheetFile(path string, attributes html.Attributes)
	AddJavaScriptFile(path string, forceHeader bool, attributes html.Attributes)
	DisplayAlert(ctx context.Context, msg string)
	AddJQueryUI()
	ChangeLocation(url string)
	PushLocation(ctx context.Context)
	PopLocation(ctx context.Context, fallback string)

	// Lifecycle calls
	Run(ctx context.Context) error
	Exit(ctx context.Context, err error)

	resetValidation()
	updateValues(ctx *Context)
	writeAllStates(ctx context.Context)
}

// ΩFormBase is a base for the FormBase struct that is in the control package.
// Normally, you should not descend your forms from here, but rather from the control.Form struct.
// It is the basic control structure for the application and also serves as the drawing mechanism for the
// <form> tag in the html output.
type ΩFormBase struct {
	Control
	response Response
	headerStyleSheets   *maps.SliceMap
	importedStyleSheets *maps.SliceMap // when refreshing, these get moved to the headerStyleSheets
	headerJavaScripts   *maps.SliceMap
	bodyJavaScripts     *maps.SliceMap
	importedJavaScripts *maps.SliceMap // when refreshing, these get moved to the bodyJavaScripts
}

// Init initializes the form control. Note that ctx might be nil if we are unit testing.
func (f *ΩFormBase) Init(ctx context.Context, self FormI, path string, id string) {
	var p = &Page{}
	p.Init()

	f.page = p
	if id == "" {
		panic("Forms must have an id assigned")
	}
	f.Control.id = id
	f.Control.Init(self, nil, id)
	f.Tag = "form"
}

func (f *ΩFormBase) this() FormI {
	return f.Self.(FormI)
}

// AddRelatedFiles adds related javascript and style sheet files. This is the default to get the minimum goradd installation working.,
// The order is important, so if you override this, be sure these files get loaded
// before other files.
func (f *ΩFormBase) AddRelatedFiles() {
	f.AddGoraddFiles()
	f.AddFontAwesome()
}

// AddJQuery adds the jquery javascript to the form
func (f *ΩFormBase) AddJQuery() {
	if !config.Release {
		f.AddJavaScriptFile(filepath.Join(config.GoraddAssets(), "js", "jquery3.js"), false, nil)
	} else {
		f.AddJavaScriptFile("https://code.jquery.com/jquery-3.3.1.min.js", false,
			html.NewAttributes().Set("integrity", "sha256-FgpCb/KJQlLNfOu91ta32o/NMZxltwRo8QtmkMRdAu8=").
				Set("crossorigin", "anonymous"))
	}
}

// AddJQueryUI adds the JQuery UI javascript to the form. This is not loaded by default, but many add-ons
// use it, so its here for convenience.
func (f *ΩFormBase) AddJQueryUI() {
	if !config.Release {
		f.AddJavaScriptFile(filepath.Join(config.GoraddAssets(), "js", "jquery-ui.js"), false, nil)
	} else {
		f.AddJavaScriptFile("https://code.jquery.com/ui/1.12.1/jquery-ui.min.js", false,
			html.NewAttributes().Set("integrity", "sha256-VazP97ZCwtekAsvgPBSUwPFKdrwD3unUfSGVYrahUqU=").
				Set("crossorigin", "anonymous"))
	}
}

// AddGoraddFiles adds the various goradd files to the form
func (f *ΩFormBase) AddGoraddFiles() {
	gr := config.GoraddAssets()
	f.AddJavaScriptFile(filepath.Join(gr, "js", "goradd.js"), false, nil)
	f.AddJavaScriptFile(filepath.Join(gr, "js", "goradd-ws.js"), false, nil)
	if !config.Release {
		f.AddJavaScriptFile(filepath.Join(gr, "js", "goradd-testing.js"), false, nil)
	}
	f.AddStyleSheetFile(filepath.Join(gr, "css", "goradd.css"), nil)
}

// AddFontAwesome adds the font-awesome files fo the form
func (f *ΩFormBase) AddFontAwesome() {
	f.AddStyleSheetFile("https://use.fontawesome.com/releases/v5.0.13/css/all.css",
		html.NewAttributes().Set("integrity", "sha384-DNOHZ68U8hZfKXOrtjWvjxusGo9WQnrNx2sqG0tfsghAvtVlRW3tvkXWZh58N9jp").Set("crossorigin", "anonymous"))
}

// Draw renders the form. Even though forms are technically controls, we use a custom drawing
// routine for performance reasons and for control.
func (f *ΩFormBase) Draw(ctx context.Context, buf *bytes.Buffer) (err error) {
	err = f.this().ΩPreRender(ctx, buf)
	buf.WriteString(`<form ` + f.this().ΩDrawingAttributes(ctx).String() + ">\n")
	if err = f.this().DrawTemplate(ctx, buf); err != nil {
		return // the template is required
	}
	// Render controls that are marked to auto render if the form did not render them
	if err = f.RenderAutoControls(ctx, buf); err != nil {
		panic(err)
	}

	f.resetDrawingFlags()

	// Render hidden controls

	// Place holder for postBack and postAjax functions to place their data
	buf.WriteString(`<input type="hidden" name="` + htmlVarParams + `" id="` + htmlVarParams + `" value="" />` + "\n")

	// CSRF prevention
	var csrf string

	csrf, _ = session.GetString(ctx, goradd.SessionCsrf)
	if csrf == "" {
		// first time
		csrf, err = crypt.GenerateRandomString(16)
		if err != nil {
			return err
		}
		session.Set(ctx, goradd.SessionCsrf, csrf)
	}
	buf.WriteString(fmt.Sprintf(`<input type="hidden" name="`+htmlCsrfToken+`" id="`+htmlCsrfToken+`" value="%s" />`+"\n", csrf))

	// Serialize and write out the pagestate
	buf.WriteString(fmt.Sprintf(`<input type="hidden" name="`+HtmlVarPagestate+`" id="`+HtmlVarPagestate+`" value="%s" />`, f.page.StateID()))

	f.drawBodyScriptFiles(ctx, buf) // Fixing a bug?

	buf.WriteString("\n</form>\n")

	// Draw things that come after the form tag

	// Write out the control scripts gathered above
	s := `goradd.initForm();` + "\n"
	s += fmt.Sprintf("goradd.initMessagingClient(%d, %d);\n", config.WebSocketPort, config.WebSocketTLSPort)
	f.GetActionScripts(&f.response) // actions assigned to form during form creation
	s += f.response.JavaScript()
	f.response = NewResponse() // clear response
	if !config.Release {
		// This code registers the form with the test harness. We do not want to do this in release mode since it is a security risk.
		s += "goradd.initFormTest();\n"
	}
	s = fmt.Sprintf(`<script>
%s
</script>`, s)
	buf.WriteString(s)

	f.this().ΩPostRender(ctx, buf)
	return
}

func (f *ΩFormBase) resetDrawingFlags() {
	f.RangeSelfAndAllChildren(func(ctrl ControlI) {
		c := ctrl.control()
		c.wasRendered = false
		c.isModified = false
	})
}

func (f *ΩFormBase) resetValidation() {
	f.RangeSelfAndAllChildren(func(ctrl ControlI) {
		c := ctrl.control()
		if c.validationMessage != "" {
			c.validationMessage = ""
		}
		if c.validationState != ValidationWaiting {
			c.validationState = ValidationWaiting
		}
	})
}

func (f *ΩFormBase) updateValues(ctx *Context) {
	f.RangeAllChildren(func(child ControlI) {
		// Parent is updated after children so that parent can read the state of the children
		// to update any internal caching of the state. Parent can then delete or recreate children
		// as needed.
		child.ΩUpdateFormValues(ctx)
	})
}

// writeAllStates is an internal function that will recursively write out the state of all the controls.
// This state is used by controls to restore the visual state of the control if the page is returned to. This is helpful
// in situations where a control is used to filter what is shown on the page, you zoom into an item, and then return to
// the parent control. In this situation, you want to see things in the same state they were in, and not have to set up
// the filter all over again.
func (f *ΩFormBase) writeAllStates(ctx context.Context) {
	f.RangeAllChildren(func(child ControlI) {
		c := child.control()
		c.writeState(ctx)
	})
}



// outputSqlProfile looks for sql profiling information and sends it to the browser if found
func (f *ΩFormBase) getDbProfile(ctx context.Context) (s string) {
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

// renderAjax assembles the ajax response for the entire form and draws it to the return buffer
func (f *ΩFormBase) renderAjax(ctx context.Context, buf *bytes.Buffer) (err error) {
	var buf2 []byte

	if !f.response.hasExclusiveCommand() { // skip drawing if we are in a high priority situation
		// gather modified controls
		err = f.DrawAjax(ctx, &f.response)
		if err != nil {
			log.Error("renderAjax error - " + err.Error())
			// savestate ???
			return
		}
	}

	// Inject any added style sheets and script files
	if f.importedStyleSheets != nil {
		f.importedStyleSheets.Range(func(k string,v interface{}) bool {
			f.response.addStyleSheet(k,v.(html.Attributes))
			return true
		})
	}

	if f.importedJavaScripts != nil {
		f.importedJavaScripts.Range(func(k string,v interface{}) bool {
			f.response.addJavaScriptFile(k,v.(html.Attributes))
			return true
		})
	}

	f.mergeInjectedFiles()

	f.resetDrawingFlags()
	buf2, err = f.response.GetAjaxResponse()
	f.response = NewResponse()
	buf.Write(buf2)
	log.FrameworkDebug("renderAjax - ", string(buf2))

	return
}

// ΩDrawingAttributes returns the attributes to add to the form tag.
func (f *ΩFormBase) ΩDrawingAttributes(ctx context.Context) html.Attributes {
	a := f.Control.ΩDrawingAttributes(ctx)
	a.SetDataAttribute("grctl", "form")
	return a
}

// ΩPreRender performs setup operations just before drawing.
func (f *ΩFormBase) ΩPreRender(ctx context.Context, buf *bytes.Buffer) (err error) {
	if err = f.Control.ΩPreRender(ctx, buf); err != nil {
		return
	}

	f.SetAttribute("method", "post")
	// Setting the "action" attribute prevents iFrame clickjacking.
	// This only works because we never ajax draw the form, only server render
	f.SetAttribute("action", GetContext(ctx).HttpContext.URL.RequestURI())

	return
}

// PageDrawingFunction returns the function used to draw the page object.
// If you want a custom drawing function for your page, implement this function in your form override.
func (f *ΩFormBase) PageDrawingFunction() PageDrawFunc {
	return PageTmpl // Returns the default
}

// AddJavaScriptFile registers a JavaScript file such that it will get loaded on the page.
// If forceHeader is true, the file will be listed in the header, which you should only do if the file has some
// preliminary javascript that needs to be executed before the dom loads.
// You can specify forceHeader and a "defer" attribute to get the effect of loading the javascript in the background.
// With forceHeader false, the file will be loaded after
// the dom is loaded, allowing the browser to show the page and then load the javascript in the background, giving the
// appearance of a more responsive website. If you add the file during an ajax operation, the file will be loaded
// dynamically by the goradd javascript. Controls generally should call this during the initial creation of the control if the control
// requires additional javascript to function.
//
// The path is either a url, or an internal path to the location of the file
// in the development environment.
func (f *ΩFormBase) AddJavaScriptFile(path string, forceHeader bool, attributes html.Attributes) {
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
		if f.headerJavaScripts != nil && f.headerJavaScripts.Has(path) ||
			f.bodyJavaScripts != nil && f.bodyJavaScripts.Has(path) {
			return // file is already on the page
		}
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

// AddMasterJavaScriptFile adds a javascript file that is a concatenation of other javascript files the system uses.
// This allows you to concatenate and minimize all the javascript files you are using without worrying about
// libraries and controls that are adding the individual files through the AddJavaScriptFile function
func (f *ΩFormBase) AddMasterJavaScriptFile(url string, attributes []string, files []string) {
	// TODO
}

// AddStyleSheetFile registers a StyleSheet file such that it will get loaded on the page.
// The file will be loaded on the page at initial draw in the header, or will be inserted into the file if the page
// is already drawn. The path is either a url, or an internal path to the location of the file
// in the development environment. AppModeDevelopment files will automatically get copied to the local assets directory for easy
// deployment and so that the MUX can find the file and serve it (This happens at draw time).
// The attributes will be extra attributes included with the tag,
// which is useful for things like crossorigin and integrity attributes.
func (f *ΩFormBase) AddStyleSheetFile(path string, attributes html.Attributes) {
	if path[:4] != "http" {
		url := GetAssetUrl(path)

		if url == "" {
			panic(path + " is not in a registered asset directory")
		}
		path = url
	}

	if f.isOnPage {
		if f.headerStyleSheets != nil && f.headerStyleSheets.Has(path) {
			return // the style sheet was already included when the form was loaded the first time
		}
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
func (f *ΩFormBase) DrawHeaderTags(ctx context.Context, buf *bytes.Buffer) {
	f.mergeInjectedFiles()

	if f.headerStyleSheets != nil {
		f.headerStyleSheets.Range(func(path string, attr interface{}) bool {
			var attributes = attr.(html.Attributes)
			if attributes == nil {
				attributes = html.NewAttributes()
			}
			attributes.Set("rel", "stylesheet")
			attributes.Set("href", path)
			buf.WriteString(html.RenderVoidTag("link", attributes))
			return true
		})
	}

	if f.headerJavaScripts != nil {
		f.headerJavaScripts.Range(func(path string, attr interface{}) bool {
			var attributes = attr.(html.Attributes)
			if attributes == nil {
				attributes = html.NewAttributes()
			}
			attributes.Set("src", path)
			buf.WriteString(html.RenderTag("script", attributes, ""))
			return true
		})
	}
}

func (f *ΩFormBase) mergeInjectedFiles() {
	if f.importedStyleSheets != nil {
		if f.headerStyleSheets == nil {
			f.headerStyleSheets = maps.NewSliceMap()
		}
		f.headerStyleSheets.Merge(f.importedStyleSheets)
		f.importedStyleSheets = nil
	}

	if f.importedJavaScripts != nil {
		if f.headerJavaScripts == nil {
			f.headerJavaScripts = maps.NewSliceMap()
		}
		f.headerJavaScripts.Merge(f.importedJavaScripts)
		f.importedJavaScripts = nil
	}
}

func (f *ΩFormBase) drawBodyScriptFiles(ctx context.Context, buf *bytes.Buffer) {
	f.bodyJavaScripts.Range(func(path string, attr interface{}) bool {
		var attributes = attr.(html.Attributes)
		if attributes == nil {
			attributes = html.NewAttributes()
		}
		attributes.Set("src", path)
		buf.WriteString(html.RenderTag("script", attributes, "") + "\n")
		return true
	})

}

// DisplayAlert will display a javascript alert with the given message.
func (f *ΩFormBase) DisplayAlert(ctx context.Context, msg string) {
	f.response.displayAlert(msg)
}

// ChangeLocation will redirect the browser to a new URL. It does this AFTER processing the return
// values sent to the browser. Generally you should use this to redirect the browser since you may
// have some data that needs to be processed first. The exception is
// if you are responding to some kind of security concern where you only want to send back an html
// redirect without revealing any goradd information, in which case you should use the Page
func (f *ΩFormBase) ChangeLocation(url string) {
	f.response.SetLocation(url)
}

// Response returns the form's response object that you can use to queue up javascript commands to the browser to be sent on
// the next ajax or server request
func (f *ΩFormBase) Response() *Response {
	return &f.response
}

// AddHeadTags is a lifecycle call that happens when a new form is created. This is where you should call
// AddHtmlHeaderTag or SetTitle on the page to set tags that appear in the <head> tag of the page.
// Head tags cannot be changed after the page is created.
func (f *ΩFormBase) AddHeadTags() {

}

// LoadControls is a lifecycle call that happens after a form is first created. It is the place to initialize the value
// of the controls in the form based on variables sent to the form.
func (f *ΩFormBase) LoadControls(ctx context.Context) {
}

// Run is a lifecycle function that gets called whenever a page is run, either by a whole page load, or an ajax call.
// Its a good place to validate that the current user should have access to the information on the page.
// Returning an error will result in the error message being displayed.
func (f *ΩFormBase) Run(ctx context.Context) error {
	return nil
}

// Exit is a lifecycle function that gets called after the form is processed, just before control is returned to the client.
// err will be set if an error response was detected.
func (f *ΩFormBase) Exit(ctx context.Context, err error) {
	return
}

func (f *ΩFormBase) Refresh() {
	panic("Do not refresh the form. It cannot be drawn in ajax.")
}

// PushLocation pushes the URL that got us to the current page on to the location stack.
func (f *ΩFormBase) PushLocation(ctx context.Context) {
	grctx := GetContext(ctx)
	location.Push(ctx, grctx.URL.RequestURI())
}

// PopLocation pops the most recent location off of the location stack and goes to that location.
// It will go to the fallback url if there is nothing on the stack
func (f *ΩFormBase) PopLocation(ctx context.Context, fallback string) {
	if loc := location.Pop(ctx); loc != "" {
		f.ChangeLocation(loc)
	} else {
		f.ChangeLocation(fallback)
	}
}

type formEncoded struct {
	HeaderSS   *maps.SliceMap
	ImportedSS *maps.SliceMap
	HeaderJS   *maps.SliceMap
	BodyJS     *maps.SliceMap
	ImportedJS *maps.SliceMap
}

func (f *ΩFormBase) Serialize(e Encoder) (err error) {
	if err = f.Control.Serialize(e); err != nil {
		return
	}

	if !config.Release {
		// The response is currently only changed between posts by the testing framework
		// If we ever need to change forms using some kind of push mechanism, we will need to serialize
		// the response.
		if err = f.response.Serialize(e); err != nil {
			return
		}
	}

	s := formEncoded{
		HeaderSS:   f.headerStyleSheets,
		ImportedSS: f.importedStyleSheets,
		HeaderJS:   f.headerJavaScripts,
		BodyJS:     f.bodyJavaScripts,
		ImportedJS: f.importedJavaScripts,
	}

	if err = e.Encode(s); err != nil {
		return
	}
	return
}

func (f *ΩFormBase) Deserialize(d Decoder) (err error) {
	if err = f.Control.Deserialize(d); err != nil {
		return
	}

	if !config.Release {
		// The response is currently only changed between posts by the testing framework
		// If we ever need to change forms using some kind of push mechanism, we will need to serialize
		// the response.
		if err = f.response.Deserialize(d); err != nil {
			return
		}
	}


	s := formEncoded{}
	if err = d.Decode(&s); err != nil {
		return
	}

	f.headerStyleSheets = s.HeaderSS
	f.importedStyleSheets = s.ImportedSS
	f.headerJavaScripts = s.HeaderJS
	f.bodyJavaScripts = s.BodyJS
	f.importedJavaScripts = s.ImportedJS

	return
}

func init() {
	gob.Register(&ΩFormBase{})
}
