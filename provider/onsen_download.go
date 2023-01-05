package provider

import (
	"fmt"
	"log"

	"github.com/pgeowng/japoto-dl/internal/entity"
	"github.com/pgeowng/japoto-dl/internal/types"
	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/pkg/worker"
	"github.com/pgeowng/japoto-dl/workdir"
	"github.com/pkg/errors"
)

var onsenGopts *model.LoaderOpts = &model.LoaderOpts{
	Headers: map[string]string{
		"Referer": "https://www.onsen.ag/",
	},
}

type OnsenUsecase struct{}

func (uc *OnsenUsecase) DownloadEpisode(loader types.Loader, hls types.AudioHLS, status types.LoadStatus, ep *OnsenEpisode, wd workdir.WorkdirHLS) (err error) {

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
	tsaudio, err := LoadPlaylist(playlistURL, onsenGopts, loader, hls)
	if err != nil {
		err = errors.Wrap(err, "onsen")
		return
	}

	if len(tsaudio) > 1 {
		err = errors.New("onsen.dl: tsaudio size > 1: not implemented")
		return
	}

	load := worker.NewBoundWorker(6)

	// Image loading
	{
		url := ep.PosterImageUrl
		if url == "" {
			url = ep.showRef.Image.Url
		}

		if url == "" {
			return errors.New("onsen.dl: image not found")
		}

		image := &entity.Entity{
			Type:    entity.FileEntity,
			Gopts:   onsenGopts,
			Loader:  loader,
			Workdir: wd,

			URL: ep.PosterImageUrl,
		}

		load.Input <- image.DownloadImage
		log.Default().Printf("Download image queued")
	}

	// Chunk loading
	for _, ts := range tsaudio {
		var keys []model.File
		var audio []model.File
		var tsaudioUrl string
		keys, audio, tsaudioUrl, err = LoadTSAudio(playlistURL, onsenGopts, ts, loader, hls)
		if err != nil {
			err = errors.Wrap(err, "onsen.dl.tsparse")
			return
		}

		filteredCount := len(keys) + len(audio)
		keys = FilterChunks(keys, hls)
		audio = FilterChunks(audio, hls)
		if count := filteredCount - (len(keys) + len(audio)); count > 0 {
			fmt.Printf("already loaded %d files: continue...\n", count)
		}

		defer fmt.Printf("\n")

		links := []model.File{}
		links = append(links, keys...)
		links = append(links, audio...)
		status.Total(len(links))

		for idx := range links {
			file := &entity.Entity{
				Type:    entity.FileEntity,
				Gopts:   onsenGopts,
				Loader:  loader,
				Workdir: wd,

				ModelFile:  &links[idx],
				TSAudioURL: tsaudioUrl,
				Filename:   links[idx].Name(),
			}

			do := func() (err error) {
				err = entity.DownloadFile(file)
				if err != nil {
					return fmt.Errorf("Load worker: %w", err)
				}

				status.Inc(1)
				return nil
			}

			select {
			case load.Input <- do:
			case <-load.Done():
				break
			}
		}
	}

	load.Close()
	load.Wait()
	if load.Err() != nil {
		err = load.Err()
		return
	}

	return
}
