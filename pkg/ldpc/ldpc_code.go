package ldpc

import (
	"errors"
	"github.com/egs33/ldpc-codes-list-decoding/pkg/node"
	"math"
	"math/rand"
)

type LDPCCode struct {
	codeLength            int
	edges                 []node.Edge
	informationBitIndexes []int
	frozenBitIndexes      []int
	variableNodes         []node.VariableNode
	checkNodes            []node.CheckNode
	decodingTreeNodes     []decodingTreeNode
}

const decodeIteration = 40

func (code LDPCCode) Decode(channelOutputs []float64) []int {
	for i := range code.decodingTreeNodes[0].variableNodes {
		code.decodingTreeNodes[0].variableNodes[i].ChannelLLR = channelOutputs[i]
	}
	code.decodingTreeNodes[0].executeBeliefPropagation(code.edges, decodeIteration)

	decoded := make([]int, 0)
	for _, index := range code.informationBitIndexes {
		decoded = append(decoded, code.decodingTreeNodes[0].variableNodes[index].EstimateSendBit())
	}

	return decoded
}

func (code LDPCCode) ListDecode(channelOutputs []float64, listSize int) [][]int {
	for i := range code.decodingTreeNodes[0].variableNodes {
		code.decodingTreeNodes[0].variableNodes[i].ChannelLLR = channelOutputs[i]
	}
	code.decodingTreeNodes[0].executeBeliefPropagation(code.edges, decodeIteration)
	ambiguousBitCount := int(math.Floor(math.Log2(float64(listSize))))

	for i := 0; i < ambiguousBitCount; i++ {
		originalDecodingNodeSize := len(code.decodingTreeNodes)
		for j := 0; j < originalDecodingNodeSize; j++ {
			minLlrIndex := 0
			minLlr := math.Inf(1)

			for i, v := range code.decodingTreeNodes[0].variableNodes {
				llr := math.Abs(v.Marginalize())
				if minLlr > llr {
					minLlr = llr
					minLlrIndex = i
				}
			}

			newNode := code.decodingTreeNodes[j].copy()
			code.decodingTreeNodes[j].variableNodes[minLlrIndex].ChannelLLR = math.Inf(1)
			newNode.variableNodes[minLlrIndex].ChannelLLR = math.Inf(-1)
			code.decodingTreeNodes = append(code.decodingTreeNodes, newNode)
		}

		for _, n := range code.decodingTreeNodes {
			n.executeBeliefPropagation(code.edges, decodeIteration)
		}
	}
	listDecoded := make([][]int, listSize)
	for i, n := range code.decodingTreeNodes {
		decoded := make([]int, len(code.informationBitIndexes))
		for j, index := range code.informationBitIndexes {
			decoded[j] = n.variableNodes[index].EstimateSendBit()
		}
		listDecoded[i] = decoded
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
*/
func ConstructCode(
	codeLength int,
	informationBitSize int,
	variableNodeDegree int,
	checkNodeDegree int) (*LDPCCode, error) {
	code := new(LDPCCode)
	checkNodes := codeLength * variableNodeDegree / checkNodeDegree
	code.decodingTreeNodes = []decodingTreeNode{newDecodingTreeNode(codeLength, checkNodes)}

	code.codeLength = codeLength
	for i := 0; i < informationBitSize; i++ {
		code.informationBitIndexes = append(code.informationBitIndexes, i)
	}
	var err error
	code.edges, err = createRandomEdges(codeLength, variableNodeDegree, checkNodeDegree)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(code.variableNodes); i++ {
		code.decodingTreeNodes[0].variableNodes[i] = node.NewVariableNode()
	}
	for i := 0; i < len(code.checkNodes); i++ {
		code.decodingTreeNodes[0].checkNodes[i] = node.NewCheckNode()
	}

	return code, nil
}

func createRandomEdges(length int, variableNodeDegree int, checkNodeDegree int) ([]node.Edge, error) {
	if length*variableNodeDegree%checkNodeDegree != 0 {
		return nil, errors.New("invalid length and degree")
	}
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
