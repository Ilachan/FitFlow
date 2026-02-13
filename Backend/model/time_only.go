package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"
)

// TimeOnly stores a time-of-day without a date.
type TimeOnly struct {
	time.Time
}

const (
	timeOnlyLayout      = "15:04:05"
	timeOnlyShortLayout = "15:04"
)

// Scan implements sql.Scanner for TIME columns stored as text.
func (t *TimeOnly) Scan(value interface{}) error {
	if value == nil {
		t.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		t.Time = v
		return nil
	case []byte:
		return t.parseString(string(v))
	case string:
		return t.parseString(v)
	default:
		return fmt.Errorf("unsupported scan type %T for TimeOnly", value)
	}
}

// Value implements driver.Valuer for writing TIME values.
func (t TimeOnly) Value() (driver.Value, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return t.Time.Format(timeOnlyLayout), nil
}

// MarshalJSON renders time-of-day as "HH:MM:SS".
func (t TimeOnly) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte("\"" + t.Time.Format(timeOnlyLayout) + "\""), nil
}

// UnmarshalJSON accepts "HH:MM" or "HH:MM:SS".
func (t *TimeOnly) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		t.Time = time.Time{}
		return nil
	}
	trimmed := strings.Trim(string(data), "\"")
	return t.parseString(trimmed)
}

func (t *TimeOnly) parseString(value string) error {
	if value == "" {
		t.Time = time.Time{}
		return nil
	}
	parsed, err := time.ParseInLocation(timeOnlyLayout, value, time.Local)
	if err != nil {
		parsed, err = time.ParseInLocation(timeOnlyShortLayout, value, time.Local)
		if err != nil {
			return errors.New("invalid time format for TimeOnly")
		}
	}
	t.Time = parsed
	return nil
}
