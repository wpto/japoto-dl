package provider

import (
	"fmt"

	"github.com/pgeowng/japoto-dl/model"
	"github.com/pkg/errors"
)

func (show *HibikiShow) LoadImage(loader model.Loader, workdir model.WorkdirBase) error {

	errs := []error{}
	filename := fmt.Sprintf("%s--%s", show.Provider(), show.ShowId())
	err := loadimg(loader, workdir, show.PcImageUrl, filename)
	if err != nil {
		errs = append(errs, errors.Wrap(err, "hibiki.pc_image"))
	}

	filename = fmt.Sprintf("%s-sp--%s", show.Provider(), show.ShowId())
	err = loadimg(loader, workdir, show.SpImageUrl, filename)
	if err != nil {
		errs = append(errs, errors.Wrap(err, "hibiki.sp_image"))
	}

	if len(errs) != 0 {
		msg := ""
		for _, e := range errs {
			msg += e.Error()
			msg += "\n"
		}

		return errors.New(msg)
	}

	return nil
}

func loadimg(loader model.Loader, workdir model.WorkdirBase, url string, filename string) error {
	if len(url) == 0 {
		return errors.New("no image")
	}

	imageBody, err := loader.Raw(url, hibikiGopts)
	if err != nil {
		return err
	}

	ext := GuessContentType(imageBody)
	return workdir.SaveRaw(filename+ext, imageBody)
}
