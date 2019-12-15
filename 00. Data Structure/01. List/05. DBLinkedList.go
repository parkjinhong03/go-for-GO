package main

import (
	"fmt"
)

type DBData interface {}

type DBLinkedList struct {
	head *DBNode
	cur *DBNode
	numOfData int
}

type DBNode struct {
	data DBData
	next *DBNode
	prev *DBNode
}

func NewDBLinkedList() *DBLinkedList {
	return &DBLinkedList{
		head:      nil,
		cur:       nil,
		numOfData: 0,
	}
}

func NewDBNode(data DBData) *DBNode {
	return &DBNode{
		data:   data,
		next:   nil,
		prev: nil,
	}
}

func (plist *DBLinkedList) LInsert(data DBData) {
	newNode := NewDBNode(data)

	newNode.next = plist.head
	if plist.head != nil {
		plist.head.prev = newNode
	}
	plist.head = newNode

	plist.numOfData++
}

func (plist *DBLinkedList) LFirst(data *DBData) bool {

}

func (plist *DBLinkedList) LNext(data *DBData) bool {

}

func (plist *DBLinkedList) LPrevious(data *DBData) bool {

}

func (plist *DBLinkedList) LRemove() DBData {

}
-