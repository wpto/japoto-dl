package model

import (
	"fmt"
)

type Episode interface {
	Artists() []string
	CanLoad() bool
	Download(dl Loader, tasks Tasks) error
	Date() (*Date, error)
	EpTitle() string
	PlaylistUrl() *string
	ShowId() string
	ShowTitle() string
}

type Show interface {
	GetEpisodes() []Episode
	Artists() []string
	ShowId() string
	ShowTitle() string
}

type TSAudio interface {
	Link(base string) string
	Name() string
}

type Date struct {
	Year  int
	Month int
	Day   int
}

func (d *Date) String() string {
	return fmt.Sprintf("%02d%02d%02d", d.Year%100, d.Month, d.Day)
}
