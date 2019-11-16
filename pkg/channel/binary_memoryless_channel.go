package channel

/*
Binary Input Stationary Memoryless Channel

input is {0, 1}, and output is Log likelihood ratio (ln(W(y|0)/W(y|1)), y is channel output).
*/
type BinaryMemorylessChannel interface {
	Channel([]int) []float64
}
