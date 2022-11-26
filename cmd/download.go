package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/pgeowng/japoto-dl/internal/model"
	"github.com/pgeowng/japoto-dl/internal/provider/hibiki"
	"github.com/pgeowng/japoto-dl/internal/provider/onsen"
	"github.com/pgeowng/japoto-dl/internal/repo/archive"
	"github.com/pgeowng/japoto-dl/internal/repo/dl"
	"github.com/pgeowng/japoto-dl/internal/repo/status"
	"github.com/pgeowng/japoto-dl/internal/repo/wd"
	"github.com/pgeowng/japoto-dl/internal/tasks"
	"github.com/pgeowng/japoto-dl/internal/usecase"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var ForceAudio bool
var OnlyMux bool
var LoadedJSON bool
var BeforeDateStr string

const (
	stageDownload        = true
	stageMux             = true
	stageHistoryWrite    = true
	stageJSONWrite       = false
	stageOldHistoryWrite = false
)

func DownloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download eps",
		Long:  "Downloads available episodes",
		Run:   downloadRun,
	}

	cmd.Flags().BoolVarP(&ForceAudio, "force-audio", "a", true, "Forces video to remove video")
	cmd.Flags().BoolVarP(&OnlyMux, "only-mux", "m", false, "Skips anything but muxing")
	cmd.Flags().StringSliceVarP(&FilterProviderList, "provider-only", "p", []string{}, "Shows only selected providers")
	cmd.Flags().StringSliceVarP(&FilterShowIdList, "show-only", "s", []string{}, "Shows only selected shows")
	cmd.Flags().BoolVarP(&LoadedJSON, "export-json", "j", true, "On save exports json about loaded show")
	cmd.Flags().StringVar(&BeforeDateStr, "before-date", "", "Load only before specific date (excluded)")

	return cmd
}

var ffwg sync.WaitGroup
var history workdir.History = workdir.NewHistory("history.txt")

func downloadRun(cmd *cobra.Command, args []string) {

	BeforeDateCallback := func(date time.Time) bool {
		return true
	}

	if BeforeDateStr != "" {
		beforeDate, err := time.Parse("2006-01-02", BeforeDateStr)
		if err != nil {
			fmt.Printf("before-date parse error: %v", err)
			return
		}
		BeforeDateCallback = func(date time.Time) bool {
			return beforeDate.After(date)
		}
	}

	d := dl.NewGrequests()
	providers := []usecase.Provider{
		onsen.NewOnsen(),
		hibiki.NewHibiki(),
	}

	// status := &printline.PrintLine{}
	status := status.New(os.Stdout)

	archiveDB, err := archive.CreateDB("./japoto-archive.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	archiveRepo, err := archive.NewRepo()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = archiveRepo.Migrate(archiveDB)
	if err != nil {
		fmt.Println(err)
		return
	}

	MapEpisode(d, providers, status, func(ep model.Episode) error {
		status.SetPrefix(ep.EpId())
		status.Status("loading ep")
		if !ep.CanDownload() {
			return errors.New("cant load - skip")
		}

		if ep.IsVideo() && !ForceAudio {
			return status.Error(errors.New("saving video not implemented"))
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
			return status.Error(fmt.Errorf("PlaylistURL is nil: %w", err))
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
			return status.Error(fmt.Errorf("check archive repo(key): %w", err))
		}

		if archiveStatus == archive.Loaded {
			return status.Error(errors.New("already downloaded"))
		}

		if !BeforeDateCallback(date.Time()) {
			return status.Error(fmt.Errorf("skip before date: %s", date.String()))
		}

		if history.Check(salt) {
			err = archiveRepo.Create(archiveDB, salt, archive.Loaded, model.ArchiveItem{})
			if err != nil {
				err = archiveRepo.SetStatus(archiveDB, salt, archive.Loaded)
				if err != nil {
					return status.Error(fmt.Errorf("migrate file history error: %w", err))
				}
			}
			return status.Error(errors.New("already downloaded"))
		}
		fmt.Println("salt", archiveKey)

		if archiveStatus == archive.NotExists {
			err = archiveRepo.Create(archiveDB, archiveKey, archive.NotLoaded, description)
			if err != nil {
				return status.Error(fmt.Errorf("predownload: %w", err))
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

		// statusRepo := status.NewLoadStatus(ep.Show().Provider(), ep.EpId())

		status.SetEp(ep.Show().Provider(), ep.EpIdx())

		t := tasks.NewTasks(wdHLS)
		if stageDownload && !OnlyMux {
			status.AddLoadedCount()
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
			status.AddLoaded()
			if err != nil {
				return status.Error(errors.Errorf("error - %s\n", err))
			}
		}

		ffwg.Add(1)
		status.AddMuxedCount()
		go func() {
			var err error
			if stageMux {
				if OnlyMux {
					err = wdHLS.ForceMux()
				} else {
					err = wdHLS.Mux()
				}
				if err != nil {
					status.Error(errors.Errorf("ffmpeg error: %v", err))
					return
				}
			}
			if stageHistoryWrite {
				if err = archiveRepo.SetStatus(archiveDB, archiveKey, archive.Loaded); err != nil {
					status.Error(errors.Errorf("archive write: %w", err))
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
			status.AddMuxed()
		}()

		return nil
	})

	fmt.Printf("\nloaded. waiting ffmpeg...\n")
	ffwg.Wait()
	fmt.Println("done")

	folders, err := os.ReadDir("./.cache")
	if err == nil {
		if len(folders) > 0 {
			fmt.Println("there are some redundant .cache folders")
			fmt.Println("they are may never be used again.")
			for _, dir := range folders {
				fmt.Println(dir.Name())
			}
		}
	}
}
