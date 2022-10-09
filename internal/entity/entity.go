package entity

import (
	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/workdir"
)

type EntityType int

const (
	FileEntity EntityType = iota
	// ShowEntity EntityType = iota
	// AudioChunkEntity
)

type Entity struct {
	Type EntityType

	Gopts   *model.LoaderOpts
	Loader  model.Loader
	Workdir workdir.WorkdirHLS

	ModelFile  *model.File
	TSAudioURL string

	URL  string
	Body []byte

	Filename string

	// // -- Show properies

	// // Name of provider from which this entity comes
	// Provider string

	// // Label within provider that identifies show
	// ShowID string

	// // Displayed full show title
	// ShowTitle string

	// // -- Episode properties

	// // Local date on which episode is released
	// EpisodeDate time.Time

	// // Whether episode can be downloaded freely
	// IsPremium bool

	// // Whether episode is a demo
	// IsPreview bool

	// // Whether there is some way to download episode
	// CanDownload bool

	// // Using this url we can download episode
	// EpisodeURL string

	// // -- HLS specific properties
	// HLSType HLSType

	// URL  string
	// Body []byte
}

// type HLSType int

// const (
// 	HLSPlaylist HLSType = iota
// 	HLSChunklist
// 	HLSKey
// 	HLSChunk
// )
