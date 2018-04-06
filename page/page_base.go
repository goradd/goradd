package page

import (
	"context"
	"bytes"
	"strconv"
	"github.com/spekary/goradd"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/util/types"
	"github.com/spekary/goradd/log"
)

type PageRenderStatus int

type PageDrawFunc func(context.Context, PageI, *bytes.Buffer) error

const (
	UNRENDERED PageRenderStatus = iota // Form has not started isRendering
	BEGUN // Form has started isRendering but has not finished
	ENDED // Form isRendering has already been started and finished
)

/*
PageI is the interface for a GoraddPage. A PageI is expected to be a sub-object of the Page struct.

Implementation of the Create function is isRequired, but all other lifecycle functions are optional

 */
type PageI interface {
	goradd.BaseI
	Init(ctx context.Context, path string)
	Restore()


	// Run is a good place to do user authentication. It runs every time the page is invoked, whether from an Ajax call
	// a Server call, or initial page creation. Note that the form controls are not set up yet if its a new form.
	Run()

	// Load is called after the page is loaded from memory, but before any other initializations are done.
	Load()

	// PreRender is called after all actions are executed, but before the form object is wasRendered.
	PreRender()


	// The Validate method validates the form after all its controls are validated. If you have form level validations,
	// implement This method. This method will be called whether the controls were valid or not. valid will be true
	// if the controls passed validation, and false if not, so you know whether validation is going to pass. Return
	// false to invalidate the form and cancel the current request.
	Validate(valid bool) bool

	// Called in response to an invalid form or control, just before handling the invalid state
	Invalid()

	// Exit is called in all situations at the end of processing a call. The form has already been drawn at This point
	// and the results saved. This is a good place to update any per-user settings.
	Exit()

	Form() FormI

	DrawI

	GenerateControlId(givenId string) string

	GetPageBase() *PageBase	// Return its page base composition
	Path() string // The url path corresponding to This page.

	GetControl(id string) ControlI
	addControl(ControlI)
	removeControl(string)

	DrawHeaderTags(context.Context, *bytes.Buffer)
	SetTitle(title string)
	StateId() string

	GoraddTranslator() Translater
	ProjectTranslator() Translater
	DrawFunction() PageDrawFunc

	Encode(e Encoder)(err error)
	Decode(e Decoder)(err error)
}

// Anything that draws into the draw buffer must implement This interface
type DrawI interface {
	Draw (context.Context, *bytes.Buffer) error
}


type PageBase struct {
	goradd.Base
	stateId      string // Id in cache of the page. Needs to be output by form.
	path         string // The path to the page. Form needs to know This so it can make the action tag
	renderStatus PageRenderStatus
	idPrefix 	 string	// For creating unique ids for the app

	controlRegistry *types.OrderedMap
	form            FormI
	idCounter       int
	drawFunc        PageDrawFunc
	title           string	// page title to draw in head tag
	htmlHeaderTags  []html.VoidTag
	responseHeader  map[string]string	// queues up anything to be sent in the response header

	goraddTranslator PageTranslator
	projectTranslator PageTranslator
}

// Initialize the page base. Should be called by a page just after creating PageBase.
func (p *PageBase) Init(ctx context.Context, self PageI, path string) {
	p.Base.Init(self)
	p.path = path
	p.drawFunc = p.this().DrawFunction()
	p.goraddTranslator = PageTranslator{Domain:GoraddDomain}
	p.projectTranslator = PageTranslator{Domain:ProjectDomain}
}

// Restore is called immediately after the page has been unserialized, to restore data that did not get serialized.
func (p *PageBase) Restore() {
	p.drawFunc = p.this().DrawFunction()
	p.form.Restore()
}



func (p *PageBase) this() PageI {
	return p.Self.(PageI)
}

// DrawFunction returns the drawing function. This implementation returns the default. Override to change it.
func (p *PageBase) DrawFunction() PageDrawFunc {
	return PageTmpl
}

