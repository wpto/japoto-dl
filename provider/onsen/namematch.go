package onsen

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Guess struct {
	DateY int
	DateM int
	DateD int
	Num   string
}

var mistakes map[string][]string = map[string][]string{
	"gaikotsukishi":   {"gaikotukishi"},
	"sega_girls":      {"segagirls"},
	"mushinobu_radio": {"mushinoburadio"},
	"seikowa_otsuge":  {"seikowaotsuge"},
	"seikowa_radio":   {"seikowaradio"},
	"tane":            {"tate"},
	"fuchigami_mai":   {"fuchigamimai"},
	"ore-ski":         {"ore.ski"},
	"techno-roid":     {"technoroid"},
}

func Extract(streamingUrl string, showId string) (*Guess, error) {
	reYM := regexp.MustCompile(`\/(\d{4})(\d{2})\/(.+)`)

	matchYM := reYM.FindStringSubmatch(streamingUrl)
	if matchYM == nil {
		return &Guess{-1, -1, -1, ""}, nil
		// return nil, errors.New("/yyyymm/ not found")
	}

	yearStr := matchYM[1]
	year, err := strconv.ParseInt(yearStr, 10, 0)
	if err != nil {
		return nil, errors.Errorf("cant parse year %s", yearStr)
	}

	monthStr := matchYM[2]
	month, err := strconv.ParseInt(monthStr, 10, 0)
	if err != nil {
		return nil, errors.Errorf("cant parse month %s", monthStr)
	}

	rest := matchYM[3]
	result := Guess{
		DateY: int(year),
		DateM: int(month),
		DateD: -1,
		Num:   "",
	}

	if !strings.HasPrefix(rest, showId) {
		bypass := false
		miss, ok := mistakes[showId]
		if ok {
			for _, m := range miss {
				if strings.HasPrefix(rest, m) {
					rest = strings.TrimPrefix(rest, showId)
					bypass = true
					break
				}
			}
		}

		if !bypass {
			return nil, errors.Errorf("showId(%s) not found: %s", showId, rest)
		}
	}

	rest = strings.TrimPrefix(rest, showId)

	reRest := regexp.MustCompile(`(\d{5,7}?)?([^-_]{4}|[^-_]{8})(?:-|_)(.+).mp4`)
	matchRest := reRest.FindStringSubmatch(rest)
	if matchRest == nil {
		return nil, errors.Errorf("unexpected rest: %s", rest)
	}

	dateStr := matchRest[1]
	result.Num = matchRest[3]

	if len(dateStr) > 0 { // assuming yymmdd
		reText := fmt.Sprintf(`^(%d)0*(%d)0*([1-9]\d*)$`, year%100, month)
		reDate := regexp.MustCompile(reText)

		matchDate := reDate.FindStringSubmatch(dateStr)
		if matchDate == nil {
			return nil, errors.Errorf("cant parse date(%s) based on %d,%d", dateStr, year%100, month)
		}

		yyStr := matchDate[1]
		_, err = strconv.ParseInt(yyStr, 10, 0)
		if err != nil {
			return nil, errors.Errorf("cant parse yy %s in %s", yyStr, dateStr)
		}

		mmStr := matchDate[2]
		_, err = strconv.ParseInt(mmStr, 10, 0)
		if err != nil {
			return nil, errors.Errorf("cant parse mm %s in %s", mmStr, dateStr)
		}

		ddStr := matchDate[3]
		dd, err := strconv.ParseInt(ddStr, 10, 0)
		if err != nil {
			return nil, errors.Errorf("cant parse dd %s in %s", ddStr, dateStr)
		}

		result.DateD = int(dd)
	}

	return &result, nil
}
