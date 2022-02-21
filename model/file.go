package model

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

type File struct {
	link string
	name string
	body []byte
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

	result := base.ResolveReference(link)
	return result.String(), nil
}

func (f *File) Name() string {
	return f.name
}

func (f *File) SetBody(data []byte) {
	fmt.Printf("%s: setraw\n", f.name)
	if len(f.body) != 0 {
		panic(errors.New("body was already set"))
	}
	f.body = data
}

func (f *File) BodyRaw() []byte {
	fmt.Printf("%s: getraw\n", f.name)
	if len(f.body) == 0 {
		panic(errors.Errorf("%s: body was not set", f.name))
	}
	return f.body
}

func (f *File) SetBodyString(data *string) {
	fmt.Printf("%s: setstr\n", f.name)
	if len(f.body) != 0 {
		panic(errors.New("body was already set"))
	}

	f.body = []byte(*data)
}

func (f *File) BodyString() string {
	fmt.Printf("%s: getstr\n", f.name)
	if len(f.body) == 0 {
		panic(errors.Errorf("%s: body was not set", f.name))
	}
	return string(f.body)
}
