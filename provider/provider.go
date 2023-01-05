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

func NewProviders() *Providers {
	return &Providers{
		Onsen:  NewOnsen(),
		Hibiki: NewHibiki(),
	}
}

func NewProvidersList() []Provider {
	return []Provider{NewOnsen(), NewHibiki()}
}
