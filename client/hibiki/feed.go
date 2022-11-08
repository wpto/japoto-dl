package hibiki

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

func (c *Client) GetFeed() (feed []FeedEntry, err error) {

	client := resty.New()

	resp, err := client.R().
		SetHeaders(c.Headers).
		SetResult(&feed).
		Get(c.FeedURL)

	if err != nil {
		err = fmt.Errorf("hibiki.feed: %w", err)
		return
	}

	if resp.IsSuccess() {
		return
	}

	err = errors.New("hibiki.feed: failed to get feed")
	return
}

func (c *Client) GetShow(showID string) (show Show, err error) {

	client := resty.New()

	resp, err := client.R().
		SetHeaders(c.Headers).
		SetResult(&show).
		Get(fmt.Sprintf(c.ShowURL, showID))

	if err != nil {
		err = fmt.Errorf("hibiki.show: %w", err)
		return
	}

	if resp.IsSuccess() {
		return
	}

	err = errors.New("hibiki.show: failed to get show")
	return
}

func (c *Client) GetShowList() (showList []string, err error) {
	client := resty.New()
	var feed []FeedEntry

	// TODO: pagination
	resp, err := client.R().
		SetHeaders(c.Headers).
		SetResult(&feed).
		Get(c.FeedURL)

	if err != nil {
		err = fmt.Errorf("hibiki.showList: %w", err)
		return
	}

	if resp.IsSuccess() {
		showList = make([]string, len(feed))
		for i := range feed {
			showList[i] = feed[i].ShowID
		}
		return
	}

	err = errors.New("hibiki.showList: failed to get feed")
	return
}
