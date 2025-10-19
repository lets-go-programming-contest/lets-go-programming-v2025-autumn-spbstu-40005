package bank

import (
	"fmt"
	"strconv"
	"strings"
)

func convertValue(valueStr string) (float64, error) {
	normalizStr := strings.ReplaceAll(valueStr, ",", ".")

	value, err := strconv.ParseFloat(normalizStr, 64)
	if err != nil {
		return 0, fmt.Errorf("parse float: %w", err)
	}

	return value, nil
}

func convertNumCode(numCodeStr string) (int, error) {
	if strings.TrimSpace(numCodeStr) == "" {
		return 0, nil
	}

	numCode, err := strconv.Atoi(numCodeStr)
	if err != nil {
		return 0, fmt.Errorf("parse int: %w", err)
	}

	return numCode, nil
}
