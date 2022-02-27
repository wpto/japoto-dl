package onsen

import (
	"fmt"

	"github.com/pgeowng/japoto-dl/helpers"
	"github.com/pgeowng/japoto-dl/model"
	"github.com/pkg/errors"
)

var gopts *model.LoaderOpts = &model.LoaderOpts{
	Headers: map[string]string{
		"Referer": "https://www.onsen.ag/",
	},
}

func (ep *OnsenEpisode) Download(loader model.Loader, tasks model.Tasks) error {
	hls := tasks.AudioHLS()

	if ep.StreamingUrl == nil {
		return errors.New("onsen.dl.plget: cant be loaded")
	}

	playlistUrl := *ep.StreamingUrl

	playlistBody, err := loader.Text(playlistUrl, gopts)
	if err != nil {
		return errors.Wrap(err, "onsen.dl.plget")
	}

	tsaudio, err := hls.Playlist(*playlistBody)
	if err != nil {
		return errors.Wrap(err, "onsen.dl.plparse")
	}

	if len(tsaudio) > 1 {
		return errors.New("onsen.dl: tsaudio size > 1: not implemented")
	}

	errcImg := make(chan error)
	go func(errc chan<- error) {
		if len(ep.PosterImageUrl) == 0 {
			fmt.Printf("onsen.dl: empty poster image for %s", ep.EpId())
		}

		if len(ep.showRef.Image.Url) == 0 {
			fmt.Printf("onsen.dl: empty show image for %s", ep.EpId())
		}

		url := ep.PosterImageUrl

		if len(url) == 0 {
			url = ep.showRef.Image.Url
		}

		if len(url) == 0 {
			errc <- errors.New("onsen.dl: image not found")
			return
		}

		imageBody, err := loader.Raw(url, gopts)
		if err != nil {
			errc <- errors.Wrap(err, "onsen.dl.img")
			return
		}

		file := model.NewFile("", "")
		file.SetBody(imageBody)

		err = hls.Image(file)
		if err != nil {
			errc <- errors.Wrap(err, "onsen.dl.img")
			return
		}

		errc <- nil
	}(errcImg)

	for _, ts := range tsaudio {
		tsaudioUrl, err := ts.Url(playlistUrl)
		if err != nil {
			return errors.Wrap(err, "onsen.dl.tsurl")
		}
		tsaudioBody, err := loader.Text(tsaudioUrl, gopts)
		if err != nil {
			return errors.Wrap(err, "onsen.dl.tsget")
		}

		ts.SetBodyString(tsaudioBody)

		keys, audio, err := hls.TSAudio(ts)
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

			links := []model.File{}
			links = append(links, keys...)
			links = append(links, audio...)

			for _, l := range links {
				loaders <- func(f model.File) func() error {
					return func() error {
						err := fetch(loader, tsaudioUrl, &f)
						if err == nil {
							files <- f
						}
						return err
					}
				}(l)
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
						return save(tasks, &f)
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

	errImg := <-errcImg
	if errImg != nil {
		return errors.Wrap(errImg, "onsen.dl:")
	}

	return nil
}

func fetch(loader model.Loader, prefix string, file *model.File) error {
	url, err := file.Url(prefix)
	if err != nil {
		return errors.Wrap(err, "fetch.url")
	}

	body, err := loader.Raw(url, gopts)
	if err != nil {
		return errors.Wrap(err, "fetch.raw")
	}

	file.SetBody(body)
	return nil
}
func save(tasks model.Tasks, file *model.File) error {
	err := tasks.AudioHLS().Validate(*file)
	if err != nil {
		return errors.Wrap(err, "save.validate")
	}
	return nil
}
