package control_base

import (
	"github.com/spekary/goradd/pkg/page"
	"context"
	"github.com/spekary/goradd/pkg/html"
)

// The local FormBase override. All framework forms descend from this one. You can change how all the forms in your
// application work by making modifications here. This struct is overridden by the one in the control package, and
// so you should descend your forms from that one.
type FormBase struct {
	page.FormBase
}


func (f *FormBase) Init(ctx context.Context, self page.FormI, path string, id string) {
	f.FormBase.Init(ctx, self, path, id)

	// additional initializations. For example, your custom page template.
	//f.Page().SetDrawFunction()
}

// You can put overrides that should apply to all your forms here.
func (f *FormBase) AddRelatedFiles() {
	f.FormBase.AddRelatedFiles() // add default files
}

func (f *FormBase) AddHeadTags() {
	f.FormBase.AddHeadTags()
	f.Page().AddHtmlHeaderTag(html.VoidTag{"link",
	html.NewAttributes().
		Set("rel", "icon").
		Set("type", "image/x-icon").
		Set("href", "/favicon.ico")})
}