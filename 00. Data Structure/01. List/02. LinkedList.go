package main

import (
	"fmt"
)

type Data interface {}

// 시작 Node와 끝 Node를 가지는 List 구조체
type LinkedList struct {
	head *Node
	tail *Node
	cur *Node
	before *Node
}

// value와 다음 Node의 주솟값을 가지는 각 노드를 표현할 구조체
type Node struct {
	data Data
	next *Node
}

// 새로운 LinkedList 생성 및 초기화 후 주소값 반환 함수
func NewLinkedList() *LinkedList {
	return &LinkedList{
		head:   nil,
		tail:   nil,
		cur:    nil,
		before: nil,
	}
}

// 새로운 Node 생성 및 인자값으로 받은 값으로 초기화 후 주소값 반환 함수
func newNode(data Data) *Node {
	return &Node{
		data: data,
		next: nil,
	}
}

// 리시버로 받은 연결 리스트의 맨 마지막 부분에 새 노드를 추가하는 메서드
func (plist *LinkedList) LInsert(data Data) {
	newNode := newNode(data)

	if plist.head == nil {
		plist.head = newNode
	} else {
		plist.tail.next = newNode
	}

	plist.tail = newNode
}

// 리시버로 받은 연결 리스트의 첫 번째 노드를 참조한 후 해당 노드의 값을 매개변수에 담는 메서드
func (plist *LinkedList) LFirst(data *Data) bool {
	if plist.head == nil {
		return false
	}

	plist.cur = plist.head
	plist.before = nil
	*data = plist.cur.data
	return true
}

// 리시버로 받은 연결 리스트의 현재 참조 노드의 다음 노드를 참조한 후 해당 노드의 값을 매개변수에 담는 메서드
func (plist *LinkedList) LNext(data *Data) bool {
	if plist.cur.next == nil {
		return false
	}

	plist.before = plist.cur
	plist.cur = plist.cur.next
	*data = plist.cur.data
	return true
}

// 리시버로 받은 연결 리스트의 현재 참조중인 노드를 삭제한 후 삭제한 값을 반환하는 메서드
func (plist *LinkedList) LRemove() Data {
	if plist.cur == nil {
		fmt.Println("참조된 노드가 없습니다")
		return nil
	}
	rpos := plist.cur.data
	plist.before.next = plist.cur.next
	plist.cur = plist.before
	return rpos
}

// 리시버로 받은 연결 리스트의 현재 참조중인 값들과 저장된 값들을 출력하는 메서드
func (plist LinkedList) LPrint() {
	data := new(Data)

	if plist.before == nil {
		fmt.Println("before: nil")
	} else {
		fmt.Println("before: ", plist.before.data)
	}
	if plist.cur == nil {
		fmt.Println("cur: nil")
	} else {
		fmt.Println("cur: ", plist.cur.data)
	}

	if plist.LFirst(data) {
		fmt.Print(*data, " ")

		for {
			if plist.LNext(data) {
				fmt.Print(*data , " ")
				continue
			}
			break
		}

		fmt.Println()
	}
}

func main() {
	list := NewLinkedList()
	data := new(Data)

	list.LInsert(1)
	list.LInsert(2)
	list.LInsert(3)

	list.LPrint()
	// before: nil
	// cur: nil
	// 1 2 3

	list.LFirst(data)
	list.LNext(data)
	list.LRemove()

	list.LPrint()
	// before:  1
	// cur:  1
	// 1 3
}
