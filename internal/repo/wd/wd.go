package wd

import (
	"fmt"
	"os"
	"path/filepath"
)

type Wd struct {
	dir string
}

func NewWd(prefix, salt string) *Wd {
	dir := filepath.Join(prefix, salt)
	return &Wd{dir}
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

func (wd *Wd) Save(filename, filebody string) error {
	return wd.SaveRaw(filename, []byte(filebody))
}

func (wd *Wd) SaveRaw(filename string, filebody []byte) error {
	wd.Ensure()

	tempfile, err := os.CreateTemp(wd.dir, filename)
	if err != nil {
		return err
	}

	targetPath := filepath.Join(wd.dir, filename)

	_, err = tempfile.Write(filebody)
	if err != nil {
		//TODO may not close. how to handle?
		tempfile.Close()
		return err
	}

	err = tempfile.Close()
	if err != nil {
		return err
	}

	err = os.Rename(tempfile.Name(), targetPath)
	if err != nil {
		return err
	}

	return nil
}

func (wd *Wd) Clean() {
	if err := os.RemoveAll(wd.dir); err != nil {
		fmt.Printf("wd(%s): clean error - %v\n", wd.dir, err)
	}
}

func (wd *Wd) CacheDir() string {
	return wd.dir
}
