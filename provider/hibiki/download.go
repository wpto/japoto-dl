package hibiki

import (
	"fmt"

	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/provider/common"
	"github.com/pkg/errors"

	"github.com/pgeowng/japoto-dl/internal/types"
)

const (
	playCheckURL = "https://vcms-api.hibiki-radio.jp/api/v1/videos/play_check?video_id=%d"
)

type HibikiUsecase struct{}

func (uc *HibikiUsecase) DownloadEpisode(loader types.Loader, hls types.AudioHLS, status types.LoadStatus, ep *HibikiEpisodeMedia) (err error) {
	// TODO: move outside download function
	// pl.SetPrefix(fmt.Sprintf("%s/%s", ep.Show().Provider(), ep.EpId()))
	// pl.SetChunk(0)

	defer func() {
		if err != nil {
			err = fmt.Errorf("HibikiUsecase.DownloadEpisode: %w", err)
		}
	}()

	var checkObj struct {
		PlaylistURL string `json:"playlist_url"`
	}
	err = loader.JSON(fmt.Sprintf(playCheckURL, ep.Id), &checkObj, gopts)
	if err != nil {
		return
	}

	playlistURL := checkObj.PlaylistURL

	tsaudio, err := common.LoadPlaylist(playlistURL, gopts, loader, hls)
	if err != nil {
		return
	}

	if len(tsaudio) > 1 {
		err = fmt.Errorf("Not implemented: tsaudio size > 1")
		return
	}

	errcImg := make(chan error)
	go func(errc chan<- error) {
		errc <- func() error {
			imageUrl := ep.showRef.PcImageUrl
			if len(imageUrl) == 0 {
				return errors.New("ImageURL not found")
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
			return
		}

		filteredCount := len(keys) + len(audio)
		keys = common.FilterChunks(keys, hls)
		audio = common.FilterChunks(audio, hls)
		if count := filteredCount - (len(keys) + len(audio)); count > 0 {
			fmt.Printf("already loaded %d files: continue...\n", count)
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
					err := common.LoadChunk(tsaudioUrl, gopts, file, loader)
					if err != nil {
						loadError <- errors.Wrap(err, "hibiki.dl.file")
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
						validateError <- errors.Wrap(err, "hibiki.dl.validate")
					}
				}
			}
		}()

		defer fmt.Printf("\n")

		links := []model.File{}
		links = append(links, keys...)
		links = append(links, audio...)
		status.Total(len(links))
		// pl.SetChunkCount(len(links))

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
		err = errImg
		return
	}

	return
}
