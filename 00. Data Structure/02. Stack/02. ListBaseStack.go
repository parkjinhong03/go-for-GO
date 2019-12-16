package main

import "fmt"

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
	newNode := NewListNode(data)

	newNode.next = ps.head
	ps.head = newNode

	ps.numOfData++
}

func (ps *ListStack) SPop() LData {
	if ps.SIsEmpty() {
		return nil
	}

	rData := ps.head.data
	ps.head = ps.head.next

	ps.numOfData--
	return rData
}

func (ps *ListStack) SPrint() {
	curNode := ps.head
	fmt.Printf("현재 데이터의 수: %d\n", ps.numOfData)

	for {
		fmt.Print(curNode.data, " ")

		if curNode.next == nil {
			break
		}
		curNode = curNode.next
	}

	fmt.Println()
}
