package workdir

import (
	"fmt"

	"github.com/pgeowng/japoto-dl/workdir/muxer"
	"github.com/pgeowng/japoto-dl/workdir/wd"
	"github.com/pkg/errors"
)

type WorkdirFile interface {
	SaveNamed(name string, fileBody string) error
	SaveNamedRaw(name string, fileBody []byte) error
	ResolveName(name string) string
	WasWritten(name string) bool
}

type WorkdirHLS interface {
	WorkdirFile
	Save(fileName, fileBody string) error
	SaveRaw(fileName string, fileBody []byte) error
}

type WorkdirHLSMuxer interface {
	WorkdirFile
	Mux() error
}

type Workdir struct {
	wd.Wd
	muxer   muxer.MuxerHLS
	namemap map[string]string
	written map[string]bool
}

func NewWorkdir(wd *wd.Wd, muxer muxer.MuxerHLS, namemap map[string]string) *Workdir {
	return &Workdir{
		Wd:      *wd,
		muxer:   muxer,
		namemap: namemap,
		written: map[string]bool{},
	}
}

func (wd *Workdir) ResolveName(name string) string {
	p, ok := wd.namemap[name]
	if !ok {
		panic(errors.Errorf("wd: named file not found for %s", name))
	}
	return p
}

func (wd *Workdir) SaveNamed(name string, fileBody string) error {
	err := wd.Save(wd.ResolveName(name), fileBody)
	if err == nil {
		wd.written[name] = true
	}
	return err
}

func (wd *Workdir) SaveNamedRaw(name string, fileBody []byte) error {
	err := wd.SaveRaw(wd.ResolveName(name), fileBody)
	if err == nil {
		wd.written[name] = true
	}
	return err
}

func (wd *Workdir) WasWritten(name string) bool {
	_, ok := wd.written[name]
	return ok
}

func (wd *Workdir) Mux() error {
	imagePath := new(string)
	if wd.WasWritten("image") {
		str := wd.ResolveName("image")
		imagePath = &str
	} else {
		fmt.Println("image wasnot written")
	}

	if !wd.WasWritten("playlist") {
		return errors.New("wd.mux: playlist was not written")
	}

	err := wd.muxer.Mux(wd.ResolveName("playlist"), imagePath)
	if err != nil {
		return errors.Wrap(err, "wd.mux")
	}
	return nil
}
