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

func NewLinkedList() *LinkedList {
	return &LinkedList{
		head:   nil,
		tail:   nil,
		cur:    nil,
		before: nil,
	}
}

func newNode(data Data) *Node {
	return &Node{
		data: data,
		next: nil,
	}
}

func (plist *LinkedList) LInsert(data Data) {
	newNode := newNode(data)

	if plist.head == nil {
		plist.head = newNode
	} else {
		plist.tail.next = newNode
	}

	plist.tail = newNode
}

func (plist *LinkedList) LFirst(data *Data) bool {
	if plist.head == nil {
		return false
	}

	plist.cur = plist.head
	plist.before = nil
	*data = plist.cur.data
	return true
}

func (plist *LinkedList) LNext(data *Data) bool {
	if plist.cur.next == nil {
		return false
	}

	plist.before = plist.cur
	plist.cur = plist.cur.next
	*data = plist.cur.data
	return true
}

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