package main

import "fmt"

type DData interface {}

type CompareFunc func(DData, DData) bool

type DLinkedList struct {
	head *DNode // 더미 노드를 가리키는 필드
	cur *DNode // 참조 및 삭제를 돕는 필드
	before *DNode // 삭제를 돕는 필드
	numOfData int // 저장된 데이터의 수를 기록하기 위한 필드
	comp CompareFunc // 정렬의 기준을 등록하기 위한 필드
}

type DNode struct {
	data DData // 해당 노드의 값을 저장하기 위한 필드
	next *DNode // 다음 노드의 주소값을 저장하기 위한 필드
}

// 더미 노드 기반의 연결 리스트 객체 생성 및 초기화 후 주소 반환
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

// 노드 생성 및 초기화 후 주소 반환
func NewDNode() *DNode {
	return &DNode{
		data: nil,
		next: nil,
	}
}

// 새 노드 생성 후 더미 노드와 연결시킴
func (plist *DLinkedList) fInsert(data DData) {
	newNode := NewDNode()
	newNode.data = data

	newNode.next = plist.head.next
	plist.head.next = newNode

	plist.numOfData++
}

// 새 노드 생성 후 정의된 정렬 함수를 기준으로 노드를 연결시킴
func (plist *DLinkedList) sInsert(data DData) {
	newNode := NewDNode()
	newNode.data = data
	predNode := plist.head

	for predNode.next != nil && plist.comp(data, predNode.next.data) {
		predNode = predNode.next
	}

	newNode.next = predNode.next
	predNode.next = newNode

	plist.numOfData++
}

// 정의된 정렬 함수가 없다면 fInsert() 실행, 있으면 sInsert() 실행
func (plist *DLinkedList) LInsert(data DData) {
	if plist.comp == nil {
		plist.fInsert(data)
	} else {
		plist.sInsert(data)
	}
}

// 현재 참조 위치를 더미 노드 바로 다음 노드로 이동시킴
func (plist *DLinkedList) LFirst(data *DData) bool {
	if plist.head.next == nil {
		return false
	}
	plist.cur = plist.head.next
	plist.before = plist.head

	*data = plist.cur.data
	return true
}

// 참조 위치를 현재 참조 위치의 다음 연결 노드로 이동시킴
func (plist *DLinkedList) LNext(data *DData) bool {
	if plist.cur.next == nil {
		return false
	}

	plist.before = plist.cur
	plist.cur = plist.cur.next
	*data = plist.cur.data
	return true
}

// 현재 참조된 위치의 노드 삭제 후 리스트 이어 붙힘
func (plist *DLinkedList) LRemove() DData {
	rData := plist.cur.data
	plist.before.next = plist.cur.next
	plist.cur = plist.before
	plist.numOfData--
	return rData
}

// 현재 더미 연결 리스트에 입력된 데이터의 수와 그 값을 출력함
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

// 정렬 함수를 매개변수로 받아 리시버 변수로 받은 객체의 필드에 등록함
func (plist *DLinkedList) SetSortRule(comp CompareFunc) {
	plist.comp = comp
}

func main() {
	list := NewDLinkedList()
	data := new(DData)

	list.SetSortRule(func(d1 DData, d2 DData) bool {
		if d1.(int) >= d1.(int) {
			return true
		} else {
			return false
		}
	})

	list.LInsert(1)
	list.LInsert(2)
	list.LInsert(3)

	list.LPrint()
	// 현재 데이터의 수: 3
	// 1 2 3

	list.LFirst(data)
	list.LNext(data)
	list.LNext(data)
	list.LRemove()

	list.LPrint()
	// 현재 데이터의 수: 2
	// 1 2
}