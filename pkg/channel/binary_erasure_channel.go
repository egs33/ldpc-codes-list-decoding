package channel

import (
	"math"
	"math/rand"
)

type BinaryErasureChannel struct {
	erasureProbability float64
}

func NewBinaryErasureChannel(erasureProbability float64) BinaryErasureChannel {
	return BinaryErasureChannel{erasureProbability: erasureProbability}
}

func (bec BinaryErasureChannel) Channel(input []int) []float64 {
	output := make([]float64, len(input))

	for index, bit := range input {
		if rand.Float64() < bec.erasureProbability {
			output[index] = 0
			continue
		}
		if bit == 0 {
			output[index] = math.Inf(1)
		} else {
			output[index] = math.Inf(-1)
		}
	}

	return output
}
