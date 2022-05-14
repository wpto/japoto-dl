package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type LoadedModel struct {
	Uid string `json:"uid"`

	Basename string `json:"base_name"`
	Filename string `json:"file_name"`

	// Duration int
	// Size     int

	Date     string `json:"date"`
	Provider string `json:"provider"`
	ShowName string `json:"show_name"`

	ShowTitle    string `json:"show_title"`
	EpisodeTitle string `json:"ep_title"`

	Artists []string `json:"artists"`
}

func (l *LoadedModel) filepath() string {
	return fmt.Sprintf("./%s.json", l.Basename)
}

func (l *LoadedModel) IsExists() bool {
	file, err := os.Open(l.filepath())
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	file.Close()
	return true
}

func (l *LoadedModel) SaveFile() error {
	// file, err := os.OpenFile(l.filepath(), os.O_RDWR|os.O_CREATE, os.ModePerm)
	// if err != nil {
	// 	return err
	// }

	jsonString, err := json.Marshal(l)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(l.filepath(), jsonString, os.ModePerm)
	// encoder := json.NewEncoder(file)
	// defer file.Close()
	// if err := encoder.Encode(l); err != nil {
	// 	return err
	// }
}
