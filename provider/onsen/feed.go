package onsen

import (
	"fmt"

	"github.com/pgeowng/japoto-dl/model"
)

const (
	feedURL    = "https://onsen.ag/web_api/programs/"
	programURL = "https://onsen.ag/web_api/programs/%s"
)

type OnsenShowAccess struct {
	DirectoryName string `json:"directory_name"`
}

func (sa *OnsenShowAccess) ShowId() string {
	return sa.DirectoryName
}

func (p *OnsenUsecase) GetFeed(loader model.Loader) (result chan *OnsenShow, errors chan error) {
	result = make(chan *OnsenShow)
	errors = make(chan error)

	go func() {
		defer func() { close(result) }()
		defer func() { close(errors) }()

		var err error
		eps := []OnsenShowAccess{}
		err = loader.JSON(feedURL, &eps, nil)
		if err != nil {
			errors <- err
			return
		}

		for _, ep := range eps {
			resObj := OnsenShow{}
			err := loader.JSON(fmt.Sprintf(programURL, ep.DirectoryName), &resObj, nil)
			if err != nil {
				return
				return nil, errors.Wrap(err, "onsen.show")
			}

			for idx := range resObj.Contents {
				resObj.Contents[idx].showRef = &resObj
			}

			result <- &resObj
		}
	}()
	return
}
