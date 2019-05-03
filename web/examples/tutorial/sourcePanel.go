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
	buttonPanel *Panel
	filePanel *FilePanel
}

func NewSourcePanel(parent page.ControlI) *SourcePanel {
	p := &SourcePanel{}
	p.Panel.Init(p, parent, "sourcePanel")
	p.buttonPanel = NewPanel(p, "buttonPanel")
	p.SetVisible(false)

	p.filePanel = NewFilePanel(p) // we will be doing our own escaping

	return p
}

// show shows the panel and loads the button bar with buttons
func (p *SourcePanel) show(files []string) {
	p.buttonPanel.RemoveChildren()

	for i,path := range files {
		base := filepath.Base(path)
		b := NewButton(p.buttonPanel, "")
		b.SetLabel(fmt.Sprintf("%d. %s", i, base))
		b.SetActionValue(path)
		b.On(event.Click(), action.Ajax(p.ID(), FileAction))
	}

	b := NewButton(p.buttonPanel, "closeButton")
	b.SetLabel("Close")
	b.On(event.Click(), action.Ajax(p.ID(), CloseAction))

	p.SetVisible(true)
}

func (p *SourcePanel) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case CloseAction:
		p.SetVisible(false)
	case FileAction:
		file := a.ControlValueString()
		p.filePanel.SetFile(file)
	}
}


func init() {
}
