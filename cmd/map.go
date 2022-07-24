package cmd

import (
	"fmt"

	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/provider"
	"github.com/pgeowng/japoto-dl/provider/hibiki"
	"github.com/pgeowng/japoto-dl/provider/onsen"
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
	// providers = FilterProvider(providers, FilterProviderList)
	if len(FilterProviderList) == 0 {
		FilterProviderList = []string{"onsen", "hibiki"}
	}

	onlyLoadShowID := map[string]struct{}{}
	for _, item := range FilterShowIdList {
		onlyLoadShowID[item] = struct{}{}
	}
	onlyLoadEnabled := len(onlyLoadShowID) > 0

	for _, provider := range FilterProviderList {
		pl.SetPrefix(provider)
		pl.Status("loading feed")

		switch provider {
		case "onsen":
			uc := onsen.OnsenUsecase{}
			shows, errors := uc.GetFeed(dl)
			go func() {
				for show := range shows {
					if onlyLoadEnabled {
						_, ok := onlyLoadShowID[show.ShowId()]
						if !ok {
							continue
						}
					}

				}
			}()
			if err != nil {
				pl.Error(errors.Errorf("err: %v", err))
				break
			}
			for _, show := range shows {
				processShow(show)
			}
		case "hibiki":
			uc := hibiki.HibikiUsecase{}
			shows, err := uc.GetFeed(dl)
			if err != nil {
				pl.Error(errors.Errorf("err: %v", err))
				break
			}
			for _, show := range shows {
				if onlyLoadEnabled {
					_, ok := onlyLoadShowID[show.ShowId()]
					if !ok {
						continue
					}
				}
				pl.SetPrefix(fmt.Sprintf("%s/%s", provider, show.ShowId()))
				pl.Status("loading show")
				processShow(show)
			}
		default:
			fmt.Println("unsupported provider", provider)
			return
		}

		// shows, err := provider.GetFeed(dl)
		// if err != nil {
		// 	pl.Error(errors.Errorf("err: %v", err))
		// 	break
		// }
		// shows = FilterShowId(shows, FilterShowIdList)

		// for _, showAccess := range shows {
		// 	pl.SetPrefix(fmt.Sprintf("%s/%s", provider, showAccess.ShowId()))
		// 	pl.Status("loading show")
		// 	show, err := showAccess.GetShow(dl)
		// 	if err != nil {
		// 		pl.Error(errors.Errorf("showAccess: %v", err))
		// 		break
		// 	}

		// 	processShow(show)
		// }
	}
}
