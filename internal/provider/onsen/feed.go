package onsen

import (
	"reflect"

	"github.com/pgeowng/japoto-dl/internal/model"
	"github.com/pkg/errors"
)

type OnsenShowAccess struct {
	DirectoryName string `json:"directory_name"`
}

func (sa *OnsenShowAccess) GetShow(loader model.Loader) (model.Show, error) {
	resObj := OnsenShow{}
	err := loader.JSON("https://onsen.ag/web_api/programs/"+sa.DirectoryName, &resObj, nil)
	if err != nil {
		return nil, errors.Wrap(err, "onsen.show")
	}

	for idx := range resObj.Contents {
		resObj.Contents[idx].showRef = &resObj
	}

	v := reflect.ValueOf(&resObj).Interface()
	c := v.(model.Show)

	return c, nil
}

func (sa *OnsenShowAccess) ShowId() string {
	return sa.DirectoryName
}

func (p *Onsen) GetFeed(loader model.Loader) ([]model.ShowAccess, error) {
	resObj := []OnsenShowAccess{}
	err := loader.JSON("https://onsen.ag/web_api/programs/", &resObj, nil)
	if err != nil {
		return nil, errors.Wrap(err, "onsen.feed.get")
	}

	result := make([]model.ShowAccess, 0)
	for i := range resObj {
		v := reflect.ValueOf(&resObj[i]).Interface()
		c := v.(model.ShowAccess)
		result = append(result, c)
	}

	return result, nil
}
