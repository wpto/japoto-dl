package wd

import (
  "fmt"
  "io/ioutil"
  "os"
  "path/filepath"

  "github.com/pkg/errors"
)

type Wd struct {
  dir string
  // Muxer Muxer
}

func NewWd(prefix, uniqueSalt string) *Wd {
  dir := filepath.Join(prefix, uniqueSalt)
  return &Wd{dir} //Muxer}
}

func (wd *Wd) Resolve(filePath string) string {
  return filepath.Join(wd.dir, filePath)
}

func (wd *Wd) Ensure() {
  err := os.MkdirAll(wd.dir, 0755)
  if err != nil {
    panic("wd: cant ensure folder " + wd.dir)
  }
}

func (wd *Wd) Save(fileName, fileBody string) error {
  wd.Ensure()
  filePath := filepath.Join(wd.dir, fileName)
  err := ioutil.WriteFile(filePath, []byte(fileBody), 0644)
  if err != nil {
    return errors.Wrap(err, "wd.save("+filePath+")")
  }
  return nil
}

func (wd *Wd) SaveRaw(fileName string, fileBody []byte) error {
  fmt.Printf("writing %s\n", fileName)
  wd.Ensure()
  filePath := filepath.Join(wd.dir, fileName)
  err := ioutil.WriteFile(filePath, fileBody, 0644)
  if err != nil {
    return errors.Wrap(err, "wd.save("+filePath+")")
  }

  return nil
}
