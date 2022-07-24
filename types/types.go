package types

type Loader interface {
	Url(url string) Loader
	Headers(headers map[string]string) Loader
	Transform(func(content []byte) (err error)) Loader

	JSON(output interface{}) error
	Save(filepath string) error
}
