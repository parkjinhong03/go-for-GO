package main

import "fmt"

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

func main() {
	bt1 := MakeBTreeNode()
	bt2 := MakeBTreeNode()
	bt3 := MakeBTreeNode()
	bt4 := MakeBTreeNode()

	SetData(bt1, 1)
	SetData(bt2, 2)
	SetData(bt3, 3)
	SetData(bt4, 4)

	MakeLeftSubTree(bt1, bt2)
	MakeRightSubTree(bt1, bt3)
	MakeLeftSubTree(bt2, bt4)

	fmt.Printf("%d\n", GetData(GetLeftSubTree(bt1)))
	// 2
	fmt.Printf("%d\n", GetData(GetLeftSubTree(GetLeftSubTree(bt1))))
	// 4
}