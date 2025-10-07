package sound

import "errors"

type Sound struct {
	id       int
	authorID int
	name     string
	album    string
	genre    string
	duration int

	fileName   string
	filePath   string
	fileSize   int
	fileFormat string

	status     string
	uploadDate string
}

func (s *Sound) ID() int       { return s.id }
func (s *Sound) AuthorID() int { return s.authorID }

func (s *Sound) Name() string  { return s.name }
func (s *Sound) Ablum() string { return s.album }
func (s *Sound) Genre() string { return s.genre }
func (s *Sound) Duration() int { return s.duration }

func (s *Sound) FileName() string   { return s.fileName }
func (s *Sound) FilePath() string   { return s.filePath }
func (s *Sound) FileSize() int      { return s.fileSize }
func (s *Sound) FileFormat() string { return s.fileFormat }

func NewSound(name, album, genre string, authorID int) (*Sound, error) {
	if name == "" {
		return nil, errors.New("sound name cannot be empty")
	}
	if album == "" {
		return nil, errors.New("album name cannot be empty")
	}
	if genre == "" {
		return nil, errors.New("genre name cannot be empty")
	}
	if authorID < 0 {
		return nil, errors.New("invalid id")
	}

	return &Sound{
		authorID: authorID,
		name:     name,
		album:    album,
		genre:    genre,
	}, nil
}

func RebuildSoundFromStorage(id, authorID, duration int, name, album, genre, fileName, filePath string, fileSize int, fileFormat, uploadDate string) *Sound {
	return &Sound{
		id:         id,
		authorID:   authorID,
		name:       name,
		album:      album,
		genre:      genre,
		duration:   duration,
		fileName:   fileName,
		filePath:   filePath,
		fileSize:   fileSize,
		fileFormat: fileFormat,
		uploadDate: uploadDate,
	}
}
