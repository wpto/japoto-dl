package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type ArchiveItem struct {
	ArchiveKey  string                  `json:"key"`
	Description *ArchiveItemDescription `json:"desc,omitempty"`
	Meta        *ArchiveItemMeta        `json:"meta,omitempty"`
	Chan        *ArchiveItemChan        `json:"chan,omitempty"`
}

type ArchiveItemDescription struct {
	Date         string   `json:"date"`
	Source       string   `json:"source"`
	ShowID       string   `json:"show_id"`
	ShowTitle    string   `json:"show_title"`
	EpisodeID    string   `json:"ep_id"`
	EpisodeTitle string   `json:"ep_title"`
	Artists      []string `json:"artists"`
}

type ArchiveItemMeta struct {
	Filename string `json:"filename"`
	Duration *int   `json:"duration,omitempty"`
	Size     *int   `json:"size,omitempty"`
}

type ArchiveItemChan struct {
	MessageId int `json:"msg_id,omitempty"`
}

func (l *ArchiveItem) filepath() string {
	return fmt.Sprintf("./%s.json", l.ArchiveKey)
}

func (l *ArchiveItem) IsExists() bool {
	file, err := os.Open(l.filepath())
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	file.Close()
	return true
}

func (l *ArchiveItem) SaveFile() error {
	jsonString, err := json.Marshal(l)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(l.filepath(), jsonString, os.ModePerm)
}
