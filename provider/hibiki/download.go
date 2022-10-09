package hibiki

import (
	"fmt"

	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/pkg/worker"
	"github.com/pgeowng/japoto-dl/provider/common"
	"github.com/pgeowng/japoto-dl/workdir"

	"github.com/pgeowng/japoto-dl/internal/entity"
	"github.com/pgeowng/japoto-dl/internal/types"
)

const (
	playCheckURL = "https://vcms-api.hibiki-radio.jp/api/v1/videos/play_check?video_id=%d"
)

type HibikiUsecase struct{}

func (uc *HibikiUsecase) DownloadEpisode(loader types.Loader, hls types.AudioHLS, status types.LoadStatus, ep *HibikiEpisodeMedia, wd workdir.WorkdirHLS) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("HibikiUsecase.DownloadEpisode: %w", err)
		}
	}()

	playlistURL := *(ep.PlaylistURL())

	tsaudio, err := common.LoadPlaylist(playlistURL, gopts, loader, hls)
	if err != nil {
		return
	}

	if len(tsaudio) > 1 {
		err = fmt.Errorf("Not implemented: tsaudio size > 1")
		return
	}

	load := worker.NewBoundWorker(6)

	// Image loading
	{
		image := &entity.Entity{
			Type:    entity.FileEntity,
			Gopts:   gopts,
			Loader:  loader,
			Workdir: wd,

			URL: ep.showRef.PcImageUrl,
		}

		load.Input <- image.DownloadImage
	}

	// Chunk loading
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
		defer fmt.Printf("\n")

		links := append([]model.File{}, keys...)
		links = append(links, audio...)
		status.Total(len(links))
		// pl.SetChunkCount(len(links))

		for idx := range links {
			file := &entity.Entity{
				Type:    entity.FileEntity,
				Gopts:   gopts,
				Loader:  loader,
				Workdir: wd,

				ModelFile:  &links[idx],
				TSAudioURL: tsaudioUrl,

				Filename: links[idx].Name(),
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
