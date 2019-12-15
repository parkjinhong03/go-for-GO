package main

type DBDData interface {}

type DBDLinkedList struct {
	head *DBDNode
	tail *DBNode
	cur *DBNode
	numOfData int
}

type DBDNode struct {
	data DBData
	next *DBNode
	prev *DBNode
}