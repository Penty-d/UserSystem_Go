package models

import "time"

type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"` //hashed password
	Role      string    `json:"role"`     //admin user
	Email     string    `json:"email"`
	FullName  string    `json:"fullname"`
	Status    string    `json:"status"` // active, inactive or deleted
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Role     string `json:"role" binding:"required,oneof=admin user"`
	Email    string `json:"email" binding:"required,email,max=100"`
	FullName string `json:"fullname" binding:"required,max=50"`
}

type UpdateUserRequest struct {
	Username string  `json:"username" binding:"required,max=50"`
	Password *string `json:"password,omitempty" binding:"omitempty,min=6,max=50"`
	Role     *string `json:"role,omitempty" binding:"omitempty,oneof=admin user"`
	Email    *string `json:"email,omitempty" binding:"omitempty,email,max=100"`
	FullName *string `json:"fullname,omitempty" binding:"omitempty,max=50"`
	Status   *string `json:"status,omitempty" binding:"omitempty,oneof=active inactive deleted"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

type Response struct {
	Message string `json:"message" binding:"required"`
	Type    int    `json:"-" binding:"required"` // HTTP status code, not included in JSON response
}
