package types

type LoaderOpts struct {
	Headers map[string]string
}

type Loader interface {
	Text(url string, opts *LoaderOpts) (*string, error)
	JSON(url string, dest interface{}, opts *LoaderOpts) error
	Raw(url string, opts *LoaderOpts) ([]byte, error)
}
