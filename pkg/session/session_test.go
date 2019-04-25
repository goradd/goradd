package session_test

import (
	"context"
	"github.com/goradd/goradd/pkg/session"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// runRequestTest will run a session test by first calling the setupHandler, and then calling the testHandler
// mimicking a process where a session variable is set in one request, and then retrieved in a later request
func runRequestTest(t *testing.T, setupHandler, testHandler http.Handler) {
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	h := session.Use(setupHandler)
	h.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v, body: %v",
			status, http.StatusOK, rec.Body)
	}

	// extract cookie
	cookie := rec.Header().Get("Set-Cookie")

	// now run it through the tester
	req = httptest.NewRequest("GET", "/", nil)
	rec = httptest.NewRecorder()
	req.Header.Set("Cookie", cookie)

	h = session.Use(testHandler)
	h.ServeHTTP(rec, req)
}

const intKey = "test.intKey"
const boolKey = "test.boolKey"
const stringKey = "test.stringKey"
const floatKey = "test.floatKey"

func setRequestHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		setupTest(ctx)
	}
	return http.HandlerFunc(fn)
}

func testRequestHandler(t *testing.T) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		runTest(t, ctx)

	}
	return http.HandlerFunc(fn)
}

func setupTest(ctx context.Context) {
	session.SetInt(ctx, intKey, 4) // testing replacing a value here
	session.SetInt(ctx, intKey, 5)
	session.SetBool(ctx, boolKey, true)
	session.SetString(ctx, stringKey, "Here")
	session.SetFloat(ctx, floatKey, 7.6)
}


func runTest(t *testing.T, ctx context.Context) {
	i,ok := session.GetInt(ctx, intKey)
	assert.Equal(t, 5, i)
	assert.True(t, ok)
	assert.True(t, session.Has(ctx, intKey))
	assert.False(t, session.Has(ctx, "randomval"))

	// test that getting the wrong kind of value produces error
	s,ok := session.GetString(ctx, intKey)
	assert.False(t, ok)
	assert.Equal(t, s, "")

	b,ok := session.GetBool(ctx, boolKey)
	assert.True(t, ok)
	assert.True(t, b)

	f,ok := session.GetFloat(ctx, floatKey)
	assert.True(t, ok)
	assert.Equal(t, 7.6, f)
	// repeat
	f,ok = session.GetFloat(ctx, floatKey)
	assert.True(t, ok)
	assert.Equal(t, 7.6, f)

	session.Clear(ctx)
	f,ok = session.GetFloat(ctx, floatKey)
	assert.False(t, ok)
	assert.Equal(t, 0.0, f)

}
