package main

import (
	"fmt"
)

type DBData interface {}

type DBLinkedList struct {
	head *DBNode
	cur *DBNode
	numOfData int
}

type DBNode struct {
	data DBData
	next *DBNode
	before *DBNode
}

func (plist *DBLinkedList) LInsert(data DBData) {

}

func (plist *DBLinkedList) LFirst(data *DBData) bool {

}

func (plist *DBLinkedList) LNext(data *DBData) bool {

}

func (plist *DBLinkedList) LPrevious(data *DBData) bool {

}

func (plist *DBLinkedList) LRemove() DBData {

}
