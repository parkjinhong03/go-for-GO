package main

type BTData interface {}

type BTreeNode struct {
	data BTData
	left *BTreeNode
	right *BTreeNode
}
