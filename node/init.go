package node

var node *Node

func SetNode(node2 *Node) {
	node = node2
}

func GetNode() *Node {
	if node == nil {
		panic("node未初始化")
	}
	return node
}
