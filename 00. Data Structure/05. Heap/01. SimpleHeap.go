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
	idx := ph.numOfData+1
	nelem := HeapElem{
		pr:   pr,
		data: data,
	}
	// 새 노드가 저장될 위치가 루트 노드의 위치가 아니라면 while문 반복
	for idx!=1 {
		// 새 노드와 부모 노드의 우선순위 비교
		if pr < ph.heapArr[getParentIDX(idx)].pr {
			// 만약 새 노드의 우선순위가 더 높다면 부모 노드 한 단계 실제로 내림
			ph.heapArr[idx] = ph.heapArr[getParentIDX(idx)]
			// 새 노드의 레벨을 한 레벨 올림, 실제로 값을 바꾸진 않고 인덱스 값만 갱신
			idx = getParentIDX(idx)
		} else {
			// 만약 부모 노드의 우선순위가 더 높다면 반복문 중단
			break
		}
	} // 반복문이 끝난 후의 idx 값을 새 노드가 대입될 heapArr의 인덱스 값이다.

	ph.heapArr[idx] = nelem // 새 노드를 배열에 저장
	ph.numOfData++
}

func (ph *Heap) HDelete() HData {
	retData := ph.heapArr[1].data // 반환을 위해서 삭제할 데이터 저장
	lastElem := ph.heapArr[ph.numOfData] // 힙의 마지막 노드 저저이

	// 아래의 변수 parentIdx에는 내림차순으로 비교 후 마지막 노드가 저장될 위치정보가 담김.
	var parentIdx = 1	// 기본적으로 루트 노드의 인덱스인 1 설정
	var childIdx int

	// 루트 노드의 우선순위가 높은 자식 노드를 시작으로 반복문 시작
	for true {
		childIdx = ph.getHiChildIDX(parentIdx)

		if lastElem.pr <= ph.heapArr[childIdx].pr {
			break // 마지막 노드의 우선순위보다 자식 노드의 우선순위가 높으면 반복문 탈출
		}

		ph.heapArr[parentIdx] = ph.heapArr[childIdx] // 마지막 노드보다 우선순위가 높으니 비교대상 노드의 위치를 한 레벨 올림
		parentIdx = childIdx // 마지막 노드가 저장될 위치 정보를 한 레벨 내림
	} // 반복문을 탈출하면 parentIdx에는 마지막 노드의 위치정보가 저장됨

	ph.heapArr[parentIdx] = lastElem // 마지막 노드 최종 저장
	ph.numOfData--
	return retData
}
