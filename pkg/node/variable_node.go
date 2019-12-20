package node

import (
	"math"
	"math/rand"
)

type VariableNode struct {
	ChannelLLR      float64
	receivedMessage map[int]float64
	//isFrozen        bool
}

func NewVariableNode() VariableNode {
	return VariableNode{receivedMessage: map[int]float64{}}
}

func (node VariableNode) CalcInitialMessage() float64 {
	//if node.isFrozen {
	//	return math.Inf(1)
	//}
	return node.ChannelLLR
}

func (node VariableNode) CalcMessage(to int) float64 {
	//if node.isFrozen {
	//	return math.Inf(1)
	//}
	sum := node.ChannelLLR
	for index, message := range node.receivedMessage {
		if index == to {
			continue
		}
		if math.IsInf(message, 0) {
			return message
		}
		sum += message
	}

	return sum
}

func (node *VariableNode) ReceiveMessage(from int, message float64) {
	node.receivedMessage[from] = message
}

// return LLR
func (node VariableNode) Marginalize() float64 {
	return node.CalcMessage(-1)
}

func (node VariableNode) EstimateSendBit() int {
	llr := node.Marginalize()
	if llr > 0 {
		return 0
	}
	if llr < 0 {
		return 1
	}
	return rand.Intn(2)
}

func (node *VariableNode) Clear() {
	node.ChannelLLR = 0
	node.receivedMessage = map[int]float64{}
}

func (node VariableNode) Copy() VariableNode {
	newNode := VariableNode{
		ChannelLLR:      node.ChannelLLR,
		receivedMessage: map[int]float64{},
		//isFrozen:        node.isFrozen,
	}

	for k, v := range node.receivedMessage {
		newNode.receivedMessage[k] = v
	}
	return newNode
}
