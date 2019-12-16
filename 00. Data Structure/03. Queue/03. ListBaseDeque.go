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

func NewDequeue() *DQueue {
	return &DQueue{
		head:      nil,
		tail:      nil,
		numOfData: 0,
	}
}

func NewDNode(data DData) *DNode {
	return &DNode{
		data: data,
		next: nil,
		prev: nil,
	}
}

func (pq *DQueue) DQIsEmpty() bool {
	if pq.head == nil {
		return true
	}
	return false
}

func (pq *DQueue) DQAddFirst(data DData) {
	newNode := NewDNode(data)

	if pq.DQIsEmpty() {
		pq.tail = newNode
		pq.head = newNode
	} else {
		newNode.next = pq.head
		newNode.prev = pq.head.prev
		pq.head.prev.next = newNode
		pq.head.prev = newNode
	}

	pq.numOfData++
}

func (pq *DQueue) DQAddLast(data DData) {

}

func (pq *DQueue) DQRemoveFirst() DData {

}

func (pq *DQueue) DQRemoveLast() DData {

}

func (pq *DQueue) DQGetFirst() DData {

}

func (pq *DQueue) DQGetLast() DData {

}