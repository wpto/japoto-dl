package onsen

import (
	"reflect"

	"github.com/pgeowng/japoto-dl/model"
)

type OnsenPersonalityGroup struct {
	Roles []struct {
		Role *string `json:"role_name"`
		Name string  `json:"name"`
	} `json:"role_of_performers"`
}

type Performer struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	AllowLike bool   `json:"allow_like"`
}

type OnsenShow struct {
	Id            int    `json:"id"`
	DirectoryName string `json:"directory_name"`
	// Display           bool   `json:"display"`
	// ShowContentsCount int    `json:"show_contents_count"`
	// BrandNew          bool   `json:"brand_new"`
	// BrandNewSp        bool   `json:"brand_new_sp"`

	ProgramInfo struct {
		Title string `json:"title"`
		Image struct {
			Url string `json:"url"`
		} `json:"image"`
	} `json:"program_info"`

	Image struct {
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

	Performers []Performer `json:"performers"`

	RelatedLinks []struct {
		LinkUrl string `json:"link_url"`
		Image   string `json:"image"`
	} `json:"related_links"`

	RelatedInfos []struct {
		Category string `json:"category"`
		Caption  string `json:"caption"`
		LinkUrl  string `json:"link_url"`
		Image    string `json:"image"`
	} `json:"related_infos"`

	RelatedPrograms []struct {
		Title         string      `json:"title"`
		DirectoryName string      `json:"directory_name"`
		Category      string      `json:"category"`
		Image         string      `json:""`
		Performers    []Performer `json:"performers"`
	} `json:"related_programs"`

	GuestInNewContent []Performer    `json:"guest_in_new_content"`
	Guests            []Performer    `json:"guests"`
	Contents          []OnsenEpisode `json:"contents"`

	PersonalityGroups []OnsenPersonalityGroup `json:"personality_groups"`
}

func (show *OnsenShow) GetEpisodes() []model.Episode {
	result := make([]model.Episode, 0)

	for i := range show.Contents {
		v := reflect.ValueOf(&show.Contents[i]).Interface()
		c := v.(model.Episode)
		result = append(result, c)
	}

	return result
}

func (show *OnsenShow) Artists() []string {
	result := []string{}

	for _, pg := range show.PersonalityGroups {
		for _, p := range pg.Roles {
			str := p.Name
			if p.Role != nil && len(*p.Role) > 0 {
				str += "(" + *p.Role + ")"
			}
			result = append(result, str)
		}
	}

	return result
}

func (show *OnsenShow) ShowId() string {
	return show.DirectoryName
}

func (show *OnsenShow) ShowTitle() string {
	return show.ProgramInfo.Title
}

func (show *OnsenShow) canDownload() bool {
	for _, c := range show.Contents {
		if c.StreamingUrl != nil {
			return true
		}
	}
	return false
}

func (show *OnsenShow) PPrint() model.PPrintRow {
	return model.PPrintRow{
		IsDir:   true,
		CanLoad: show.canDownload(),
		IsVid:   false,
		Date:    model.Date{},
		Ref:     show.ShowId(),
		Note:    show.ShowTitle(),
		Cast:    show.Artists(),
	}
}
func (show *OnsenShow) Provider() string {
	return "onsen"
}
