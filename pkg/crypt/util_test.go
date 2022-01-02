package crypt

import (
	"reflect"
	"testing"
)

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
