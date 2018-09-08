package page

import (
	"bytes"
	"context"
	"fmt"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/log"
	"github.com/spekary/goradd/util/types"
	"strconv"
	"strings"
)

type PageRenderStatus int

type PageDrawFunc func(context.Context, *Page, *bytes.Buffer) error

const (
	UNRENDERED PageRenderStatus = iota // FormBase has not started isRendering
	BEGUN                              // FormBase has started isRendering but has not finished
	ENDED                              // FormBase isRendering has already been started and finished
)

// Anything that draws into the draw buffer must implement This interface
type DrawI interface {
	Draw(context.Context, *bytes.Buffer) error
}

type Page struct {
	stateId      string // Id in cache of the override. Needs to be output by form.
	path         string // The path to the override. FormBase needs to know this so it can make the action tag
	renderStatus PageRenderStatus
	idPrefix     string // For creating unique ids for the app

	controlRegistry *types.OrderedMap
	form            FormI
	idCounter       int
	drawFunc        PageDrawFunc
	title           string // override title to draw in head tag
	htmlHeaderTags  []html.VoidTag
	responseHeader  map[string]string // queues up anything to be sent in the response header
	responseError   int

	goraddTranslator  PageTranslator
	projectTranslator PageTranslator
}

// Initialize the override base. Should be called by a override just after creating PageBase.
func (p *Page) Init(ctx context.Context, path string) {
	p.path = path
	p.drawFunc = p.DrawFunction()
	p.goraddTranslator = PageTranslator{Domain: GoraddDomain}
	p.projectTranslator = PageTranslator{Domain: ProjectDomain}
}

// Restore is called immediately after the override has been unserialized, to restore data that did not get serialized.
func (p *Page) Restore() {
	p.drawFunc = p.DrawFunction()
	p.form.Restore(p.form)
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

	if grCtx.err != nil {
		panic(grCtx.err)	// An error occurred during unpacking of the context, so report that now
	}

	if err = p.Form().Run(ctx); err != nil {
		return err
	}

	p.renderStatus = UNRENDERED

	log.FrameworkDebugf("Run: %s", grCtx)

	// TODO: Lifecycle calls - push them to the form

	if isNew {
		p.Form().AddHeadTags()
		p.Form().LoadControls(ctx)
	} else {
		p.Form().control().updateValues(grCtx) // Tell all the controls to update their values.
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
		//p.SetResponseHeader("Content-Type", "text/html")	// default for web override. Response can change This if drawing something else.
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

// Draws from the override template. The default should be fine for most situations.
// You can replace the template function with your own
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
		for _, tag := range p.htmlHeaderTags {
			buf.WriteString(tag.Render())
		}
	}

	p.Form().DrawHeaderTags(ctx, buf)
}

// Sets the prefix for control ids. Some javascript frameworks (i.e. jQueryMobile) require that control ids
// be unique across the application, vs just in the override, because they create internal caches of control ids. This
// allows you to set a per override prefix that will be added on to all control ids to make them unique across the whole
// application. However, its up to you to make sure the names are unique per override.
func (p *Page) SetControlIdPrefix(prefix string) *Page {
	p.idPrefix = prefix
	return p
}

// Generates unique control ids. If you want to do your own id generation, or modifying of given ids, implement that
// in an override to the control.Init function. The given id is one that the user supplies. User provided ids and
// generated ids can be further munged by providing an id prefix through SetControlIdPrefix().
func (p *Page) GenerateControlID(id string) string {
	if id != "" {
		if strings.Contains(id, "_") {
			// underscores are used by the action system to route actions to sub items of the control.
			panic ("You cannot add a control with an underscore in the name. Use a hyphen instead.")
		}
		if p.idPrefix != "" {
			if !strings.HasPrefix(id, p.idPrefix) {	// subcontrols might already have this prefix
				id = p.idPrefix + id
			}
		}
		if p.GetControl(id) != nil {
			panic (fmt.Sprintf(`A control with id "%s" is being added a second time to the override. Ids must be unique on the override.`, id))
		} else {
			return id
		}
	} else {
		var trialid string
		for trialid == "" || p.GetControl(trialid) != nil { // checks to make sure user did not previously add a control that might match our generation pattern
			p.idCounter++
			trialid = p.idPrefix + "c" + strconv.Itoa(p.idCounter)
		}
		return trialid
	}
}

func (p *Page) GetControl(id string) ControlI {
	if id == "" || p.controlRegistry == nil {
		return nil
	}
	i := p.controlRegistry.Get(id)
	if c, ok := i.(ControlI); ok {
		return c
	} else {
		return nil
	}
}

// Add the given control to the controlRegistry. Called by the control code whenever a control is created or restored
func (p *Page) addControl(control ControlI) {
	id := control.ID()

	if id == "" {
		panic("Control must have an id before being added.")
	}

	if p.controlRegistry == nil {
		p.controlRegistry = types.NewOrderedMap()
	}

	if p.controlRegistry.Has(id) {
		panic("Control id already exists. Control must have a unique id on the override before being added.")
	}

	p.controlRegistry.Set(id, control)

	if control.Parent() == nil {
		if f, ok := control.(FormI); ok {
			if p.form != nil {
				panic("The Form object for the override has already been set.")
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
		panic(fmt.Errorf("This control id is already defined on the override: %s", newId))
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

func (p *Page) Title() string {
	return p.title
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

// TODO: Move these to the session object, since language is likely the same on a session basis
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

// MarshalBinary will binary encode the override for the purpose of saving the override in the formstate.
func (p *Page) Encode(e Encoder) (err error) {
	if err = e.Encode(p.stateId); err != nil {
		return
	}
	if err = e.Encode(p.path); err != nil {
		return
	}
	if err = e.Encode(p.idPrefix); err != nil {
		return
	}
	if err = e.Encode(p.form); err != nil {
		return
	}
	if err = e.Encode(p.idCounter); err != nil {
		return
	}
	if err = e.Encode(p.title); err != nil {
		return
	}
	if err = e.Encode(p.htmlHeaderTags); err != nil {
		return
	}
	if err = e.Encode(p.goraddTranslator); err != nil {
		return
	}
	if err = e.Encode(p.projectTranslator); err != nil {
		return
	}
	return
}

func (p *Page) Decode(d Decoder) (err error) {
	if err = d.Decode(&p.stateId); err != nil {
		return
	}
	if err = d.Decode(&p.path); err != nil {
		return
	}
	if err = d.Decode(&p.idPrefix); err != nil {
		return
	}
	if err = d.Decode(&p.form); err != nil {
		return
	}
	if err = d.Decode(&p.idCounter); err != nil {
		return
	}
	if err = d.Decode(&p.title); err != nil {
		return
	}
	if err = d.Decode(&p.htmlHeaderTags); err != nil {
		return
	}
	if err = d.Decode(&p.goraddTranslator); err != nil {
		return
	}
	if err = d.Decode(&p.projectTranslator); err != nil {
		return
	}
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

