package loader

import (
	"time"

	"github.com/levigross/grequests"
	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/types"
	"github.com/pkg/errors"
)

type grequestsLoader struct {
	session *grequests.Session

	url       string
	headers   map[string]string
	transform func(content []byte) (err error)
}

func NewGrequestsLoader(session *grequests.Session) types.Loader {
	return grequestsLoader{session: session}
}

func (l grequestsLoader) Url(url string) types.Loader {
	l.url = url
	return l
}

func (l grequestsLoader) Headers(headers map[string]string) types.Loader {
	l.headers = headers
	return l
}

func (l grequestsLoader) Transform(transform func(content []byte) (err error)) types.Loader {
	l.transform = transform
	return l
}

func (l grequestsLoader) retryGet(url string, opts *model.LoaderOpts) (*grequests.Response, error) {
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

func (l grequestsLoader) JSON(output interface{}) error {

}
