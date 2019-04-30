package orm

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/web/examples/tutorial"
)

type IntroPanel struct {
	Panel
}

func NewIntroPanel(ctx context.Context, parent page.ControlI) page.ControlI {
	p := &IntroPanel{}
	p.Panel.Init(p, parent, "")
	return p
}


func init() {
	tutorial.RegisterTutorialPage("orm", 0, "Introduction to the ORM", NewIntroPanel)
}

