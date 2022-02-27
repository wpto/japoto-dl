package hibiki

import (
	"reflect"

	"github.com/pgeowng/japoto-dl/model"
)

type HibikiShow struct {
	AccessId string `json:"access_id"`
	Name     string `json:"name"`
	Casts    []struct {
		Name     string  `json:"name"`
		RollName *string `json:"roll_name"`
	} `json:"casts"`
	Episode HibikiEpisode `json:"episode"`
}

func (show *HibikiShow) ShowId() string {
	return show.AccessId
}

func (show *HibikiShow) Artists() []string {
	result := []string{}
	for _, c := range show.Casts {
		str := c.Name
		if c.RollName != nil && len(*c.RollName) > 0 {
			str += "(" + *c.RollName + ")"
		}
		result = append(result, str)
	}
	return result
}

func (show *HibikiShow) ShowTitle() string {
	return show.Name
}

func (show *HibikiShow) GetEpisodes() []model.Episode {
	result := []model.Episode{}

	if show.Episode.Video != nil {
		v := reflect.ValueOf(show.Episode.Video).Interface()
		c := v.(model.Episode)
		result = append(result, c)
	}
	if show.Episode.AdditionalVideo != nil {
		v := reflect.ValueOf(show.Episode.AdditionalVideo).Interface()
		c := v.(model.Episode)
		result = append(result, c)
	}
	return result
}

func (show *HibikiShow) canDownload() bool {
	return false
}

func (show *HibikiShow) latestDate() model.Date {
	return model.Date{}
}

func (show *HibikiShow) PPrint() model.PPrintRow {
	return model.PPrintRow{
		IsDir:   true,
		CanLoad: show.canDownload(),
		IsVid:   false,
		Date:    show.latestDate(),
		Ref:     show.ShowId(),
		Note:    show.ShowTitle(),
		Cast:    show.Artists(),
	}
}
