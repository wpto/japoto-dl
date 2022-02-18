package types

import "github.com/levigross/grequests"

type Loader interface {
	Text(url string, opts *grequests.RequestOptions) (*string, error)
	JSON(url string, dest interface{}, opts *grequests.RequestOptions) error
	Raw(url string, opts *grequests.RequestOptions) ([]byte, error)
}
