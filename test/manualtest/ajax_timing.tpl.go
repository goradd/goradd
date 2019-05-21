//** This file was code generated by got. DO NOT EDIT. ***

package manualtest

import (
	"bytes"
	"context"

	"github.com/goradd/goradd/pkg/page"
)

func (form *AjaxTimingForm) AddHeadTags() {
	form.FormBase.AddHeadTags()
	if "Ajax Timing" != "" {
		form.Page().SetTitle("Ajax Timing")
	}

	// double up to deal with body attributes if they exist
	form.Page().BodyAttributes = `
buf.WriteString(fmt.Sprintf("%v", bodyAttributes))
`
}

func (form *AjaxTimingForm) DrawTemplate(ctx context.Context, buf *bytes.Buffer) (err error) {

	buf.WriteString(`
`)

	buf.WriteString(`
	<div>
        `)

	buf.WriteString(`
`)

	{
		err := form.Txt1.With(page.NewLabelWrapper()).Draw(ctx, buf)
		if err != nil {
			return err
		}
	}

	buf.WriteString(`
        `)

	buf.WriteString(`
`)

	{
		err := form.Txt1ChangeLabel.With(page.NewLabelWrapper()).Draw(ctx, buf)
		if err != nil {
			return err
		}
	}

	buf.WriteString(`
        `)

	buf.WriteString(`
`)

	{
		err := form.Txt1KeyUpLabel.With(page.NewLabelWrapper()).Draw(ctx, buf)
		if err != nil {
			return err
		}
	}

	buf.WriteString(`
	</div>
	<div>
	    `)

	buf.WriteString(`
`)

	{
		err := form.Chk.With(page.NewLabelWrapper()).Draw(ctx, buf)
		if err != nil {
			return err
		}
	}

	buf.WriteString(`
	    `)

	buf.WriteString(`
`)

	{
		err := form.ChkLabel.With(page.NewLabelWrapper()).Draw(ctx, buf)
		if err != nil {
			return err
		}
	}

	buf.WriteString(`
	</div>
    <div>
        `)

	buf.WriteString(`
`)

	{
		err := form.Txt2.With(page.NewLabelWrapper()).Draw(ctx, buf)
		if err != nil {
			return err
		}
	}

	buf.WriteString(`
        `)

	buf.WriteString(`
`)

	{
		err := form.Btn.Draw(ctx, buf)
		if err != nil {
			return err
		}
	}

	buf.WriteString(`
    </div>

`)

	buf.WriteString(`
`)

	return
}