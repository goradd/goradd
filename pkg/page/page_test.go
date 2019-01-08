package page

import (
	"context"
	"github.com/goradd/goradd/test/browser"
	"testing"
)

// Page unit testing, just testing the public interface of the page part of a form at this point.

const TestPath = "/test/PageTest"
const TestId = "LoginForm"

type testPageForm struct {
	FormBase
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
	f := newTestForm(ctx)
}
