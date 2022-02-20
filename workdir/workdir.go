package workdir

import (
	"github.com/pgeowng/japoto-dl/workdir/muxer"
	"github.com/pgeowng/japoto-dl/workdir/wd"
	"github.com/pkg/errors"
)

type Muxer interface {
	Mux(errc chan<- error)
}

type WorkdirHLS interface {
	SavePlaylist(playlistBody string) error
	Save(fileName, fileBody string) error
	SaveRaw(fileName string, fileBody []byte) error
}

type WorkdirHLSMuxer interface {
	Mux() error
}

type WorkdirHLSImpl struct {
	wd.Wd
	playlistName string
	muxer        muxer.MuxerHLS
}

func NewWorkdirHLSImpl(wd *wd.Wd, muxer muxer.MuxerHLS, playlistName string) *WorkdirHLSImpl {
	return &WorkdirHLSImpl{
		Wd:           *wd,
		playlistName: playlistName,
		muxer:        muxer,
	}
}

func (wd *WorkdirHLSImpl) SavePlaylist(playlistBody string) error {
	err := wd.Save(wd.playlistName, playlistBody)
	return err
}

func (wd *WorkdirHLSImpl) Mux() error {
	err := wd.muxer.Mux(wd.Resolve(wd.playlistName))
	if err != nil {
		return errors.Wrap(err, "wdhlsimpl.mux")
	}
	return nil
}
