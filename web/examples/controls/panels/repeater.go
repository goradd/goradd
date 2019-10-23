package panels

import (
	"bytes"
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/page"
	. "github.com/goradd/goradd/pkg/page/control"
)

type RepeaterPanel struct {
	Panel
}


func NewRepeaterPanel(ctx context.Context, parent page.ControlI) {
	p := &RepeaterPanel{}
	p.Panel.Init(p, parent, "repeaterPanel")
	p.AddControls(ctx,
		RepeaterCreator{
			ID: "repeater1",
			ItemHtmler: p,
			DataProvider: p,
			PageSize:5,
		},
		// A DataPager can be a standalone control, which you draw manually
		DataPagerCreator{
			ID:           "pager1",
			PagedControl: "repeater1",
		},
	)
}

// BindData satisfies the data provider interface so that the parent panel of the table
// is the one that is providing the table.
func (f *RepeaterPanel) BindData(ctx context.Context, s DataManagerI) {
	switch s.ID() {
	case "repeater1":
		t := s.(PagedControlI)
		t.SetTotalItems(uint(len(tableMapData)))
		start, end := t.SliceOffsets()
		s.SetData(tableMapData[start:end])
	}
}

func (f *RepeaterPanel) RepeaterHtml(ctx context.Context, r RepeaterI, i int, data interface{}, buf *bytes.Buffer) error {
	d := data.(TableMapData)
	buf.WriteString(fmt.Sprintf(`<div>ID: %s, Name: %s</div>`, d["id"], d["name"]))
	return nil
}

func init() {
	page.RegisterControl(RepeaterPanel{})
}

