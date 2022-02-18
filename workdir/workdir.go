package workdir

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type Workdir struct {
	dir string
}

func NewWorkdir(prefix, uniqueSalt string) *Workdir {
	dir := filepath.Join(prefix, uniqueSalt)
	return &Workdir{dir}
}

func (wd *Workdir) Ensure() {
	err := os.MkdirAll(wd.dir, 0755)
	if err != nil {
		panic("wd: cant ensure folder " + wd.dir)
	}
}

func (wd *Workdir) Save(fileName, fileBody string) error {
	wd.Ensure()
	filePath := filepath.Join(wd.dir, fileName)
	err := ioutil.WriteFile(filePath, []byte(fileBody), 0644)
	if err != nil {
		return errors.Wrap(err, "wd.save("+filePath+")")
	}
	return nil
}

func (wd *Workdir) SaveRaw(fileName string, fileBody []byte) error {
	wd.Ensure()
	filePath := filepath.Join(wd.dir, fileName)
	err := ioutil.WriteFile(filePath, fileBody, 0644)
	if err != nil {
		return errors.Wrap(err, "wd.save("+filePath+")")
	}

	return nil
}
