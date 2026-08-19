package util

import (
	"math"
	"math/rand/v2"
)

type RandomGenerator struct {
	random *rand.Rand
}

func NewRandomGenerator(random *rand.Rand) *RandomGenerator {
	if random == nil {
		src := rand.NewPCG(rand.Uint64(), rand.Uint64())
		random = rand.New(src)
	}

	return &RandomGenerator{random: random}
}

func (r *RandomGenerator) RollInt(valueMin, valueMax int) int {
	randomMaxExclusive := valueMax - valueMin + 1
	return valueMin + r.IntN(randomMaxExclusive)
}

func (r *RandomGenerator) RollChance(chanceYes float64) bool {
	return r.random.Float64() < chanceYes
}

func (r *RandomGenerator) GainPercent(value, gainNegativePercCap, gainPositivePercCap, clampMin, clampMax int) int {

	randomMaxExclusive := gainNegativePercCap + gainPositivePercCap + 1
	percentChange := r.IntN(randomMaxExclusive) - gainNegativePercCap

	multiplier := 1.0 + (float64(percentChange) / 100.0)
	newValue := int(math.RoundToEven(float64(value) * multiplier))

	return Clamp(newValue, clampMin, clampMax)
}

func (r *RandomGenerator) IntN(maxExclusive int) int {
	return r.random.IntN(maxExclusive)
}
