package onsen

import (
	"fmt"

	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/provider/common"
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

	// TODO remove pointers from ep
	if ep.StreamingUrl == nil {
		return fmt.Errorf("onsen.dl: streaming url is null. this cant happened.")
	}

	// TODO rewrites playlist file in any case. should be like that?
	tsaudio, err := common.LoadPlaylist(playlistUrl, gopts, loader, hls)
	if err != nil {
		return errors.Wrap(err, "onsen")
	}

	if len(tsaudio) > 1 {
		return errors.New("onsen.dl: tsaudio size > 1: not implemented")
	}

	errcImg := make(chan error)
	go func(errc chan<- error) {
		errc <- func() error {
			imageUrl := ep.PosterImageUrl

			if len(imageUrl) == 0 {
				imageUrl = ep.showRef.Image.Url
			}

			if len(imageUrl) == 0 {
				return errors.New("onsen.dl: image not found")
			}

			return common.LoadImage(imageUrl, gopts, loader, hls)
		}()

	}(errcImg)

	for _, ts := range tsaudio {
		keys, audio, tsaudioUrl, err := common.LoadTSAudio(playlistUrl, gopts, ts, loader, hls)
		if err != nil {
			return errors.Wrap(err, "onsen.dl.tsparse")
		}

		filteredCount := len(keys) + len(audio)
		keys = common.FilterChunks(keys, hls)
		audio = common.FilterChunks(audio, hls)
		if count := filteredCount - (len(keys) + len(audio)); count > 0 {
			fmt.Printf("already loaded %d files: continue...\n", count)
		}

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
					err := common.LoadChunk(tsaudioUrl, gopts, file, loader)
					if err != nil {
						loadError <- err
						return
					}
					pl.AddChunk()
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
