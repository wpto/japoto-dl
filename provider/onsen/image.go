package onsen

import (
	"fmt"

	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/provider/common"
	"github.com/pgeowng/japoto-dl/workdir"
	"github.com/pkg/errors"
)

func (ep *OnsenEpisode) Image(loader model.Loader, workdir workdir.Workdir) error {
	if len(ep.PosterImageUrl) == 0 {
		fmt.Printf("onsen.img: note empty poster image for %s\n", ep.EpId())
	}

	if len(ep.showRef.Image.Url) == 0 {
		fmt.Printf("onsen.img: note empty show image for %s\n", ep.EpId())
	}

	url := ep.PosterImageUrl

	if len(url) == 0 {
		url = ep.showRef.Image.Url
	}

	if len(url) == 0 {
		return errors.New("onsen.img: not found")
	}

	imageBody, err := loader.Raw(url, gopts)
	if err != nil {
		return errors.Wrap(err, "onsen.img")
	}

	ext := common.GuessContentType(imageBody)
	filename := fmt.Sprintf("%s-%s%s", ep.Show().Provider(), ep.ShowId(), ext)

	return errors.Wrap(workdir.SaveRaw(filename, imageBody), "onsen.img")
}
