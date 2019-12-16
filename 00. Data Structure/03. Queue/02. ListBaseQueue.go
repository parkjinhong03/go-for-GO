package main

type LData interface {}

type LQueue struct {
	front *LNode
	rear *LNode
	numOfData int
}

type LNode struct {
	data LData
	next *LNode
}