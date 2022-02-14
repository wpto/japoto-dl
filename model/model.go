package model

type Episode interface {
	Title() string
	Date() string
}

type Show interface {
	GetEpisodes()
}
