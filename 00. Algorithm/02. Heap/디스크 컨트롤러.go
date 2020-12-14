package algorithm2

type JobMinHeap struct {
	Nodes []JobNode
	NodeNum int
}

type JobNode struct {
	comeTime int // 요청 시점
	workTime int // 작업 시간 (비교 기준)
}

func NewJobMinHeap(idxLen int) JobMinHeap {
	return JobMinHeap{
		Nodes: make([]JobNode, idxLen),
	}
}

func (heap *JobMinHeap) HInsert(comeTime, workTime int) {
	heap.NodeNum++
	insertIdx := heap.NodeNum
	for ; insertIdx!=1 && heap.Nodes[insertIdx/2].workTime >= workTime; {
		heap.Nodes[insertIdx] = heap.Nodes[insertIdx/2]
		insertIdx /= 2
	}
	heap.Nodes[insertIdx] = JobNode{
		comeTime: comeTime,
		workTime: workTime,
	}
}
