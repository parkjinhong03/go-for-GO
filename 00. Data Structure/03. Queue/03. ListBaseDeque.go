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

}

func NewDNode(data DData) *DNode {

}

func (pq *DQueue) DQIsEmpty() bool {

}

func (pq *DQueue) DQAddFirst(data DData) {

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