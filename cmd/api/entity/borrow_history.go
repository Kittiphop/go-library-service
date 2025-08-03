package entity

import "time"

// BorrowHistory is a model for borrow history table
type BorrowHistory struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	BookID     uint       `gorm:"not null" json:"bookId"`
	UserID     uint       `gorm:"not null" json:"userId"`
	BorrowedAt *time.Time `gorm:"default:null" json:"borrowedAt"`
	ReturnedAt *time.Time `gorm:"default:null" json:"returnedAt,omitempty"`
	Status     string     `gorm:"type:varchar(20);not null" json:"status"` //  "borrowed", "returned"
	CreatedAt  *time.Time `gorm:"default:now()" json:"createdAt"`
	UpdatedAt  *time.Time `gorm:"default:now()" json:"updatedAt"`
}

// BorrowBookRequest is a request for borrow a book
type BorrowBookRequest struct {
	BookID uint `json:"bookId" validate:"required"`
	UserID uint `json:"-"`
}

// ReturnBookRequest is a request for return a book
type ReturnBookRequest struct {
	HistoryID uint `json:"historyId" validate:"required"`
	BookID 	  uint `json:"bookId" validate:"required"`
	UserID    uint `json:"-"`
}

// BorrowHistoryResponse represents the response for borrow history
type BorrowHistoryResponse struct {
	ID         uint       `json:"id"`
	BookID     uint       `json:"bookId"`
	UserID     uint       `json:"userId"`
	BorrowedAt time.Time  `json:"borrowedAt"`
	ReturnedAt *time.Time `json:"returnedAt,omitempty"`
	Status     string     `json:"status"`
}
