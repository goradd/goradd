package panels

import "github.com/goradd/goradd/web/examples/controls"

func init() {
	controls.RegisterPanel("", "Home", NewDefaultPanel)
	controls.RegisterPanel("textbox", "Textboxes", NewTextboxPanel)
	controls.RegisterPanel("checkbox", "Checkboxes and Radio Buttons", NewCheckboxPanel)
	controls.RegisterPanel("selectlist", "Selection Lists", NewSelectListPanel)
	controls.RegisterPanel("table", "Tables", NewTablePanel)
	controls.RegisterPanel("tablecheckbox", "Tables - Checkbox Column", NewTableCheckboxPanel)
	controls.RegisterPanel("tabledb", "Tables - Database Columns", NewTableDbPanel)
	controls.RegisterPanel("tableproxy", "Tables - Proxy Column", NewTableProxyPanel)
	controls.RegisterPanel("tableselect", "Tables - Select Row", NewTableSelectPanel)
	controls.RegisterPanel("hlist", "Nested Lists", NewHListPanel)
	controls.RegisterPanel("repeater", "Repeaters", NewRepeaterPanel)
	controls.RegisterPanel("dialogs", "Dialogs", NewDialogsPanel)
	controls.RegisterPanel("imagecapture", "Image Capture Widget", NewImageCapturePanel)
}

