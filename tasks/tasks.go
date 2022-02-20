package tasks

import (
	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/tasks/audiohls"
	"github.com/pgeowng/japoto-dl/workdir"
)

type AudioHLS interface {
	Playlist(body string) (tsaudio []model.File, err error)
	TSAudio(tsaudio model.File) (keys []model.File, audio []model.File, err error)
	Validate(file model.File) error
}

type Tasks struct {
	AudioHLS AudioHLS
}

func NewTasks(workdirHLS workdir.WorkdirHLS) *Tasks {
	return &Tasks{
		AudioHLS: audiohls.NewAudioHLSImpl(workdirHLS),
	}
}
