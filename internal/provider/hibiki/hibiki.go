package hibiki

import "github.com/pgeowng/japoto-dl/internal/model"

type Hibiki struct{}

var gopts *model.LoaderOpts = &model.LoaderOpts{
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
