package crypt

import (
	"context"
	"github.com/goradd/goradd/pkg/session"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newMockContext() (ctx context.Context) {
		s := session.NewMock()
		session.SetSessionManager(s)
		ctx = s.With(context.Background())
		return
}

func TestSessionEncryptUrlValue(t *testing.T) {
	ctx := newMockContext()
	v := SessionEncryptUrlValue(ctx, "abc?")
	assert.NotEqual(t, "abc?", v)
	v2 := SessionDecryptUrlValue(ctx, v)
	assert.Equal(t, "abc?", v2)

	// test using same key
	v = SessionEncryptUrlValue(ctx, "def*")
	assert.NotEqual(t, "def*", v)
	v2 = SessionDecryptUrlValue(ctx, v)
	assert.Equal(t, "def*", v2)

}

func TestSessionEncryptErrors(t *testing.T) {
	ctx := newMockContext()

	v := SessionDecrypt(ctx, []byte("abc"))
	assert.Empty(t, v)

	v2 := SessionDecryptUrlValue(ctx, "abc")
	assert.Empty(t, v2)
}

func TestSessionEncryptErrors2(t *testing.T) {
	ctx := newMockContext()
	_ = SessionEncryptUrlValue(ctx, "abc?")
	v2 := SessionDecryptUrlValue(ctx, "abc?")
	assert.Empty(t, v2)
}



func TestSessionEncryptAttributeValue(t *testing.T) {
	ctx := newMockContext()
	v := SessionEncryptAttributeValue(ctx, "abc?")
	assert.NotEqual(t, "abc?", v)
	v2 := SessionDecryptAttributeValue(ctx, v)
	assert.Equal(t, "abc?", v2)
}

