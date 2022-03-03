package cmd

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/pgeowng/japoto-dl/dl"
	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/provider"
	"github.com/pgeowng/japoto-dl/tasks"
	"github.com/pgeowng/japoto-dl/workdir"
	"github.com/pgeowng/japoto-dl/workdir/muxer"
	"github.com/pgeowng/japoto-dl/workdir/wd"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var ForceAudio bool
var OnlyMux bool
var FilterProviderList []string
var FilterShowIdList []string

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

	return cmd
}

var ffwg sync.WaitGroup
var history workdir.History = workdir.NewHistory("history.txt")

func FilterProvider(src []provider.Provider, filter []string) []provider.Provider {
	if len(filter) == 0 {
		return src
	}

	result := []provider.Provider{}
	for _, provider := range src {
		for _, label := range filter {
			if label == provider.Label() {
				result = append(result, provider)
				break
			}
		}
	}
	return result
}

func FilterShowId(src []model.ShowAccess, filter []string) []model.ShowAccess {
	if len(filter) == 0 {
		return src
	}

	result := []model.ShowAccess{}
	for _, sa := range src {
		for _, label := range filter {
			if label == sa.ShowId() {
				result = append(result, sa)
				break
			}
		}
	}
	return result
}

func downloadRun(cmd *cobra.Command, args []string) {

	d := dl.NewGrequests()
	providers := provider.NewProvidersList()

	pl := &PrintLine{}

	MapEpisode(d, providers, pl, func(ep model.Episode) error {
		pl.SetPrefix(ep.EpId())
		pl.Status("loading ep")
		if !ep.CanDownload() {
			return pl.Error(errors.New("cant load - skip"))
		}

		if ep.IsVideo() && !ForceAudio {
			return pl.Error(errors.New("saving video not implemented"))
		}

		date := ep.Date()

		if !date.IsGood() {
			return pl.Error(errors.Errorf("bad date - %s\n", date.String()))
		}

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

		salt := fmt.Sprintf("%s-%s--%s-u%d", date.Filename(), ep.Show().ShowId(), ep.Show().Provider(), ep.EpIdx())

		if history.Check(salt) {
			return pl.Error(errors.New("already downloaded"))
		}

		destPath := fmt.Sprintf("./%s.mp3", salt)

		// fmt.Printf("%s: loading\n", salt)

		ffm := muxer.NewFFMpegHLS(destPath, tags)
		wd1 := wd.NewWd("./.cache", salt)

		wdHLS := workdir.NewWorkdir(wd1, ffm, map[string]string{
			"playlist": "playlist.m3u8",
			"image":    "image",
		})

		t := tasks.NewTasks(wdHLS)
		if !OnlyMux {
			pl.AddLoadedCount()
			err := ep.Download(d, t, pl)
			pl.AddLoaded()
			if err != nil {
				return pl.Error(errors.Errorf("error - %s\n", err))
			}
		}

		ffwg.Add(1)
		pl.AddMuxedCount()
		go func() {
			var err error
			if OnlyMux {
				err = wdHLS.ForceMux()
			} else {
				err = wdHLS.Mux()
			}
			if err != nil {
				pl.Error(errors.Errorf("ffmpeg error: %v", err))
			} else {
				history.Write(salt)
				wd1.Clean()
			}
			ffwg.Done()
			pl.AddMuxed()
		}()

		return nil
	})

	fmt.Printf("loaded. waiting ffmpeg...\n")
	ffwg.Wait()
	fmt.Printf("muxed")
}

func MapEpisode(dl model.Loader, providers []provider.Provider, pl model.PrintLine, processEpisode func(ep model.Episode) error) {
	providers = FilterProvider(providers, FilterProviderList)

	for _, prov := range providers {
		pl.SetPrefix(prov.Label())
		pl.Status("loading feed")
		shows, err := prov.GetFeed(dl)
		if err != nil {
			pl.Error(errors.Errorf("err: %v", err))
			break
		}

		shows = FilterShowId(shows, FilterShowIdList)

		for _, showAccess := range shows {
			pl.SetPrefix(fmt.Sprintf("%s/%s", prov.Label(), showAccess.ShowId()))
			pl.Status("loading show")
			show, err := showAccess.GetShow(dl)
			if err != nil {
				pl.Error(errors.Errorf("showAccess: %v", err))
				break
			}

			eps := show.GetEpisodes()
			for _, ep := range eps {
				err := processEpisode(ep)
				if err != nil {
					break
				}
			}
		}
	}
}
