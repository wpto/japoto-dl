package model

import "time"

type ListEntry struct {
	IsShow      bool
	CanDownload bool
	IsAudio     bool
	IsPaid      bool
	Date        time.Time
	ShowID      string
	EpisodeID   string
	Title       string
}
