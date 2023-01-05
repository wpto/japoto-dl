package provider

import (
	"fmt"
	"log"
	"reflect"

	"github.com/pgeowng/japoto-dl/internal/types"
	"github.com/pgeowng/japoto-dl/model"
)

type HibikiShow struct {
	AccessId string `json:"access_id"`
	Name     string `json:"name"`
	Casts    []struct {
		Name     string  `json:"name"`
		RollName *string `json:"roll_name"`
	} `json:"casts"`
	PcImageUrl string        `json:"pc_image_url"`
	SpImageUrl string        `json:"sp_image_url"`
	Episode    HibikiEpisode `json:"episode"`
}

func (show *HibikiShow) ShowId() string {
	return show.AccessId
}

func (show *HibikiShow) Artists() []string {
	result := []string{}
	for _, c := range show.Casts {
		result = append(result, c.Name)
		if c.RollName != nil && len(*c.RollName) > 0 {
			result = append(result, *c.RollName)
		}
	}
	return result
}

func (show *HibikiShow) ShowTitle() string {
	return show.Name
}

func (show *HibikiShow) GetEpisodes(loader model.Loader) (result []model.Episode, err error) {
	if show.Episode.Video != nil {
		var url string
		var err error
		url, err = loadCheckPlaylistURL(loader, show.Episode.Video.Id)
		if err != nil {
			log.Printf("WARN: HibikiShow.GetEpisodes: cant get episode playlist url: %v, accessID=%v\n", err, show.AccessId)
			show.Episode.Video.URL = nil
		} else {
			show.Episode.Video.URL = &url
		}

		v := reflect.ValueOf(show.Episode.Video).Interface()
		c := v.(model.Episode)
		if err == nil {
			result = append(result, c)
		}
	}

	if show.Episode.AdditionalVideo != nil {
		var url string
		var err error
		url, err = loadCheckPlaylistURL(loader, show.Episode.AdditionalVideo.Id)
		if err != nil {
			log.Printf("WARN: HibikiShow.GetEpisodes: cant get episode playlist url: %v, accessID=%v\n", err, show.AccessId)
			show.Episode.AdditionalVideo.URL = nil
		} else {
			show.Episode.AdditionalVideo.URL = &url
		}

		v := reflect.ValueOf(show.Episode.AdditionalVideo).Interface()
		c := v.(model.Episode)
		if err == nil {
			result = append(result, c)
		}
	}
	return result, nil
}

func loadCheckPlaylistURL(loader types.Loader, id int) (url string, err error) {
	var checkObj struct {
		PlaylistURL string `json:"playlist_url"`
	}

	err = loader.JSON(fmt.Sprintf(playCheckURL, id), &checkObj, hibikiGopts)
	if err != nil {
		return
	}

	return checkObj.PlaylistURL, nil
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

func (show *HibikiShow) Provider() string {
	return "hibiki"
}
