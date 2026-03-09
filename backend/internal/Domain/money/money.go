package money

import (
	"errors"
	"strconv"
	"strings"
)

func ParseToCents(raw string) (int64, error) {
	s := strings.TrimSpace(strings.ReplaceAll(raw, ",", "."))
	if s == "" {
		return 0, errors.New("invalid price format")
	}
	if strings.HasPrefix(s, "-") {
		return 0, errors.New("price must be positive")
	}

	parts := strings.Split(s, ".")
	if len(parts) > 2 {
		return 0, errors.New("invalid price format")
	}
	if parts[0] == "" {
		return 0, errors.New("invalid price format")
	}

	whole, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, errors.New("invalid price format")
	}

	frac := int64(0)
	if len(parts) == 2 {
		if len(parts[1]) == 0 || len(parts[1]) > 2 {
			return 0, errors.New("invalid price format")
		}
		fracPart := parts[1]
		if len(fracPart) == 1 {
			fracPart += "0"
		}
		frac, err = strconv.ParseInt(fracPart, 10, 64)
		if err != nil {
			return 0, errors.New("invalid price format")
		}
	}

	cents := whole*100 + frac
	if cents <= 0 {
		return 0, errors.New("price must be positive")
	}
	return cents, nil
}
