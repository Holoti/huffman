package models

type Forest struct {
	Weight int
	Root   int
}

type Node struct {
	Left   int
	Right  int
	Parent int
	Symbol byte
}

func NewNode() *Node {
	return &Node{Left: -1, Right: -1, Parent: -1}
}
