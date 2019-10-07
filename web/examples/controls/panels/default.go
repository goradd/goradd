package panels

import (
	"context"
	"encoding/gob"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
)

// shared
const controlsFormPath = "/goradd/examples/controls.g"

const (
	AjaxSubmit int = iota + 1
	ServerSubmit
	ButtonSubmit
	ResetStateSubmit
	ProxyClick
)


type DefaultPanel struct {
	Panel
}

func NewDefaultPanel(ctx context.Context, parent page.ControlI) {
	p := &DefaultPanel{}
	p.Panel.Init(p, parent, "defaultPanel")
}

func init() {
	gob.Register(DefaultPanel{})
}