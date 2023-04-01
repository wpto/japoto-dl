package provider

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pgeowng/japoto-dl/model"
	"github.com/pkg/errors"
)

type OnsenShowAccess struct {
	DirectoryName string `json:"directory_name"`
}

func (p *Onsen) GetShow(showName string) (model.Show, error) {
	resObj := OnsenShow{}
	err := p.loader.JSON("https://onsen.ag/web_api/programs/"+showName, &resObj, nil)
	if err != nil {
		return nil, errors.Wrap(err, "onsen.show")
	}

	b, _ := json.MarshalIndent(resObj, "", "  ")
	fmt.Println(string(b))

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

func (p *Onsen) GetFeedW() ([]OnsenShowAccess, error) {
	result := []OnsenShowAccess{}
	err := p.loader.JSON("https://onsen.ag/web_api/programs/", &result, nil)
	if err != nil {
		return nil, errors.Wrap(err, "onsen.feed.get")
	}

	return result, nil
}
