package node

import (
	"math"
)

type CheckNode struct {
	receivedMessage map[int]float64
}

func NewCheckNode() CheckNode {
	return CheckNode{receivedMessage: map[int]float64{}}
}

func (node CheckNode) CalcMessage(to int) float64 {
	product := 1.0
	for index, message := range node.receivedMessage {
		if index == to {
			continue
		}
		if math.IsInf(message, 1) {
			continue
		}
		if math.IsInf(message, -1) {
			product *= -1
			continue
		}
		product *= math.Tanh(message / 2)
	}

	return 2 * math.Atanh(product)
}

func (node *CheckNode) ReceiveMessage(from int, message float64) {
	node.receivedMessage[from] = message
}

func (node *CheckNode) Clear() {
	node.receivedMessage = map[int]float64{}
}
