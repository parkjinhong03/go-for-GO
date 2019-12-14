package main

import (
	"fmt"
)

type CData interface {}

type CLinkedList struct {
	tail *CNode
	cur *Node
	before *Node
	numOfData int
}

type CNode struct {
	data CData
	next *CNode
}

func NewCLinkedList() *CLinkedList {

}

func (plist *CLinkedList) LInsert(data CData) {

}

func (plist *CLinkedList) LInsertFront(data CData) {

}

func (plist *CLinkedList) LFirst(data *CData) bool {

}

func (plist *CLinkedList) LNext(data *CData) bool {

}

func (plist *CLinkedList) LRemove() CData {

}
