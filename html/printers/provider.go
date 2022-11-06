package printers

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pgeowng/japoto-dl/html/config"
	"github.com/pgeowng/japoto-dl/html/types"
)

type Provider struct {
	Name string
}

func (p Provider) Print(entries []types.Entry) {
	files := []string{
		config.TemplateDir + "base.layout.tmpl",
		config.TemplateDir + "provider.content.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
		return
	}

	fpath := filepath.Join(config.Dest, fmt.Sprintf("%s.html", p.Name))
	f, err := os.Create(fpath)
	if err != nil {
		log.Fatalf("%s create error: %v\n", fpath, err)
	}
	defer f.Close()

	entries = FilterEntries(entries, func(entry types.Entry) bool {
		return entry.Provider == p.Name
	})
	entries = UniqueRecentShows(entries)

	alphabet := make(map[string][]types.Entry)
	for _, ep := range entries {
		programName := ep.ShowId
		if len(programName) == 0 {
			log.Fatalf("programName zero\n %v", ep)
		}

		letter := string(programName[0])
		letter = strings.ToLower(letter)
		if "0" <= letter && letter <= "9" {
			letter = "0"
		}
		alphabet[letter] = append(alphabet[letter], ep)
	}

	for _, eps := range alphabet {
		sort.Slice(eps, func(i, j int) bool {
			return strings.ToLower(eps[i].ShowId) < strings.ToLower(eps[j].ShowId)
		})
	}

	err = ts.Execute(f, map[string]interface{}{
		"PublicURL":  config.PublicURL,
		"CreateTime": config.CreateTime,
		"Alphabet":   &alphabet,
	})
	if err != nil {
		log.Fatalf("%s. write error: %v\n", fpath, err)
	}
}
