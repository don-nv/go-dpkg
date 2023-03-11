package djson

import (
	"bytes"
	"fmt"
	"time"
)

/*
Duration - is used to parse duration value (string) having raw json format. Standard json unmarshal doesn't support it
natively.
*/
type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	data = bytes.Map(
		func(r rune) rune {
			switch r {
			case '"', ' ', '\t', '\n', '\r':
				return -1
			}

			return r
		},
		data,
	)

	duration, err := time.ParseDuration(string(data))
	if err != nil {
		return fmt.Errorf("parsing duration: %w", err)
	}

	d.Duration = duration

	return nil
}
