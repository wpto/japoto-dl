package workdir

import (
	"fmt"
	"os"

	"github.com/pgeowng/japoto-dl/internal/repo/ffmpeg"
	"github.com/pgeowng/japoto-dl/internal/repo/wd"
	"github.com/pkg/errors"
)

type Workdir struct {
	wd.Wd
	muxer   ffmpeg.MuxerHLS
	namemap map[string]string
	written map[string]bool
}

func NewWorkdir(wd *wd.Wd, muxer ffmpeg.MuxerHLS, namemap map[string]string) *Workdir {
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

func (wd *Workdir) checkMuxFiles() (playlistPath string, imagePath *string, err error) {
	if wd.WasWritten("image") {
		str := wd.Resolve(wd.ResolveName("image"))
		imagePath = &str
	} else {
		fmt.Println("image wasn't written")
	}

	playlistPath = wd.Resolve(wd.ResolveName("playlist"))
	if !wd.WasWritten("playlist") {
		err = errors.New("wd.mux: playlist was not written")
		str := wd.Resolve(wd.ResolveName("image"))
		imagePath = &str
	}

	return
}

func (wd *Workdir) Mux() error {
	playlistPath, imagePath, err := wd.checkMuxFiles()
	if err != nil {
		return err
	}

	err = wd.muxer.Mux(playlistPath, imagePath)
	if err != nil {
		return errors.Wrap(err, "wd.mux")
	}
	return nil
}

func (wd *Workdir) ForceMux() error {
	playlistPath, imagePath, _ := wd.checkMuxFiles()
	err := wd.muxer.Mux(playlistPath, imagePath)
	if err != nil {
		return errors.Wrap(err, "wd.mux")
	}
	return nil
}

func (wd *Workdir) AlreadyLoadedChunks() map[string]bool {
	result := map[string]bool{}
	files, err := os.ReadDir(wd.CacheDir())
	if err == nil {
		for _, file := range files {
			result[file.Name()] = true
		}
	}

	return result
}
