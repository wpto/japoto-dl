package model

import "strings"

type Episode interface {
	Artists() []string
	CanDownload() bool
	Date() Date
	Download(loader Loader, tasks Tasks) error
	EpId() string
	EpIdx() int
	EpTitle() string
	// PlaylistUrl() *string
	// ShowId() string
	// ShowTitle() string
	IsVideo() bool

	Show() Show

	PPrint() PPrintRow
}

type Show interface {
	GetEpisodes() []Episode
	Artists() []string
	ShowId() string
	ShowTitle() string

	PPrint() PPrintRow
	Provider() string
}

type ShowAccess interface {
	GetShow(loader Loader) (Show, error)
	ShowId() string
}

type TSAudio interface {
	Link(base string) string
	Name() string
}

type PPrintRow struct {
	IsDir   bool
	CanLoad bool
	IsVid   bool
	Date    Date
	Ref     string
	Note    string
	Cast    []string
}

func (a PPrintRow) Attrs() string {
	result := []string{"-", "-", "-", "-"}

	if a.IsDir {
		result[0] = "d"
	} else {
		if a.IsVid {
			result[2] = "v"
		} else {
			result[2] = "a"
		}
	}

	if a.CanLoad {
		result[1] = "l"
	} else {
		result[3] = "$"
	}

	return strings.Join(result, "")
}

func (a PPrintRow) Pprint() []interface{} {
	cast := ""
	if len(a.Cast) != 0 {
		cast = "# " + strings.Join(a.Cast, " ")
	}

	return []interface{}{a.Attrs(), a.Date.String(), a.Ref, a.Note, cast}
}
