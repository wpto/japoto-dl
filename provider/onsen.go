package provider

import (
	"log"

	"github.com/levigross/grequests"
	"github.com/pgeowng/japoto-dl/model"
	"github.com/pkg/errors"
)

type Onsen struct{}

func NewOnsen() *Onsen {
	return &Onsen{}
}

type FeedRawResponse []FeedRawShow
type FeedRawShow []struct {
	Id                int    `json:"id"`
	DirectoryName     string `json:"directory_name"`
	Display           bool   `json:"display"`
	ShowContentsCount int    `json:"show_contents_count"`
	BrandNew          bool   `json:"brand_new"`
	BrandNewSp        bool   `json:"brand_new_sp"`
	Title             string `json:"title"`
	Image             struct {
		Url string `json:"url"`
	} `json:"image"`
	New               bool     `json:"new"`
	List              bool     `json:"list"`
	DeliveryInterval  string   `json:"delivery_interval"`
	DeliveryDayOfWeek []int    `json:"delivery_day_of_week"`
	CategoryList      []string `json:"category_list"`
	Copyright         string   `json:"copyright"`
	SponsorName       string   `json:"sponsor_name"`
	Updated           string   `json:"updated"`
	Performers        []struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		AllowLike bool   `json:"allow_like"`
	} `json:"performers"`
	RelatedLinks []struct {
		LinkUrl string `json:"link_url"`
		Image   string `json:"image"`
	} `json:"related_links"`

	RelatedInfos      []map[string]interface{} `json:"related_infos"`        // !
	RelatedPrograms   []map[string]interface{} `json:"related_programs"`     // !
	GuestInNewContent int                      `json:"guest_in_new_content"` // !
	Guests            int                      `json:"guest_in_new_content"` // !

	Contents []FeedRawContent `json:"contents"`
}

type FeedRawContent struct {
	Id             int     `json:"id"`
	Title          string  `json:"title"`
	Latest         bool    `json:"latest"`
	MediaType      string  `json:"media_type"`
	ProgramId      int     `json:"program_id"`
	New            bool    `json:"new"`
	Event          bool    `json:"event"`
	Block          bool    `json:"block"`
	OngenId        int     `json:"ongen_id"`
	Premium        bool    `json:"premium"`
	Free           bool    `json:"free"`
	DeliveryDate   string  `json:"delivery_date"`
	Movie          bool    `json:"movie"`
	PosterImageUrl string  `json:"poster_image_url"`
	StreamingUrl   *string `json:"streaming_url"`

	TagImage struct {
		Url *string `json:"url"`
	} `json:"tag_image"` // !

	Guests   int  `json:"guests"` //!
	Expiring bool `json:"expiring"`
}

func (f *FeedRawContent) GetEps() {}

func (p *Onsen) GetFeed() ([]model.Show, error) {
	res, err := grequests.Get("https://onsen.ag/web_api/programs/", nil)
	if err != nil {
		return nil, errors.Wrap(err, "onsen.feed.get")
	}

	apiRes := make(FeedRawResponse, 0)
	err = res.JSON(apiRes)
	if err != nil {
		log.Println(resp.String())
		return nil, errors.Wrap(err, "onsen.feed.parse")
	}

	return apiRes, nil
}
