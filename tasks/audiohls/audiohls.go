package audiohls

import "github.com/pgeowng/japoto-dl/workdir"

type AudioHLSImpl struct {
	workdir workdir.WorkdirHLS
}

func NewAudioHLSImpl(workdir workdir.WorkdirHLS) *AudioHLSImpl {
	return &AudioHLSImpl{workdir}
}
