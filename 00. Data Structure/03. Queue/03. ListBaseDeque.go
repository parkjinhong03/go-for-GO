package main

type DData interface {}

type DQueue struct {
	head *DNode
	tail *DNode
	numOfData int
}

type DNode struct {
	data DData
	next *DNode
	prev *DNode
}