func (p *PageBase) GetPageBase() *PageBase {
	return p
}


func (p *PageBase) setStateId(stateId string) {
	p.stateId = stateId
}

// The following are the lifecycle events of the page. The only isRequired one to implement is Create.

func (p *PageBase) runPage(ctx context.Context, buf *bytes.Buffer, isNew bool) (err error) {
	grCtx := GetContext(ctx)

	if grCtx.err != nil {
		panic(grCtx.err)	// If we received an error during the unpacking process, let the deferred code above handle the error.
	}

	p.renderStatus = UNRENDERED

	log.FrameworkDebugf("Run page: %s", grCtx)

	// TODO: Lifecycle calls

	if !isNew {
		p.Form().control().updateValues(grCtx)	// Tell all the controls to update their values.
		// if This is an event response, do the actions associated with the event
		if c := p.GetControl(grCtx.actionControlId); c != nil {
			c.control().doAction(ctx)
		}
	}

	p.ClearResponseHeaders()
	//p.SetResponseHeader("charset", "utf-8")
	if grCtx.RequestMode() == Ajax {
		err = p.DrawAjax(ctx, buf)
		p.SetResponseHeader("Content-Type", "application/json")
	} else if grCtx.RequestMode() == Server || grCtx.RequestMode() == Http {
		//p.SetResponseHeader("Content-Type", "text/html")	// default for web page. Response can change This if drawing something else.
		err = p.Draw(ctx, buf)
	} else {
		// TODO: Implement a hook for the CustomAjax call and/or Rest API calls?
	}

	p.Form().control().writeState(ctx)
	return
}

func (p *PageBase) Run() {
}

func (p *PageBase) Load() {
}

func (p *PageBase) PreRender() {
}

func (p *PageBase) Validate(valid bool) bool {
	return true
}

func (p *PageBase) Invalid()  {
}

func (p *PageBase) Exit()  {
}

// Returns the form for pages that only have one form
func (p *PageBase) Form() FormI {
	//return p.forms.GetAt(0).(FormI)
	return p.form
}

func (p *PageBase) SetForm(f FormI) {
	p.form = f
}

// For pages that have multiple forms, get the form by id
/*
func (p *PageBase) FormById(id string) FormI {
	if id == "" {
		panic("Can't get a form by a blank id")
	} else if !p.forms.Has(id) {
		panic("Unknown form, id: " + id)
	} else {
		return p.forms.Get(id).(FormI)
	}
}
*/

// Draws from the page template. The default should be fine for most situations.
// You can replace the template function with your own, or override This for even more control
func (p *PageBase) Draw(ctx context.Context, buf *bytes.Buffer) (err error) {
	return p.drawFunc(ctx, p.Self.(PageI), buf)
}

func (p *PageBase) DrawHeaderTags(ctx context.Context, buf *bytes.Buffer) {
	if p.title != "" {
		buf.WriteString("  <title>")
		buf.WriteString(p.title)
		buf.WriteString("  </title>\n")
	}

	// draw things like additional meta tags, etc
	if p.htmlHeaderTags != nil {
		for _,tag := range p.htmlHeaderTags {
			buf.WriteString(tag.Render())
		}
	}

	p.Form().DrawHeaderTags(ctx, buf)
}


// Sets the prefix for control ids. Some javascript frameworks (i.e. jQueryMobile) require that control ids
// be unique across the application, vs just in the page, because they create internal caches of control ids. This
// allows you to set a per page prefix that will be added on to all control ids to make them unique across the whole
// application. However, its up to you to make sure the names are unique per page.
func (p *PageBase) SetControlIdPrefix(prefix string) *PageBase {
	p.idPrefix = prefix
	return p
}


// Overridable generator for control ids. This is called through the PageI interface, meaning you can change how This
// is done by simply implementing it in a subclass.
func (p *PageBase) GenerateControlId(given string) string {
	if given != "" {
		return p.idPrefix + given
	}
	p.idCounter++	// id counter defaults to zero, so pre-increment
	return p.idPrefix + "c" + strconv.Itoa(p.idCounter)
}

