package main

type BTData interface {}

type BTreeNode struct {
	data BTData
	left *BTreeNode
	right *BTreeNode
}

func MakeBTreeNode() *BTreeNode {
	return &BTreeNode{
		data:  nil,
		left:  nil,
		right: nil,
	}
}

func GetData(bt *BTreeNode) BTData {
	return bt.data
}

func SetData(bt *BTreeNode, data BTData) {
	bt.data = data
}

func GetLeftSubTree(bt *BTreeNode) *BTreeNode {
	return bt.left
}

func GetRightSubTree(bt *BTreeNode) *BTreeNode {
	return bt.right
}

func MakeLeftSubTree(main *BTreeNode, sub *BTreeNode) {
	main.left = sub
}

func MakeRightSubTree(main *BTreeNode, sub *BTreeNode) {
	main.right = sub
}
