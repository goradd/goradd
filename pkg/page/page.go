// Package page is the user-interface layer of goradd, and implements state management and rendering
// of an html page, as well as the framework for rendering controls.
//
// To use the page package, you start by creating a form object, and then adding controls to that form.
// You also should add a drawing template to define additional html for the form.
package page

import (
	"bytes"
	"context"
	"fmt"
	"github.com/goradd/gengen/pkg/maps"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/i18n"
	"github.com/goradd/goradd/pkg/messageServer"
	"strconv"
	"strings"
)

// PageRenderStatus keeps track of whether we are rendering the page or not
type PageRenderStatus int

// Future note. Below is for general information but should NOT be used to synchronize multiple drawing routines.
// An architecture using channels to synchronize page changes and drawing would be better.
// For now, except for testing, we should not get in a situation where multiple copies of a form
// are being used.

const (
	PageIsNotRendering PageRenderStatus = iota // FormBase has started rendering but has not finished
	PageIsRendering
)

// PageDrawFunc is the type of the page drawing function. This is implemented by the page drawing template.
type PageDrawFunc func(context.Context, *Page, *bytes.Buffer) error

const EncodingVersion = 1

// DrawI is the interface for items that draw into the draw buffer
type DrawI interface {
	Draw(context.Context, *bytes.Buffer) error
}

// The Page object is the top level drawing object, and is essentially a wrapper for the form. The Page draws the
// html, head and body tags, and includes the one Form object on the page. The page also maintains a record of all
// the controls included on the form.
type Page struct {
	// BodyAttributes contains the attributes that will be output with the body tag. It should be set before the
	// form draws, like in the AddHeadTags function.
	BodyAttributes  string

	stateId      string // Id in cache of the pagestate. Needs to be output by form.
	renderStatus PageRenderStatus
	idPrefix     string // For creating unique ids for the app

	controlRegistry *maps.SliceMap
	form            FormI
	idCounter       int
	title           string // page title to draw in head tag
	htmlHeaderTags  []html.VoidTag
	responseHeader  map[string]string // queues up anything to be sent in the response header
	responseError   int

	language 	    int		// Don't serialize this. This is a cached version of what the session holds.
}

// Init initializes the page. Should be called by a form just after creating Page.
func (p *Page) Init() {
}

// Restore is called immediately after the page has been unserialized, to restore data that did not get serialized.
func (p *Page) Restore() {
	p.form.Restore(p.form)
}

func (p *Page) runPage(ctx context.Context, buf *bytes.Buffer, isNew bool) (err error) {
	grCtx := GetContext(ctx)

	if grCtx.err != nil {
		panic(grCtx.err)	// An error occurred during unpacking of the context, so report that now
	}

	if err = p.Form().Run(ctx); err != nil {
		return err
	}

	// TODO: Lifecycle calls - push them to the form

	// cache the language tags so we only need to look them up once for every call
	p.language = i18n.SetDefaultLanguage(ctx, grCtx.Header.Get("accept-language"))

	if isNew {
		p.Form().AddHeadTags()
		p.Form().LoadControls(ctx)

	} else {
		p.Form().control().updateValues(grCtx) // Tell all the controls to update their values.
		// if this is an event response, do the actions associated with the event
		if c := p.GetControl(grCtx.actionControlID); c != nil {
			c.control().doAction(ctx)
		}
	}

	p.ClearResponseHeaders()
	if grCtx.RequestMode() == Ajax {
		err = p.DrawAjax(ctx, buf)
		p.SetResponseHeader("Content-Type", "application/json")
	} else if grCtx.RequestMode() == Server || grCtx.RequestMode() == Http {
		//p.url = grCtx.HttpContext.URL. We might want a record of the original URL to be used during ajax calls someday. Until we have a reason, this will remain commented out.
		err = p.Draw(ctx, buf)
	} else {
		// TODO: Implement a hook for the CustomAjax call and/or Rest API calls?
	}

	p.Form().control().writeState(ctx)
	p.Form().Exit(ctx, err)
	return
}

