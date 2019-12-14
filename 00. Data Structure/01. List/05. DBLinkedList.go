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
