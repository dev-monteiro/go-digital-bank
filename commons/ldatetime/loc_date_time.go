package ldatetime

import (
	"strconv"
	"time"
)

type LocDateTime struct {
	tim time.Time
}

const (
	ISO = "2006-01-02T15:04:05"
)

func Now() *LocDateTime {
	return &LocDateTime{
		tim: time.Now(),
	}
}

func (date *LocDateTime) String() string {
	return date.tim.Format(ISO)
}

func (date *LocDateTime) UnmarshalJSON(bytArr []byte) error {
	strDateTime, err := strconv.Unquote(string(bytArr))
	if err != nil {
		return err
	}

	tim, err := time.Parse(ISO, strDateTime)
	if err != nil {
		return err
	}

	date.tim = tim
	return nil
}

func (date *LocDateTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(date.String())), nil
}
