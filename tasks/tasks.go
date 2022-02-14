package tasks

type AudioHLS interface {
	Playlist() error
}

type Tasks struct {
	Audio AudioHLS
}

func NewTasks() *Tasks {
	return &Tasks{
		Audio: NewAudioHLSImpl(),
	}
}
