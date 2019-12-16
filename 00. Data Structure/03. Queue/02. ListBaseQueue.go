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

func NewLQueue() *LQueue {

}

func NewLNode(data LData) *LNode {
	
}

func (pq *LQueue) Enqueue(data LData) {

}

func (pq *LQueue) Dequeue() LData {

}

func (pq *LQueue) QPrint() {

}
