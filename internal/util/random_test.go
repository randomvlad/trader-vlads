package util

import (
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomGenerator_RandomGain(t *testing.T) {
	// given
	generator := createRandomFixed()
	valueStart := 100

	// when
	actualGain := generator.RandomGain(valueStart, 50, 50, 75, 125)

	// then
	assert.Equal(t, 125, actualGain) // fixed rand rolls 127 and gets clamped to 125
}

func TestRandomGenerator_RollChance(t *testing.T) {
	// given
	generator := createRandomFixed()
	coinFlipChance := 0.5

	// when
	actualFalse := generator.RollChance(coinFlipChance)

	// then
	assert.Equal(t, false, actualFalse)
}

func TestRandomGenerator_IntN(t *testing.T) {
	// given
	generator := createRandomFixed()
	coinFlipChance := 0.5

	// when
	actualFalse := generator.RollChance(coinFlipChance)

	// then
	assert.Equal(t, false, actualFalse)
}

func createRandomFixed() *RandomGenerator {
	randomFixed := rand.New(rand.NewPCG(1, 2))
	return NewRandomGenerator(randomFixed)
}
