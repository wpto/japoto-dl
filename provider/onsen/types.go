package onsen

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"github.com/pgeowng/japoto-dl/model"
	"github.com/pkg/errors"
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
	DeliveryDate   *string `mapstructure:"delivery_date"`
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

func (show *FeedRawShow) GetEpisodes() []model.Episode {
	result := make([]model.Episode, 0)

	for i := range show.Contents {
		v := reflect.ValueOf(&show.Contents[i]).Interface()
		c := v.(model.Episode)
		result = append(result, c)
	}

	return result
}

func (show *FeedRawShow) Artists() []string {
	result := []string{}

	for _, p := range show.Performers {
		result = append(result, p.Name)
	}

	return result
}

func (show *FeedRawShow) ShowId() string {
	return show.DirectoryName
}

func (show *FeedRawShow) ShowTitle() string {
	return show.Title
}

func (ep *FeedRawEp) Artists() []string {
	result := []string{}

	result = append(result, ep.showRef.Artists()...)

	for _, p := range ep.Guests {
		result = append(result, p.Name)
	}

	return result
}

func (ep *FeedRawEp) CanLoad() bool {
	return ep.StreamingUrl != nil
}

func (ep *FeedRawEp) Date() (*model.Date, error) {
	if ep.StreamingUrl == nil {
		return nil, errors.New("you can't call ep.Date() on episode that is can't be loaded. onsen specific limitation")
	}

	result, err := Extract(*ep.StreamingUrl, ep.ShowId())

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("onsen.date: bad format for (%s, %s)", *ep.StreamingUrl, ep.ShowId()))
	}

	if result.DateD != 0 && result.DateM != 0 && result.DateY != 0 {
		return &model.Date{
			Year:  result.DateY,
			Month: result.DateM,
			Day:   result.DateD,
		}, nil
	}

	if result.DateY == 0 {
		return nil, errors.New("onsen.date: cant resolve without year")
	}

	if ep.DeliveryDate == nil {
		return nil, errors.New("onsen.date: undefined")
	}

	reText := `(\d+)/(\d+)`
	reDate := regexp.MustCompile(reText)

	match := reDate.FindStringSubmatch(*ep.DeliveryDate)
	if match == nil {
		return nil, errors.Errorf("onsen.date: bad delivery_date format: %s", *ep.DeliveryDate)
	}

	mm, err := strconv.ParseInt(match[1], 10, 0)
	if err != nil {
		return nil, errors.Wrap(err, "onsen.date.delivery_month")
	}

	if int(mm) != result.DateM {
		return nil, errors.Errorf("onsen.date: month mismatch %s and %s", *ep.DeliveryDate, *ep.StreamingUrl)
	}

	dd, err := strconv.ParseInt(match[2], 10, 0)
	if err != nil {
		return nil, errors.Wrap(err, "onsen.date.delivery_month")
	}

	return &model.Date{
		Year:  result.DateY,
		Month: result.DateM,
		Day:   int(dd),
	}, nil
}

func (ep *FeedRawEp) EpTitle() string {
	return ep.Title
}

func (ep *FeedRawEp) PlaylistUrl() *string {
	return ep.StreamingUrl
}

func (ep *FeedRawEp) ShowId() string {
	return ep.showRef.ShowId()
}

func (ep *FeedRawEp) ShowTitle() string {
	return ep.showRef.ShowTitle()
}
