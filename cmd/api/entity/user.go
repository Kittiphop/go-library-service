package entity

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID       	uint      	`gorm:"primaryKey" json:"id"`
	Name      	string    	`gorm:"not null" json:"name"`
	Username 	string    	`gorm:"unique;not null" json:"username"`
	Password 	string    	`gorm:"not null" json:"password"`
	Role 		string 		`gorm:"not null" json:"role"`
	CreatedAt 	*time.Time 	`gorm:"default:now()" json:"createdAt"`
	UpdatedAt 	*time.Time 	`gorm:"default:now()" json:"updatedAt"`
	DeletedAt 	*time.Time 	`gorm:"default:null" json:"deletedAt"`

	BorrowHistories []BorrowHistory `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
}

// UserLoginRequest is a request for log in
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represent a response for user
type LoginResponse struct {
	Token string `json:"token"`
}


// UserCreateRequest is a request for create a user
type UserCreateRequest struct {
	Name  		string 	`json:"name" binding:"required"`
	Username 	string 	`json:"username" binding:"required"`
	Password 	string 	`json:"password" binding:"required"`
}

// UserUpdateRequest is a request for update a user
type UserUpdateRequest struct {
	ID 			uint 		`json:"-"`
	Name 		string 		`json:"name" binding:"required"`
}

// UserResponse represents a response for user
type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// UserDeleteRequest is a request for delete a user
type UserDeleteRequest struct {
	ID uint `form:"id" binding:"required"`
}
