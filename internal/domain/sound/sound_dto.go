package sound

type SoundDTO struct {
	ID         int    `json:"id"`
	AuthorID   int    `json:"author_id"`
	Name       string `json:"name"`
	Album      string `json:"album"`
	Genre      string `json:"genre"`
	Duration   int    `json:"duration"`
	FileName   string `json:"file_name"`
	FilePath   string `json:"file_path"`
	FileSize   int    `json:"file_size"`
	FileFormat string `json:"file_format"`
	Status     string `json:"status"`
	UploadDate string `json:"upload_date"`
}

func (s *Sound) ToDTO() *SoundDTO {
	return &SoundDTO{
		ID:         s.id,
		AuthorID:   s.authorID,
		Name:       s.name,
		Album:      s.album,
		Genre:      s.genre,
		Duration:   s.duration,
		FileName:   s.fileName,
		FilePath:   s.filePath,
		FileSize:   s.fileSize,
		FileFormat: s.fileFormat,
		Status:     s.status,
		UploadDate: s.uploadDate,
	}
}

func SoundsToDTO(sounds []*Sound) []*SoundDTO {
	dtos := make([]*SoundDTO, len(sounds))
	for i, s := range sounds {
		dtos[i] = s.ToDTO()
	}
	return dtos
}
