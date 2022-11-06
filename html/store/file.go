package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/pgeowng/japoto-dl/html/types"
)

type FileStore struct {
	dir string
}

func NewFileStore(dir string) *FileStore {
	return &FileStore{dir}
}

func (fs *FileStore) Read() []types.Entry {
	files, err := os.ReadDir(fs.dir)
	if err != nil {
		log.Fatalf("error: read %s: %s\n", fs.dir, err)
	}

	result := make([]types.Entry, 0)

	for _, file := range files {
		pp := filepath.Join(fs.dir, file.Name())
		data, err := ioutil.ReadFile(pp)
		if err != nil {
			log.Fatalf("error: read %s: %s\n", pp, err)
		}

		filejson := make([]types.Entry, 0)
		err = json.Unmarshal(data, &filejson)
		if err != nil {
			log.Fatalf("error: parse %s: %s\n", pp, err)
		}
		result = append(result, filejson...)
	}

	return result
}

func (fs *FileStore) Write(eps []types.Entry) {

	pardir := filepath.Dir(fs.dir)
	tmpdir := filepath.Join(pardir, "__")

	err := os.Mkdir(tmpdir, 0755)
	if err != nil {
		log.Fatalf("new dir mk", err)
	}

	sort.Slice(eps, func(i, j int) bool {
		return eps[i].MessageId < eps[j].MessageId
	})

	step := 1000
	idx := 1
	for left := 0; left < len(eps); left += step {

		right := left + step
		if right > len(eps) {
			right = len(eps)
		}

		fmt.Println(left, right)
		sl := eps[left:right]
		text, err := json.MarshalIndent(sl, "", "  ")
		if err != nil {
			log.Fatalf("error when marshaling: %v\n", err)
		}

		pp := filepath.Join(tmpdir, fmt.Sprintf("ep_%03d.json", idx))
		idx += 1
		err = ioutil.WriteFile(pp, text, 0644)
		if err != nil {
			log.Fatalf("error when writing: %v\n", err)
		}
	}

	err = os.RemoveAll(fs.dir)
	if err != nil {
		log.Fatalf("old dir rm", err)
	}
	err = os.Rename(tmpdir, fs.dir)
	if err != nil {
		log.Fatalf("new dir rename", err)
	}
}
