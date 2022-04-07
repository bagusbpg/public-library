package helper

import "time"

func TimeFormatter(oldtime interface{}) (t time.Time, err error) {
	_t := oldtime.(time.Time)
	t, err = time.Parse("2006-01-02 15:04:05", _t.Format("2006-01-02 15:04:05"))
	return
}
