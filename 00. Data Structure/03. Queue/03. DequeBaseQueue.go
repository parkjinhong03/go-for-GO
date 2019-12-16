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
	} else {
		newNode.next = pq.head
		pq.head.prev = newNode
	}

	pq.head = newNode
	pq.numOfData++
}

func (pq *DQueue) DQAddLast(data DData) {
	newNode := NewDNode(data)

	if pq.DQIsEmpty() {
		pq.head = newNode
	} else {
		newNode.prev = pq.tail
		pq.tail.next = newNode
	}

	pq.tail = newNode
	pq.numOfData++
}

func (pq *DQueue) DQRemoveFirst() DData {
	rData := pq.head.data

	if pq.DQIsEmpty() {
		return nil
	}

	pq.head = pq.head.next
	if pq.head == nil {
		pq.tail = nil
	} else {
		pq.head.prev = nil
	}

	return rData
}

func (pq *DQueue) DQRemoveLast() DData {
	rData := pq.tail.data

	if pq.DQIsEmpty() {
		return nil
	}

	pq.tail = pq.tail.prev
	if pq.tail == nil {
		pq.head = nil
	} else {
		pq.tail.next = nil
	}

	return rData
}

func (pq *DQueue) DQGetFirst() DData {

}

func (pq *DQueue) DQGetLast() DData {

}