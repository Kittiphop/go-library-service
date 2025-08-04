package handler

import (
	"errors"
	"go-library-service/cmd/api/entity"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

//go:generate mockgen -destination=./mock/handler.go -source=./handler.go -package=mock

// Handler manages the handler
type Handler struct {
	deps   *Dependencies
	config *Config
}

// Dependencies is a dependencies of handler
type Dependencies struct {
	Service   Service
	Validator *validator.Validate
}

// Config is a configuration of handler
type Config struct{}

// Service manages the service layer
type Service interface {
	// Book
	CreateBook(request entity.BookCreateRequest) error
	ListBook(req entity.ListBookRequest) ([]entity.BookResponse, error)
	GetBookByID(bookID uint) (*entity.BookResponse, error)
	UpdateBook(req entity.BookUpdateRequest) error
	ListLatestBooks() ([]entity.BookResponse, error)

	// User
	CreateUser(user entity.UserCreateRequest) (*uint, error)
	GetUserByID(userID uint) (*entity.UserResponse, error)
	LoginUser(username, password string) (*entity.User, error)
	UpdateUser(user entity.UserUpdateRequest) error
	DeleteUser(userID uint) error

	// Borrow
	BorrowBook(req entity.BorrowBookRequest) (*entity.BorrowHistory, error)
	ReturnBook(req entity.ReturnBookRequest) error
	GetBookBorrowHistory(bookID uint) ([]entity.BorrowHistoryResponse, error)
}

// NewHandler creates a new handler
func NewHandler(deps *Dependencies, config *Config) *Handler {
	return &Handler{
		deps:   deps,
		config: config,
	}
}

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
// InitRoute initialize the route
func InitRoute(router *gin.RouterGroup, handler *Handler) error {
	if handler == nil {
		return errors.New("unable to initialize route, handler is nil")
	}

	publicRoute := router.Group("")
	{
		publicRoute.POST("/login", handler.Login)
		publicRoute.POST("/register", handler.CreateUser)
		publicRoute.GET("/books/latest", handler.ListLatestBooks)
	}
	

	RegisterUserRoutes(router, handler)
	RegisterBookRoutes(router, handler)
	
	return nil
}

// getJWTInfo get jwt info
func (h *Handler) getJWTInfo(c *gin.Context) (userID uint) {
	userID = c.MustGet("userID").(uint)

	return userID
}
