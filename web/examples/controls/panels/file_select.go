package panels

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	. "github.com/goradd/goradd/pkg/page/control/button"
	. "github.com/goradd/goradd/pkg/page/control/dialog"
)

const (
	FileUploadAction = iota + 10120
)

type FileSelectPanel struct {
	Panel
}

func NewFileSelectPanel(ctx context.Context, parent page.ControlI) {
	p := new(FileSelectPanel)
	p.Init(p, ctx, parent, "FileSelectPanel")
}

func (p *FileSelectPanel) Init(self any, ctx context.Context, parent page.ControlI, id string) {
	p.Panel.Init(self, parent, "FileSelectPanel")
	p.AddControls(ctx,
		PanelCreator{
			ID: "result",
		},
		FileSelectCreator{
			ID:       "uploadButton",
			Multiple: true,
			OnUpload: action.Do(p.ID(), FileUploadAction),
		},
	)

	// This is really specific to this demo because we are switching back and forth between this and bootstrap dialogs.
	// You do not normally need to do this.
	RestoreNewDialogFunction()
}

func (p *FileSelectPanel) DoAction(ctx context.Context, a action.Params) {
	switch a.ID {
	case FileUploadAction:
		Alert(p, "Alert", "Look out!", true)
	}
}

func init() {
	page.RegisterControl(&FileSelectPanel{})
}
