package commons

import (
	"strconv"
	"time"
)

type LocalDate struct {
	tim time.Time
}

const (
	ISO              = "2006-01-02"
	MonLitCapsDayNum = "Jan 02"
)

func NewLocalDate(y int, m int, d int) *LocalDate {
	return &LocalDate{
		tim: time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Now().Location()),
	}
}

func Today() *LocalDate {
	return &LocalDate{
		tim: time.Now(),
	}
}

func (date *LocalDate) Year() int {
	return date.tim.Year()
}

func (date *LocalDate) Month() int {
	return int(date.tim.Month())
}

func (date *LocalDate) Day() int {
	return date.tim.Day()
}

func (date *LocalDate) String() string {
	return date.Format(ISO)
}

func (date *LocalDate) Format(layout string) string {
	return date.tim.Format(layout)
}

func (date *LocalDate) After(othDate *LocalDate) bool {
	return date.tim.After(othDate.tim.Truncate(24 * time.Hour))
}

func (date *LocalDate) UnmarshalJSON(bytArr []byte) error {
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

func (date *LocalDate) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(date.String())), nil
}
