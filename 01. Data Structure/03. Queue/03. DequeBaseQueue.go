package main

import "fmt"

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
	if pq.DQIsEmpty() {
		return nil
	}
	rData := pq.head.data


	pq.head = pq.head.next
	if pq.head == nil {
		pq.tail = nil
	} else {
		pq.head.prev = nil
	}
	pq.numOfData--
	return rData
}

func (pq *DQueue) DQRemoveLast() DData {
	if pq.DQIsEmpty() {
		return nil
	}
	rData := pq.tail.data


	pq.tail = pq.tail.prev
	if pq.tail == nil {
		pq.head = nil
	} else {
		pq.tail.next = nil
	}

	pq.numOfData--
	return rData
}

func (pq *DQueue) DQGetFirst() DData {
	if pq.DQIsEmpty() {
		return nil
	}

	return pq.head.data
}

func (pq *DQueue) DQGetLast() DData {
	if pq.DQIsEmpty() {
		return nil
	}

	return pq.tail.data
}

func (pq *DQueue) DQPrint() {
	fmt.Printf("현재 데이터의 갯수: %d\n", pq.numOfData)

	if pq.DQIsEmpty() {
		fmt.Println("데이터가 존재하지 않습니다.")
		return
	}

	cur := pq.head; for {
		fmt.Print(cur.data, " ")
		if cur.next == nil {
			break
		}
		cur = cur.next
	}

	fmt.Println()
}

func main() {
	deque := NewDequeue()

	deque.DQAddFirst(1)
	deque.DQAddFirst(2)
	deque.DQAddLast(3)
	deque.DQAddLast(4)

	deque.DQPrint()
	// 현재 데이터의 갯수: 4
	// 2 1 3 4

	deque.DQRemoveFirst()
	deque.DQRemoveLast()

	deque.DQPrint()
	// 현재 데이터의 갯수: 2
	// 1 3

	deque.DQRemoveLast()
	deque.DQRemoveFirst()
	deque.DQRemoveLast()
	deque.DQRemoveFirst()

	deque.DQPrint()
	// 현재 데이터의 갯수: 0
	// 데이터가 존재하지 않습니다.
}