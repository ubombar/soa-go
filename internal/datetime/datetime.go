package datetime

import (
	"time"

	"gopkg.in/yaml.v3"
)

const customDateFormat = "2006-01-02"

type Date struct {
	time.Time
}

func (ct Date) String() string {
	return ct.Format(customDateFormat)
}

func (ct Date) MarshalYAML() (interface{}, error) {
	return ct.Format(customDateFormat), nil
}

func (ct *Date) UnmarshalYAML(value *yaml.Node) error {
	t, err := time.Parse(customDateFormat, value.Value)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

const customDateTimeFormat = "2006-01-02 15:04:05"

type DateTime struct {
	time.Time
}

func (ct DateTime) String() string {
	return ct.Format(customDateTimeFormat)
}

func (ct DateTime) MarshalYAML() (interface{}, error) {
	return ct.Format(customDateTimeFormat), nil
}

func (ct *DateTime) UnmarshalYAML(value *yaml.Node) error {
	t, err := time.Parse(customDateTimeFormat, value.Value)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

func CurrentDateTime() DateTime {
	return DateTime{time.Now()}
}

func CurrentDate() Date {
	return Date{time.Now()}
}
