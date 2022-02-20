package onsen

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Guess struct {
	DateY   int
	DateM   int
	DateD   int
	EpNum   int
	Special bool
}

func Extract(streamingUrl string, showId string) (*Guess, error) {
	// re := `\/(\d{4})(\d{2})\/(.+?)(\d{5,6}?)?([^-_]{4}|[^-_]{8})(-|_)(\d+|sp\d*)`
	reYM := regexp.MustCompile(`\/(\d{4})(\d{2})\/(.+)`)

	matchYM := reYM.FindStringSubmatch(streamingUrl)
	if matchYM == nil {
		return nil, errors.New("/yyyymm/ not found")
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
	if !strings.HasPrefix(rest, showId) {
		return nil, errors.Errorf("showId(%s) not found: %s", showId, rest)
	}

	result := Guess{
		DateY:   int(year),
		DateM:   int(month),
		DateD:   0,
		EpNum:   0,
		Special: false,
	}

	rest = strings.TrimPrefix(rest, showId)

	reRest := regexp.MustCompile(`(\d{5,7}?)?([^-_]{4}|[^-_]{8})(-|_)(sp)?(\d*)`)
	matchRest := reRest.FindStringSubmatch(rest)
	if matchRest == nil {
		return nil, errors.Errorf("unexpected rest: %s", rest)
	}

	dateStr := matchRest[1]
	spStr := matchRest[4]
	epNumStr := matchRest[5]

	if len(dateStr) > 0 { // assuming yymmdd
		reText := fmt.Sprintf(`^(%d)0*(%d)0*([1-9]+)$`, year%100, month)
		reDate := regexp.MustCompile(reText)

		matchDate := reDate.FindStringSubmatch(dateStr)
		if matchDate == nil {
			return nil, errors.Errorf("cant parse date(%s) based on %d,%d", dateStr, year, month)
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

	if len(spStr) > 0 {
		if spStr == "sp" {
			result.Special = true
		} else {
			return nil, errors.New("sp parse error")
		}
	}

	var epNum int64 = 0

	if len(epNumStr) > 0 {
		epNum, err = strconv.ParseInt(epNumStr, 10, 0)
		if err != nil {
			return nil, errors.Wrap(err, "epNum:")
		}
	}

	result.EpNum = int(epNum)

	return &result, nil
}
