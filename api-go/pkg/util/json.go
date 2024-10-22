package util

import (
	"encoding/json"
	"time"
)

type OhDuration struct {
	time.Duration
}

func (d *OhDuration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		d.Duration = 0
		return nil
	}
	duration, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	d.Duration = duration
	return nil
}

type OhTime struct {
	time.Time
}

func (t *OhTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		t.Time = time.Time{}
		return nil
	}
	layout := "2006-01-02T15:04:05Z07:00"
	parsedTime, err := time.Parse(layout, s)
	if err != nil {
		return err
	}
	t.Time = parsedTime
	return nil
}
