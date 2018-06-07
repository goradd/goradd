package control_base

import (
	"bytes"
	"context"
	"fmt"
	"github.com/spekary/goradd/page"
	localPage "goradd/page"
	"text/template"
)

// The interface is used to add template drawing functions that are created by the got template engine. Its important
// to embed the template drawing function in an interface so that we can properly serialize and unserialize it.
type TemplateDrawer interface {
	DrawTemplate(ctx context.Context, c page.ControlI, buf *bytes.Buffer) (err error)
}

// Panel is the base type for all panel-type controls. These controls are surrounded by a span or div, and can
// have a template draw into the inner-html of the template.
type Panel struct {
	localPage.Control
	goTemplate     string             // a locally defined go template, to be parsed
	parsedTemplate *template.Template // a parsed template to use for drawing
	goTemplateName string             // a named template from the global template pool
	gotTemplate    TemplateDrawer     // a TemplateDrawer implementation
}

func (c *Panel) SetGoTemplate(t string) {
	c.goTemplate = t
	c.parsedTemplate = nil
	c.Refresh()
}

func (c *Panel) SetNamedGoTemplate(name string) {
	c.goTemplateName = name
	c.Refresh()
}

func (c *Panel) SetTemplateDrawer(t TemplateDrawer) {
	c.gotTemplate = t
}

func (c *Panel) DrawTemplate(ctx context.Context, buf *bytes.Buffer) (err error) {
	if c.gotTemplate != nil {
		c.gotTemplate.DrawTemplate(ctx, c.Self.(page.ControlI), buf)
		return
	} else if c.goTemplateName != "" {
		if t := page.GetTemplate(c.goTemplateName); t != nil {
			return t.Execute(buf, c)
		} else {
			return fmt.Errorf("Could not find template %s", c.goTemplateName)
		}
	} else if c.goTemplate != "" {
		if c.parsedTemplate == nil {
			c.parsedTemplate, err = template.New("").Parse(c.goTemplate)
			if err != nil {
				return
			}
		}
		return c.parsedTemplate.Execute(buf, c)
	}
	return page.NewAppErr(page.AppErrNoTemplate)
}
