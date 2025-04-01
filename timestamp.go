package yaml

import (
	"time"
)

func ParseTimestamp(input string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, input); err == nil {
		return t, nil
	}
	return time.Parse(time.RFC3339Nano, input)
}
