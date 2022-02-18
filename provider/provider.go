package provider

import (
	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/provider/onsen"
	"github.com/pgeowng/japoto-dl/tasks"
	"github.com/pgeowng/japoto-dl/types"
)

type ProviderInfo interface {
	GetFeed() ([]model.Show, error)
	// GetShow()
}
type Provider interface {
	Download(playlistUrl string) error
}

type ProvidersInfo struct {
	Onsen ProviderInfo
}

type Providers struct {
	Onsen Provider
}

// func NewProviders(tasks *tasks.Tasks) *Providers {
func NewProviders(loader types.Loader, tasks *tasks.Tasks) *Providers {
	return &Providers{
		Onsen: onsen.NewOnsen(loader, tasks),
	}
}

func NewProvidersInfo(loader types.Loader) *ProvidersInfo {
	return &ProvidersInfo{
		Onsen: onsen.NewOnsenInfo(loader),
	}
}
