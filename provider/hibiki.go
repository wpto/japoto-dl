package provider

import "github.com/pgeowng/japoto-dl/model"

type Hibiki struct {
	loader model.Loader
}

var hibikiGopts *model.LoaderOpts = &model.LoaderOpts{
	Headers: map[string]string{
		"X-Requested-With": "XMLHttpRequest",
	},
}

var HibikiGopts = hibikiGopts

func NewHibiki(loader model.Loader) *Hibiki {
	return &Hibiki{
		loader: loader,
	}
}

func (h *Hibiki) Label() string {
	return "hibiki"
}
