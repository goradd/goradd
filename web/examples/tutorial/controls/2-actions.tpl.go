//** This file was code generated by GoT. DO NOT EDIT. ***

package controls

import (
	"context"
	"io"
)

// DrawTemplate draws the content of the matching control's template file.
func (ctrl *ActionsPanel) DrawTemplate(ctx context.Context, _w io.Writer) (err error) {

	if _, err = io.WriteString(_w, `<h1>Actions</h1>
<h2>Intro</h2>
<p>
As mentioned previously, <i>Events</i> are attached to controls and trigger <i>Actions</i>.
</p>
<p>
There are two types of actions:
<ul>
<li>Javascript Actions
<li>Callback Actions
</ul>
</p>
<h2>Javascript Actions</h2>
<p>
Javascript Actions are snippets of Javascript code that execute in the client browser in response to an event.
</p>
<p>
For example, the following code will convert any text to uppercase while it is typed.
</p>
<label>To Uppercase</label>`); err != nil {
		return
	}

	if _, err = io.WriteString(_w, `
`); err != nil {
		return
	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("textbox1").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("textbox1").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if _, err = io.WriteString(_w, `<code>
textbox1 := NewTextbox(p, "textbox1")
textbox1.On(event.Input().Action(action.Javascript("event.target.value = event.target.value.toUpperCase()")))
</code>
<p>
GoRADD predefines some Javascript actions to do common tasks. You can find the complete list of predefined Javascript
actions in the <a href="https://pkg.go.dev/github.com/goradd/goradd/pkg/page/action">Actions Documentation</a> under the ActionI heading.
</p>
<h2>Callback Actions</h2>
<p>
Callback actions invoke the control.DoAction() that is in every GoRADD control. By default, if you do not specify an
action to an event, the event will invoke the DoAction function on the receiving control using an Ajax call from the client browser.
</p>
<p>
The two buttons below use actions to get both the server's time and browser's time.
Click on the "2-actions.go" button under the View Source list to see how these buttons are
created and how the DoAction() function responds.
Note that the DoAction() function just sets the text of the span and the span is automatically redrawn
in the browser with the new value.
</p>
<p>
`); err != nil {
		return
	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("serverTimeButton").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("serverTimeButton").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if _, err = io.WriteString(_w, ` `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, `
`); err != nil {
		return
	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("clientTimeButton").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("clientTimeButton").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if _, err = io.WriteString(_w, ` `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, `
`); err != nil {
		return
	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("timeSpan").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("timeSpan").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if _, err = io.WriteString(_w, `</p>

`); err != nil {
		return
	}

	return
}
