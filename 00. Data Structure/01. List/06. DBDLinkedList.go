package main

import "fmt"

type DBDData interface {}

// 더미 노드 기반의 양방향 연결 리스트의 구현을 위한 구조체 정의
type DBDLinkedList struct {
	head *DBDNode
	tail *DBDNode
	cur *DBDNode
	numOfData int
}

// 값과 다음 노드, 이전 노드의 주소를 가지고 있는 단위인 노드의 구조체 정의
type DBDNode struct {
	data DBDData
	next *DBDNode
	prev *DBDNode
}

// 새로운 DBDLinkedList 생성 및 헤더와 테일에 더미 노드 연결 후 주소값 반환
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

// 인자값으로 받은 data값을 가지고 있는 노드 생성 후 주소값 반환
func NewDBDNode(data DBDData) *DBDNode {
	return &DBDNode{
		data: data,
		next: nil,
		prev: nil,
	}
}

// DBDLinkedList의 tail 부분에 새 노드 생성 및 연결
func (plist *DBDLinkedList) LInsert(data DBDData) {
	newNode := NewDBDNode(data)

	newNode.next = plist.tail
	newNode.prev = plist.tail.prev
	plist.tail.prev.next = newNode
	plist.tail.prev = newNode

	plist.numOfData++
}

// DBDLinkedList의 현재 참조 값을 head의 더미 노드의 다음 노드로 이동시킴
func (plist *DBDLinkedList) LFirst(data *DBDData) bool {
	if plist.head.next == nil {
		return false
	}

	plist.cur = plist.head.next
	*data = plist.cur.data
	return true
}

// DBDLinkedList의 참조 값을 현재 참조 중인 노드의 다음 노드로 이동시킴
func (plist *DBDLinkedList) LNext(data *DBDData) bool {
	if plist.cur.next == plist.tail {
		return false
	}

	plist.cur = plist.cur.next
	*data = plist.cur.data
	return true
}

// DBDLinkedList의 참조 값을 현재 참조 중인 노드의 이전 노드로 이동시킴
func (plist *DBDLinkedList) LPrev(data *DBDData) bool {
	if plist.cur.prev == plist.head {
		return false
	}

	plist.cur = plist.cur.prev
	*data = plist.cur.data
	return true
}

// 현재 참조 중인 노드를 삭제 및 재연결 후 삭제한 해당 노드의 data값을 반환함
func (plist *DBDLinkedList) LRemove() DBDData {
	rData := plist.cur.data

	plist.cur.prev.next = plist.cur.next
	plist.cur.next.prev = plist.cur.prev
	plist.cur = plist.cur.prev

	plist.numOfData--
	return rData
}

// 리시버 변수로 받은 DBDlinkedList의 현재 데이터의 수와 그 값들을 출력해줌
func (plist *DBDLinkedList) LPrint() {
	data := new(DBDData)

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
	list := NewDBDLinkedList()
	data := new(DBDData)

	list.LInsert(1)
	list.LInsert(2)
	list.LInsert(3)

	list.LPrint()
	// 현재 데이터의 수: 3
	// 1 2 3

	list.LFirst(data)
	list.LNext(data)
	list.LPrev(data)
	list.LRemove()

	list.LPrint()
	// 현재 데이터의 수: 2
	// 2 3
}