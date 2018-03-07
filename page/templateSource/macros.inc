{{# Macros added to template files for drawing forms. }}

{{# draw draws a goradd control. It expects one parameter, the control itself. You can get the control from the form. }}
{{define draw}}
{{e $1.Draw(ctx, buf) }}
{{end}}
