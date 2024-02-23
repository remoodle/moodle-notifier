package repository

import (
	"github.com/fernet/fernet-go"
)

// CryptoRepository defines the interface for cryptographic operations
type CryptoRepository interface {
	DecryptToken(token string) (string, error)
}

type cryptoRepositoryImpl struct {
	key string // Fernet key
}

// NewCryptoRepository creates a new instance of CryptoRepository with the given Fernet key
func NewCryptoRepository(key string) CryptoRepository {
	return &cryptoRepositoryImpl{key}
}

func (c *cryptoRepositoryImpl) DecryptToken(token string) (string, error) {
	key, err := fernet.DecodeKey(c.key)
	if err != nil {
		return "", err
	}

	message := fernet.VerifyAndDecrypt([]byte(token), 0, []*fernet.Key{key})

	return string(message), nil
}
