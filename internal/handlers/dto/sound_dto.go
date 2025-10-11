// sound_dto.go
package handlers

// Sound represents the request body for create sound
type CreateSoundRequest struct {
	Name  string `json:"name" example:"My Awesome Sound"`
	Album string `json:"album" example:"Summer Vibes"`
	Genre string `json:"genre" example:"Electronic"`
}

// Sound represents the request body for update sound
type UpdateSoundRequest struct {
	Name  string `json:"name" example:"Updated Sound Name"`
	Album string `json:"album" example:"Updated Album"`
	Genre string `json:"genre" example:"Updated Genre"`
}

// Sound represents the request body for delete sound
type DeleteSoundRequest struct {
	Name string `json:"name" example:"Sound to delete"`
}
