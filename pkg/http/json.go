package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/goradd/goradd/pkg/log"
	"io"
	"net/http"
	"strings"
)

// ParseJsonBody will look for json in the request, parse it into the given dest, and handle errors.
//
// The dest should be a pointer to a structure or some other value you want filled with the data.
// Errors will result in an appropriate error response written to the response writer, and an error response.
// If maxBytes is reached, it will close the connection and error.
func ParseJsonBody(w http.ResponseWriter, r *http.Request, maxBytes int64, dest any) error {
	v, _ := ParseValueAndParams(r.Header.Get("Content-Type"))

	if v != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Content-Type must be application/json")
		return fmt.Errorf("content-type must be application/json")
	}

	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxBytes))
	dec := json.NewDecoder(bytes.NewReader(body))

	if err == nil {
		err = dec.Decode(dest)
	}

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		// Catch any syntax errors in the JSON and send an error message
		// which interpolates the location of the problem to make it
		// easier for the client to fix.
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(w, msg)
			return fmt.Errorf(msg)

		// In some circumstances Decode() may also return an
		// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
		// is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("request body contains badly-formed JSON")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(w, msg)
			return fmt.Errorf(msg)

		// Catch any type errors, like trying to assign a string in the
		// JSON request body to a int field in our Person struct. We can
		// interpolate the relevant field name and position into the error
		// message to make it easier for the client to fix.
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(w, msg)
			return fmt.Errorf(msg)

		// Catch the error caused by extra unexpected fields in the request
		// body. We extract the field name from the error message and
		// interpolate it in our custom error message. There is an open
		// issue at https://github.com/golang/go/issues/29035 regarding
		// turning this into a sentinel error.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(w, msg)
			return fmt.Errorf(msg)

		// An io.EOF error is returned by Decode() if the request body is
		// empty.
		case errors.Is(err, io.EOF):
			msg := "request body must not be empty"
			w.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(w, msg)
			return fmt.Errorf(msg)

		// Catch the error caused by the request body being too large. Again
		// there is an open issue regarding turning this into a sentinel
		// error at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			msg := fmt.Sprintf("request body must not be larger than %d bytes", maxBytes)
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			_, _ = io.WriteString(w, msg)
			return fmt.Errorf(msg)

		// Otherwise default to logging the error and sending a 500 Internal
		// Server Error response.
		default:
			log.Debug(err.Error())
			msg := http.StatusText(http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = io.WriteString(w, msg)
			return fmt.Errorf(msg)
		}
	}

	// Call decode again, using a pointer to an empty anonymous struct as
	// the destination. If the request body only contained a single JSON
	// object this will return an io.EOF error. So if we get anything else,
	// we know that there is additional data in the request body.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "request body must only contain a single JSON object"
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, msg)
		return fmt.Errorf(msg)
	}
	return nil
}
