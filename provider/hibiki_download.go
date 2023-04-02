package provider

import (
	"context"
	"fmt"

	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/pkg/worker"
	"github.com/pgeowng/japoto-dl/workdir"

	"github.com/pgeowng/japoto-dl/internal/entity"
	"github.com/pgeowng/japoto-dl/internal/types"
	"github.com/pgeowng/japoto-dl/repo/status"
)

const (
	playCheckURL = "https://vcms-api.hibiki-radio.jp/api/v1/videos/play_check?video_id=%d"
)

type GeneralLoader struct {
	loader types.Loader
	hls    types.AudioHLS
	metric status.Metric
	wd     workdir.WorkdirHLS
}

func NewGeneralLoader(loader types.Loader, hls types.AudioHLS, metric status.Metric, wd workdir.WorkdirHLS) *GeneralLoader {
	return &GeneralLoader{
		loader: loader,
		hls:    hls,
		metric: metric,
		wd:     wd,
	}
}

type DownloadEpisodeParams struct {
	PlaylistURL    string
	ImageURL       string
	RequestOptions *model.LoaderOpts
}

func (uc *GeneralLoader) DownloadEpisode(ctx context.Context, ep DownloadEpisodeParams) (err error) {

	ropts := ep.RequestOptions

	loader := uc.loader
	hls := uc.hls
	metric := uc.metric
	wd := uc.wd

	defer func() {
		if err != nil {
			err = fmt.Errorf("GeneralLoader.DownloadEpisode: %w", err)
		}
	}()

	playlistURL := ep.PlaylistURL
	fmt.Println("Using playlist url:", playlistURL)
	tsaudio, err := LoadPlaylist(playlistURL, ropts, loader, hls)
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
			Gopts:   ropts,
			Loader:  loader,
			Workdir: wd,

			URL: ep.ImageURL,
		}

		load.Input <- image.DownloadImage
	}

	// Chunk loading
	for _, ts := range tsaudio {
		var keys []model.File
		var audio []model.File
		var tsaudioUrl string
		keys, audio, tsaudioUrl, err = LoadTSAudio(playlistURL, ropts, ts, loader, hls)
		if err != nil {
			return
		}

		filteredCount := len(keys) + len(audio)
		keys = FilterChunks(keys, hls)
		audio = FilterChunks(audio, hls)
		if count := filteredCount - (len(keys) + len(audio)); count > 0 {
			fmt.Printf("already loaded %d files: continue...\n", count)
		}

		// total := len(keys) + len(audio)
		defer fmt.Printf("\n")

		links := append([]model.File{}, keys...)
		links = append(links, audio...)

		metric.Set("total", float32(len(links)))
		// pl.SetChunkCount(len(links))

		for idx := range links {
			file := &entity.Entity{
				Type:    entity.FileEntity,
				Gopts:   ropts,
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

				metric.Inc("progress")
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
