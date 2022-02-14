package provider

import (
	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/provider/onsen"
)

type ProviderInfo interface {
	GetFeed() ([]model.Show, error)
	// GetShow()
}
type Provider interface {
	Download(model.Show) error
}

type Providers struct {
	OnsenInfo ProviderInfo
	Onsen     Provider
}

// func NewProviders(tasks *tasks.Tasks) *Providers {
func NewProviders() *Providers {
	return &Providers{
		OnsenInfo: onsen.NewOnsenInfo(),
		// Onsen:     onsen.NewOnsen(tasks),
	}
}
