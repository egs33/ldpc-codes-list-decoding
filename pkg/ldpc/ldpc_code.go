package ldpc

import (
	"errors"
	"github.com/egs33/ldpc-codes-list-decoding/pkg/node"
	"math"
	"math/rand"
	"sort"
	"time"
)

type LDPCCode struct {
	codeLength            int
	edges                 []node.Edge
	informationBitIndexes []int
	frozenBitIndexes      []int
	variableNodes         []node.VariableNode
	checkNodes            []node.CheckNode
}

const decodeIteration = 40

// return whether all check node
func (code LDPCCode) isSatisfyAllChecks() bool {
	estimates := make([]int, code.codeLength)
	for i, variableNode := range code.variableNodes {
		llr := variableNode.Marginalize()
		switch {
		case llr == 0:
			return false
		case llr < 0:
			estimates[i] = 0
		case llr > 0:
			estimates[i] = 1
		}
	}
	checks := make([]int, len(code.checkNodes))
	for _, edge := range code.edges {
		checks[edge.CheckNodeIndex] += estimates[edge.VariableNodeIndex]
	}

	for _, c := range checks {
		if c%2 != 0 {
			return false
		}
	}
	return true

}

func (code LDPCCode) executeMessagePassing(channelOutputs []float64) {
	for i, _ := range code.variableNodes {
		code.variableNodes[i].ChannelLLR = channelOutputs[i]
	}

	for i, edge := range code.edges {
		message := code.variableNodes[edge.VariableNodeIndex].CalcInitialMessage()
		code.checkNodes[edge.CheckNodeIndex].ReceiveMessage(i, message)
	}

	for i := 0; i < decodeIteration; i++ {
		for edgeIndex, edge := range code.edges {
			message := code.checkNodes[edge.CheckNodeIndex].CalcMessage(edgeIndex)
			code.variableNodes[edge.VariableNodeIndex].ReceiveMessage(edgeIndex, message)
		}
		for edgeIndex, edge := range code.edges {
			message := code.variableNodes[edge.VariableNodeIndex].CalcMessage(edgeIndex)
			code.checkNodes[edge.CheckNodeIndex].ReceiveMessage(edgeIndex, message)
		}
		if code.isSatisfyAllChecks() {
			return
		}
	}
}

func (code LDPCCode) Decode(channelOutputs []float64) []int {
	code.executeMessagePassing(channelOutputs)

	decoded := make([]int, 0)
	for _, index := range code.informationBitIndexes {
		decoded = append(decoded, code.variableNodes[index].EstimateSendBit())
	}

	return decoded
}

func (code LDPCCode) ListDecode(channelOutputs []float64, listSize int) [][]int {
	code.executeMessagePassing(channelOutputs)
	ambiguousBitCount := int(math.Floor(math.Log2(float64(listSize))))

	llrs := make([]struct {
		index int
		llr   float64
	}, 0)

	for _, index := range code.informationBitIndexes {
		llr := code.variableNodes[index].Marginalize()
		llrs = append(llrs, struct {
			index int
			llr   float64
		}{index: index, llr: llr})
	}

	sort.Slice(llrs, func(i, j int) bool {
		return math.Abs(llrs[i].llr) < math.Abs(llrs[j].llr)
	})

	uniqueDecoded := make([]int, 0)
	for _, index := range code.informationBitIndexes {
		uniqueDecoded = append(uniqueDecoded, code.variableNodes[index].EstimateSendBit())
	}

	listDecoded := make([][]int, 1)
	listDecoded[0] = uniqueDecoded

	for i := 0; i < ambiguousBitCount; i++ {
		llr := llrs[i]
		temp := make([][]int, 0)
		for _, v := range listDecoded {
			inverted := make([]int, len(v))
			copy(inverted, v)
			inverted[llr.index] = 1 - inverted[llr.index]
			temp = append(temp, inverted)
		}
		listDecoded = append(listDecoded, temp...)
	}

	return listDecoded
}

func (code LDPCCode) GetRate() float64 {
	return float64(len(code.informationBitIndexes)) / float64(code.GetRealCodeLength())
}

func (code LDPCCode) GetListRate(listSize int) float64 {
	ambiguousBitCount := int(math.Floor(math.Log2(float64(listSize))))
	return float64(len(code.informationBitIndexes)-ambiguousBitCount) / float64(code.GetRealCodeLength())
}

func (code LDPCCode) GetRealCodeLength() int {
	return code.codeLength - len(code.frozenBitIndexes)
}

/*
Construct Random Regular LDPC code.
Freeze and puncture bits if smaller information bit size
*/
func ConstructCode(
	originalCodeLength int,
	informationBitSize int,
	variableNodeDegree int,
	checkNodeDegree int) (*LDPCCode, error) {
	code := new(LDPCCode)

	code.codeLength = originalCodeLength
	for i := 0; i < informationBitSize; i++ {
		code.informationBitIndexes = append(code.informationBitIndexes, i)
	}

	originalInfoBitSize := originalCodeLength - originalCodeLength*variableNodeDegree/checkNodeDegree

	for i := informationBitSize; i < originalInfoBitSize; i++ {
		code.frozenBitIndexes = append(code.frozenBitIndexes, i)
	}
	var err error
	code.edges, err = createRandomEdges(originalCodeLength, variableNodeDegree, checkNodeDegree)
	if err != nil {
		return nil, err
	}

	code.variableNodes = make([]node.VariableNode, originalCodeLength)
	code.checkNodes = make([]node.CheckNode, originalCodeLength*variableNodeDegree/checkNodeDegree)
	for i := 0; i < len(code.variableNodes); i++ {
		code.variableNodes[i] = node.NewVariableNode()
	}
	for i := 0; i < len(code.checkNodes); i++ {
		code.checkNodes[i] = node.NewCheckNode()
	}

	for _, index := range code.frozenBitIndexes {
		code.variableNodes[index].SetIsFrozen(true)
	}

	return code, nil
}

func createRandomEdges(length int, variableNodeDegree int, checkNodeDegree int) ([]node.Edge, error) {
	if length*variableNodeDegree%checkNodeDegree != 0 {
		return nil, errors.New("invalid length and degree")
	}
	rand.Seed(time.Now().UnixNano())
	ret := make([]node.Edge, length*variableNodeDegree)

	temp := make([]int, length*variableNodeDegree)
	for i := range temp {
		temp[i] = i / variableNodeDegree
	}
	for i := 0; i < len(temp); i++ {
		i1 := rand.Intn(length * variableNodeDegree)
		temp[i], temp[i1] = temp[i1], temp[i]
	}
	for k, v := range temp {
		ret[k] = node.Edge{
			VariableNodeIndex: v,
			CheckNodeIndex:    k / checkNodeDegree,
		}
	}

	return ret, nil
}
