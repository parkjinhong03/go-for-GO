package main

import "fmt"

type DData interface {}

type SortingFunc func(DData, DData) int

type DLinkedList struct {
	head *DNode // 더미 노드를 가리키는 필드
	cur *DNode // 참조 및 삭제를 돕는 필드
	before *DNode // 삭제를 돕는 필드
	numOfData int // 저장된 데이터의 수를 기록하기 위한 필드
	comp SortingFunc // 정렬의 기준을 등록하기 위한 필드
}

type DNode struct {
	data DData // 해당 노드의 값을 저장하기 위한 필드
	next *DNode // 다음 노드의 주소값을 저장하기 위한 필드
}

func NewDLinkedList() *DLinkedList {
	return &DLinkedList{
		head:      &DNode{
			data: nil,
			next: nil,
		},
		cur:       nil,
		before:    nil,
		numOfData: 0,
		comp:      nil,
	}
}

func NewDNode() *DNode {
	return &DNode{
		data: nil,
		next: nil,
	}
}

func (plist *DLinkedList) fInsert(data DData) {
	newNode := NewDNode()
	newNode.data = data

	newNode.next = plist.head.next
	plist.head.next = newNode

	plist.numOfData++
}

func (plist *DLinkedList) sInsert(data DData) {

}


func (plist *DLinkedList) LInsert(data DData) {
	if plist.comp == nil {
		plist.fInsert(data)
	} else {
		plist.sInsert(data)
	}
}

func (plist *DLinkedList) LFirst(data *DData) bool {
	if plist.head.next == nil {
		return false
	}
	plist.cur = plist.head.next
	plist.before = plist.head

	*data = plist.cur.data
	return true
}

func (plist *DLinkedList) LNext(data *DData) bool {
	if plist.cur.next == nil {
		return false
	}

	plist.before = plist.cur
	plist.cur = plist.cur.next
	*data = plist.cur.data
	return true
}

func (plist *DLinkedList) LRemove() DData {
	rData := plist.cur.data
	plist.before.next = plist.cur.next
	plist.cur = plist.before
	plist.numOfData--
	return rData
}

func (plist DLinkedList) LPrint() {
	fmt.Printf("현재 데이터의 수: %d\n", plist.numOfData)
	data := new(DData)

	if plist.LFirst(data) {
		fmt.Print(*data, " ")

		for {
			if plist.LNext(data) {
				fmt.Print(*data, " ")
				continue
			}
			break
		}
	}

	fmt.Println()
}

func (plist *DLinkedList) SetSortRule(sf SortingFunc) {

}

func main() {
	list := NewDLinkedList()
	data := new(DData)

	list.LInsert(1)
	list.LInsert(2)
	list.LInsert(3)

	list.LPrint()
	// 현재 데이터의 수: 3
	// 3 2 1

	list.LFirst(data)
	list.LNext(data)
	list.LNext(data)
	list.LRemove()

	list.LPrint()
	// 현재 데이터의 수: 2
	// 3 2
}