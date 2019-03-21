package crypt

import (
	"context"
	"encoding/ascii85"
	"encoding/base64"
	"github.com/goradd/goradd/pkg/goradd"
	"github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/session"
)

// SessionEncrypt encrypts the value so that as long as the session does not change, it can be decrypted using
// SessionDecrypt. This is useful for any time you might need to send something to the browser that would pose
// a security risk, like a database id, or a pagestate in a URL.
// The resulting value is a raw byte stream and is not suitable for sending to the browser. See the other
// utilities for helpers depending on where you are sending the data.
func SessionEncrypt(ctx context.Context, value []byte) []byte {
	var salt []byte
	var err error
	if i := session.Get(ctx, goradd.SessionSalt); i == nil {
		if salt,err = GenerateRandomBytes(16); err != nil {
			panic (err)
		}
		session.Set(ctx, goradd.SessionSalt, salt)
	} else {
		salt = i.([]byte)
		if len(salt) != 16 {
			panic("Salt is the wrong length")
		}
	}
	return Encrypt(value, salt)
}

// SessionDecrypt will decrypt a value that was encrypted by the SessionEncrypt.
// If there was a problem with the decryption process, the result will be nil.
func SessionDecrypt(ctx context.Context, encryptedValue []byte) []byte {
	var salt []byte
	if i := session.Get(ctx, goradd.SessionSalt); i == nil {
		log.Warning("Attempted to decrypt a value without a key in the session.")
		return nil // we don't have an encryption key
	} else {
		salt = i.([]byte)
		if len(salt) != 16 {
			log.Error("The session has an encryption key that is the wrong length.")
			return nil
		}
	}

	v, err := Decrypt(encryptedValue, salt)
	if err != nil {
		log.Warningf("Decryption error %s", err.Error())
		return nil
	}
	return v
}

// SessionEncryptUrlValue encrypts a value such that the resulting text can be used as a value in a URL.
// Decrypt the value using SessionDecryptUrlValue
func SessionEncryptUrlValue(ctx context.Context, value string) string {
	v := SessionEncrypt(ctx, []byte(value))
	return base64.URLEncoding.EncodeToString(v)
}

// SessionDecryptUrlValue decrypts the give value that was encrypted using SessionEncryptUrlValue.
// If there was a problem, an empty string is returned.
func SessionDecryptUrlValue(ctx context.Context, encryptedValue string) string {
	v,err := base64.URLEncoding.DecodeString(encryptedValue)
	if err != nil {
		log.Warning("Bad URL Encoding in SessionDecryptUrlValue")
		return ""
	}
	return string( SessionDecrypt(ctx, []byte(v)))
}

// SessionEncryptAttributeValue encrypts a value so that it can be used as a value in an html attribute.
// Decrypt this value with SessionDecryptAttributeValue. It uses the ascii85 encoding algorithm.
func SessionEncryptAttributeValue(ctx context.Context, value string) string {
	v := SessionEncrypt(ctx, []byte(value))
	dst := make([]byte, ascii85.MaxEncodedLen(len(v)))
	count := ascii85.Encode(dst, v)
	return string(dst[:count])
}

// SessionDecryptAttributeValue decrypts a value encrypted with SessionEncryptAttributeValue.
// If there was a problem, an empty string is returned.
func SessionDecryptAttributeValue(ctx context.Context, encryptedValue string) string {
	dst := make([]byte, len(encryptedValue))
	ndst,_,err := ascii85.Decode(dst, []byte(encryptedValue), true)
	if err != nil {
		log.Warning(err)
		return ""
	}

	v := SessionDecrypt(ctx, dst[:ndst])
	if v == nil {
		return ""
	}
	return string(v)
}