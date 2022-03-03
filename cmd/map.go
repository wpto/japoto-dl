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
func MapEpisode(dl model.Loader, providers []provider.Provider, pl model.PrintLine, processEpisode func(ep model.Episode) error) {
	MapShow(dl, providers, pl, func(show model.Show) error {
		eps := show.GetEpisodes()
		for _, ep := range eps {
			processEpisode(ep)
		}
		return nil
	})
}

func MapShow(dl model.Loader, providers []provider.Provider, pl model.PrintLine, processShow func(show model.Show) error) {
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
