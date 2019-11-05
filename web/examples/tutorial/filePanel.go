package tutorial

import (
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
	"path/filepath"
)


type FilePanel struct {
	Panel
	File string
	Base string
}

func NewFilePanel(parent page.ControlI) *FilePanel {
	p := &FilePanel{}
	p.Panel.Init(p, parent, "filePanel")
	return p
}

func (p *FilePanel) SetFile(f string) {
	p.File = f
	p.Base = filepath.Base(f)
	p.Refresh()
}

func init() {
	page.RegisterControl(FilePanel{})
}

// TODO: Serialize