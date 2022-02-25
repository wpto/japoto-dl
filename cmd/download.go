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
	prov := provider.NewProviders(d)

	shows, err := prov.Onsen.GetFeed()
	if err != nil {
		log.Fatal(err)
		return
	}

	logger := log.Default()

	for _, show := range shows {
		eps := show.GetEpisodes()
		for _, ep := range eps {
			canLoad := ep.CanLoad()

			if !canLoad {
				logger.Printf("%s - cant load - skipping...", ep.ShowId())
				break
			}

			date, err := ep.Date()
			if err != nil {
				fmt.Println(err)
				fmt.Printf("%#v\n", ep)
				return
			}

			tags := map[string]string{
				"title":  strings.Join([]string{date.String(), ep.ShowId(), ep.EpTitle(), ep.ShowTitle()}, " "),
				"artist": strings.Join(ep.Artists(), " "),
				"album":  ep.ShowTitle(),
				"track":  date.String(),
			}

			ffm := muxer.NewFFMpegHLS("./output.mp3", tags)
			wd1 := wd.NewWd("./.cache", "salt")

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

// func LoadImage(dl model.Loader, wd workdir.Workdir, url string, fileName string) (string, error) {
// 	body, err := dl.Raw(url, nil)
// 	if err != nil {
// 		return "", errors.Wrap(err, "loadimage")
// 	}

// 	wd.SaveRaw(fileName, body)
// 	if err != nil {
// 		return "", errors.Wrap(err, "loadimage")
// 	}
// }
