// register_dto.go
package handlers

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Username string `json:"username" example:"johndoe"`
	Password string `json:"password" example:"securepassword123"`
	Email    string `json:"email" example:"john@example.com"`
}
