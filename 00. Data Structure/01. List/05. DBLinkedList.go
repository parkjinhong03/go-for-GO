package main

import (
	"fmt"
	"net/http"
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
	if plist.head == nil {
		return false
	}

	plist.cur = plist.head
	*data = plist.cur.data
	return true
}

func (plist *DBLinkedList) LNext(data *DBData) bool {
	if plist.cur.next == nil {
		return false
	}

	plist.cur = plist.cur.next
	*data = plist.cur.data
	return true
}

func (plist *DBLinkedList) LPrevious(data *DBData) bool {
	if plist.cur.prev == nil {
		return false
	}

	plist.cur = plist.cur.prev
	*data = plist.cur.data
	return true
}

func (plist *DBLinkedList) LRemove() DBData {
	rData := plist.cur.data

	if plist.cur == plist.head {
		plist.cur.next.prev = nil
		plist.head = plist.cur.next
		plist.cur = plist.head
	} else if plist.cur.next == nil {
		plist.cur.prev.next = nil
		plist.cur = plist.cur.prev
	} else {
		plist.cur.prev.next = plist.cur.next
		plist.cur.next.prev = plist.cur.prev
		plist.cur = plist.cur.next
	}

	return rData
}