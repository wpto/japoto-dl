package provider

import "github.com/pgeowng/japoto-dl/model"

type Hibiki struct{}

var hibikiGopts *model.LoaderOpts = &model.LoaderOpts{
	Headers: map[string]string{
		"X-Requested-With": "XMLHttpRequest",
	},
}

func NewHibiki() *Hibiki {
	return &Hibiki{}
}

func (h *Hibiki) Label() string {
	return "hibiki"
}
