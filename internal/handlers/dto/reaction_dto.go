// reaction_dto.go
package handlers

// SetReactionRequest represents the request body for setting a reaction
type SetReactionRequest struct {
	Type string `json:"type" example:"like" enums:"like,dislike"`
}
