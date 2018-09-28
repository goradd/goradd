{{# Macros added to template files for drawing forms. }}

{{# draw draws a goradd control. It expects one parameter, the control itself. You can get the control from the form. }}
{{define draw}}
{{e $1.Draw(ctx, buf) }}
{{end}}

{{define drawWithLabel}}
{{e $1.With(page.NewLabelWrapper()).Draw(ctx, buf) }}
{{end}}

{{define drawWithError}}
{{e $1.With(page.NewErrorWrapper()).Draw(ctx, buf) }}
{{end}}

{{define renderTemplate}}
package {{package }}

import (
    "context"
    "github.com/spekary/goradd/page"
    "bytes"
    {{imports}}
)


func (form *{{form}}) AddHeadTags()  {
    form.FormBase.AddHeadTags()
    if "{{title}}" != "" {
        form.Page().SetTitle("{{title}}")
    }

	{{headTags}}

	// double up to deal with body attributes if they exist
    form.Page().BodyAttributes = `{{bodyAttributes}}`
}


func (form *{{form}}) DrawTemplate(ctx context.Context, buf *bytes.Buffer) (err error) {
{{
{{template}}
}}
	return
}
{{end}}