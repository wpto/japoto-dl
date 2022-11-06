package printers

import (
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/pgeowng/japoto-dl/html/config"
	"github.com/pgeowng/japoto-dl/html/types"
)

type BasePage struct {
	PublicURL string
	CreatedAt time.Time
}

type ShowPage struct {
	Provider string
	Show     string
	Entries  []types.Entry
}

type Show struct {
	Provider string
	Name     string
}

func (show Show) Print(entries []types.Entry) {

	files := []string{
		config.TemplateDir + "base.layout.tmpl",
		config.TemplateDir + "show.content.tmpl",
	}

	err := os.MkdirAll(filepath.Join(config.Dest, show.Provider), fs.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
	}

	fpath := filepath.Join(config.Dest, show.Provider, show.Name+".html")
	f, err := os.Create(fpath)
	if err != nil {
		log.Fatalf("%s create error: %v\n", fpath, err)
	}
	defer f.Close()

	entries = FilterEntries(entries, func(entry types.Entry) bool {
		return entry.ShowId == show.Name && entry.Provider == show.Provider
	})

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Title > entries[j].Title
	})

	err = ts.Execute(f, map[string]interface{}{
		"PublicURL":  config.PublicURL,
		"CreateTime": config.CreateTime,
		"Provider":   show.Provider,
		"ShowName":   show.Name,
		"Entries":    entries,
	})

	if err != nil {
		log.Fatalf("%s write error: %v\n", fpath, err)
	}
}
