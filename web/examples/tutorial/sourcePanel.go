package tutorial

import (
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/event"
	"path/filepath"
)

const (
	FileAction = iota + 1
	CloseAction
)

type SourcePanel struct {
	Panel
	ButtonPanel *Panel
	FilePanel *FilePanel
}

func NewSourcePanel(parent page.ControlI, id string) *SourcePanel {
	p := &SourcePanel{}
	p.Self = p
	p.Init(parent, id)
	return p
}

func (p *SourcePanel) Init(parent page.ControlI, id string) {
	p.Panel.Init(parent, id)
	p.ButtonPanel = NewPanel(p, "buttonPanel")
	p.FilePanel = NewFilePanel(p) // we will be doing our own escaping
}

// show shows the panel and loads the button bar with buttons
func (p *SourcePanel) show(files []string) {
	p.ButtonPanel.RemoveChildren()

	for i,path := range files {
		base := filepath.Base(path)
		b := NewButton(p.ButtonPanel, "")
		b.SetLabel(fmt.Sprintf("%d. %s", i, base))
		b.SetActionValue(path)
		b.On(event.Click(), action.Ajax(p.ID(), FileAction))
	}
}

func (p *SourcePanel) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case FileAction:
		file := a.ControlValueString()
		p.FilePanel.SetFile(file)
	}
}


func init() {
	page.RegisterControl(&SourcePanel{})
}

func GetSourcePanel(p page.ControlI) *SourcePanel {
	return p.Page().GetControl("sourcePanel").(*SourcePanel)
}