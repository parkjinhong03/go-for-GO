package main

type LData interface {}

type ListStack struct {
	head *ListNode
}

type ListNode struct {
	data LData
	next *ListNode
}