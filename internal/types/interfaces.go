package types

import "github.com/pgeowng/japoto-dl/model"

type (
	LoadStatus interface {
		Inc(step int)
		Total(total int)
	}

	Loader interface {
		Text(url string, opts *model.LoaderOpts) (*string, error)
		JSON(url string, dest interface{}, opts *model.LoaderOpts) error
		Raw(url string, opts *model.LoaderOpts) ([]byte, error)
	}

	LoaderNew interface {
		Url(url string) LoaderNew
		Headers(headers map[string]string) LoaderNew
		Transform(func(content []byte) (err error)) LoaderNew

		JSON(output interface{}) error
		Save(filepath string) error
	}

	AudioHLS interface {
		Playlist(body string) (tsaudio []model.File, err error)
		TSAudio(tsaudio model.File) (keys []model.File, audio []model.File, err error)
		CheckAlreadyLoaded(filename string) bool
	}
)