// Returns the form for the page
func (p *Page) Form() FormI {
	return p.form
}

// Draw draws the page.
func (p *Page) Draw(ctx context.Context, buf *bytes.Buffer) (err error) {
	f := p.form.PageDrawingFunction()
	return f(ctx, p, buf)
}

// DrawHeaderTags draws all the inner html for the head tag
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

// SetControlIdPrefix sets the prefix for control ids. Some javascript frameworks (i.e. jQueryMobile) require that control ids
// be unique across the application, vs just in the page, because they create internal caches of control ids. This
// allows you to set a per page prefix that will be added to all control ids to make them unique across the whole
// application. However, its up to you to make sure the names are unique per page.
func (p *Page) SetControlIdPrefix(prefix string) *Page {
	p.idPrefix = prefix
	return p
}

// GenerateControlID generates unique control ids. If you want to do your own id generation, or modifying of given ids, implement that
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
			panic (fmt.Sprintf(`A control with id "%s" is being added a second time to the page. Ids must be unique on the page.`, id))
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

// GetControl returns the control with the given id. If not found, it returns nil.
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

// addControl adds the given control to the controlRegistry. It is called by the control code whenever a control is created.
func (p *Page) addControl(control ControlI) {
	id := control.ID()

	if id == "" {
		panic("Control must have an id before being added.")
	}

	if p.controlRegistry == nil {
		p.controlRegistry = maps.NewSliceMap()
	}

	if p.controlRegistry.Has(id) {
		panic("Control id already exists. Control must have a unique id on the page before being added.")
	}

	p.controlRegistry.Set(id, control)

	if control.Parent() == nil {
		if f, ok := control.(FormI); ok {
			if p.form != nil {
				panic("The Form object for the page has already been set.")
			} else {
				p.form = f
			}
		} else {
			panic("Controls must have a parent.")
		}
	}
}

/* Remove?
func (p *Page) changeControlID(oldId string, newId string) {
	if p.GetControl(newId) != nil {
		panic(fmt.Errorf("this control id is already defined on the page: %s", newId))
	}
	ctrl := p.GetControl(oldId)
	p.controlRegistry.Delete(oldId)
	p.controlRegistry.Set(newId, ctrl)
}
*/

func (p *Page) removeControl(id string) {
	// Execute the javascript to remove the control from the dom if we are in ajax mode
	// TODO: Application::ExecuteSelectorFunction('#' . $objControl->getWrapperID(), 'remove');
	// TODO: Make This a direct command in the ajax renderer

	p.controlRegistry.Delete(id)
}

// Title returns the content of the <title> tag that will be output in the head of the page.
func (p *Page) Title() string {
	return p.title
}

// Call SetTitle to set the content of the <title> tag to be output in the head of the page.
func (p *Page) SetTitle(title string) {
	p.title = title
}

func (p *Page) setStateID(stateId string) {
	p.stateId = stateId
}

// StateID returns the page state id. This is output by the form so that we can recover the saved state of the page
// each time we call into the application.
func (p *Page) StateID() string {
	return p.stateId
}

// DrawAjax renders the page during an ajax call. Since the page itself is already rendered, it simply hands off this
// responsibility to the form.
func (p *Page) DrawAjax(ctx context.Context, buf *bytes.Buffer) (err error) {
	err = p.Form().renderAjax(ctx, buf)
	return
}

// GobEncode here is implemented to intercept the GobSerializer to only encode an empty structure. We use this as part
// of our overall serialization stratgey for forms. Controls still need to be registered with gob.
func (p *Page) GobEncode() (data []byte, err error) {
	return
}

func (p *Page) GobDecode(data []byte) (err error) {
	return
}

func (p *Page) MarshalJSON() (data []byte, err error) {
	return
}

func (p *Page) UnmarshalJSON(data []byte) (err error) {
	return
}

