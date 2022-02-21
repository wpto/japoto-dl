package model

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"
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

type File struct {
	link    string
	name    string
	body    []byte
	bodyStr *string
}

func NewFile(link, name string) File {
	return File{link: link, name: name, body: []byte{}}
}

func (f *File) Url(baseUrl string) (string, error) {
	base, err := url.Parse(baseUrl)
	if err != nil {
		return "", errors.Wrap(err, "ts.url.base")
	}

	link, err := url.Parse(f.link)
	if err != nil {
		return "", errors.Wrap(err, "ts.url.link")
	}

	// fmt.Printf("file.url: %s, %s", f.link, baseUrl)

	result := base.ResolveReference(link)
	return result.String(), nil
}

func (f *File) Name() string {
	return f.name
}

func (f *File) SetBody(data []byte) {
	if f.bodyStr != nil {
		panic(errors.New("bodystr was already set"))
	}
	f.body = data
}

func (f *File) BodyRaw() []byte {
	if len(f.body) == 0 {
		panic(errors.New("body was not set"))
	}
	return f.body
}

func (f *File) SetBodyString(data *string) {
	if len(f.body) != 0 {
		panic(errors.New("bodyraw was already set"))
	}
	f.bodyStr = data
}

func (f *File) BodyString() string {
	if f.bodyStr == nil {
		panic(errors.New("bodystr was not set"))
	}
	return *f.bodyStr
}

type Date struct {
	Year  int
	Month int
	Day   int
}

func (d *Date) String() string {
	return fmt.Sprintf("%02d%02d%02d", d.Year%100, d.Month, d.Day)
}
