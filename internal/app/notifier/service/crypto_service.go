package service

import (
	"github.com/remoodle/notifier/internal/app/notifier/repository"
)

// CryptoService provides an interface for cryptographic operations
type CryptoService interface {
	DecryptToken(token string) (string, error)
}

type cryptoServiceImpl struct {
	repo repository.CryptoRepository
}

// NewCryptoService creates a new CryptoService
func NewCryptoService(repo repository.CryptoRepository) CryptoService {
	return &cryptoServiceImpl{repo}
}

func (s *cryptoServiceImpl) DecryptToken(token string) (string, error) {
	return s.repo.DecryptToken(token)
}
