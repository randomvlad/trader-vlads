package util

import (
	"errors"
	"fmt"
	"image/color"
	"strconv"
	"strings"

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

func CreateFormInputFields(items []string, getTitle func(item string) string) []gimme.Field {
	fields := make([]gimme.Field, len(items))

	for index, name := range items {
		fields[index] = gimme.NewInput().
			Key(name).
			Title(getTitle(name)).
			Inline(true).
			Prompt(": ").
			Validate(func(inputValue string) error {
				if strings.TrimSpace(inputValue) == "" {
					return nil
				}

				_, err := strconv.Atoi(inputValue)
				if err != nil {
					return errors.New("Please enter a valid number or leave blank")
				}
				return err
			})
	}

	return fields
}
