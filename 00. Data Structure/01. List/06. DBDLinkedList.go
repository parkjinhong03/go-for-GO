package main

type DBDData interface {}

type DBDLinkedList struct {
	head *DBDNode
	tail *DBDNode
	cur *DBDNode
	numOfData int
}

type DBDNode struct {
	data DBDData
	next *DBDNode
	prev *DBDNode
}

func NewDBDLinkedList() *DBDLinkedList {
	DBDLinkedList := DBDLinkedList{
		head:      NewDBDNode(nil),
		tail:      NewDBDNode(nil),
		cur:       nil,
		numOfData: 0,
	}

	DBDLinkedList.head.next = DBDLinkedList.tail
	DBDLinkedList.tail.prev = DBDLinkedList.head

	return &DBDLinkedList
}

func NewDBDNode(data DBDData) *DBDNode {
	return &DBDNode{
		data: data,
		next: nil,
		prev: nil,
	}
}

func (plist *DBDLinkedList) LInsert(data DBDData) {
	newNode := NewDBDNode(data)

	newNode.next = plist.tail
	newNode.prev = plist.tail.prev
	plist.tail.prev.next = newNode
	plist.tail.prev = newNode

	plist.numOfData++
}

func (plist *DBDLinkedList) LFirst(data *DBDData) bool {
	if plist.head.next == nil {
		return false
	}

	plist.cur = plist.head.next
	*data = plist.cur.data
	return true
}

func (plist *DBDLinkedList) LNext(data *DBDData) bool {
	if plist.cur.next == plist.tail {
		return false
	}

	plist.cur = plist.cur.next
	*data = plist.cur.data
	return true
}

func (plist *DBDLinkedList) LPrev(data *DBDData) bool {
	if plist.cur.prev == plist.head {
		return false
	}

	plist.cur = plist.cur.prev
	*data = plist.cur.data
	return true
}

func (plist *DBDLinkedList) LRemove() DBDData {
	rData := plist.cur.data

	plist.cur.prev.next = plist.cur.next
	plist.cur.next.prev = plist.cur.prev
	plist.cur = plist.cur.prev

	plist.numOfData--
	return rData
}