package provider

import (
	"github.com/pgeowng/japoto-dl/internal/provider/hibiki"
	"github.com/pgeowng/japoto-dl/internal/provider/onsen"
)

type Providers struct {
	Onsen  Provider
	Hibiki Provider
}

func NewProviders() *Providers {
	return &Providers{
		Onsen:  onsen.NewOnsen(),
		Hibiki: hibiki.NewHibiki(),
	}
}
