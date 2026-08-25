package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatCurrency(t *testing.T) {
	// given
	amount := 100

	// when
	actualFormatted := FormatCurrency(amount)

	// then
	assert.Equal(t, "💰100", actualFormatted)
}
