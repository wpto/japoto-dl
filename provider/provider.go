package provider

import (
	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/provider/onsen"
)

type Provider interface {
	GetFeed(loader model.Loader) ([]model.Show, error)
}

type Providers struct {
	Onsen Provider
}

func NewProviders() *Providers {
	return &Providers{
		Onsen: onsen.NewOnsen(),
	}
}
