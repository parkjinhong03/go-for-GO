package main

type DData interface {}

type SortingFunc func(Data, Data) int

type DLinkedList struct {
	head *DNode // 더미 노드를 가리키는 필드
	cur *DNode // 참조 및 삭제를 돕는 필드
	before *DNode // 삭제를 돕는 필드
	numOfData int // 저장된 데이터의 수를 기록하기 위한 필드
	comp SortingFunc // 정렬의 기준을 등록하기 위한 필드
}

type DNode struct {
	data DData // 해당 노드의 값을 저장하기 위한 필드
	next *DNode // 다음 노드의 주소값을 저장하기 위한 필드
}

func NewDLinkedList() *DLinkedList {

}

func (plist *DLinkedList) LInset(data DData) {

}

func (plist *DLinkedList) LFirst(data *DData) bool {

}

func (plist *DLinkedList) LSecond(data *DData) bool {

}

func (plist *DLinkedList) LRemove() DData {

}

func (plist *DLinkedList) SetSortRule(sf SortingFunc) {
	
}