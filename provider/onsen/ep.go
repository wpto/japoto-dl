package onsen

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"github.com/pgeowng/japoto-dl/model"
	"github.com/pkg/errors"
)

// !
type OnsenEpisode struct {
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
	DeliveryDate   *string `json:"delivery_date"`
	Movie          bool    `json:"movie"`
	PosterImageUrl string  `json:"poster_image_url"`
	StreamingUrl   *string `json:"streaming_url"`

	TagImage struct {
		Url *string `json:"url"`
	} `json:"tag_image"`

	Guests   []Performer `json:"guests"`
	Expiring bool        `json:"expiring"`

	showRef *OnsenShow
}

func (ep *OnsenEpisode) Artists() []string {
	result := []string{}
	for _, p := range ep.Guests {
		if len(p.Name) > 0 {
			result = append(result, p.Name)
		}
	}

	return result
}

func (ep *OnsenEpisode) CanDownload() bool {
	return ep.StreamingUrl != nil
}

func (ep *OnsenEpisode) Date() model.Date {
	result := model.Date{Year: -1, Month: -1, Day: -1}

	if ep.StreamingUrl != nil {
		other, err := parseStreamingUrlDate(*ep.StreamingUrl, ep.ShowId())
		if err != nil {
			fmt.Println(errors.Wrap(err, "onsen.ep.date"))
		} else {
			if result.Year == -1 {
				result.Year = other.Year
			}
			if result.Month == -1 {
				result.Month = other.Month
			}
			if result.Day == -1 {
				result.Day = other.Day
			}
		}
	}

	if ep.DeliveryDate != nil {
		other, err := parseDeliveryDate(*ep.DeliveryDate)
		if err != nil {
			fmt.Println(errors.Wrap(err, "onsen.ep.date"))
		} else {
			if result.Year == -1 {
				result.Year = other.Year
			}
			if result.Month == -1 {
				result.Month = other.Month
			}
			if result.Day == -1 {
				result.Day = other.Day
			}
		}
	}

	return result
}

func parseDeliveryDate(date string) (*model.Date, error) {
	re := regexp.MustCompile(`(\d+)/(\d+)`)

	match := re.FindStringSubmatch(date)
	if match == nil {
		return nil, errors.Errorf("dont match pattern m/d: %s", date)
	}

	mm, _ := strconv.ParseInt(match[1], 10, 0)
	dd, _ := strconv.ParseInt(match[2], 10, 0)

	return &model.Date{
		Year:  -1,
		Month: int(mm),
		Day:   int(dd),
	}, nil
}

func parseStreamingUrlDate(url string, showId string) (*model.Date, error) {
	result, err := Extract(url, showId)
	if err != nil {
		return nil, err
	}
	return &model.Date{
		Year:  result.DateY,
		Month: result.DateM,
		Day:   result.DateD,
	}, nil
}

func (ep *OnsenEpisode) EpId() string {
	return fmt.Sprintf("%s/%d", ep.ShowId(), ep.Id)
}

func (ep *OnsenEpisode) EpTitle() string {
	return ep.Title
}

func (ep *OnsenEpisode) PlaylistUrl() *string {
	return ep.StreamingUrl
}

func (ep *OnsenEpisode) ShowId() string {
	return ep.showRef.ShowId()
}

func (ep *OnsenEpisode) ShowTitle() string {
	return ep.showRef.ShowTitle()
}

func (ep *OnsenEpisode) IsVideo() bool {
	return ep.MediaType != "sound"
}

func (ep *OnsenEpisode) PPrint() model.PPrintRow {
	return model.PPrintRow{
		IsDir:   false,
		CanLoad: ep.CanDownload(),
		IsVid:   ep.IsVideo(),
		Date:    ep.Date(),
		Ref:     ep.EpId(),
		Note:    ep.EpTitle(),
		Cast:    ep.Artists(),
	}
}

func (ep *OnsenEpisode) Show() model.Show {
	v := reflect.ValueOf(ep.showRef).Interface()
	c := v.(model.Show)
	return c
}
