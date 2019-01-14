package page_test

import (
	"github.com/goradd/goradd/pkg/page"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWrapperInterface(t *testing.T) {
	e := page.NewErrorWrapper()
	var w page.WrapperI

	w = e
	assert.Equal(t, w.ΩNewI().TypeName(), page.ErrorWrapper)

	l := page.NewLabelWrapper()
	w = l
	assert.Equal(t, w.ΩNewI().TypeName(), page.LabelWrapper)

	d := page.NewDivWrapper()

	w = d
	assert.Equal(t, w.ΩNewI().TypeName(), page.DivWrapper)
}

