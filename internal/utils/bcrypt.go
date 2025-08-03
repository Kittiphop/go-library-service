package utils

import "golang.org/x/crypto/bcrypt"

type Bcrypt interface {
	Hash(password string) (string, error)
	Compare(password, hash string) error
}

type BcryptService struct {
	Bcrypt Bcrypt
}

func NewBcryptService() *BcryptService {
	return &BcryptService{}
}

func (b *BcryptService) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (b *BcryptService) Compare(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
