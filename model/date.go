package model

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type Date struct {
	Year  int
	Month int
	Day   int
}

var monthMap []string = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

func (d *Date) String() string {
	result := []string{"----", "---", "--"}

	if d.Year > 0 {
		result[0] = fmt.Sprintf("%4d", d.Year)
	}

	if d.Month > 0 {
		if d.Month > 12 {
			panic(errors.New("bad month"))
		}
		result[1] = monthMap[d.Month-1]
	}

	if d.Day > 0 {
		result[2] = fmt.Sprintf("%2d", d.Day)
	}

	return strings.Join(result, " ")
}

func (d *Date) Filename() string {
	result := []string{"00", "00", "00"}
	if d.Year > 0 {
		result[0] = fmt.Sprintf("%02d", d.Year%100)
	}
	if d.Month > 0 {
		result[1] = fmt.Sprintf("%02d", d.Month)
	}
	if d.Day > 0 {
		result[2] = fmt.Sprintf("%02d", d.Day)
	}
	return strings.Join(result, "")
}

func intStr(i int) string {
	if i < 1 {
		return "--"
	} else {
		return fmt.Sprintf("%02d", i)
	}
}

func (d *Date) IsGood() bool {
	return d.Year > 0 && d.Month > 0 && d.Day > 0
}
