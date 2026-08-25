package util

import (
	"fmt"
	"image/color"
	"strconv"

	"charm.land/lipgloss/v2"
	"github.com/randomvlad/trader-vlads/internal/component/gimme"
)

const CurrencySymbol = "💰"

func FormatCurrency(gold int) string {
	return fmt.Sprintf("%s%d", CurrencySymbol, gold)
}

func ViewTurnsLeft(permanent bool, turnsLeft int) string {
	if permanent {
		return "each week"
	} else {
		var expiresSoonIcon string
		if turnsLeft == 1 {
			expiresSoonIcon = " ⌛"
		}

		return "for " + FormatCountPluralized(turnsLeft, "week") + expiresSoonIcon
	}
}

func FormatCountPluralized(count int, nameSingular string) string {
	var nameFormatted string
	if count == 1 {
		nameFormatted = nameSingular
	} else {
		nameFormatted = nameSingular + "s"
	}

	return strconv.Itoa(count) + " " + nameFormatted
}

func ToColors(names ...string) []color.Color {
	var objs []color.Color
	for _, name := range names {
		objs = append(objs, lipgloss.Color(name))
	}
	return objs
}

func CreateFormInputFields(
	items []string,
	getTitle func(item string) string,
	validate func(value string, input *gimme.Input) error,
) []gimme.Field {

	fields := make([]gimme.Field, len(items))

	for index, name := range items {
		fields[index] = gimme.NewInput().
			Key(name).
			Title(getTitle(name)).
			Inline(true).
			Prompt(": ").
			Validate(validate)
	}

	return fields
}

func Clamp(value, low, high int) int {
	return max(low, min(value, high))
}

func GetOrDefault[T comparable](val, def T) T {
	var zero T
	if val == zero {
		return def
	}
	return val
}
