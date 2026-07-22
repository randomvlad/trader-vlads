package gimme

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func ValidateNumberStringInRange(value string, minRange, maxRange int) error {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	quantity, err := strconv.Atoi(value)
	if err != nil || quantity < minRange || quantity > maxRange {
		return errors.New(fmt.Sprintf("Invalid value. Enter a number between %v - %v or leave blank", minRange, maxRange))
	}

	return nil
}
