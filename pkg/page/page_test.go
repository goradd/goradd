package page

import (
	"context"
	"testing"
)

const TestPath = "/test/PageTest"
const TestId = "LoginForm"

type testPageForm struct {
	ΩFormBase
}

func newTestForm(ctx context.Context) FormI {
	f := &testPageForm{}
	f.Init(ctx, f, TestPath, TestId)
	f.AddRelatedFiles()
	return f
}


func loadPageValues(f FormI) *Page {
	p := f.Page()

	return p
}

func TestPageValues(t *testing.T)  {
	f := newTestForm(nil)
	f
}
