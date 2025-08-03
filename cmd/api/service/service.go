package service

import (
	"go-library-service/cmd/api/entity"
	"time"
)

//go:generate mockgen -destination=./mock/service.go -source=./service.go -package=mock

// Service manages the service
type Service struct {
	deps *Dependencies
	conf *Config
}

// Dependencies is a dependencies of service
type Dependencies struct {
	BcryptService BcryptService
	PostgresRepo PostgresRepository
	RedisRepo RedisRepository
}

// Config is a configuration of service
type Config struct {
	
}

// PostgresRepository is a repository for postgres
type PostgresRepository interface{
	// User
	CreateUser(user entity.User) (*uint, error)
	GetUserByID(userID uint) (*entity.UserResponse, error)
	GetUserByUsername(username string) (*entity.User, error)
	UpdateUser(user entity.User) error
	DeleteUser(userID uint) error

	// Book
	CreateBook(book entity.Book) error
	GetBookByID(bookID uint) (*entity.BookResponse, error)
	UpdateBook(book entity.Book) error
	ListBook(req entity.ListBookRequest) ([]entity.BookResponse, error)
	ListLatestBooks() ([]entity.BookResponse, error)


	// BorrowHistory
	BorrowBook(history *entity.BorrowHistory) (*entity.BorrowHistory, error)
	ReturnBook(historyID, BookID uint, returnedAt time.Time) error
	GetBorrowHistoryByBookID(bookID uint) ([]entity.BorrowHistory, error)
	GetBorrowHistoryByID(id uint) (*entity.BorrowHistory, error)
}

// RedisRepository is a repository for redis
type RedisRepository interface{
	Get(key string) (string, error)
	Set(key string, value interface{}, expiration uint) error
	Delete(key string) error
}

// BcryptService is a service for bcrypt
type BcryptService interface{
	Hash(password string) (string, error)
	Compare(password, hash string) error
}

// NewService creates a new service
func NewService(deps *Dependencies, conf *Config) *Service {
	return &Service{
		deps: deps,
		conf: conf,
	}
}