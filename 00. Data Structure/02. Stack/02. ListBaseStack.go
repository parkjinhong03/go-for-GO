package main

type LData interface {}

type ListStack struct {
	head *ListNode
	numOfData int
}

type ListNode struct {
	data LData
	next *ListNode
}

func NewListStack() *ListStack {
	return &ListStack{
		head:      nil,
		numOfData: 0,
	}
}

func NewListNode(data LData) *ListNode {
	return &ListNode{
		data: data,
		next: nil,
	}
}

func (ps *ListStack) SIsEmpty() bool {
	if ps.head == nil {
		return true
	}
	return false
}

func (ps *ListStack) SPush(data LData) {
	if ps.SIsEmpty() {
		return
	}

	newNode := NewListNode(data)

	newNode.next = ps.head
	ps.head = newNode
	
	ps.numOfData++
}

func (ps *ListStack) SPop() LData {

}

func (ps *ListStack) SPrint() {

}
