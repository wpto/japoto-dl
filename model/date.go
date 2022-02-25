package model

import (
	"fmt"
)

type Date struct {
	Year  int
	Month int
	Day   int
}

func (d *Date) String() string {
	return fmt.Sprintf("%s%s%s", intStr(d.Year%100), intStr(d.Month), intStr(d.Day))
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
