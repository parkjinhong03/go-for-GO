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

func (plist *DBDLinkedList) LFirst(data DBDData) bool {

}

func (plist *DBDLinkedList) LNext(data DBDData) bool {

}

func (plist *DBDLinkedList) LPrev(data DBDData) bool {

}

func (plist *DBDLinkedList) LRemove() DBDData {
	rData := plist.cur.data

	plist.cur.prev.next = plist.cur.next
	plist.cur.next.prev = plist.cur.prev
	plist.cur = plist.cur.prev

	plist.numOfData--
	return rData
}