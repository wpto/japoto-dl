package hibiki

import (
	"fmt"

	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/provider/common"
	"github.com/pkg/errors"
)

var defaultHeaders = map[string]string{
	"X-Requested-With": "XMLHttpRequest",
}

type HibikiUsecase struct {
}

func NewHibikiUsecase() (*HibikiUsecase, error) {
	return &HibikiUsecase{}, nil
}

func (ep *HibikiEpisodeMedia) Download(loader model.Loader, tasks model.Tasks, pl model.PrintLine) error {
	pl.SetPrefix(fmt.Sprintf("%s/%s", ep.Show().Provider(), ep.EpId()))
	pl.SetChunk(0)
	hls := tasks.AudioHLS()

	var checkObj struct {
		PlaylistUrl string `json:"playlist_url"`
	}
	var err error
	err = loader.
		Url(fmt.Sprintf("https://vcms-api.hibiki-radio.jp/api/v1/videos/play_check?video_id=%d", ep.Id)).
		Headers(defaultHeaders).
		JSON(&checkObj)

	err := loader.JSON("https://vcms-api.hibiki-radio.jp/api/v1/videos/play_check?video_id="+fmt.Sprint(ep.Id), &checkObj, gopts)
	if err != nil {
		return errors.Wrap(err, "hibiki.dl.check")
	}

	playlistUrl := checkObj.PlaylistUrl

	err = loader.
		Url(playlistUrl).
		Headers(defaultHeaders).
		Transform(func(content string) (err error) {
			return nil
		}).
		Save(wd.Permanent("playlist.m3u8"))

	tsaudio, err := common.LoadPlaylist(playlistUrl, gopts, loader, hls)
	if err != nil {
		return errors.Wrap(err, "hibiki.dl.playlist")
	}

	if len(tsaudio) > 1 {
		return errors.New("hibiki.dl.playlist: tsaudio size > 1: not implemented")
	}

	errcImg := make(chan error)
	go func(errc chan<- error) {
		errc <- func() error {
			imageUrl := ep.showRef.PcImageUrl
			if len(imageUrl) == 0 {
				return errors.New("hibiki.dl.image: not found")
			}

			err := loader.
				Url(imageUrl).
				Headers(defaultHeaders).
				Save(wd.Permanent("image"))

			return common.LoadImage(imageUrl, gopts, loader, hls)
		}()
	}(errcImg)

	for _, ts := range tsaudio {
		err = loader.
			Url(tsaudio).
			Headers(defaultHeaders).
			Transform(func(content []byte) (err error) { return }).
			Save(wd.Permanent("tsaudio"))

		keys, audio, tsaudioUrl, err := common.LoadTSAudio(playlistUrl, gopts, ts, loader, hls)
		if err != nil {
			return errors.Wrap(err, "hibiki.dl.ts")
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
						validateError <- errors.Wrap(err, "hibiki.dl.validate")
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
				return errors.Wrap(err, "hibiki.dl")
			case err := <-validateError:
				close(done)
				return errors.Wrap(err, "hibiki.dl")
			}
		}
		close(loadChan)

		select {
		case err := <-loadError:
			close(done)
			return errors.Wrap(err, "hibiki.dl")
		case err := <-validateError:
			if err != nil {
				close(done)
				return errors.Wrap(err, "hibiki.validate")
			}
		}
	}

	errImg := <-errcImg
	if errImg != nil {
		return errors.Wrap(errImg, "hibiki.dl:")
	}

	return nil
}
