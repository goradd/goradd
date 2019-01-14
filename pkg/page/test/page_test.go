package test

import (
	"context"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const TestPath = "/test/PageTest"
const TestFormId = "PageTest"
const TestTxtId = "TxtTest"

type testPageForm struct {
	ΩFormBase
	txt *control.Textbox
}

func newTestForm(ctx context.Context) *testPageForm {
	f := &testPageForm{}
	f.Init(ctx, f, TestPath, TestFormId)
	f.AddRelatedFiles()

	f.txt = control.NewTextbox(f, TestTxtId)
	return f
}


func loadTestFormValues(f *testPageForm) *Page {
	p := f.Page()

	return p
}

// Focuses on exercising page and form values
func TestFormPageValues(t *testing.T)  {
	f := newTestForm(nil)
	loadTestFormValues(f)
	checkTestFormValues(t,f)
}

func checkTestFormValues(t *testing.T, f *testPageForm) {
	p := f.Page()
	require.NotNil(t, p)

	c := p.GetControl(TestTxtId)
	assert.NotNil(t, c)
	assert.IsType(t, &control.Textbox{}, c)

}