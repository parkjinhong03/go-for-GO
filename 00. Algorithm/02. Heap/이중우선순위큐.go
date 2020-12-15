package algorithm2

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
