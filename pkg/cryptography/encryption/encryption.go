package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

var (
	errInvalidEncryptionKey = errors.New("invalid encryption key")
	errInvalidCypher        = errors.New("invalid cypher")
)

type encryption struct {
	gcm cipher.AEAD
}

func New(key []byte) (*encryption, error) {
	// validate key
	if len(key) != 32 {
		return nil, errInvalidEncryptionKey
	}

	// generate a new aes cipher using the 32 byte long key
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// gcm or Galois/Counter Mode, is a mode of operation for
	// symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	// construct object
	e := &encryption{
		gcm: gcm,
	}

	return e, nil
}

func (e *encryption) Encrypt(plain []byte) ([]byte, error) {
	// creates a new byte array the size of the nonce which
	// must be passed to Seal
	nonce := make([]byte, e.gcm.NonceSize())

	// populates nonce with a cryptographically secure random
	// sequence
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Seal encrypts and authenticates plaintext, authenticates
	// the additional data and appends the result to dst,
	// returning the updated slice.
	bytes := e.gcm.Seal(nonce, nonce, plain, nil)

	return bytes, nil
}

func (e *encryption) Decrypt(cipher []byte) ([]byte, error) {
	nonceSize := e.gcm.NonceSize()
	if len(cipher) < nonceSize {
		return nil, errInvalidCypher
	}

	nonce, c := cipher[:nonceSize], cipher[nonceSize:]
	plain, err := e.gcm.Open(nil, nonce, c, nil)
	if err != nil {
		return nil, err
	}

	return plain, nil
}
