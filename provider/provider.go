package provider

import (
	"github.com/pgeowng/japoto-dl/provider/hibiki"
	"github.com/pgeowng/japoto-dl/provider/onsen"
)

type Provider interface {
	// GetFeed(loader model.Loader) ([]model.ShowAccess, error)
	Label() string
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

func NewProvidersList() []Provider {
	return []Provider{onsen.NewOnsen(), hibiki.NewHibiki()}
}
