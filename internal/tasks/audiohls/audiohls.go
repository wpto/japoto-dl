package audiohls

type AudioHLSImpl struct {
	workdir workdir.WorkdirHLS
}

func NewAudioHLSImpl(workdir workdir.WorkdirHLS) *AudioHLSImpl {
	return &AudioHLSImpl{workdir}
}
