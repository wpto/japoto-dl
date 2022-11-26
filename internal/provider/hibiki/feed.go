package hibiki

import (
	"reflect"

	"github.com/pgeowng/japoto-dl/internal/model"
	"github.com/pkg/errors"
)

type HibikiShowAccess struct {
	AccessId string `json:"access_id"`
}

func (sa *HibikiShowAccess) GetShow(loader model.Loader) (model.Show, error) {
	resObj := HibikiShow{}
	err := loader.JSON("https://vcms-api.hibiki-radio.jp/api/v1/programs/"+sa.AccessId, &resObj, gopts)
	if err != nil {
		return nil, errors.Wrap(err, "hibiki.getshow")
	}

	if resObj.Episode.Video != nil {
		resObj.Episode.Video.showRef = &resObj
		resObj.Episode.Video.epRef = &resObj.Episode
	}

	if resObj.Episode.AdditionalVideo != nil {
		resObj.Episode.AdditionalVideo.IsAdditional = true
		resObj.Episode.AdditionalVideo.showRef = &resObj
		resObj.Episode.AdditionalVideo.epRef = &resObj.Episode
	}

	v := reflect.ValueOf(&resObj).Interface()
	c := v.(model.Show)

	return c, nil
}

func (sa *HibikiShowAccess) ShowId() string {
	return sa.AccessId
}

func (p *Hibiki) GetFeed(loader model.Loader) ([]model.ShowAccess, error) {
	resObj := []HibikiShowAccess{}
	err := loader.JSON("https://vcms-api.hibiki-radio.jp/api/v1//programs?limit=99&page=1", &resObj, gopts)
	if err != nil {
		return nil, errors.Wrap(err, "hibiki.feed")
	}

	result := make([]model.ShowAccess, 0)
	for i := range resObj {
		v := reflect.ValueOf(&resObj[i]).Interface()
		c := v.(model.ShowAccess)
		result = append(result, c)
	}

	return result, nil
}
