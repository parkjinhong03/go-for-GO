package main

import (
	"fmt"
)

type CData interface {}

type CLinkedList struct {
	tail *CNode
	cur *Node
	before *Node
	numOfData int
}

type CNode struct {
	data CData
	next *CNode
}
