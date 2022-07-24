package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/pgeowng/japoto-dl/dl"
	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/provider"
	"github.com/pgeowng/japoto-dl/provider/hibiki"
	"github.com/pgeowng/japoto-dl/provider/onsen"
	"github.com/pgeowng/japoto-dl/repo/archive"
	"github.com/pgeowng/japoto-dl/repo/status"
	"github.com/pgeowng/japoto-dl/tasks"
	"github.com/pgeowng/japoto-dl/workdir"
	"github.com/pgeowng/japoto-dl/workdir/muxer"
	"github.com/pgeowng/japoto-dl/workdir/wd"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var ForceAudio bool
var OnlyMux bool
var LoadedJSON bool

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

	return cmd
}

var ffwg sync.WaitGroup
var history workdir.History = workdir.NewHistory("history.txt")

func downloadRun(cmd *cobra.Command, args []string) {

	d := dl.NewGrequests()
	providers := provider.NewProvidersList()

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
		historyKey := salt
		description := LoadedModel{
			HistoryKey: historyKey,
			Basename:   salt,
			Filename:   filename,

			Provider: ep.Show().Provider(),
			Uid:      ep.EpIdx(),

			Date:     date.Filename(),
			ShowName: ep.Show().ShowId(),

			ShowTitle:    ep.Show().ShowTitle(),
			EpisodeTitle: ep.EpTitle(),

			Artists: artists,
		}

		archiveStatus, err := archiveRepo.IsLoaded(archiveDB, historyKey)
		if err != nil {
			return status.Error(fmt.Errorf("check archive repo: %w", err))
		}

		if loaded {
			return status.Error(errors.New("description already exists"))
		}

		if history.Check(salt) {
			return status.Error(errors.New("already downloaded"))
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
		if !OnlyMux {
			status.AddLoadedCount()
			var ii interface{} = ep
			var err error
			switch v := ii.(type) {
			case *hibiki.HibikiEpisodeMedia:
				hibikiUC := hibiki.HibikiUsecase{}
				err = hibikiUC.DownloadEpisode(d, t.AudioHLS(), status, v)
			case *onsen.OnsenEpisode:
				onsenUC := onsen.OnsenUsecase{}
				err = onsenUC.DownloadEpisode(d, t.AudioHLS(), status, v)
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
			if OnlyMux {
				err = wdHLS.ForceMux()
			} else {
				err = wdHLS.Mux()
			}
			if err == nil && LoadedJSON {
				err = description.SaveFile()
				if err != nil {
					err = fmt.Errorf("save description file: %s", err)
				}
			}
			if err != nil {
				status.Error(errors.Errorf("ffmpeg error: %v", err))
			} else {
				history.Write(salt)
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