type pageEncoded struct {
	StateId      string // Id in cache of the pagestate. Needs to be output by form.
	Path         string // The path to the page. FormBase needs to know this so it can make the action tag
	IdPrefix     string // For creating unique ids for the app
	IdCounter       int
	Title           string // page title to draw in head tag
	HtmlHeaderTags  []html.VoidTag
	BodyAttributes  string

	FormID string // to record the form

}

// Encode is called by the framework to serialize the page state.
// TODO: serialization is not completely implemented yet
func (p *Page) Encode(e Encoder) (err error) {
	s := pageEncoded{
		StateId:           p.stateId,
		IdPrefix:          p.idPrefix,
		Title:             p.title,
		HtmlHeaderTags:    p.htmlHeaderTags,
		BodyAttributes:    p.BodyAttributes,
		FormID:			   p.form.ID(),
	}

	if err = e.Encode(s); err != nil {
		return
	}

	if err = e.EncodeControl(p.form); err != nil {
		return
	}

	// Add the items from the control registry that were not serialized as part of serializing the form.
	// This might happen if the item had no parent, like dialogs or other objects that are automatically drawn.
	var count int
	p.controlRegistry.Range(func(key string, value interface{}) bool {
		if !value.(ControlI).control().encoded {
			count++
		}
		return true
	})
	if err = e.Encode(count); err != nil {
		return
	}
	p.controlRegistry.Range(func(key string, value interface{}) bool {
		c := value.(ControlI)
		if !c.control().encoded {
			if err = e.EncodeControl(c); err != nil {
				return false
			}
		}
		return true
	})

	return
}

// Decode is called by the framework to serialize the page state.
func (p *Page) Decode(d Decoder) (err error) {
	s := pageEncoded{}
	if err = d.Decode(&s); err != nil {
		return
	}
	p.controlRegistry = maps.NewSliceMap()
	p.stateId = s.StateId
	p.idPrefix = s.IdPrefix
	p.title = s.Title
	p.htmlHeaderTags = s.HtmlHeaderTags
	p.BodyAttributes = s.BodyAttributes

	var ci ControlI
	if ci,err = d.DecodeControl(p); err != nil {
		return
	}
	p.form = ci.(FormI)

	// Deserialize the controls that were not part of the form structure, like dialogs
	var count int
	if err = d.Decode(&count); err != nil {
		return
	}

	for i:=0; i<count;i++ {
		if ci,err = d.DecodeControl(p); err != nil { // the process of decoding will automatically add to the control registry, so no need to do anything with the result.
			return
		}
	}

	return err
}

// AddHtmlHeaderTag adds the given tag to the head section of the page.
func (p *Page) AddHtmlHeaderTag(t html.VoidTag) {
	p.htmlHeaderTags = append(p.htmlHeaderTags, t)
}

// SetResponseHeader sets a value in the html response header. You generally would only need to do this if your are outputting
// custom content, like a pdf file.
func (p *Page) SetResponseHeader(key, value string) {
	if p.responseHeader == nil {
		p.responseHeader = map[string]string{}
	}
	p.responseHeader[key] = value
}

// ClearResponseHeaders removes all the current response headers.
func (p *Page) ClearResponseHeaders() {
	p.responseHeader = nil
}

// PushRedraw will cause the form to refresh in between events. This will cause the client to pull
// the ajax response. Its possible that this will happen while drawing. We avoid the race condition
// by sending the message anyways, and allowing the client to send an event back to us, essentially
// using the javascript event mechanism to synchronize us. We might get an unnecessary redraw, but
// that is not a big deal.
func (p *Page) PushRedraw() {
	channel := "form-" + p.stateId
	if messageServer.HasChannel(channel) {	// If we call this while launching a page, the channel isn't created yet, but the page is going to be drawn, so its ok.
		messageServer.SendMessage(channel, map[string]interface{}{"grup":true})
	}
}

// LanguageCode returns the language code that will be put in the lang attribute of the html tag.
// It is taken from the i18n package.
func (p *Page) LanguageCode() string {
	return i18n.CanonicalValue(p.language)
}