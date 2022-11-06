package expanddb

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pgeowng/japoto-dl/html/types"
)

func GuessPerformers(artistField string) ([]types.Person, error) {
	if len(artistField) < 1 {
		return make([]types.Person, 0), nil
	}

	re := regexp.MustCompile(`(\[object Object\]|(?:.+?)( 役)?)\s/g`)
	match := re.FindAllStringSubmatch(artistField, -1)
	if len(match) == 0 {
		return nil, fmt.Errorf("empty match for %s", artistField)
	}

	result, err := TryCharacter(match)
	if err != nil {
		fmt.Printf("err: %s is not character", artistField)
	}

	result, err = TryCommon(match)
	if err != nil {
		return nil, fmt.Errorf("unknown match for %s", artistField)
	}

	return result, nil

}

func TryCharacter(match [][]string) ([]types.Person, error) {
	result := make([]types.Person, 0)
	last := types.Person{}
	for idx, p := range match {
		word := p[1]
		if CheckObjectObject(word) {
			fmt.Printf("ch: has obj: %s", word[0])
		}
		if idx%2 == 0 {
			last.Name = word
		}
		if idx%2 == 1 {
			if strings.HasSuffix(word, " 役") {
				last.Character = &word
			} else {
				return nil, fmt.Errorf("seems not like character")
			}
			result = append(result, last)
			last = types.Person{}
		}
	}

	return result, nil
}

func TryCommon(match [][]string) ([]types.Person, error) {
	return nil, fmt.Errorf("not implemented")
}

func CheckObjectObject(field string) bool {
	if field == "[object Object]" {
		return true
	}
	return false
}
