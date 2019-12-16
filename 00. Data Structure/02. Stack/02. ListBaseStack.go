package main

type LData interface {}

type ListStack struct {
	head *ListNode
}

type ListNode struct {
	data LData
	next *ListNode
}

func NewListStack() *ListStack {

}

func NewListNode() *ListNode {

}

func (ps *ListStack) SIsEmpty() bool {

}

func (ps *ListStack) SPush(data LData) {

}

func (ps *ListStack) SPop() LData {

}

func (ps *ListStack) SPrint() {

}
