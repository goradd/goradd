package http

import (
	"io"
	"net/http"
)


type FileProcessorFunc func(r io.Reader, w http.ResponseWriter, req *http.Request) error

// fileProcessors is a map that connects file endings to processors that will process the content and return it
// to the output stream, bypassing other means of processing static files.
var fileProcessors map[string]FileProcessorFunc

// RegisterFileProcessor registers a processor function for static files that have a particular extension.
// Do this at init time.
func RegisterFileProcessor(extension string, processorFunc FileProcessorFunc) {
	if fileProcessors == nil {
		fileProcessors = make(map[string]FileProcessorFunc)
	}
	fileProcessors[extension] = processorFunc
}
