package main

import (
	"fmt"
)

const QueLen = 10
type CData interface {}

type CQueue struct {
	queArr [QueLen]CData
	front int
	rear int
	numOfData int
}

func NewCQueue() *CQueue {
	return &CQueue{
		queArr: [QueLen]CData{},
		front:  0,
		rear:   0,
	}
}

func nextPosIdx(pos int) int {
	if pos == QueLen-1 {
		return 0
	}
	return pos+1
}

func (pq *CQueue) QIsEmpty() bool {
	if pq.rear == pq.front {
		return true
	}
	return false
}

func (pq *CQueue) QIsFull() bool {
	if nextPosIdx(pq.rear) == pq.front {
		return true
	}
	return false
}

func (pq *CQueue) Enqueue(data CData) {
	if pq.QIsFull() {
		return
	}

	pq.rear = nextPosIdx(pq.rear)
	pq.queArr[pq.rear] = data
	pq.numOfData++
}

func (pq *CQueue) Dequeue() CData {
	if pq.QIsEmpty() {
		return nil
	}

	pq.front = nextPosIdx(pq.front)
	pq.numOfData--
	return pq.queArr[pq.front]
}

func (pq CQueue) QPrint() {
	fmt.Printf("현재 데이터의 수: %d\n", pq.numOfData)

	for i:=nextPosIdx(pq.front); i<=pq.rear; i=nextPosIdx(i) {
		fmt.Print(pq.queArr[i], " ")
	}
	fmt.Println()
}

func main() {
	queue := NewCQueue()

	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)

	queue.QPrint()
	// 현재 데이터의 수: 3
	// 1 2 3

	queue.Dequeue()

	queue.QPrint()
	// 현재 데이터의 수: 2
	// 2 3
}