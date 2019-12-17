// 02. ArrayBaseHeap을 import 했다는 조건으로 시작, min heap을 이용한 우선순위 큐 구현 코드
package main

import "fmt"

func NewPQueue() *Heap {
	return NewHeap()
}

func PQIsEmpty(ppq *Heap) bool {
	return ppq.HIsEmpty()
}

func PEnqueue(ppq *Heap, data HData) {
	ppq.HInsert(data, Priority(ppq.numOfData))
}

func PDequeue(ppq *Heap) HData {
	return ppq.HDelete()
}

func main() {
	pq := NewPQueue()

	PEnqueue(pq, 'A')
	PEnqueue(pq, 'B')
	PEnqueue(pq, 'C')

	for !PQIsEmpty(pq) {
		fmt.Println(PDequeue(pq))
	}
}