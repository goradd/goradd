package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestParseJsonBody(t *testing.T) {
	t.Run("valid_json", func(t *testing.T) {
		// Setup test
		jsonBody := `{"name": "John Doe", "age": 30}`
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader([]byte(jsonBody)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Call the function being tested
		var user person
		_ = ParseJsonBody(w, req, int64(len(jsonBody)), &user)

		// Verify result
		if user.Name != "John Doe" || user.Age != 30 {
			t.Errorf("ParseJsonBody() returned incorrect result: got %+v, expected %+v", user, person{Name: "John Doe", Age: 30})
		}
	})

	t.Run("invalid_content_type", func(t *testing.T) {
		// Setup test
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set("Content-Type", "text/plain")
		w := httptest.NewRecorder()

		// Call the function being tested
		ParseJsonBody(w, req, int64(10), nil)

		// Verify result
		if w.Result().StatusCode != http.StatusBadRequest {
			t.Errorf("ParseJsonBody() returned incorrect status code: got %d, expected %d", w.Result().StatusCode, http.StatusBadRequest)
		}
	})

	t.Run("badly_formed_json", func(t *testing.T) {
		// Setup test
		jsonBody := `{"name": "John Doe", "age": "30}`
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader([]byte(jsonBody)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Call the function being tested
		ParseJsonBody(w, req, int64(len(jsonBody)), nil)

		// Verify result
		if w.Result().StatusCode != http.StatusBadRequest {
			t.Errorf("ParseJsonBody() returned incorrect status code: got %d, expected %d", w.Result().StatusCode, http.StatusBadRequest)
		}
	})
}
