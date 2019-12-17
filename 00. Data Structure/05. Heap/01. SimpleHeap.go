package main

const HeapLen = 100

type HData int
type Priority int

type HeapElem struct {
	pr Priority	// 값이 작을수록 높은 우선순위
	data HData
}

type Heap struct {
	numOfData int
	heapArr [HeapLen]HeapElem
}
