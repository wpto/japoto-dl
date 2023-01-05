package archive

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type OnsenFeed struct {
	storage FeedStorage
	http    http.Client
}

func NewOnsenFeed(storage FeedStorage) *OnsenFeed {
	return &OnsenFeed{
		storage: storage,
		http:    http.Client{},
	}
}

func (f *OnsenFeed) Run(ctx context.Context) (err error) {
	shows, err := f.GetShowFeed(ctx)
	if err != nil {
		log.Printf("OnsenFeed: Run: GetShowFeed: %v", err)
		return
	}

	err = f.storage.AddShows(ctx, shows)
	if err != nil {
		log.Printf("OnsenFeed: Run: AddShows: %v", err)
		return
	}

	return nil
}

func (f *OnsenFeed) GetShowFeed(ctx context.Context) ([]ShowEntry, error) {
	req, err := http.NewRequest("GET", "https://onsen.ag/web_api/programs/", nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create feed")
	}

	res, err := f.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request: %w")
	}
	defer res.Body.Close()

	var result []struct {
		DirectoryName string `json:"directory_name"`
	}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode response: %w", err)
	}

	shows := make([]ShowEntry, len(result))
	for i, v := range result {
		shows[i] = ShowEntry{
			Source: "onsen",
			ShowID: v.DirectoryName,
		}
	}

	return shows, nil
}
