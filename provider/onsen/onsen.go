package onsen

import "github.com/pgeowng/japoto-dl/model"

type Onsen struct {
	loader model.Loader
}

func NewOnsen(loader model.Loader) *Onsen {
	return &Onsen{loader}
}
