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

type MaxHeap struct {
	count int
	Nodes []HeapNode
}

func NewMaxHeap(queueLen int) MaxHeap {
	return MaxHeap{
		Nodes: make([]HeapNode, queueLen),
	}
}
