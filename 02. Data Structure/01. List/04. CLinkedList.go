package main

import (
	"fmt"
)

type CData interface {}

// 원형 연결 리스트를 구현하기 위한 구조체 정의
type CLinkedList struct {
	tail *CNode
	cur *CNode
	before *CNode
	numOfData int
}

// 연결 리스트에 각각 연결한 노드 구조체 정의
type CNode struct {
	data CData
	next *CNode
}

// 새 ClinkedList 객체 생성 및 초기화 후 주소 반환
func NewCLinkedList() *CLinkedList {
	return &CLinkedList{
		tail:      nil,
		cur:       nil,
		before:    nil,
		numOfData: 0,
	}
}

// 새 노드 생성 및 초기화 후 주소 반환
func NewCNode(data CData) *CNode {
	return &CNode{
		data: data,
		next: nil,
	}
}

// 새 노드 성성 후 ClinkedList의 tail에 해당 노드 연결
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

// 새 노드 생성 후 CLinkedList의 head(tail Node의 next 필드 값)에 해당 노드 연결
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

// CLinkedList의 before은 tail로, cur은 head로 참조 값을 설정해준다.
func (plist *CLinkedList) LFirst(data *CData) bool {
	if plist.tail == nil {
		return false
	}

	plist.before = plist.tail
	plist.cur = plist.tail.next
	*data = plist.cur.data
	return true
}

// CLickedList의 before과 cur의 참조 값을 각각 다음 노드로 이동시킨다.
func (plist *CLinkedList) LNext(data *CData) bool {
	if plist.cur == plist.tail {
		return false
	}

	plist.before = plist.cur
	plist.cur = plist.cur.next
	*data = plist.cur.data
	return true
}

// 현재 참조 중인 노드를 삭제하고, 만약 현재 참조 값이 tail 이라면 tail을 before의 참조 값으로 바꾼다.
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

// 리시버 번수로 받은 CLinkedList의 현재 데이터 수와 그 값들을 head에서 tail 순서로 출력해준다.
func (plist *CLinkedList) LPrint() {
	data := new(CData)

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
	list := NewCLinkedList()
	data := new(CData)

	list.LInsert(1)
	list.LInsert(2)
	list.LInsertFront(3)
	list.LInsertFront(4)
	
	list.LPrint()
	// 현재 데이터의 수: 4
	// 4 3 1 2

	list.LFirst(data)
	list.LNext(data)
	list.LRemove()

	list.LPrint()
	// 현재 데이터의 수: 3
	// 4 1 2
}