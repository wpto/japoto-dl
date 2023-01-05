package cmd

import (
	"fmt"
	"log"

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
func (m *ShowMapper) MapEpisodes() (result []model.Episode, err error) {
	shows, err := m.MapShows()
	if err != nil {
		log.Println("Map Episodes: map shows failed")
		return
	}

	for _, show := range shows {
		var eps []model.Episode
		eps, err := show.GetEpisodes(m.dl)
		if err != nil {
			fmt.Println("show.GetEpisodes: error: %v", err)
			return nil, err
		}

		result = append(result, eps...)
	}

	return
}

type ShowMapper struct {
	dl        model.Loader
	providers []provider.Provider
	pl        model.PrintLine
}

func (m *ShowMapper) MapShows() (result []model.Show, err error) {
	result = make([]model.Show, 0, 10)

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
			return nil, err
		}
		shows = FilterShowId(shows, FilterShowIdList)

		for _, showAccess := range shows {
			pl.SetPrefix(fmt.Sprintf("%s/%s", prov.Label(), showAccess.ShowId()))
			pl.Status("loading show")

			show, err := showAccess.GetShow(dl)
			if err != nil {
				pl.Error(errors.Errorf("showAccess: %v", err))
				return nil, err
			}

			result = append(result, show)
		}
	}

	return
}
