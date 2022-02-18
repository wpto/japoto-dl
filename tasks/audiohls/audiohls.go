package audiohls

import "github.com/pgeowng/japoto-dl/workdir"

type AudioHLSImpl struct {
	workdir *workdir.Workdir
}

func NewAudioHLSImpl(workdir *workdir.Workdir) *AudioHLSImpl {
	return &AudioHLSImpl{workdir}
}
