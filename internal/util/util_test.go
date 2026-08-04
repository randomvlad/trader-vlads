package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatMoney(t *testing.T) {
	// given
	amount := 100

	// when
	actualFormatted := FormatMoney(amount)

	// then
	assert.Equal(t, "💰100", actualFormatted)
}
