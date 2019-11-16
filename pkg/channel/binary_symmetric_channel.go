package channel

import (
	"math"
	"math/rand"
	"time"
)

type BinarySymmetricChannel struct {
	crossoverProbability float64
}

func NewBinarySymmetricChannel(crossoverProbability float64) BinarySymmetricChannel {
	return BinarySymmetricChannel{crossoverProbability: crossoverProbability}
}

func (bsc BinarySymmetricChannel) Channel(input []int) []float64 {
	rand.Seed(time.Now().UnixNano())
	output := make([]float64, len(input))
	zeroLLR := math.Log((1 - bsc.crossoverProbability) / bsc.crossoverProbability)

	for index, bit := range input {
		if rand.Float64() < bsc.crossoverProbability {
			bit = 1 - bit
		}
		if bit == 0 {
			output[index] = zeroLLR
		} else {
			output[index] = -zeroLLR
		}
	}

	return output
}
