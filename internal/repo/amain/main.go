package main

import (
	"context"
	"log"
	"time"

	"github.com/pgeowng/japoto-dl/internal/repo/archive"
)

func main() {
	db, err := archive.NewFileSQLite("data.db")
	if err != nil {
		log.Fatal(err)
	}

	storage := archive.NewStorage(db)
	if err := storage.MigrateShowFeed(context.Background()); err != nil {
		log.Fatal(err)
		return
	}

	feedHibiki := archive.NewHibikiFeed(storage)
	feedOnsen := archive.NewOnsenFeed(storage)

	for i := 0; i < 5; i++ {
		err := feedHibiki.Run(context.Background())
		if err != nil {
			log.Println(err)
		} else {
			log.Println("hibiki success!")
		}

		err = feedOnsen.Run(context.Background())
		if err != nil {
			log.Println(err)
		} else {
			log.Println("onsen success!")
		}

		entries, err := storage.GetShows(context.Background())
		if err != nil {
			log.Println(err)
		}

		for _, entry := range entries {
			log.Printf("%#v\n", entry)
		}

		<-time.After(time.Second)

	}
}
