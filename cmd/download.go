package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/pgeowng/japoto-dl/dl"
	"github.com/pgeowng/japoto-dl/provider"
	"github.com/pgeowng/japoto-dl/tasks"
	"github.com/pgeowng/japoto-dl/workdir"
	"github.com/pgeowng/japoto-dl/workdir/muxer"
	"github.com/pgeowng/japoto-dl/workdir/wd"
	"github.com/spf13/cobra"
)

func DownloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download eps",
		Long:  "Downloads available episodes",
		Run:   downloadRun,
	}

	return cmd
}

func downloadRun(cmd *cobra.Command, args []string) {

	d := dl.NewGrequests()
	prov := provider.NewProviders()
	shows, err := prov.Hibiki.GetFeed(d)
	if err != nil {
		log.Fatal(err)
		return
	}

	logger := log.Default()

	for _, showAccess := range shows {
		show, err := showAccess.GetShow(d)
		if err != nil {
			fmt.Println(err)
			continue
		}
		eps := show.GetEpisodes()
		for _, ep := range eps {
			canLoad := ep.CanDownload()

			if !canLoad {
				logger.Printf("%s - cant load - skipping...", ep.EpId())
				break
			}

			date := ep.Date()

			if !date.IsGood() {
				fmt.Printf("err: bad date %s", date.String())
				break
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

			destPath := fmt.Sprintf("./%s-%s--%s.mp3", date.Filename(), ep.Show().ShowId(), ep.Show().Provider())

			ffm := muxer.NewFFMpegHLS(destPath, tags)
			wd1 := wd.NewWd("./.cache", "salt1")

			wdHLS := workdir.NewWorkdir(wd1, ffm, map[string]string{
				"playlist": "playlist.m3u8",
				"image":    "image",
			})

			t := tasks.NewTasks(wdHLS)
			err = ep.Download(d, t)
			if err != nil {
				logger.Printf("skipping... loading error: %v", err)
				break
			}

			// err = wdHLS.Mux()
			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }

			fmt.Printf("%s ", ep.EpTitle())
		}
	}
}
