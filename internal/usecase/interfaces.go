package usecase

import "github.com/pgeowng/japoto-dl/internal/model"

type PrintLine interface {
	SetChunk(count int)
	AddChunk()
	SetChunkCount(count int)
	SetPrefix(prefix string)
	Status(str string)
	Error(err error) error
}

type Provider interface {
	GetFeed(loader model.Loader) ([]model.ShowAccess, error)
	Label() string
}

type History interface {
	Check(key string) bool
	Write(key string) error
}

type WorkdirHLS interface {
	model.WorkdirFile
}

type WorkdirHLSMuxer interface {
	model.WorkdirFile
	Mux() error
	ForceMux() error
}

type DL interface {
	Text(url string, opts *model.LoaderOpts) (*string, error)
	JSON(url string, dest interface{}, opts *model.LoaderOpts) error
	Raw(url string, opts *model.LoaderOpts) ([]byte, error)
}
