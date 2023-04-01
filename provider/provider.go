package provider

import (
	"github.com/pgeowng/japoto-dl/model"
)

type Provider interface {
	GetFeed(loader model.Loader) ([]model.ShowAccess, error)
	Label() string
}

type Providers struct {
	Onsen  Provider
	Hibiki Provider
}

func NewProviders(loader model.Loader) *Providers {
	return &Providers{
		Onsen:  NewOnsen(loader),
		Hibiki: NewHibiki(loader),
	}
}

func NewProvidersList(loader model.Loader) []Provider {
	return []Provider{NewOnsen(loader), NewHibiki(loader)}
}
