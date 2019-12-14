package main

import (
	"fmt"
)

type Data interface {}

// 시작 Node와 끝 Node를 가지는 List 구조체
type LinkedList struct {
	head *Node
	tail *Node
}

// value와 다음 Node의 주솟값을 가지는 각 노드를 표현할 구조체
type Node struct {
	data Data
	next *Node
}

func NewLinkedList() *LinkedList {

}

func newNode(data Data) *Node {

}

func (plist *LinkedList) LInsert(data Data) {

}

func (plist *LinkedList) LFirst(data *Data) {

}

func (plist *LinkedList) LNext(data *Data) {

}

func (plist *LinkedList) LRemove() Data {

}
