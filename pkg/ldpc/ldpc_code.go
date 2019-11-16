package ldpc

import (
	"../node"
	"math/rand"
	"time"
) // TODO: replace with absolute path

type LDPCCode struct {
	codeLength            int
	edges                 []node.Edge
	informationBitIndexes []int
	frozenBitIndexes      []int
}

const decodeIteration = 20

func (code LDPCCode) Decode(channelOutputs []float64) []int {
	variableNodes := make([]node.VariableNode, code.codeLength)
	checkNodes := make([]node.CheckNode, code.codeLength/2)
	for i := 0; i < len(variableNodes); i++ {
		variableNodes[i] = node.NewVariableNode()
	}
	for i := 0; i < len(checkNodes); i++ {
		checkNodes[i] = node.NewCheckNode()
	}

	for i, _ := range variableNodes {
		variableNodes[i].ChannelLLR = channelOutputs[i]
	}

	for _, index := range code.frozenBitIndexes {
		variableNodes[index].SetIsFrozen(true)
	}

	for i, edge := range code.edges {
		message := variableNodes[edge.VariableNodeIndex].CalcInitialMessage()
		checkNodes[edge.CheckNodeIndex].ReceiveMessage(i, message)
	}

	for i := 0; i < decodeIteration; i++ {
		for edgeIndex, edge := range code.edges {
			message := checkNodes[edge.CheckNodeIndex].CalcMessage(edgeIndex)
			variableNodes[edge.VariableNodeIndex].ReceiveMessage(edgeIndex, message)
		}
		for edgeIndex, edge := range code.edges {
			message := variableNodes[edge.VariableNodeIndex].CalcMessage(edgeIndex)
			checkNodes[edge.CheckNodeIndex].ReceiveMessage(edgeIndex, message)
		}
	}

	decoded := make([]int, 0)
	for _, index := range code.informationBitIndexes {
		decoded = append(decoded, variableNodes[index].EstimateSendBit())
	}

	return decoded
}

func (code LDPCCode) GetRate() float64 {
	return float64(len(code.informationBitIndexes)) / float64(code.codeLength-len(code.frozenBitIndexes))
}

/*
Construct Random LDPC code.
- (3, 6) regular
- freeze bits if smaller information bit size
*/
func ConstructCode(codeLength int, informationBitSize int) *LDPCCode {
	code := new(LDPCCode)

	code.codeLength = codeLength
	for i := 0; i < informationBitSize; i++ {
		code.informationBitIndexes = append(code.informationBitIndexes, i)
	}

	for i := informationBitSize; i < codeLength/2; i++ {
		code.frozenBitIndexes = append(code.frozenBitIndexes, i)
	}
	code.edges = createRandomEdges(codeLength)

	return code
}

func createRandomEdges(length int) []node.Edge {
	rand.Seed(time.Now().UnixNano())
	ret := make([]node.Edge, length*3)

	temp := make([]int, length*3)
	for i := range temp {
		temp[i] = i / 3
	}
	for i := 0; i < len(temp); i++ {
		i1 := rand.Intn(length * 3)
		temp[i], temp[i1] = temp[i1], temp[i]
	}
	for k, v := range temp {
		ret[k] = node.Edge{
			VariableNodeIndex: v,
			CheckNodeIndex:    k / 6,
		}
	}

	return ret
}