func (p *PageBase) GetControl(id string) ControlI {
	if id == "" || p.controlRegistry == nil {
		return nil
	}
	return p.controlRegistry.Get(id).(ControlI)
}

// Add the given control to the pathRegistry. Call by the control code whenever a control is created or restored
func (p *PageBase) addControl(control ControlI) {
	id := control.Id()

	if id == "" {
		panic("Control must have an id before being added.")
	}

	if p.controlRegistry == nil {
		p.controlRegistry = types.NewOrderedMap()
	}

	if p.controlRegistry.Has(id) {
		panic("Control id already exists. Control must have a unique id on the page before being added.")
	}

	p.controlRegistry.Set(id, control)

	if control.Parent() == nil {
		if f,ok := control.(FormI); ok {
			if p.form != nil {
				panic ("The Form object for the page has already been set.")
			} else {
				p.form = f
			}
		} else {
			panic("Controls must have a parent.")
		}
	}
}

func (p *PageBase) removeControl(id string) {
	// Execute the javascript to remove the control from the dom if we are in ajax mode
	// TODO: Application::executeSelectorFunction('#' . $objControl->getWrapperId(), 'remove');
	// TODO: Make This a direct command in the ajax renderer


	p.controlRegistry.Remove(id)
}


func (p *PageBase) SetTitle(title string) {
	p.title = title
}

func (p *PageBase) Path() string {
	return p.path
}

func (p *PageBase) StateId() string {
	return p.stateId
}

func (p *PageBase) DrawAjax(ctx context.Context, buf *bytes.Buffer) (err error) {
	err = p.Form().renderAjax(ctx, buf)
	return
}

func (p *PageBase) GoraddTranslator() Translater {
	return &p.goraddTranslator
}

func (p *PageBase) ProjectTranslator() Translater {
	return &p.projectTranslator
}


func (p *PageBase) SetLanguage(l string) {
	p.goraddTranslator.Language = l
	p.projectTranslator.Language = l
}

// MarshalBinary will binary encode the page for the purpose of saving the page in the formstate.
func (p *PageBase) Encode(e Encoder)(err error) {
	if err = e.Encode(p.stateId); err != nil {return}
	if err = e.Encode(p.path); err != nil {return}
	if err = e.Encode(p.idPrefix); err != nil {return}
	if err = e.Encode(p.form); err != nil {return}
	if err = e.Encode(p.idCounter); err != nil {return}
	if err = e.Encode(p.title); err != nil {return}
	if err = e.Encode(p.htmlHeaderTags); err != nil {return}
	if err = e.Encode(p.goraddTranslator); err != nil {return}
	if err = e.Encode(p.projectTranslator); err != nil {return}
	return
}

func (p *PageBase) Decode(d Decoder)(err error) {
	if err = d.Decode(&p.stateId); err != nil {return}
	if err = d.Decode(&p.path); err != nil {return}
	if err = d.Decode(&p.idPrefix); err != nil {return}
	if err = d.Decode(&p.form); err != nil {return}
	if err = d.Decode(&p.idCounter); err != nil {return}
	if err = d.Decode(&p.title); err != nil {return}
	if err = d.Decode(&p.htmlHeaderTags); err != nil {return}
	if err = d.Decode(&p.goraddTranslator); err != nil {return}
	if err = d.Decode(&p.projectTranslator); err != nil {return}
	return
}

func (p *PageBase) AddHtmlHeaderTag(t html.VoidTag) {
	p.htmlHeaderTags = append(p.htmlHeaderTags, t)
}

func (p *PageBase) SetResponseHeader(key, value string) {
	if p.responseHeader == nil {
		p.responseHeader = map[string]string{}
	}
	p.responseHeader[key] = value
}

func (p *PageBase) ClearResponseHeaders() {
	p.responseHeader = nil
}


