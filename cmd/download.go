package cmd

import (
	"fmt"
	"strings"
	"sync"

	"github.com/pgeowng/japoto-dl/dl"
	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/provider"
	"github.com/pgeowng/japoto-dl/tasks"
	"github.com/pgeowng/japoto-dl/workdir"
	"github.com/pgeowng/japoto-dl/workdir/muxer"
	"github.com/pgeowng/japoto-dl/workdir/wd"
	"github.com/spf13/cobra"
)

var ForceAudio bool
var OnlyMux bool
var FilterProvider []string
var FilterShowId []string

func DownloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download eps",
		Long:  "Downloads available episodes",
		Run:   downloadRun,
	}

	cmd.Flags().BoolVarP(&ForceAudio, "force-audio", "a", true, "Forces video to remove video")
	cmd.Flags().BoolVarP(&OnlyMux, "only-mux", "m", false, "Skips anything but muxing")
	cmd.Flags().StringSliceVarP(&FilterProvider, "provider-only", "p", []string{}, "Shows only selected providers")
	cmd.Flags().StringSliceVarP(&FilterShowId, "show-only", "s", []string{}, "Shows only selected shows")

	return cmd
}

var ffwg sync.WaitGroup
var history workdir.History = workdir.NewHistory("history.txt")

func downloadRun(cmd *cobra.Command, args []string) {
	d := dl.NewGrequests()
	providers := provider.NewProvidersList()

	for _, provider := range providers {
		pass := len(FilterProvider) == 0
		for _, label := range FilterProvider {
			if label == provider.Label() {
				pass = true
				break
			}
		}

		if !pass {
			fmt.Println("skipping " + provider.Label())
			continue
		}

		if err := loadProvider(d, provider); err != nil {
			fmt.Printf("err: %v", err)
		}
	}

	fmt.Printf("loaded. waiting ffmpeg...\n")
	ffwg.Wait()
	fmt.Printf("muxed")
}

func loadProvider(d model.Loader, p provider.Provider) error {
	fmt.Printf("%s: %s\n", p.Label(), "getting feed")
	shows, err := p.GetFeed(d)
	if err != nil {
		return err
	}

	for _, showAccess := range shows {
		pass := len(FilterShowId) == 0
		for _, showId := range FilterShowId {
			if showId == showAccess.ShowId() {
				pass = true
				break
			}
		}
		if !pass {
			continue
		}

		fmt.Printf("%s/%s: %s\n", p.Label(), showAccess.ShowId(), "loading show")
		show, err := showAccess.GetShow(d)
		if err != nil {
			return err
		}

		err = loadEpisodes(d, show)
	}
	return nil
}

func loadEpisodes(d model.Loader, show model.Show) error {
	eps := show.GetEpisodes()
	for _, ep := range eps {
		canLoad := ep.CanDownload()

		if !canLoad {
			fmt.Printf("%s: %s\n", ep.EpId(), "cant load - skip")
			continue
		}

		if ep.IsVideo() && !ForceAudio {
			fmt.Printf("%s: saving video not implemented\n", ep.EpId())
			continue
		}

		date := ep.Date()

		if !date.IsGood() {
			fmt.Printf("%s: bad date - %s\n", ep.EpId(), date.String())
			continue
		}

		artists := []string{}

		artists = append(artists, ep.Show().Artists()...)
		artists = append(artists, ep.Artists()...)

		tags := map[string]string{
			"title":  strings.Join([]string{date.String(), ep.Show().ShowId(), ep.EpTitle(), ep.Show().ShowTitle()}, " "),
			"artist": strings.Join(ep.Artists(), " "),
			"album":  ep.Show().ShowTitle(),
			"track":  date.Filename(),
		}

		salt := fmt.Sprintf("%s-%s--%s-u%d", date.Filename(), ep.Show().ShowId(), ep.Show().Provider(), ep.EpIdx())

		if history.Check(salt) {
			fmt.Printf("%s: already downloaded\n", ep.EpId())
			continue
		}

		destPath := fmt.Sprintf("./%s.mp3", salt)

		fmt.Printf("%s: loading\n", salt)

		ffm := muxer.NewFFMpegHLS(destPath, tags)
		wd1 := wd.NewWd("./.cache", salt)

		wdHLS := workdir.NewWorkdir(wd1, ffm, map[string]string{
			"playlist": "playlist.m3u8",
			"image":    "image",
		})

		t := tasks.NewTasks(wdHLS)
		if !OnlyMux {
			err := ep.Download(d, t)
			if err != nil {
				fmt.Printf("%s: error - %s\n", ep.EpId(), err)
				break
			}
		}

		ffwg.Add(1)
		go func() {
			var err error
			if OnlyMux {
				err = wdHLS.ForceMux()
			} else {
				err = wdHLS.Mux()
			}
			if err != nil {
				fmt.Printf("ffmpeg error: %v", err)
			} else {
				history.Write(salt)
				wd1.Clean()
			}
			ffwg.Done()
		}()

		fmt.Printf("%s ", ep.EpTitle())
	}

	return nil
}
