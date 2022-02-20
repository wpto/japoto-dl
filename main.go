package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/pgeowng/japoto-dl/dl"
	"github.com/pgeowng/japoto-dl/provider"
)

func main() {
	ddl := dl.NewGrequests()
	providersInfo := provider.NewProvidersInfo(ddl)

	shows, err := providersInfo.Onsen.GetFeed()
	if err != nil {
		log.Fatal(err)
		return
	}

	// fmt.Println("len", len(shows))
	for _, show := range shows {
		eps := show.GetEpisodes()
		for _, ep := range eps {
			url := ep.PlaylistUrl()
			if url != nil {
				// wd1 := wd.NewWd("./.cache/", "fasqwge")
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
				fmt.Printf("%#v\n", tags)
				// mux1 := muxer.NewMuxerHLS("./file.mp3", nil, tags)
				// workdirHLS := workdir.NewWorkdirHLSImpl(wd1, mux1, "playlist.m3u8")
				// tasks := tasks.NewTasks(workdirHLS)
				// providers := provider.NewProviders(ddl, tasks)
				// err = providers.Onsen.Download(*url)
				// if err != nil {
				// 	fmt.Println(err)
				// 	return
				// }
				// err = workdirHLS.Mux()
				// if err != nil {
				// 	fmt.Println(err)
				// 	return
				// }
				// return
			}
			fmt.Printf("%s ", ep.EpTitle())
		}
	}

	// log.Printf("%#v\n", result)

}
