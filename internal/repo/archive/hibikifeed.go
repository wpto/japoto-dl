package archive

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type FeedStorage interface {
	AddShows(ctx context.Context, entries []ShowEntry) error
}

type HibikiFeed struct {
	storage FeedStorage
	http    http.Client
}

func NewHibikiFeed(storage FeedStorage) *HibikiFeed {
	return &HibikiFeed{
		storage: storage,
		http:    http.Client{},
	}
}

func (f *HibikiFeed) Run(ctx context.Context) (err error) {
	shows, err := f.GetShowFeed(ctx)
	if err != nil {
		log.Printf("HibikiFeed: Run: GetShowFeed: %v", err)
		return
	}

	err = f.storage.AddShows(ctx, shows)
	if err != nil {
		log.Printf("HibikiFeed: Run: AddShows: %v", err)
		return
	}

	return
}

func (f *HibikiFeed) GetShowFeed(ctx context.Context) ([]ShowEntry, error) {
	req, err := http.NewRequest("GET", "https://vcms-api.hibiki-radio.jp/api/v1//programs?limit=99&page=1", nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %w", err)
	}

	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	res, err := f.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request: %w", err)
	}
	defer res.Body.Close()

	var result []struct {
		AccessID string `json:"access_id"`
	}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode response: %w", err)
	}

	shows := make([]ShowEntry, len(result))
	for i, v := range result {
		shows[i] = ShowEntry{
			Source: "hibiki",
			ShowID: v.AccessID,
		}
	}

	return shows, nil
}
