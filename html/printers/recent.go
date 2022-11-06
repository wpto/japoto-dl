package printers

import (
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/pgeowng/japoto-dl/html/config"
	"github.com/pgeowng/japoto-dl/html/types"
)

type Recent struct{}

func (r Recent) Print(entries []types.Entry) {
	files := []string{
		config.TemplateDir + "base.layout.tmpl",
		config.TemplateDir + "recent.content.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
		return
	}

	f, err := os.Create(filepath.Join(config.Dest, "index.html"))
	if err != nil {
		log.Fatalln("index.html create error:", err)
	}

	defer f.Close()

	recent := make(map[string][]types.Entry)
	for _, ep := range entries {
		provider := ep.Provider
		recent[provider] = append(recent[provider], ep)
	}

	for provider, eps := range recent {
		recent[provider] = UniqueRecentShows(eps)

		currLimit := cap(recent[provider])
		if currLimit > config.RecentLimit {
			currLimit = config.RecentLimit
		}
		recent[provider] = recent[provider][:currLimit]
	}

	err = ts.Execute(f, map[string]interface{}{
		"PublicURL":  config.PublicURL,
		"CreateTime": config.CreateTime,
		"Recent":     &recent,
	})

	if err != nil {
		log.Fatalln("index.html write error:", err)
	}
}
