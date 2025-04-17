package yaml

import (
	"fmt"
	"strings"
	"time"
)

// ParseTimestamp parses standard ISO-8601 and YAML 1.1/1.2 timestamp scalar representations.
// Usage example:
//   t, err := ParseTimestamp("2025-06-01T12:00:00Z")
func ParseTimestamp(input string) (time.Time, error) {
	val := strings.TrimSpace(input)
	if val == "" {
		return time.Time{}, fmt.Errorf("yaml: empty timestamp scalar, expected ISO-8601 or RFC3339 formatted date")
	}
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, val); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("yaml: invalid timestamp %q, expected ISO8601 or RFC3339 format", input)
}

// IsTimestampScalar checks if a string scalar matches common date patterns.
// Usage example:
//   ok := IsTimestampScalar("2025-06-01")
func IsTimestampScalar(input string) bool {
	_, err := ParseTimestamp(input)
	return err == nil
}

