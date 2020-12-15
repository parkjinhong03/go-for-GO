package algorithm2

import (
	"strconv"
	"strings"
)

type HeapNode struct {
	data int
}

type MinHeap struct {
	count int
	Nodes []HeapNode
}

func NewMinHeap(queueLen int) MinHeap {
	return MinHeap{
		Nodes: make([]HeapNode, queueLen),
	}
}

func (heap *MinHeap) getMinChildIdx(parentIdx int) int {
	lChildIdx := parentIdx * 2
	rChildIdx := parentIdx * 2 + 1
	if rChildIdx > heap.count {
		return lChildIdx
	}
	if heap.Nodes[rChildIdx].data > heap.Nodes[lChildIdx].data {
		return lChildIdx
	}
	return rChildIdx
}

func (heap *MinHeap) HInsert(data int) {
	heap.count++
	idx := heap.count
	for ; idx != 1 && heap.Nodes[idx/2].data > data; {
		heap.Nodes[idx] = heap.Nodes[idx/2]
		idx /= 2
	}
	heap.Nodes[idx] = HeapNode{
		data: data,
	}
}

func (heap *MinHeap) HDelete() (deletedNode HeapNode) {
	idx := 1
	deletedNode = heap.Nodes[idx]
	heap.Nodes[idx] = heap.Nodes[heap.count]
	heap.Nodes[heap.count] = HeapNode{}
	heap.count--

	for {
		if idx * 2 > heap.count {
			break
		}
		childIdx := heap.getMinChildIdx(idx)
		if heap.Nodes[childIdx].data >= heap.Nodes[idx].data {
			break
		}
		agentNode := heap.Nodes[idx]
		heap.Nodes[idx] = heap.Nodes[childIdx]
		heap.Nodes[childIdx] = agentNode
		idx = childIdx
	}

	return
}

type MaxHeap struct {
	count int
	Nodes []HeapNode
}

func NewMaxHeap(queueLen int) MaxHeap {
	return MaxHeap{
		Nodes: make([]HeapNode, queueLen),
	}
}

func (heap *MaxHeap) getMaxChildIdx(parentIdx int) int {
	lChildIdx := parentIdx * 2
	rChildIdx := parentIdx * 2 + 1
	if rChildIdx > heap.count {
		return lChildIdx
	}
	if heap.Nodes[lChildIdx].data > heap.Nodes[rChildIdx].data {
		return lChildIdx
	}
	return rChildIdx
}

func (heap *MaxHeap) HInsert(data int) {
	heap.count++
	idx := heap.count
	for ; idx != 1 && data > heap.Nodes[idx/2].data; {
		heap.Nodes[idx] = heap.Nodes[idx/2]
		idx /= 2
	}
	heap.Nodes[idx] = HeapNode{
		data: data,
	}
}

func (heap *MaxHeap) HDelete() (deletedNode HeapNode) {
	idx := 1
	deletedNode = heap.Nodes[idx]
	heap.Nodes[idx] = heap.Nodes[heap.count]
	heap.Nodes[heap.count] = HeapNode{}
	heap.count--

	for {
		if idx * 2 > heap.count {
			break
		}
		childIdx := heap.getMaxChildIdx(idx)
		if heap.Nodes[idx].data >= heap.Nodes[childIdx].data {
			break
		}
		agentNode := heap.Nodes[idx]
		heap.Nodes[idx] = heap.Nodes[childIdx]
		heap.Nodes[childIdx] = agentNode
		idx = childIdx
	}

	return
}


func solution2(operations []string) []int {
	var opts [][]string
	isInsertFirst := false
	var insertCnt, deleteCnt int
	for _, operation := range operations {
		cmd := strings.Split(operation, " ")
		switch cmd[0] {
		case "I":
			isInsertFirst = true
			insertCnt++
		case "D":
			if !isInsertFirst {
				continue
			}
			deleteCnt++
		}
		opts = append(opts, cmd)
	}

	if deleteCnt >= insertCnt {
		return []int{0, 0}
	}

	minHeap := NewMinHeap(insertCnt + 1)
	maxHeap := NewMaxHeap(insertCnt + 1)
	popsMap := map[int]int{}
	for _, opt := range opts {
		switch opt[0] {
		case "I":
			insertInt, _ := strconv.Atoi(opt[1])
			minHeap.HInsert(insertInt)
			maxHeap.HInsert(insertInt)
		case "D":
			switch opt[1] {
			case "1":
				deletedNode := maxHeap.HDelete()
				popsMap[deletedNode.data]++
			case "-1":
				deletedNode := minHeap.HDelete()
				popsMap[deletedNode.data]++
			}
		}
	}

	var maxNode, minNode HeapNode
	for {
		maxNode = maxHeap.HDelete()
		if _, ok := popsMap[maxNode.data]; !ok { break }
	}
	for {
		minNode = minHeap.HDelete()
		if _, ok := popsMap[minNode.data]; !ok { break }
	}
	return []int{maxNode.data, minNode.data}
}
