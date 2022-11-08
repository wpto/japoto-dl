package hibiki

type FeedEntry struct {
	ShowID string `json:"access_id"`
}

type Show struct {
	AccessId string `json:"access_id"`
	Name     string `json:"name"`
	Casts    []struct {
		Name     string  `json:"name"`
		RollName *string `json:"roll_name"`
	} `json:"casts"`
	PcImageUrl string  `json:"pc_image_url"`
	SpImageUrl string  `json:"sp_image_url"`
	Episode    Episode `json:"episode"`
}

type Episode struct {
	Id   int    `json:"id"`
	Name string `json:"name"`

	UpdatedAt string `json:"updated_at"`

	Video           *EpisodeMedia `json:"video"`
	AdditionalVideo *EpisodeMedia `json:"additional_video"`
}

type EpisodeMedia struct {
	Id        int `json:"id"`
	MediaType int `json:"media_type"`
	URL       *string

	IsAdditional bool
}
