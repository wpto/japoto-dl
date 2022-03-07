package model

import (
	"fmt"
	"strings"
)

type Episode interface {
	Artists() []string
	CanDownload() bool
	Date() Date
	Download(loader Loader, tasks Tasks, pl PrintLine) error
	EpId() string
	EpIdx() string // base62
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

	LoadImage(loader Loader, workdir WorkdirBase) error

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

func (a PPrintRow) String() string {
	cast := ""
	if len(a.Cast) != 0 {
		cast = "# " + strings.Join(a.Cast, " ")
	}
	return fmt.Sprintf("%s %s %s %s %s", a.Attrs(), a.Date.String(), a.Ref, a.Note, cast)
}

type WorkdirBase interface {
	Save(fileName, fileBody string) error
	SaveRaw(fileName string, fileBody []byte) error
}

type WorkdirFile interface {
	SaveNamed(name string, fileBody string) error
	SaveNamedRaw(name string, fileBody []byte) error
	ResolveName(name string) string
	WasWritten(name string) bool
	WorkdirBase
}
