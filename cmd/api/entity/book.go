package entity

import "time"

// Book is a model for book table
type Book struct {
	ID     uint 		`gorm:"primaryKey" json:"id"`
	Title  string 		`gorm:"not null" json:"title"`
	Author string		`gorm:"not null" json:"author"`
	Price  float64		`gorm:"not null" json:"price"`
	Stock  uint			`gorm:"not null" json:"stock"`
	CreatedAt *time.Time	`gorm:"default:now()" json:"createdAt"`
	UpdatedAt *time.Time	`gorm:"default:now()" json:"updatedAt"`
	DeletedAt *time.Time	`gorm:"index" json:"deletedAt"`

	BorrowHistories []BorrowHistory `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`	
}

// BookCreateRequest is a request for creating a book
type BookCreateRequest struct {
	Title  string 	`json:"title" validate:"required"`
	Author string	`json:"author" validate:"required"`
	Price  float64	`json:"price" validate:"required"`
	Stock  uint		`json:"stock" validate:"required,min=1"`
}

// ListBookRequest is a request for listing books
type ListBookRequest struct {
	Page   int 		`form:"page" validate:"required,min=1"`
	Size   int 		`form:"size" validate:"required,min=1"`
	Search *string 	`form:"search"`
}

// BookUpdateRequest is a request for updating a book
type BookUpdateRequest struct {
	ID     uint 		`json:"id" validate:"required"`
	Title  string 		`json:"title" validate:"required"`
	Author string		`json:"author" validate:"required"`
	Price  float64		`json:"price" validate:"required"`
	Stock  uint			`json:"stock" validate:"required,min=1"`
}

// BookResponse represents a response for book
type BookResponse struct {
	ID     uint 		`json:"id"`
	Title  string 		`json:"title"`
	Author string		`json:"author"`
	Price  float64		`json:"price"`
	Stock  uint			`json:"stock"`
	CreatedAt *time.Time	`json:"createdAt"`
	UpdatedAt *time.Time	`json:"updatedAt"`
}

