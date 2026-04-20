package service

import (
	"math"
	"time"
)

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func daysLeft(raw string) int {
	due, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return 0
	}
	d := int(math.Ceil(time.Until(due).Hours() / 24))
	if d < 0 {
		return 0
	}
	return d
}

func formatDate(raw string) string {
	t, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return raw
	}
	return t.Format("Jan 2, 2006")
}

