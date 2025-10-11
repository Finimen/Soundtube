// login_dto.go
package handlers

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Username string `json:"username" example:"johndoe"`
	Password string `json:"password" example:"securepassword123"`
}

// LogoutRequest represents the request body for user logout
type LogoutRequest struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}
