package model

type AudioHLS interface {
	Playlist(body string) (tsaudio []File, err error)
	TSAudio(tsaudio File) (keys []File, audio []File, err error)
	Validate(file File) error
}

type Tasks interface {
	AudioHLS() AudioHLS
}
