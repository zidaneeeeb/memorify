package cryptography

// Encryption denotes encryption functionalites.
type Encryption interface {
	// Encrypt encrypts the given plain bytes and returns the
	// encrypted bytes.
	Encrypt(plain []byte) ([]byte, error)

	// Decrypt decrypts the given encrypted bytes and returns
	// the plain bytes.
	Decrypt(cipher []byte) ([]byte, error)
}
