package ldate

import (
	"strconv"
	"time"
)

type LocDate struct {
	tim time.Time
}

const (
	ISO  = "2006-01-02"
	MMdd = "Jan 02"
)

func NewLocDate(y int, m int, d int) *LocDate {
	return &LocDate{
		tim: time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Now().Location()),
	}
}

func Today() *LocDate {
	return &LocDate{
		tim: time.Now().Truncate(24 * time.Hour),
	}
}

func (date *LocDate) Year() int {
	return date.tim.Year()
}

func (date *LocDate) Month() int {
	return int(date.tim.Month())
}

func (date *LocDate) Day() int {
	return date.tim.Day()
}

func (date *LocDate) String() string {
	return date.Format(ISO)
}

func (date *LocDate) Format(layout string) string {
	return date.tim.Format(layout)
}

func (date *LocDate) After(othDate *LocDate) bool {
	return date.tim.After(othDate.tim.Truncate(24 * time.Hour))
}

func (date *LocDate) UnmarshalJSON(bytArr []byte) error {
	strDate, err := strconv.Unquote(string(bytArr))
	if err != nil {
		return err
	}

	tim, err := time.Parse(time.DateOnly, strDate)
	if err != nil {
		return err
	}

	date.tim = tim
	return nil
}

func (date *LocDate) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(date.String())), nil
}
