package sound

import "context"

type ISoundRepository interface {
	ISoundRepositoryReader
	ISoundRepositoryWriter
}

type ISoundRepositoryReader interface {
	GetSounds(ctx context.Context) ([]*Sound, error)
	GetSoundByID(ctx context.Context, id int) (*Sound, error)
	GetSoundByName(ctx context.Context, name string) (*Sound, error)
}

type ISoundRepositoryWriter interface {
	CreateSound(ctx context.Context, sound *Sound) error
	DeleteSound(ctx context.Context, sound *Sound) error
	UpdateSoundFile(ctx context.Context, name, filename, filepath string, fileSize int64) error
}
