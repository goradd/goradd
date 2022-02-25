package http

import (
	"io/fs"
	"os"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func makeFs() fs.FS {
	fsys := fstest.MapFS{
		"test1":    {Data: []byte("test"), Mode: os.ModePerm},
		"test2.js": {Data: []byte("test2"), Mode: os.ModePerm},
	}
	return fsys
}

func TestCacheBuster(t *testing.T) {
	RegisterAssetDirectory("/assets", makeFs())
	url := GetAssetUrl("/assets/test2.js")
	s := StripCacheBusterPath(url)
	assert.Equal(t, s, "/assets/test2.js")
}
