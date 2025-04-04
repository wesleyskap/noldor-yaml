package yaml

import (
	"time"
)

func ParseTimestamp(input string) (time.Time, error) {
	formats := []string{time.RFC3339, time.RFC3339Nano, "2006-01-02 15:04:05", "2006-01-02"}
	for _, f := range formats {
		if t, err := time.Parse(f, input); err == nil {
			return t, nil
		}
	}
	return time.Time{}, nil
}
