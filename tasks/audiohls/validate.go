package audiohls

import (
	"github.com/pgeowng/japoto-dl/model"
	"github.com/pkg/errors"
)

func (a *AudioHLSImpl) Validate(file model.File) error {
	data := file.BodyRaw()
	if len(data) == 0 {
		return errors.Errorf("empty file: %s", file.Name())
	}

	if len(data) == 4 &&
		data[0] == 110 &&
		data[1] == 117 &&
		data[2] == 108 &&
		data[3] == 108 {
		return errors.Errorf("file containing null: %s", file.Name())
	}

	err := a.workdir.SaveRaw(file.Name(), file.BodyRaw())
	if err != nil {
		return errors.Wrap(err, "ahls.validate")
	}

	return nil
}
