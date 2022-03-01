package control

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
)

// The FormBase is the control that all Form objects should include,
// and is the master container for all other goradd controls.
// This is here for future expansion.
type FormBase struct {
	page.FormBase
}


// Init initializes the form.
func (f *FormBase) Init(ctx context.Context, id string) {
	f.FormBase.Init(ctx, id)
}


