package dl

import (
	"fmt"

	"github.com/levigross/grequests"
	"github.com/pgeowng/japoto-dl/types"
	"github.com/pkg/errors"
)

type dlGrequests struct{}

func NewGrequests() *dlGrequests {
	return &dlGrequests{}
}

// options should me separated from code...
func (dl *dlGrequests) Text(url string, opts *types.LoaderOpts) (*string, error) {
	gopts := parseOpts(opts)
	res, err := grequests.Get(url, gopts)
	if err != nil {
		return nil, errors.Wrap(err, "get")
	}
	fmt.Printf("%s: %v\n", url, res.StatusCode)
	if !res.Ok {
		return nil, errors.Errorf("get: bad status code(%d): %s", res.StatusCode, url)
	}
	str := res.String()
	return &str, nil
}

func (dl *dlGrequests) JSON(url string, dest interface{}, opts *types.LoaderOpts) error {
	gopts := parseOpts(opts)
	res, err := grequests.Get(url, gopts)
	if err != nil {
		return errors.Wrap(err, "get")
	}
	fmt.Printf("%s: %v\n", url, res.StatusCode)
	if !res.Ok {
		return errors.Errorf("get: bad status code(%d): %s", res.StatusCode, url)
	}

	// mapObj := make([]map[string]interface{}, 0)
	err = res.JSON(dest)
	if err != nil {
		return errors.Wrap(err, "get.json")
	}

	return nil
}

func (dl *dlGrequests) Raw(url string, opts *types.LoaderOpts) ([]byte, error) {
	gopts := parseOpts(opts)
	res, err := grequests.Get(url, gopts)
	if err != nil {
		return nil, errors.Wrap(err, "get")
	}
	fmt.Printf("%s: %v\n", url, res.StatusCode)
	if !res.Ok {
		return nil, errors.Errorf("get: bad status code(%d): %s", res.StatusCode, url)
	}

	return res.Bytes(), nil
}

func parseOpts(opts *types.LoaderOpts) *grequests.RequestOptions {
	if opts == nil {
		return nil
	}

	return &grequests.RequestOptions{Headers: opts.Headers}
}
