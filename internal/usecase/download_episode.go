package usecase

import (
	"fmt"
	"sort"
	"strings"

	"errors"

	"github.com/pgeowng/japoto-dl/internal/model"
	"github.com/pgeowng/japoto-dl/internal/provider/hibiki"
	"github.com/pgeowng/japoto-dl/internal/provider/onsen"
	"github.com/pgeowng/japoto-dl/internal/repo/archive"
	"github.com/pgeowng/japoto-dl/internal/repo/wd"
	"github.com/pgeowng/japoto-dl/internal/repo/workdir"
	"github.com/pgeowng/japoto-dl/internal/tasks"
)

var (
	ErrNoPlaylistURL = errors.New("no playlist url")
	ErrCannotLoad    = errors.New("cannot load")
)

type DownloadEpisode struct {
	metrics Metrics
}

func (d DownloadEpisode) Download(ep model.Episode) (err error) {
	d.d.metrics.SetPrefix(ep.EpId())
	d.metrics.Status("loading ep")
	if !ep.CanDownload() {
		return errors.New("cant load - skip")
	}

	if ep.IsVideo() && !ForceAudio {
		return d.metrics.Error(errors.New("saving video not implemented"))
	}

	date := ep.Date()

	artists := []string{}

	artists = append(artists, ep.Show().Artists()...)
	sort.Strings(artists)
	guests := ep.Artists()
	sort.Strings(guests)
	artists = append(artists, guests...)

	tags := map[string]string{
		"title":  strings.Join([]string{date.Filename(), ep.Show().ShowId(), ep.EpTitle(), ep.Show().ShowTitle()}, " "),
		"artist": strings.Join(artists, " "),
		"album":  ep.Show().ShowTitle(),
		"track":  date.Filename(),
	}

	salt := fmt.Sprintf("%s-%s--%s-u%s", date.Filename(), ep.Show().ShowId(), ep.Show().Provider(), ep.EpIdx())

	filename := fmt.Sprintf("%s.mp3", salt)
	playlistURL := ep.PlaylistURL()
	if playlistURL == nil {
		err = ErrNoPlaylistURL
		d.metrics.Error(err)
		return
	}

	/*
		u, err := url.Parse(*playlistURL)
		if err != nil {
			fmt.Println("playlisturl parse error", *playlistURL, err)
			return err
		}

			query := u.Query()
			query.Del("token")
			u.RawQuery = query.Encode()

			archiveKey := u.String()
	*/
	archiveKey := salt

	description := model.ArchiveItem{
		ArchiveKey: archiveKey,

		Description: &model.ArchiveItemDescription{
			Date:         date.Filename(),
			Source:       ep.Show().Provider(),
			ShowID:       ep.Show().ShowId(),
			ShowTitle:    ep.Show().ShowTitle(),
			EpisodeID:    ep.EpIdx(),
			EpisodeTitle: ep.EpTitle(),
			Artists:      artists,
		},

		Meta: &model.ArchiveItemMeta{
			Filename: filename,
		},
	}

	archiveStatus, err := archiveRepo.IsLoaded(archiveDB, archiveKey)
	if err != nil {
		return d.metrics.Error(fmt.Errorf("check archive repo(key): %w", err))
	}

	if archiveStatus == archive.Loaded {
		return d.metrics.Error(errors.New("already downloaded"))
	}

	if !BeforeDateCallback(date.Time()) {
		return d.metrics.Error(fmt.Errorf("skip before date: %s", date.String()))
	}

	if history.Check(salt) {
		err = archiveRepo.Create(archiveDB, salt, archive.Loaded, model.ArchiveItem{})
		if err != nil {
			err = archiveRepo.SetStatus(archiveDB, salt, archive.Loaded)
			if err != nil {
				return d.metrics.Error(fmt.Errorf("migrate file history error: %w", err))
			}
		}
		return d.metrics.Error(errors.New("already downloaded"))
	}
	fmt.Println("salt", archiveKey)

	if archiveStatus == archive.NotExists {
		err = archiveRepo.Create(archiveDB, archiveKey, archive.NotLoaded, description)
		if err != nil {
			return d.metrics.Error(fmt.Errorf("predownload: %w", err))
		}
	}

	destPath := fmt.Sprintf("./%s", filename)

	// fmt.Printf("%s: loading\n", salt)

	ffm := muxer.NewFFMpegHLS(destPath, tags)
	wd1 := wd.NewWd("./.cache", salt)

	wdHLS := workdir.NewWorkdir(wd1, ffm, map[string]string{
		"playlist": "playlist.m3u8",
		"image":    "image",
	})

	// statusRepo := d.metrics.NewLoadStatus(ep.Show().Provider(), ep.EpId())

	d.metrics.SetEp(ep.Show().Provider(), ep.EpIdx())

	t := tasks.NewTasks(wdHLS)
	if stageDownload && !OnlyMux {
		d.metrics.AddLoadedCount()
		var ii interface{} = ep
		var err error
		switch v := ii.(type) {
		case *hibiki.HibikiEpisodeMedia:
			hibikiUC := hibiki.HibikiUsecase{}
			err = hibikiUC.DownloadEpisode(d, t.AudioHLS(), status, v, wdHLS)
		case *onsen.OnsenEpisode:
			onsenUC := onsen.OnsenUsecase{}
			err = onsenUC.DownloadEpisode(d, t.AudioHLS(), status, v, wdHLS)
		default:
			fmt.Printf("Provider %s download is not implemented\n", ep.Show().Provider())
		}
		// err := ep.Download(d, t, status)
		d.metrics.AddLoaded()
		if err != nil {
			return d.metrics.Error(errors.Errorf("error - %s\n", err))
		}
	}

	ffwg.Add(1)
	d.metrics.AddMuxedCount()
	go func() {
		var err error
		if stageMux {
			if OnlyMux {
				err = wdHLS.ForceMux()
			} else {
				err = wdHLS.Mux()
			}
			if err != nil {
				d.metrics.Error(errors.Errorf("ffmpeg error: %v", err))
				return
			}
		}
		if stageHistoryWrite {
			if err = archiveRepo.SetStatus(archiveDB, archiveKey, archive.Loaded); err != nil {
				d.metrics.Error(errors.Errorf("archive write: %w", err))
				return
			}
		}
		if stageJSONWrite {
			if err == nil && LoadedJSON {
				err = description.SaveFile()
				if err != nil {
					err = fmt.Errorf("save description file: %s", err)
				}
			}
		}
		if err == nil {
			if stageOldHistoryWrite {
				history.Write(salt)
			}
			wd1.Clean()
		}

		ffwg.Done()
		d.metrics.AddMuxed()
	}()

	return nil
}
