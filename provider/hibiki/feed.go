package hibiki

import (
	"fmt"

	"github.com/pgeowng/japoto-dl/internal/types"
	"github.com/pkg/errors"
)

const (
	accessURL = "https://vcms-api.hibiki-radio.jp/api/v1/programs/%s"
	feedURL   = "https://vcms-api.hibiki-radio.jp/api/v1//programs?limit=99&page=1"
)

type hibikiAccessResponse struct {
	AccessId string `json:"access_id"`
}

func (sa *hibikiAccessResponse) ShowId() string {
	return sa.AccessId
}

func (uc *HibikiUsecase) GetFeed(loader types.Loader) (result []*HibikiShow, err error) {
	accessEps := []hibikiAccessResponse{}
	err = loader.JSON(feedURL, &accessEps, gopts)
	if err != nil {
		return
	}

	for _, access := range accessEps {
		resObj := HibikiShow{}
		err := loader.JSON(fmt.Sprintf(accessURL, access.AccessId), &resObj, gopts)
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

		result = append(result, &resObj)
	}

	return
}
