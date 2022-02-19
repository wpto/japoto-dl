package main

import (
	"fmt"
	"log"

	"github.com/pgeowng/japoto-dl/dl"
	"github.com/pgeowng/japoto-dl/provider"
	"github.com/pgeowng/japoto-dl/tasks"
	"github.com/pgeowng/japoto-dl/workdir"
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
				wd := workdir.NewWorkdir("./.cache/", "fasqwge")
				tasks := tasks.NewTasks(wd)
				providers := provider.NewProviders(ddl, tasks)
				err := providers.Onsen.Download(*url)
				fmt.Println(err)
				return
			}
			fmt.Printf("%s ", ep.EpTitle())
		}
	}

	// log.Printf("%#v\n", result)

}
