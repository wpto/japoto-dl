package model

type LoaderOpts struct {
	Headers  map[string]string
	Timeouts []int
}

var LoaderOptsDefault *LoaderOpts = &LoaderOpts{
	Headers:  map[string]string{},
	Timeouts: []int{5, 10, 20},
}

type Loader interface {
	Text(url string, opts *LoaderOpts) (*string, error)
	JSON(url string, dest interface{}, opts *LoaderOpts) error
	Raw(url string, opts *LoaderOpts) ([]byte, error)
}
