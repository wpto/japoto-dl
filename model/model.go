package model

type Episode interface {
	Artists() []string
	CanLoad() bool
	Date() (*Date, error)
	Download(loader Loader, tasks Tasks) error
	EpId() string
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
