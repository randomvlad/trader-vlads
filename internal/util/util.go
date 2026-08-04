package util

import (
	"fmt"
	"image/color"

	"charm.land/lipgloss/v2"
	"github.com/randomvlad/trader-vlads/internal/component/gimme"
)

const CurrencySymbol = "💰"

func FormatMoney(amount int) string {
	return fmt.Sprintf("%s%d", CurrencySymbol, amount)
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
