package html

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pgeowng/japoto-dl/html/config"
	"github.com/pgeowng/japoto-dl/html/entity"
	"github.com/pgeowng/japoto-dl/html/expanddb"
	"github.com/pgeowng/japoto-dl/html/printers"
	"github.com/pgeowng/japoto-dl/html/types"
	"github.com/pgeowng/japoto-dl/repo/archive"
	"github.com/spf13/cobra"
)

const destinationPath = "./public"

func run1() (err error) {

	return
}

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "html",
		Short: "Generates html pages",
		Long:  `read database and output html pages`,
		Run:   run,
	}

	cmd.Flags().StringVarP(&config.PublicURL, "public-url", "p", "", "Prefix for urls")

	return cmd
}

func run(cmd *cobra.Command, args []string) {

	a, err := archive.NewRepo()
	if err != nil {
		return
	}

	pool, err := archive.CreateDB("japoto-archive.db")
	if err != nil {
		return
	}

	defer pool.Close()

	err = a.Migrate(pool)
	if err != nil {
		return
	}

	items, err := a.GetArchiveEntries(pool)
	if err != nil {
		return
	}

	entries := make([]types.Entry, 0, len(items))
	for _, item := range items {
		if item.ArchiveKey == "" {
			continue
		}

		var (
			filename string
			duration int
			size     int
		)

		if item.Meta != nil {
			filename = item.Meta.Filename
			if item.Meta.Duration != nil {
				duration = *item.Meta.Duration
			}
			if item.Meta.Size != nil {
				size = *item.Meta.Size
			}
		}

		var (
			date      string
			provider  string
			showID    string
			showTitle string
			performer string
		)

		if item.Description != nil {
			date = item.Description.Date
			provider = item.Description.Source
			showID = item.Description.ShowID
			showTitle = item.Description.ShowTitle
			performer = strings.Join(item.Description.Artists, ", ")
		}

		var (
			messageID  int
			messageURL string
		)

		if item.Chan != nil {
			messageID = item.Chan.MessageId
			messageURL = config.ChannelPrefix + strconv.Itoa(item.Chan.MessageId)
		}

		entry := types.Entry{
			HasImage: false,

			MessageId: messageID,
			URL:       messageURL,

			Date:      date,
			Provider:  provider,
			ShowId:    showID,
			Title:     showTitle,
			Performer: performer,

			Filename:      filename,
			Duration:      duration,
			DurationHuman: entity.FormatDurationHuman(duration),
			Size:          size,
			SizeHuman:     entity.FormatSizeHuman(size),
		}

		entries = append(entries, entry)
	}

	//store := store.NewFileStore(config.FileStorePath)
	//entries := store.Read()

	entries = expanddb.ExtendContent(entries)

	err = os.MkdirAll(config.Dest, fs.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	staticFiles, err := os.ReadDir(config.Static)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range staticFiles {
		fp := filepath.Join(config.Static, file.Name())
		ff, err := ioutil.ReadFile(fp)
		if err != nil {
			fmt.Println(err)
			continue
		}

		dp := filepath.Join(config.Dest, file.Name())
		err = ioutil.WriteFile(dp, ff, 0644)
		if err != nil {
			fmt.Println("Error creating", dp)
			fmt.Println(err)
			return
		}
	}

	arranged := make(map[string]map[string]bool)
	for _, ep := range entries {
		if _, ok := arranged[ep.Provider]; !ok {
			arranged[ep.Provider] = make(map[string]bool)
		}
		arranged[ep.Provider][ep.ShowId] = true
	}

	r := printers.Recent{}
	r.Print(entries)

	for provider := range arranged {
		pr := printers.Provider{Name: provider}
		pr.Print(entries)
	}

	for provider := range arranged {
		for name := range arranged[provider] {
			sh := printers.Show{
				Provider: provider,
				Name:     name,
			}
			sh.Print(entries)
		}
	}

	// renderIndex(db)
	// renderAll(db)

	// // sc := presenters.ShowContent(s)
	// // presenters.RenderShowContent(sc)
	// for provider := range db {
	// 	for showName := range db[provider] {
	// 		renderPage(provider, showName, db[provider][showName])
	// 	}
	// }

}
