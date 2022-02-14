package onsen

import (
	"reflect"

	"github.com/pgeowng/japoto-dl/model"
)

type Performer struct {
	Id        int    `mapstructure:"id"`
	Name      string `mapstructure:"name"`
	AllowLike bool   `mapstructure:"allow_like"`
}

type FeedRawShow struct {
	Id                int    `mapstructure:"id"`
	DirectoryName     string `mapstructure:"directory_name"`
	Display           bool   `mapstructure:"display"`
	ShowContentsCount int    `mapstructure:"show_contents_count"`
	BrandNew          bool   `mapstructure:"brand_new"`
	BrandNewSp        bool   `mapstructure:"brand_new_sp"`
	Title             string `mapstructure:"title"`

	Image struct {
		Url string `mapstructure:"url"`
	} `mapstructure:"image"`

	New               bool     `mapstructure:"new"`
	List              bool     `mapstructure:"list"`
	DeliveryInterval  string   `mapstructure:"delivery_interval"`
	DeliveryDayOfWeek []int    `mapstructure:"delivery_day_of_week"`
	CategoryList      []string `mapstructure:"category_list"`
	Copyright         string   `mapstructure:"copyright"`
	SponsorName       string   `mapstructure:"sponsor_name"`
	Updated           string   `mapstructure:"updated"`

	Performers []Performer `mapstructure:"performers"`

	RelatedLinks []struct {
		LinkUrl string `mapstructure:"link_url"`
		Image   string `mapstructure:"image"`
	} `mapstructure:"related_links"`

	RelatedInfos []struct {
		Category string `mapstructure:"category"`
		Caption  string `mapstructure:"caption"`
		LinkUrl  string `mapstructure:"link_url"`
		Image    string `mapstructure:"image"`
	} `mapstructure:"related_infos"`

	RelatedPrograms []struct {
		Title         string      `mapstructure:"title"`
		DirectoryName string      `mapstructure:"directory_name"`
		Category      string      `mapstructure:"category"`
		Image         string      `mapstructure:""`
		Performers    []Performer `mapstructure:"performers"`
	} `mapstructure:"related_programs"`

	GuestInNewContent []Performer `mapstructure:"guest_in_new_content"`
	Guests            []Performer `mapstructure:"guests"`
	Contents          []FeedRawEp `mapstructure:"contents"`
}

// !
type FeedRawEp struct {
	Id             int     `mapstructure:"id"`
	Title          string  `mapstructure:"title"`
	Latest         bool    `mapstructure:"latest"`
	MediaType      string  `mapstructure:"media_type"`
	ProgramId      int     `mapstructure:"program_id"`
	New            bool    `mapstructure:"new"`
	Event          bool    `mapstructure:"event"`
	Block          bool    `mapstructure:"block"`
	OngenId        int     `mapstructure:"ongen_id"`
	Premium        bool    `mapstructure:"premium"`
	Free           bool    `mapstructure:"free"`
	DeliveryDate   string  `mapstructure:"delivery_date"`
	Movie          bool    `mapstructure:"movie"`
	PosterImageUrl string  `mapstructure:"poster_image_url"`
	StreamingUrl   *string `mapstructure:"streaming_url"`

	TagImage struct {
		Url *string `mapstructure:"url"`
	} `mapstructure:"tag_image"`

	Guests   []Performer `mapstructure:"guests"`
	Expiring bool        `mapstructure:"expiring"`

	showRef *FeedRawShow
}

func (show *FeedRawShow) ShowId() string {
	return show.DirectoryName
}

func (show *FeedRawShow) GetEpisodes() []model.Episode {
	result := make([]model.Episode, 0)

	for i := range show.Contents {
		v := reflect.ValueOf(&show.Contents[i]).Interface()
		c := v.(model.Episode)
		result = append(result, c)
	}

	return result
}

func (ep *FeedRawEp) ShowId() string {
	return ep.showRef.ShowId()
}

func (ep *FeedRawEp) EpTitle() string {
	return ep.Title
}
