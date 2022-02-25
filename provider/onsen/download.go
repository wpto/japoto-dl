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

func (ep *FeedRawEp) Download(dl model.Loader, tasks model.Tasks) error {
	hls := tasks.AudioHLS()

	if ep.StreamingUrl == nil {
		return errors.New("onsen.dl.plget: cant be loaded")
	}

	playlistUrl := *ep.StreamingUrl

	playlistBody, err := dl.Text(playlistUrl, gopts)
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

		imageBody, err := dl.Raw(url, gopts)
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
		tsaudioBody, err := dl.Text(tsaudioUrl, gopts)
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

			for _, key := range keys {
				loaders <- loader(dl, tsaudioUrl, key, files)
			}

			for _, au := range audio {
				loaders <- loader(dl, tsaudioUrl, au, files)
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
						return saver(hls, f)
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

func loader(dl model.Loader, prefix string, file model.File, files chan model.File) func() error {
	return func() error {
		url, err := file.Url(prefix)
		if err != nil {
			return errors.Wrap(err, "loader.url")
		}

		body, err := dl.Raw(url, gopts)
		if err != nil {
			return errors.Wrap(err, "loader.raw")
		}

		file.SetBody(body)
		files <- file

		return nil
	}
}
func saver(hls model.AudioHLS, file model.File) error {
	err := hls.Validate(file)
	if err != nil {
		return errors.Wrap(err, "loader.validate")
	}
	return nil
}
