package cmd

import (
	"context"
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

type ShowMapper struct {
	dl        model.Loader
	providers []provider.Provider
	pl        model.PrintLine

	epChan chan model.Episode
}

func NewShowMapper(dl model.Loader, providers []provider.Provider, pl model.PrintLine) *ShowMapper {

	return &ShowMapper{
		dl:        dl,
		providers: providers,
		pl:        pl,
	}
}

func (m *ShowMapper) MapShows() (result []model.Show, err error) {
	result = make([]model.Show, 0, 10)

	ch := m.StreamShows(context.Background())

	for show := range ch {
		result = append(result, show)
	}

	return result, nil
}

func (m *ShowMapper) StreamShows(ctx context.Context) <-chan model.Show {
	ch := make(chan model.Show, 10)

	go func() {
		defer close(ch)
		err := m.loadShows(ctx, ch)
		if err != nil {
			fmt.Println("load shows error:", err)
		}
	}()

	return ch
}

func (m *ShowMapper) loadShows(ctx context.Context, out chan<- model.Show) error {
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
			return err
		}
		shows = FilterShowId(shows, FilterShowIdList)

		for _, showAccess := range shows {
			pl.SetPrefix(fmt.Sprintf("%s/%s", prov.Label(), showAccess.ShowId()))
			pl.Status("loading show")

			show, err := showAccess.GetShow(dl)
			if err != nil {
				pl.Error(errors.Errorf("showAccess: %v", err))
				return err
			}

			out <- show
		}
	}

	return nil
}

func (m *ShowMapper) GetEpisodeChan(ctx context.Context) <-chan model.Episode {
	if m.epChan != nil {
		log.Fatal("RunEpisodeWorker started twice")
		return nil
	}

	m.epChan = make(chan model.Episode, 20)
	return m.epChan
}

func (m *ShowMapper) RunEpisodeWorker(ctx context.Context) {
	defer close(m.epChan)

	err := m.loadEpisodes(ctx, m.epChan)
	if err != nil {
		fmt.Println("run episode worker error:", err)
	}
}

func (m *ShowMapper) loadEpisodes(ctx context.Context, output chan<- model.Episode) error {
	shows := m.StreamShows(ctx)
	for show := range shows {
		var eps []model.Episode
		eps, err := show.GetEpisodes(m.dl)
		if err != nil {
			fmt.Println("show.GetEpisodes: error: %v", err)
			return err
		}

		for _, ep := range eps {
			output <- ep
		}
	}

	return nil
}
