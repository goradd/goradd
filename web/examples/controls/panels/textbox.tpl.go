//** This file was code generated by GoT. DO NOT EDIT. ***

package panels

import (
	"context"
	"io"
)

func (ctrl *TextboxPanel) DrawTemplate(ctx context.Context, _w io.Writer) (err error) {

	if _, err = io.WriteString(_w, `<h2>Textboxes</h2>
<p>
Textboxes create html input tags, or textarea tags. Some textbox flavors are for entering certain
kinds of data that would be found in a database, like strings, integers and floats. Others are for
validating certain kinds of text input, like URLs or email addresses.
</p>
<p>
Textboxes may be assigned a validator to validate their input. The simplest kind of validator tests
to see if a value has been entered, and can be added by using <i>SetRequired(true)</i>. Other validators
can be added using functions for specific types of controls, or you can create a custom validator.
To see the results of validation in the samples below, scroll to the bottom of this page and click
one of the submit buttons.
</p>
<h3>Database Related Textboxes</h3>
<h4>Plain Textbox</h4>
<p>
By default, the code generator will generate a *Textbox for standard text items in a database, like a VARCHAR in
a sql database. To make a <i>textarea</i> instead of an <i>input</i> tag, set the RowCount to a value that is not zero.
`); err != nil {
		return
	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("plainText-ff").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("plainText-ff").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("multiText-ff").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("multiText-ff").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if _, err = io.WriteString(_w, `</p>
<h4>IntegerTextbox</h4>
<p>
An *IntegerTextbox corresponds to an integer item in a database, like an INT in
a sql database. Integer textboxes are validated to make sure they contain an integer.
`); err != nil {
		return
	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("intText-ff").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("intText-ff").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if _, err = io.WriteString(_w, `</p>
<h4>DateTextbox</h4>
<p>
A *DateTextbox corresponds to a Date, Time or DateTime in a database.
Timestamps generally are not editable, so they usually generate a DateTimeSpan (as in html span).
These textboxes validate to make sure they match a particular format.
`); err != nil {
		return
	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("dateTimeText-ff").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("dateTimeText-ff").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("dateText-ff").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("dateText-ff").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("timeText-ff").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("timeText-ff").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if _, err = io.WriteString(_w, `</p>

<h4>FloatTextbox</h4>
<p>
A *FloatTextbox corresponds to a floating point number item in a database, like a FLOAT in
a sql database. Float textboxes are validated to make sure they contain a numeric value.
Click on one of the Submit buttons below to cause the controls to validate.
`); err != nil {
		return
	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("floatText-ff").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("floatText-ff").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if _, err = io.WriteString(_w, `</p>
<h3>Validating Textboxes</h3>
<h4>Email Textbox</h4>
<p>
The EmailTextbox accepts email addresses only. It is capable of accepting multiple email addresses separated
by commas. If it is set up to only accept one email address, it will also set its "type" attribute to "email"
so that the browser can potential help with entering and validating an email address. This is particularly
useful for mobile browsers, as they sometimes change the virtual keyboard to make it easier to enter an `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, `@ `); err != nil {
		return
	}

	if _, err = io.WriteString(_w, ` symbol
or provide a shortcut key to enter ".com".
`); err != nil {
		return
	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("emailText-ff").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("emailText-ff").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if _, err = io.WriteString(_w, `</p>
<h3>Textbox Types</h3>
<p>HTML offers a number of different types to give browsers a hint of what kind of data the server
is expecting in a particular textbox. Not all are supported on all browsers, but below are some examples
of ones that are commonly supported.
`); err != nil {
		return
	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("passwordText-ff").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("passwordText-ff").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("searchText-ff").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("searchText-ff").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if _, err = io.WriteString(_w, `</p>


`); err != nil {
		return
	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("ajaxButton").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("ajaxButton").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	if `` == "" {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("serverButton").Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	} else {

		if _, err = io.WriteString(_w, `    `); err != nil {
			return
		}
		ctrl.Page().GetControl("serverButton").ProcessAttributeString(``).Draw(ctx, _w)
		if _, err = io.WriteString(_w, `
`); err != nil {
			return
		}

	}

	return
}
