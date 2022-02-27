package provider

import (
	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/provider/hibiki"
	"github.com/pgeowng/japoto-dl/provider/onsen"
)

type Provider interface {
	GetFeed(loader model.Loader) ([]model.ShowAccess, error)
}

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
