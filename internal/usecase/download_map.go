package usecase

import (
	"fmt"

	"github.com/pgeowng/japoto-dl/internal/model"
	"github.com/pgeowng/japoto-dl/internal/provider"
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
func MapEpisode(dl model.Loader, providers []provider.Provider, pl model.PrintLine, processEpisode func(ep model.Episode) error) {
	MapShow(dl, providers, pl, func(show model.Show) error {
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

func (m *Mapper) MapShow(processShow func(show model.Show) error) {
	providers = FilterProvider(m.providers, FilterProviderList)
	for _, prov := range providers {
		m.pl.SetPrefix(prov.Label())
		m.pl.Status("loading feed")
		shows, err := prov.GetFeed(m.dl)
		if err != nil {
			err = fmt.Errorf("provider.GetFeed: error: %v", err)
			m.pl.Error(err)
			break
		}
		shows = FilterShowId(shows, FilterShowIdList)

		for _, showAccess := range shows {
			m.pl.SetPrefix(fmt.Sprintf("%s/%s", prov.Label(), showAccess.ShowId()))
			m.pl.Status("loading show")
			show, err := showAccess.GetShow(dl)
			if err != nil {
				err = fmt.Errorf("showAccess.GetShow: error: %v", err)
				m.pl.Error(err)
				break
			}

			processShow(show)
		}
	}
}
