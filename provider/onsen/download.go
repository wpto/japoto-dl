package onsen

import (
	"fmt"

	"github.com/pgeowng/japoto-dl/model"
	"github.com/pkg/errors"
)

var gopts *model.LoaderOpts = &model.LoaderOpts{
	Headers: map[string]string{
		"Referer": "https://www.onsen.ag/",
	},
}

func (ep *OnsenEpisode) Download(loader model.Loader, tasks model.Tasks, pl model.PrintLine) error {
	pl.SetPrefix(fmt.Sprintf("%s/%s", ep.Show().Provider(), ep.EpId()))
	pl.SetChunk(0)
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
			pl.Error(errors.New("empty poster(1) image"))
		}

		if len(ep.showRef.Image.Url) == 0 {
			pl.Error(errors.New("empty show(2) image"))
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

		// total := len(keys) + len(audio)

		done := make(chan struct{})

		loadChan := make(chan *model.File, 10)
		loadError := make(chan error)

		validateChan := make(chan *model.File, 20)
		validateError := make(chan error)

		go func() {
			defer close(validateChan)
			for file := range loadChan {
				select {
				case <-done:
					return
				default:
					url, err := file.Url(tsaudioUrl)
					if err != nil {
						loadError <- errors.Wrap(err, "onsen.dl.file")
						return
					}

					body, err := loader.Raw(url, gopts)
					if err != nil {
						loadError <- errors.Wrap(err, "onsen.dl.file")
						return
					}

					pl.AddChunk()
					file.SetBody(body)
					validateChan <- file
				}
			}
		}()

		go func() {
			defer close(validateError)
			for file := range validateChan {
				select {
				case <-done:
					return
				default:
					err = tasks.AudioHLS().Validate(*file)
					if err != nil {
						validateError <- errors.Wrap(err, "onsen.dl.validate")
					}
				}
			}
		}()

		defer fmt.Printf("\n")

		links := []model.File{}
		links = append(links, keys...)
		links = append(links, audio...)
		pl.SetChunkCount(len(links))

		for idx := range links {
			select {
			case loadChan <- &links[idx]:
			case err := <-loadError:
				close(done)
				return errors.Wrap(err, "onsen.dl")
			case err := <-validateError:
				close(done)
				return errors.Wrap(err, "onsen.dl")
			}
		}
		close(loadChan)

		select {
		case err := <-loadError:
			close(done)
			return errors.Wrap(err, "onsen.dl")
		case err := <-validateError:
			if err != nil {
				close(done)
				return errors.Wrap(err, "onsen.validate")
			}
		}
	}

	errImg := <-errcImg
	if errImg != nil {
		return errors.Wrap(errImg, "onsen.dl")
	}

	return nil
}
