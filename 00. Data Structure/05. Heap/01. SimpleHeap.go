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

func NewHeap() *Heap {
	return &Heap{
		numOfData: 0,
		heapArr:   [HeapLen]HeapElem{},
	}
}

func (ph *Heap) HIsEmpty() bool {
	if ph.numOfData == 0 {
		return true
	}
	return false
}

func getParentIDX(idx int) int {
	return idx/2
}

func getLChildIDX(idx int) int {
	return idx*2
}

func getRChildIDX(idx int) int {
	return idx*2+1
}

func (ph *Heap)getHiChildIDX(idx int) int {
	if getLChildIDX(idx) > ph.numOfData {
		// 자식 노드가 존재하지 않는다면
		return 0
	} else if getLChildIDX(idx) == ph.numOfData {
		// 자식 노드가 왼쪽 자식 노드 하나만 존재한다면
		return getLChildIDX(idx)
	} else {
		// 오른쪽 자식 노드의 우선순위가 높다면,
		if ph.heapArr[getLChildIDX(idx)].pr > ph.heapArr[getRChildIDX(idx)].pr {
			// 오른쪽 자식 노드의 인덱스 값 반환
			return getRChildIDX(idx)
		// 왼쪽 자식 노드의 우선순위가 높다면,
		} else {
			// 왼쪽 자식 노드의 인덱스 값 반환
			return getLChildIDX(idx)
		}
	}
}

func (ph *Heap) HInsert(data HData, pr Priority) {
}

func (ph *Heap) HDelete() HData {

}
