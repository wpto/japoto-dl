package onsen

import (
	"fmt"

	"github.com/pgeowng/japoto-dl/internal/types"
	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/provider/common"
	"github.com/pkg/errors"
)

var gopts *model.LoaderOpts = &model.LoaderOpts{
	Headers: map[string]string{
		"Referer": "https://www.onsen.ag/",
	},
}

type OnsenUsecase struct{}

func (uc *OnsenUsecase) DownloadEpisode(loader types.Loader, hls types.AudioHLS, status types.LoadStatus, ep *OnsenEpisode) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("OnsenUsecase.DownloadEpisode: %w", err)
		}
	}()

	if ep.StreamingUrl == nil {
		err = fmt.Errorf("StreamingURL is not presented")
		return
	}

	playlistURL := *ep.StreamingUrl

	// TODO rewrites playlist file in any case. should be like that?
	tsaudio, err := common.LoadPlaylist(playlistURL, gopts, loader, hls)
	if err != nil {
		err = errors.Wrap(err, "onsen")
		return
	}

	if len(tsaudio) > 1 {
		err = errors.New("onsen.dl: tsaudio size > 1: not implemented")
		return
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
		var keys []model.File
		var audio []model.File
		var tsaudioUrl string
		keys, audio, tsaudioUrl, err = common.LoadTSAudio(playlistURL, gopts, ts, loader, hls)
		if err != nil {
			err = errors.Wrap(err, "onsen.dl.tsparse")
			return
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
					status.Inc(1)
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
					err = hls.Validate(*file)
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
		status.Total(len(links))

		for idx := range links {
			select {
			case loadChan <- &links[idx]:
			case err = <-loadError:
				close(done)
				return
			case err = <-validateError:
				close(done)
				return
			}
		}
		close(loadChan)

		select {
		case err = <-loadError:
			close(done)
			return
		case err = <-validateError:
			if err != nil {
				close(done)
				err = fmt.Errorf("validate: %w", err)
				return
			}
		}
	}

	errImg := <-errcImg
	if errImg != nil {
		return
	}

	return
}
