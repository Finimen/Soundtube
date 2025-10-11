// upload_dto.go
package handlers

// Upload represents the request body for upload sound
type UploadRequest struct {
	Message  string `json:"message" example:"file uploaded successfully"`
	Filename string `json:"filename" example:"mysound.mp3"`
	Path     string `json:"path" example:"uploads/mysound.mp3"`
	FullPath string `json:"full_path" example:"/app/static/uploads/mysound.mp3"`
}
