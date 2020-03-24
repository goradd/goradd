package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a base64 encoded
// securely generated random string. The given value is the length of the initial random bytes represented by
// the base64 string. The actual length of the string will be bigger.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.StdEncoding.EncodeToString(b), err
}

// Encrypt will encrypt the data using AES-GCM and return the encrypted data.
// Using Decrypt to decrypt the data.
// key is recommended to be 16 bytes long, and will produce an AES-128 encryption using GCM
// which is currently considered secure. Note that this is due to our use of GCM. When using CBC,
// AES-256 is considered secure, but GCM is better and faster.
// Errors will panic, since they are caused by some kind of system wide failure or a bad key size.
func Encrypt(data []byte, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil) // Here we append the encrypted text to the nonce.
	return ciphertext
}

// Decrypt will decrypt something created by Encrypt. If the data is not decryptable, either because it is corrupt,
// or the key changed, an error is returned. Otherwise, if there is a system failure, it will panic.
func Decrypt(data []byte, key []byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
