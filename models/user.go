package models

type UserInfo struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type Response struct {
	Message string `json:"message" binding:"required"`
	Type    int    `json:"-" binding:"required"` // HTTP status code, not included in JSON response
}
