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

	plist.numOfData--
	return rData
}

func (plist DBLinkedList) LPrint() {
	data := new(DBData)
	fmt.Printf("현재 데이터의 수: %d\n", plist.numOfData)

	if plist.LFirst(data) {
		fmt.Print(*data, " ")

		for {
			if plist.LNext(data) {
				fmt.Print(*data, " ")
				continue
			}
			break
		}
		fmt.Println()
	}
}

func main() {
	list := NewDBLinkedList()
	data := new(DBData)

	list.LInsert(1)
	list.LInsert(2)
	list.LInsert(3)

	list.LPrint()
	// 현재 데이터의 수: 3
	// 3 2 1

	list.LFirst(data)
	list.LNext(data)
	list.LPrevious(data)
	list.LRemove()

	list.LPrint()
	// 현재 데이터의 수: 2
	// 2 1
}