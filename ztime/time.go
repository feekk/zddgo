package ztime

import(
	"time"
	"github.com/feekk/zddgo/errors"
)

type Duration time.Duration

func (d *Duration) UnmarshalText(text []byte) error {
	timeDura, err := time.ParseDuration(string(text))
	if err == nil {
		*d = Duration(timeDura)
	}
	return errors.With(err)
}


type JsonTime time.Time

func (j JsonTime) MarshalJSON() ([]byte, error) {
	t := time.Time(j)
	return []byte(t.Format("2006-01-02 15:04:05")), nil
}

func (j *JsonTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var timeStr string = string(data)[1:20]
	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	*j = JsonTime(t)
	return errors.With(err)
}

