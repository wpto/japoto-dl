package archive

import (
	"context"
	"fmt"

	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/provider/onsen"
)

type OnsenProvider struct {
	Storage
}

func NewOnsenProvider(storage Storage) *OnsenProvider {
	return &OnsenProvider{Storage: storage}
}

func (p *OnsenProvider) GetFeed(_ model.Loader) ([]model.ShowAccess, error) {

	shows, err := p.GetShows(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("OnsenProvider.GetFeed: %w", err)
	}

	result := make([]model.ShowAccess, 0, len(shows))
	for i := range shows {
		if shows[i].Source == "onsen" {
			result = append(result, &onsen.OnsenShowAccess{DirectoryName: shows[i].ShowID})
		}

		// TODO: dynamic typecast. remove
		var v interface{} = shows[i]
		c := v.(model.ShowAccess)
		result = append(result, c)
	}

	return result, nil
}

func (p *OnsenProvider) Label() string {
	return "onsen"
}
