package http

// contentTypes is a map that connects file endings to content types for files served
// out of file systems, like embedded files and statically served files.
var contentTypes map[string]string

// RegisterContentType registers a content type that associates a file extension
// with a specific content type.
// You do not need to do this for all content served, as Go's http handler
// will try to guess the content type by the name of the file or the content itself.
// This is for those situations where Go's default is not working.
//
// The extension must begin with a dot and have only one dot in it.
func RegisterContentType(extension string, contentType string) {
	if contentTypes == nil {
		contentTypes = make(map[string]string)
	}
	contentTypes[extension] = contentType
}
