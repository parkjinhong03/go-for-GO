package main

import (
	"errors"
	"fmt"
)

const QueLen = 10
type CData interface {}

type CQueue struct {
	queArr [QueLen]CData
	front int
	rear int
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

func (pq *CQueue) Enqueue(data CData) error {
	if pq.QIsEmpty() {
		return errors.New("Queue is full!! ")
	}

	pq.rear = nextPosIdx(pq.rear)
	pq.queArr[pq.rear] = data
	return nil
}

func (pq *CQueue) Dequeue() (CData, error) {
	if pq.QIsFull() {
		return nil, errors.New("Queue is empty!! ")
	}

	pq.front = nextPosIdx(pq.front)
	return pq.queArr[pq.front], nil
}

func (pq *CQueue) QPrint() {
	for i:=nextPosIdx(pq.front); i<=pq.rear; i=nextPosIdx(i) {
		fmt.Print(pq.queArr[i], " ")
	}
}