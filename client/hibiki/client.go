package hibiki

type Client struct {
	Headers map[string]string
	FeedURL string
	ShowURL string
}

func New() *Client {
	return &Client{
		Headers: map[string]string{
			"X-Requested-With": "XMLHttpRequest",
		},
		FeedURL: "https://vcms-api.hibiki-radio.jp/api/v1//programs?limit=99&page=1",
		ShowURL: "https://vcms-api.hibiki-radio.jp/api/v1/programs/%s",
	}
}
