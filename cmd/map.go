package cmd

import (
	"fmt"

	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/provider"
	"github.com/pkg/errors"
)

var FilterProviderList []string
var FilterShowIdList []string

func FilterProvider(src []provider.Provider, filter []string) []provider.Provider {
	if len(filter) == 0 {
		return src
	}

	result := []provider.Provider{}
	for _, provider := range src {
		for _, label := range filter {
			if label == provider.Label() {
				result = append(result, provider)
				break
			}
		}
	}
	return result
}

func FilterShowId(src []model.ShowAccess, filter []string) []model.ShowAccess {
	if len(filter) == 0 {
		return src
	}

	result := []model.ShowAccess{}
	for _, sa := range src {
		for _, label := range filter {
			if label == sa.ShowId() {
				result = append(result, sa)
				break
			}
		}
	}
	return result
}
func (m *ShowMapper) MapEpisodes(processEpisode func(ep model.Episode) error) {
	dl := m.dl

	m.MapShows(func(show model.Show) error {
		eps, err := show.GetEpisodes(dl)
		if err != nil {
			fmt.Println("show.GetEpisodes: error: %v", err)
			return err
		}
		for _, ep := range eps {
			processEpisode(ep)
		}
		return nil
	})
}

type ShowMapper struct {
	dl        model.Loader
	providers []provider.Provider
	pl        model.PrintLine
}

func (m *ShowMapper) MapShows(processShow func(ep model.Show) error) {
	dl := m.dl
	providers := m.providers
	pl := m.pl

	providers = FilterProvider(providers, FilterProviderList)
	for _, prov := range providers {
		pl.SetPrefix(prov.Label())
		pl.Status("loading feed")
		shows, err := prov.GetFeed(dl)
		if err != nil {
			pl.Error(errors.Errorf("err: %v", err))
			break
		}
		shows = FilterShowId(shows, FilterShowIdList)

		for _, showAccess := range shows {
			pl.SetPrefix(fmt.Sprintf("%s/%s", prov.Label(), showAccess.ShowId()))
			pl.Status("loading show")
			show, err := showAccess.GetShow(dl)
			if err != nil {
				pl.Error(errors.Errorf("showAccess: %v", err))
				break
			}

			processShow(show)
		}
	}
}
