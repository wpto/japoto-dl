package dl

import (
	"time"

	"github.com/levigross/grequests"
	"github.com/pgeowng/japoto-dl/model"
	"github.com/pkg/errors"
)

type Grequests struct {
	session *grequests.Session
}

func NewGrequests() *Grequests {
	session := grequests.NewSession(nil)
	return &Grequests{session}
}

// options should me separated from code...
func (dl *Grequests) Text(url string, opts *model.LoaderOpts) (*string, error) {
	res, err := dl.retryGet(url, opts)
	if err != nil {
		return nil, errors.Wrap(err, "get")
	}
	if !res.Ok {
		return nil, errors.Errorf("get: bad status code(%d): %s", res.StatusCode, url)
	}
	str := res.String()
	return &str, nil
}

func (dl *Grequests) JSON(url string, dest interface{}, opts *model.LoaderOpts) error {
	res, err := dl.retryGet(url, opts)
	if err != nil {
		return errors.Wrap(err, "get")
	}

	if !res.Ok {
		return errors.Errorf("get: bad status code(%d): %s", res.StatusCode, url)
	}

	err = res.JSON(dest)
	if err != nil {
		return errors.Wrap(err, "get.json")
	}

	return nil
}

func (dl *Grequests) Raw(url string, opts *model.LoaderOpts) ([]byte, error) {
	res, err := dl.retryGet(url, opts)
	if err != nil {
		return nil, errors.Wrap(err, "get")
	}

	if !res.Ok {
		return nil, errors.Errorf("get: bad status code(%d): %s", res.StatusCode, url)
	}

	return res.Bytes(), nil
}

func (dl *Grequests) retryGet(url string, opts *model.LoaderOpts) (*grequests.Response, error) {
	gopts := &grequests.RequestOptions{}
	times := model.LoaderOptsDefault.Timeouts

	if opts != nil {
		if opts.Headers != nil {
			gopts.Headers = opts.Headers
		}
		if opts.Timeouts != nil {
			times = opts.Timeouts
		}
	}

	var res *grequests.Response
	var err error
	for _, t := range times {
		gopts.RequestTimeout = time.Duration(t) * time.Second
		res, err = dl.session.Get(url, gopts)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, errors.Wrap(err, "retry failed")
	}

	// fmt.Printf("%s: %v\n", url, res.StatusCode)

	return res, nil
}
