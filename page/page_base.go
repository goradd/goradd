package page

import (
	"context"
	"bytes"
	"strconv"
	"github.com/spekary/goradd"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/util/types"
)

type PageRenderStatus int

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
	// Run is a good place to do user authentication. It runs every time the page is invoked, whether from an Ajax call
	// a Server call, or initial page creation. Note that the form controls are not set up yet if its a new form.
	Run()

	// Load is called after the page is loaded from memory, but before any other initializations are done.
	Load()

	// PreRender is called after all actions are executed, but before the form object is wasRendered.
	PreRender()


	// The Validate method validates the form after all its controls are validated. If you have form level validations,
	// implement this method. This method will be called whether the controls were valid or not. valid will be true
	// if the controls passed validation, and false if not, so you know whether validation is going to pass. Return
	// false to invalidate the form and cancel the current request.
	Validate(valid bool) bool

	// Called in response to an invalid form or control, just before handling the invalid state
	Invalid()

	// Exit is called in all situations at the end of processing a call. The form has already been drawn at this point
	// and the results saved. This is a good place to update any per-user settings.
	Exit()

	Form() FormI

	DrawI

	GenerateControlId(givenId string) string

	GetPageBase() *PageBase	// Return its page base composition
	Path() string // The url path corresponding to this page.

	GetControl(id string) ControlI
	addControl(ControlI)
	removeControl(string)

	DrawHeaderTags(context.Context, *bytes.Buffer)
	SetTitle(title string)
}

type PageDrawFunc func(context.Context, PageI, *bytes.Buffer) error

// Anything that draws into the draw buffer must implement this interface
type DrawI interface {
	Draw (context.Context, *bytes.Buffer) error
}


type PageBase struct {
	goradd.Base
	stateId      string // Id in cache of the page. Needs to be output by form.
	path         string // The path to the page. Form needs to know this so it can make the action tag
	renderStatus PageRenderStatus
	idPrefix 	 string	// For creating unique ids for the app

	controlRegistry  *types.OrderedMap
	form FormI
	idCounter	 int
	drawFunc PageDrawFunc
	title	string	// page title to draw in head tag
	headerTags []html.VoidTag
}

// Initialize the page base. Should be called by a page just after creating PageBase
func (p *PageBase) Init(self interface{}, path string) {
	p.controlRegistry = types.NewOrderedMap()
	p.Base.Init(self)
	p.path = path
	p.drawFunc = PageTmpl	// The default draw function. Replace it if you want something else.
	p.headerTags = []html.VoidTag{}
}

func (p *PageBase) this() PageI {
	return p.Self.(PageI)
}

func (p *PageBase) GetPageBase() *PageBase {
	return p
}


func (p *PageBase) setStateId(stateId string) {
	p.stateId = stateId
}

// The following are the lifecycle events of the page. The only isRequired one to implement is Create.

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
// You can replace the template function with your own, or override this for even more control
func (p *PageBase) Draw(ctx context.Context, buf *bytes.Buffer) (err error) {
	return p.drawFunc(ctx, p, buf)
}

func (p *PageBase) DrawHeaderTags(ctx context.Context, buf *bytes.Buffer) {
	if p.title != "" {
		buf.WriteString("  <title>")
		buf.WriteString(p.title)
		buf.WriteString("  </title>\n")
	}

	// draw things like additional meta tags, etc
	for _,tag := range p.headerTags {
		buf.WriteString(tag.Render())
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


// Overridable generator for control ids. This is called through the PageI interface, meaning you can change how this
// is done by simply implementing it in a subclass.
func (p *PageBase) GenerateControlId(given string) string {
	if given != "" {
		return p.idPrefix + given
	}
	p.idCounter++	// id counter defaults to zero, so pre-increment
	return p.idPrefix + "c" + strconv.Itoa(p.idCounter)
}

func (p *PageBase) GetControl(id string) ControlI {
	return p.controlRegistry.Get(id).(ControlI)
}
/*
// Gets and draws the named control
func (p *PageBase) DrawControl(ctx context.Context, id string, buf *bytes.Buffer) (err error) {
	if c := p.GetControl(id); c != nil {
		err = c.Draw(ctx, buf)
	} else {
		// TODO: issue warning
	}
	return
}

*/

// Add the given control to the registry
func (p *PageBase) addControl(control ControlI) {
	id := control.Id()

	if id == "" {
		panic("Control must have an id before being added.")
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
	// TODO: Make this a direct command in the ajax renderer


	p.controlRegistry.Remove(id)
}


func (p *PageBase) SetTitle(title string) {
	p.title = title
}

func (p *PageBase) Path() string {
	return p.path
}

