package expanddb

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type EpInfo struct {
	Date     string
	ShowId   string
	Provider string
}

type Matcher struct {
	RE     *regexp.Regexp
	Action func(match []string) EpInfo
}

var matchers = []Matcher{
	{
		regexp.MustCompile(`(\d{6})-(.+?)--(.+?).mp3`),
		func(match []string) EpInfo {
			tags := match[3]
			provider := "unknown"
			if strings.Contains(tags, "onsen") {
				provider = "onsen"
			}

			if strings.Contains(tags, "hibiki") {
				provider = "hibiki"
			}
			return EpInfo{match[1], match[2], provider}
		},
	},
	{
		regexp.MustCompile(`(\d{6})-(.+?)-(\d+?|SP\d*?)-(onsen|hibiki)(-p\d+?)?\.mp3`),
		func(match []string) EpInfo {
			return EpInfo{match[1], match[2], match[4]}
		},
	},
	{
		regexp.MustCompile(`(\d{6})-(\d*?)-(.+?)-(onsen|hibiki)(-p\d+?)?\.mp3`),
		func(match []string) EpInfo {
			return EpInfo{match[1], match[3], match[4]}
		},
	},
	{
		regexp.MustCompile(`(\d{6})-(.+?)-(\d*?|SP\d*?)?\.mp3`),
		func(match []string) EpInfo {
			return EpInfo{match[1], match[2], "onsen"}
		},
	},
	{
		regexp.MustCompile(`(\d{3})(.+?)(\d{6}?)(.{4})\.mp3`),
		func(match []string) EpInfo {
			return EpInfo{match[3], match[2], "onsen"}
		}},
	{
		regexp.MustCompile(`(210508)(-100ma)?`),

		func(match []string) EpInfo {
			return EpInfo{"210508", "100man", "onsen"}
		},
	},
	{
		regexp.MustCompile(`(\d{3})radista_ex_(\d{2})`),
		func(match []string) EpInfo {
			return EpInfo{"0000" + match[2], "radista_ex", "onsen"}
		},
	},
	{
		regexp.MustCompile(`(\d)_(\d+)生肉_PsyChe`),
		func(match []string) EpInfo {
			month, err := strconv.Atoi(match[1])
			if err != nil {
				log.Fatalf("parse error %v", match[1])
			}
			day, err := strconv.Atoi(match[2])
			if err != nil {
				log.Fatalf("parse error %v", match[2])
			}
			return EpInfo{
				fmt.Sprintf("20%02d%02d", month, day),
				"watahana",
				"onsen",
			}
		},
	},
	{
		regexp.MustCompile(`【_桐生ココ】あさココ(?:LIVE|ライブ)(?:100回目)?(?:ニュース！|NEWS初回放送)(\d{1,2})\D(\d{1,2})`),
		func(match []string) EpInfo {
			month, err := strconv.Atoi(match[1])
			if err != nil {
				log.Fatalf("parse error %v", match[1])
			}
			day, err := strconv.Atoi(match[2])
			if err != nil {
				log.Fatalf("parse error %v", match[2])
			}

			year := 20
			if month > 11 {
				year = 19
			}

			return EpInfo{
				fmt.Sprintf("%02d%02d%02d", year, month, day),
				"asacoco",
				"youtube",
			}
		},
	},
	{
		regexp.MustCompile(`(-{6})-(.+?)--(onsen|hibiki)-(.{4})?\.mp3`),
		func(match []string) EpInfo {
			return EpInfo{"000000", match[2], match[3]}
		},
	},
	{
		regexp.MustCompile(`--(\d{4})-(.+?)--(onsen|hibiki)-(.{4})?\.mp3`),
		func(match []string) EpInfo {
			return EpInfo{"00" + match[1], match[2], match[3]}
		},
	},
}

func GuessMeta(filename string) (info EpInfo, err error) {
	for _, matcher := range matchers {
		match := matcher.RE.FindStringSubmatch(filename)
		if len(match) != 0 {
			info = matcher.Action(match)
			return
		}
	}
	err = fmt.Errorf("no matcher for %s", filename)
	return
}
