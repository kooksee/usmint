package node


var n *Node

func GetNode() *Node {
	if n == nil {
		panic("please init node")
	}
	return n
}

func InitNode(n1 *Node) {
	n = n1
}