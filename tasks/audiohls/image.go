package audiohls

import (
	"github.com/pgeowng/japoto-dl/model"
	"github.com/pkg/errors"
)

func (a *AudioHLSImpl) Image(file model.File) error {
	err := a.workdir.SaveNamedRaw("image", file.BodyRaw())
	if err != nil {
		return errors.Wrap(err, "ahls.img")
	}
	return nil
}
