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
