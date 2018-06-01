package page

import (
	"context"
	"bytes"
	"strconv"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/util/types"
	"github.com/spekary/goradd/log"
	"fmt"
)

type PageRenderStatus int

type PageDrawFunc func(context.Context, *Page, *bytes.Buffer) error

const (
	UNRENDERED PageRenderStatus = iota // Form has not started isRendering
	BEGUN // Form has started isRendering but has not finished
	ENDED // Form isRendering has already been started and finished
)

// Anything that draws into the draw buffer must implement This interface
type DrawI interface {
	Draw (context.Context, *bytes.Buffer) error
}


type Page struct {
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
func (p *Page) Init(ctx context.Context, path string) {
	p.path = path
	p.drawFunc = p.DrawFunction()
	p.goraddTranslator = PageTranslator{Domain:GoraddDomain}
	p.projectTranslator = PageTranslator{Domain:ProjectDomain}
}

// Restore is called immediately after the page has been unserialized, to restore data that did not get serialized.
func (p *Page) Restore() {
	p.drawFunc = p.DrawFunction()
	p.form.Restore()
}

// DrawFunction returns the drawing function. This implementation returns the default. Override to change it.
func (p *Page) DrawFunction() PageDrawFunc {
	return PageTmpl
}

func (p *Page) SetDrawFunction(f PageDrawFunc) {
	p.drawFunc = f
}

func (p *Page) GetPageBase() *Page {
	return p
}

func (p *Page) setStateID(stateId string) {
	p.stateId = stateId
}


func (p *Page) runPage(ctx context.Context, buf *bytes.Buffer, isNew bool) (err error) {
	grCtx := GetContext(ctx)
	grCtx.WasHandled = true		// Notify listeners that the app handled the page

	if err = p.Form().Run(ctx); err != nil {
		return err
	}

	if grCtx.err != nil {
		panic(grCtx.err)
	}

	p.renderStatus = UNRENDERED

	log.FrameworkDebugf("Run page: %s", grCtx)

	// TODO: Lifecycle calls - push them to the form

	if isNew {
		p.Form().AddHeadTags()
	} else {
		p.Form().control().updateValues(grCtx)	// Tell all the controls to update their values.
		// if This is an event response, do the actions associated with the event
		if c := p.GetControl(grCtx.actionControlID); c != nil {
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
	p.Form().Exit(ctx, err)
	return
}


// Returns the form for pages that only have one form
func (p *Page) Form() FormI {
	//return p.forms.GetAt(0).(FormI)
	return p.form
}

func (p *Page) SetForm(f FormI) {
	p.form = f
}

// For pages that have multiple forms, get the form by id
/*
func (p *Page) FormByID(id string) FormI {
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
func (p *Page) Draw(ctx context.Context, buf *bytes.Buffer) (err error) {
	return p.drawFunc(ctx, p, buf)
}

func (p *Page) DrawHeaderTags(ctx context.Context, buf *bytes.Buffer) {
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
func (p *Page) SetControlIdPrefix(prefix string) *Page {
	p.idPrefix = prefix
	return p
}


// Overridable generator for control ids. This is called through the PageI interface, meaning you can change how This
// is done by simply implementing it in a subclass.
func (p *Page) GenerateControlID() string {
	p.idCounter++	// id counter defaults to zero, so pre-increment
	return p.idPrefix + "c" + strconv.Itoa(p.idCounter)
}

func (p *Page) GetControl(id string) ControlI {
	if id == "" || p.controlRegistry == nil {
		return nil
	}
	i := p.controlRegistry.Get(id)
	if c,ok := i.(ControlI); ok {
		return c
	} else {
		return nil
	}
}

// Add the given control to the pathRegistry. Called by the control code whenever a control is created or restored
func (p *Page) addControl(control ControlI) {
	id := control.ID()

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

func (p *Page) changeControlID(oldId string, newId string) {
	if p.GetControl(newId) != nil {
		panic(fmt.Errorf("This control id is already defined on the page: %s", newId))
	}
	ctrl := p.GetControl(oldId)
	p.controlRegistry.Remove(oldId)
	p.controlRegistry.Set(newId, ctrl)
}

func (p *Page) removeControl(id string) {
	// Execute the javascript to remove the control from the dom if we are in ajax mode
	// TODO: Application::ExecuteSelectorFunction('#' . $objControl->getWrapperID(), 'remove');
	// TODO: Make This a direct command in the ajax renderer


	p.controlRegistry.Remove(id)
}


func (p *Page) SetTitle(title string) {
	p.title = title
}

func (p *Page) Path() string {
	return p.path
}

func (p *Page) StateID() string {
	return p.stateId
}

func (p *Page) DrawAjax(ctx context.Context, buf *bytes.Buffer) (err error) {
	err = p.Form().renderAjax(ctx, buf)
	return
}

func (p *Page) GoraddTranslator() Translater {
	return &p.goraddTranslator
}

func (p *Page) ProjectTranslator() Translater {
	return &p.projectTranslator
}


func (p *Page) SetLanguage(l string) {
	p.goraddTranslator.Language = l
	p.projectTranslator.Language = l
}

// MarshalBinary will binary encode the page for the purpose of saving the page in the formstate.
func (p *Page) Encode(e Encoder)(err error) {
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

func (p *Page) Decode(d Decoder)(err error) {
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

func (p *Page) AddHtmlHeaderTag(t html.VoidTag) {
	p.htmlHeaderTags = append(p.htmlHeaderTags, t)
}

func (p *Page) SetResponseHeader(key, value string) {
	if p.responseHeader == nil {
		p.responseHeader = map[string]string{}
	}
	p.responseHeader[key] = value
}

func (p *Page) ClearResponseHeaders() {
	p.responseHeader = nil
}


