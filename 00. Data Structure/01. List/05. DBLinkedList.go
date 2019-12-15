package main

import (
	"fmt"
)

type DBData interface {}

// 양방향 연결 리스트를 구현하기 위한 구조체 정의
type DBLinkedList struct {
	head *DBNode
	cur *DBNode
	numOfData int
}

// 양방향 연결 리스트의 각각의 데이터가 저장될 노드 구조체 정의
type DBNode struct {
	data DBData
	next *DBNode
	prev *DBNode
}

// 새로운 DBLinkedList 객체를 생성 및 초기화 후 반환
func NewDBLinkedList() *DBLinkedList {
	return &DBLinkedList{
		head:      nil,
		cur:       nil,
		numOfData: 0,
	}
}

// 세로운 DBNode 객체를 생성 및 초기화 후 반환
func NewDBNode(data DBData) *DBNode {
	return &DBNode{
		data:   data,
		next:   nil,
		prev: nil,
	}
}

// 새로운 Node를 생성하여 DBLinkedList의 head 부분에 연결시킴
func (plist *DBLinkedList) LInsert(data DBData) {
	newNode := NewDBNode(data)

	newNode.next = plist.head
	if plist.head != nil {
		plist.head.prev = newNode
	}
	plist.head = newNode

	plist.numOfData++
}

// DBLinkedList의 참조값을 head인 노드로 이동시킴
func (plist *DBLinkedList) LFirst(data *DBData) bool {
	if plist.head == nil {
		return false
	}

	plist.cur = plist.head
	*data = plist.cur.data
	return true
}

// DBLinkedList의 참조값을 현재 참조중인 노드의 다음 연결 노드로 이동시킴
func (plist *DBLinkedList) LNext(data *DBData) bool {
	if plist.cur.next == nil {
		return false
	}

	plist.cur = plist.cur.next
	*data = plist.cur.data
	return true
}

// DBLinkedList의 참조값을 현재 참조중인 노드의 이전 연결 노드로 이동시킴
func (plist *DBLinkedList) LPrevious(data *DBData) bool {
	if plist.cur.prev == nil {
		return false
	}

	plist.cur = plist.cur.prev
	*data = plist.cur.data
	return true
}

// DBLinkedList의 현재 참조 중인 노드를 삭제한 후 그 삭제한 값을 반환함
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

// 리시버 변수로 받은 DBLinkedList의 현재 총 데이터의 수와 그 값들을 출력함
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