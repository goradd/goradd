package crypt

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestGenerateRandomString(t *testing.T) {
	s,err := GenerateRandomBase64String(10)
	assert.NoError(t, err)
	b,err2 := base64.StdEncoding.DecodeString(s)
	assert.NoError(t, err2)
	assert.Equal(t, 10, len(b))
}

func TestEncryptDecrypt(t *testing.T) {
	const key = "0123456789abcdef"
	tests := []struct {
		name string
		data []byte
	}{
		{"Basic encryption", []byte("Here and there")},
		{"Empty encryption", []byte("")},
		{"Byte encryption", []byte{1,2,3,4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := Encrypt(tt.data, []byte(key))
			got,err := Decrypt(enc, []byte(key))
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(tt.data, got) {
				t.Errorf("Encrypt() = %v, want %v", got, tt.data)
			}
		})
	}
}
