package main

type BTData interface {}

type BTreeNode struct {
	data BTData
	left *BTreeNode
	right *BTreeNode
}

func MakeBTreeNode() *BTreeNode {

}

func GetData(bt *BTreeNode) {

}

func SetData(bt *BTreeNode, data BTData) {

}

func GetLeftSubTree(bt *BTreeNode) *BTreeNode {

}

func GetRightSubTree(bt *BTreeNode) *BTreeNode {

}

func MakeLeftSubTree(main *BTreeNode, sub *BTreeNode) {

}

func MakeRightSubTree(main *BTreeNode, sub *BTreeNode) {

}
