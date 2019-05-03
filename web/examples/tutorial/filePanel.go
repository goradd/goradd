package tutorial

import (
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
	"path/filepath"
)


type FilePanel struct {
	Panel
	file string
	base string
	content string
}

func NewFilePanel(parent page.ControlI) *FilePanel {
	p := &FilePanel{}
	p.Panel.Init(p, parent, "filePanel")
	return p
}

func (p *FilePanel) SetFile(f string) {
	p.file = f
	p.base = filepath.Base(f)
	p.Refresh()
}
