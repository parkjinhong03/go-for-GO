package main

import (
	"fmt"
)

type CData interface {}

type CLinkedList struct {
	tail *CNode
	cur *CNode
	before *CNode
	numOfData int
}

type CNode struct {
	data CData
	next *CNode
}

func NewCLinkedList() *CLinkedList {
	return &CLinkedList{
		tail:      nil,
		cur:       nil,
		before:    nil,
		numOfData: 0,
	}
}

func NewCNode(data CData) *CNode {
	return &CNode{
		data: data,
		next: nil,
	}
}

func (plist *CLinkedList) LInsert(data CData) {
	newNode := NewCNode(data)

	if plist.tail == nil {
		plist.tail = newNode
		newNode.next = newNode
	} else {
		newNode.next = plist.tail.next
		plist.tail.next = newNode
		plist.tail = newNode
	}

	plist.numOfData++
}

func (plist *CLinkedList) LInsertFront(data CData) {
	newNode := NewCNode(data)

	if plist.tail == nil {
		plist.tail = newNode
		newNode.next = newNode
	} else {
		newNode.next = plist.tail.next
		plist.tail.next = newNode
	}

	plist.numOfData++
}

func (plist *CLinkedList) LFirst(data *CData) bool {
	if plist.tail == nil {
		return false
	}

	plist.before = plist.tail
	plist.cur = plist.tail.next
	*data = plist.cur.data
	return true
}

func (plist *CLinkedList) LNext(data *CData) bool {
	if plist.tail == nil {
		return false
	}

	plist.before = plist.cur
	plist.cur = plist.cur.next
	*data = plist.cur.data
	return true
}

func (plist *CLinkedList) LRemove() CData {
	rData := plist.cur.data

	if plist.cur == plist.tail {
		if plist.tail == plist.tail.next {
			plist.tail = nil
		} else {
			plist.tail = plist.before
		}
	}

	plist.before.next = plist.cur.next
	plist.cur = plist.before
	plist.numOfData--

	return rData
}

func (plist *CLinkedList) LPrint() {
	data := new(CData)

	if plist.LFirst(data) {
		fmt.Print(data, " ")

		for {
			if plist.LNext(data) {
				fmt.Println(data, " ")
				continue
			}
			break
		}

		fmt.Println()
	}
}