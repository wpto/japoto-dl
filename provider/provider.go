package provider

import (
	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/provider/onsen"
)

type Provider interface {
	GetFeed() ([]model.Show, error)
}

type Providers struct {
	Onsen Provider
}

func NewProviders(loader model.Loader) *Providers {
	return &Providers{
		Onsen: onsen.NewOnsen(loader),
	}
}
