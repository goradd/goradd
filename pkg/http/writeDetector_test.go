package http

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/http"
	"io"
	"testing"
)

type Rewinder interface {
	Rewind(n int)
}

func TestWriteDetector(t *testing.T) {
	h := http.TestResponseWriter{}
	buf := bytes.Buffer{}
	bwriter := &bufferedResponseWriter{&h, &buf, 0, false, 0}
	w := WriteDetector{bwriter, false}
	w2 := io.Writer(&w)
	_,ok := w2.(Rewinder)
	assert.True(t, ok)
}

