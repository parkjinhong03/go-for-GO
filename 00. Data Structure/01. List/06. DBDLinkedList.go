package main

type DBDData interface {}

type DBDLinkedList struct {
	head *DBDNode
	tail *DBNode
	cur *DBNode
	numOfData int
}

type DBDNode struct {
	data DBData
	next *DBNode
	prev *DBNode
}

func NewDBDLinkedList() *DBDLinkedList {

}

func NewDBDNode(data DBDData) *DBDNode {
	
}

func (plist *DBDLinkedList) LInsert(data DBDData) {

}

func (plist *DBDLinkedList) LFirst(data DBDData) bool {

}

func (plist *DBDLinkedList) LNext(data DBDData) bool {

}

func (plist *DBDLinkedList) LPrev(data DBDData) bool {

}

func (plist *DBDLinkedList) LRemove() DBDData {

}