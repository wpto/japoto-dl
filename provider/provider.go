package provider

import (
	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/tasks"
)

type Provider interface {
	GetFeed() ([]model.Show, error)
	GetShow()
	Download(model.Show) error
}

type Providers struct {
	Onsen Provider
}

func NewProviders(tasks *tasks.Tasks) *Providers {
	return Providers{
		Onsen: NewOnsen(tasks),
	}
}
