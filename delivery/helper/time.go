package helper

import "time"

func TimeFormatter(oldtime time.Time) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", oldtime.Format("2006-01-02 15:04:05"))
}
