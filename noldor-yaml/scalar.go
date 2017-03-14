package yaml

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ParseBool resolves YAML 1.1 and 1.2 boolean scalars.
// Usage example:
//   b, err := ParseBool("yes")
func ParseBool(input string) (bool, error) {
	val := strings.TrimSpace(strings.ToLower(input))
	switch val {
	case "true", "yes", "on", "1":
		return true, nil
	case "false", "no", "off", "0":
		return false, nil
	default:
		return false, fmt.Errorf("yaml: invalid boolean scalar %q, expected one of [true, false, yes, no, on, off]", input)
	}
}

// ParseInt resolves decimal, hex, binary, and octal integer scalars.
// Usage example:
//   n, err := ParseInt("0755")
func ParseInt(input string) (int64, error) {
	val := strings.TrimSpace(input)
	if val == "" {
		return 0, fmt.Errorf("yaml: empty scalar string, expected valid integer representation")
	}
	val = strings.ReplaceAll(val, "_", "")
	return resolveIntBase(val)
}

// resolveIntBase inspects integer prefixes for octal, hex, binary, and decimal.
func resolveIntBase(val string) (int64, error) {
	if strings.HasPrefix(val, "0o") || strings.HasPrefix(val, "0O") {
		return strconv.ParseInt(val[2:], 8, 64)
	}
	if strings.HasPrefix(val, "0x") || strings.HasPrefix(val, "0X") {
		return strconv.ParseInt(val[2:], 16, 64)
	}
	if strings.HasPrefix(val, "0b") || strings.HasPrefix(val, "0B") {
		return strconv.ParseInt(val[2:], 2, 64)
	}
	return parseLegacyOctalOrDecimal(val)
}

// parseLegacyOctalOrDecimal resolves YAML 1.1 legacy octals (0755) or standard base 10 integers.
func parseLegacyOctalOrDecimal(val string) (int64, error) {
	clean := val
	sign := int64(1)
	if strings.HasPrefix(clean, "-") {
		sign = -1
		clean = clean[1:]
	} else if strings.HasPrefix(clean, "+") {
		clean = clean[1:]
	}
	if len(clean) > 1 && clean[0] == '0' && clean[1] >= '0' && clean[1] <= '7' {
		res, err := strconv.ParseInt(clean, 8, 64)
		return res * sign, err
	}
	res, err := strconv.ParseInt(clean, 10, 64)
	return res * sign, err
}

// ParseFloat resolves scalar representation into float64.
// Usage example:
//   f, err := ParseFloat(".inf")
func ParseFloat(input string) (float64, error) {
	val := strings.TrimSpace(strings.ToLower(input))
	switch val {
	case ".inf", "+.inf":
		return math.Inf(1), nil
	case "-.inf":
		return math.Inf(-1), nil
	case ".nan":
		return math.NaN(), nil
	default:
		return strconv.ParseFloat(val, 64)
	}
}

// IsNullScalar determines if a scalar string represents a null value.
// Usage example:
//   isNull := IsNullScalar("null")
func IsNullScalar(input string) bool {
	val := strings.TrimSpace(input)
	return val == "" || val == "~" || strings.ToLower(val) == "null"
}
