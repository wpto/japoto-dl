package entity

import "fmt"

func FormatSizeHuman(bytes int) string {
	unit := "B"
	size := float64(bytes)
	if size*10 >= 1024 {
		unit = "KB"
		size = size / 1024
	}

	if size*10 >= 1024 {
		unit = "MB"
		size = size / 1024
	}

	if size < 10 {
		return fmt.Sprintf("%.1f%s", size, unit)
	} else {
		return fmt.Sprintf("%.f%s", size, unit)
	}
}

func FormatDurationHuman(seconds int) string {
	minutes := seconds / 60
	return fmt.Sprintf("%d:%02d", minutes, seconds%60)
}
