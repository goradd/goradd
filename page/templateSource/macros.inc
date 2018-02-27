{{# Macros added to template files for drawing form. }}

{{# draw draws a goradd control. It expects one paramter, the control itself. You can get the control from the form. }}
{{define draw}}
if err = $1.Draw(ctx, buf); err != nil {return}
{{end}}
