package onsen

import (
	"github.com/pgeowng/japoto-dl/helpers"
	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/tasks"
	"github.com/pgeowng/japoto-dl/types"
	"github.com/pkg/errors"
)

type Onsen struct {
	dl    types.Loader
	tasks *tasks.Tasks
}

func NewOnsen(dl types.Loader, tasks *tasks.Tasks) *Onsen {
	return &Onsen{dl, tasks}
}

var gopts *types.LoaderOpts = &types.LoaderOpts{
	Headers: map[string]string{
		"Referer": "https://www.onsen.ag/",
	},
}

func (p *Onsen) Download(playlistUrl string) error {

	playlistBody, err := p.dl.Text(playlistUrl, gopts)
	if err != nil {
		return errors.Wrap(err, "onsen.dl.plget")
	}

	tsaudio, err := p.tasks.AudioHLS.Playlist(*playlistBody)
	if err != nil {
		return errors.Wrap(err, "onsen.dl.plparse")
	}

	if len(tsaudio) > 1 {
		return errors.New("onsen.dl: tsaudio size > 1: not implemented")
	}

	for _, ts := range tsaudio {
		tsaudioUrl, err := ts.Url(playlistUrl)
		if err != nil {
			return errors.Wrap(err, "onsen.dl.tsurl")
		}
		tsaudioBody, err := p.dl.Text(tsaudioUrl, gopts)
		if err != nil {
			return errors.Wrap(err, "onsen.dl.tsget")
		}

		ts.SetBodyString(tsaudioBody)

		keys, audio, err := p.tasks.AudioHLS.TSAudio(ts)
		if err != nil {
			return errors.Wrap(err, "onsen.dl.tsparse")
		}

		total := len(keys) + len(audio)

		files := make(chan model.File, total)
		loaders := make(chan func() error)
		loaded := make(chan error)
		go helpers.EachLimit(loaders, loaded, 10)
		go func() {
			defer close(loaders)

			for _, key := range keys {
				loaders <- p.loader(tsaudioUrl, key, files)
			}

			for _, au := range audio {
				loaders <- p.loader(tsaudioUrl, au, files)
			}
		}()

		savers := make(chan func() error)
		saved := make(chan error)
		go helpers.EachLimit(savers, saved, 100)
		go func() {
			defer close(savers)
			for file := range files {
				savers <- func(f model.File) func() error {
					return func() error {
						return p.saver(f)
					}
				}(file)
			}
		}()

		errc := make(chan error)
		go func() {
			err = <-loaded
			close(loaded)
			close(files)
			if err != nil {
				err = errors.Wrap(err, "onsen.dl.files_load")
			}
			errc <- err
		}()

		go func() {
			err = <-saved
			close(saved)
			if err != nil {
				err = errors.Wrap(err, "onsen.dl.files_save")
			}
			errc <- err
		}()

		for i := 0; i < 2; i++ {
			err = <-errc
			if err != nil {
				return err
			}
		}

	}

	err = p.tasks.AudioHLS.Mux()
	if err != nil {
		return errors.Wrap(err, "onsen.dl.mux")
	}

	return nil
}

func (p *Onsen) loader(prefix string, file model.File, files chan model.File) func() error {
	return func() error {
		url, err := file.Url(prefix)
		if err != nil {
			return errors.Wrap(err, "loader.url")
		}

		body, err := p.dl.Raw(url, gopts)
		if err != nil {
			return errors.Wrap(err, "loader.raw")
		}

		file.SetBody(body)
		files <- file

		return nil
	}
}
func (p *Onsen) saver(file model.File) error {
	err := p.tasks.AudioHLS.Validate(file)
	if err != nil {
		return errors.Wrap(err, "loader.validate")
	}
	return nil
}
