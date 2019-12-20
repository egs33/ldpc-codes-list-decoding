package ldpc

import "github.com/egs33/ldpc-codes-list-decoding/pkg/node"

type decodingTreeNode struct {
	variableNodes []node.VariableNode
	checkNodes    []node.CheckNode
}

func newDecodingTreeNode(variableNodes int, checkNodes int) decodingTreeNode {
	newNode := decodingTreeNode{
		variableNodes: make([]node.VariableNode, variableNodes),
		checkNodes:    make([]node.CheckNode, checkNodes),
	}
	for i := 0; i < variableNodes; i++ {
		newNode.variableNodes[i] = node.NewVariableNode()
	}
	for i := 0; i < checkNodes; i++ {
		newNode.checkNodes[i] = node.NewCheckNode()
	}
	return newNode
}

// return whether all check node
func (decodingNode *decodingTreeNode) isSatisfyAllChecks(edges []node.Edge) bool {
	estimates := make([]int, len(decodingNode.variableNodes))
	for i, variableNode := range decodingNode.variableNodes {
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
	checks := make([]int, len(decodingNode.checkNodes))
	for _, edge := range edges {
		checks[edge.CheckNodeIndex] += estimates[edge.VariableNodeIndex]
	}

	for _, c := range checks {
		if c%2 != 0 {
			return false
		}
	}
	return true
}

func (decodingNode decodingTreeNode) executeBeliefPropagation(edges []node.Edge, decodeIteration int) {
	for i, edge := range edges {
		message := decodingNode.variableNodes[edge.VariableNodeIndex].CalcInitialMessage()
		decodingNode.checkNodes[edge.CheckNodeIndex].ReceiveMessage(i, message)
	}

	for i := 0; i < decodeIteration; i++ {
		for edgeIndex, edge := range edges {
			message := decodingNode.checkNodes[edge.CheckNodeIndex].CalcMessage(edgeIndex)
			decodingNode.variableNodes[edge.VariableNodeIndex].ReceiveMessage(edgeIndex, message)
		}
		for edgeIndex, edge := range edges {
			message := decodingNode.variableNodes[edge.VariableNodeIndex].CalcMessage(edgeIndex)
			decodingNode.checkNodes[edge.CheckNodeIndex].ReceiveMessage(edgeIndex, message)
		}
		if decodingNode.isSatisfyAllChecks(edges) {
			return
		}
	}
}

func (decodingNode decodingTreeNode) copy() decodingTreeNode {
	newNode := decodingTreeNode{
		variableNodes: make([]node.VariableNode, len(decodingNode.variableNodes)),
		checkNodes:    make([]node.CheckNode, len(decodingNode.checkNodes)),
	}

	for i, v := range decodingNode.variableNodes {
		newNode.variableNodes[i] = v.Copy()
	}

	for i, v := range decodingNode.checkNodes {
		newNode.checkNodes[i] = v.Copy()
	}

	return newNode
}
