package main

const QueLen = 10
type CData interface {}

type CQueue struct {
	queArr [QueLen]CData
	front int
	rear int
}

func NewCQueue() *CQueue {

}

func (pq *CQueue) QIsEmpty() bool {

}

func (pq *CQueue) QIsFull() bool {

}

func (pq *CQueue) Enqueue(data CData) {

}

func (pq *CQueue) Dequeue() CData {

}

func (pq *CQueue) QPrint() {

}