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
	return &LQueue{
		front:     nil,
		rear:      nil,
		numOfData: 0,
	}
}

func NewLNode(data LData) *LNode {
	return &LNode{
		data: data,
		next: nil,
	}
}

func (pq *LQueue) QIsEmpty() bool {
	if pq.front == nil {
		return true
	}
	return false
}

func (pq *LQueue) Enqueue(data LData) {
	newNode := NewLNode(data)

	if pq.QIsEmpty() {
		pq.rear = newNode
		pq.front = newNode
	} else {
		pq.rear.next = newNode
		pq.rear = newNode
	}

	pq.numOfData++
}

func (pq *LQueue) Dequeue() LData {

}

func (pq *LQueue) QPrint() {

}
