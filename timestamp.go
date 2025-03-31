package yaml

import (
	"time"
)

func ParseTimestamp(input string) (time.Time, error) {
	return time.Parse(time.RFC3339, input)
}
