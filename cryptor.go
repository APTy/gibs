package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

const (
	blockSize = 32
)

// Cryptor handles encryption/decryption with a shared secret
type Cryptor interface {
	Encrypt(b []byte) []byte
	Decrypt(b []byte) ([]byte, error)
}

func NewCryptor(b64Key string) (Cryptor, error) {
	key, err := base64.StdEncoding.DecodeString(b64Key)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &cryptor{aead: aesgcm}, nil
}

type cryptor struct {
	aead cipher.AEAD
}

func (c *cryptor) Encrypt(b []byte) []byte {
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}
	return c.aead.Seal(nonce, nonce, b, nil)
}

func (c *cryptor) Decrypt(b []byte) ([]byte, error) {
	nonce := b[:c.aead.NonceSize()]
	ciphertext := b[c.aead.NonceSize():]

	plaintext, err := c.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
