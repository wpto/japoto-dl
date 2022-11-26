package historyfile

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

type HistoryImpl struct {
	filepath string
}

func NewHistory(path string) *HistoryImpl {
	h := &HistoryImpl{filepath: path}
	c := reflect.ValueOf(h).Interface()
	v := c.(History)
	return v
}

func read(path string) (lines []string, err error) {
	if _, err = os.Stat(path); os.IsNotExist(err) {
		lines = []string{}
		err = nil
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		err = errors.Wrap(err, "history")
		return
	}

	lines = strings.Split(string(data), "\n")
	return
}

func write(path string, lines []string) (err error) {
	err = os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		return errors.Wrap(err, "history")
	}
	return nil
}

func (h *HistoryImpl) Write(key string) error {
	lines, err := read(h.filepath)
	if err != nil {
		return err
	}

	lines = append(lines, key)
	sort.Strings(lines)

	err = write(h.filepath, lines)
	if err != nil {
		return err
	}
	return nil
}

func (h *HistoryImpl) Check(key string) bool {
	lines, err := read(h.filepath)
	if err != nil {
		fmt.Printf("history: cant read %v\n", err)
		return false
	}

	m := make(map[string]bool, len(lines))
	for _, l := range lines {
		m[l] = true
	}

	_, ok := m[key]
	return ok
}
