package tasks

import (
	"github.com/pgeowng/japoto-dl/internal/model"
	"github.com/pgeowng/japoto-dl/internal/tasks/audiohls"
)

type Tasks struct {
	ahls model.AudioHLS
	file model.File
}

func (t *Tasks) AudioHLS() model.AudioHLS {
	return t.ahls
}

func NewTasks(workdirHLS workdir.WorkdirHLS) model.Tasks {
	return &Tasks{
		ahls: audiohls.NewAudioHLSImpl(workdirHLS),
	}
}
