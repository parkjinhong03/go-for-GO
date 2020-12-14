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

func (heap *JobMinHeap) GetMinChildIdx(parentIdx int) int {
	lChildIdx := parentIdx * 2
	rChildIdx := parentIdx * 2 + 1
	if rChildIdx > heap.NodeNum {
		return lChildIdx
	}
	if heap.Nodes[lChildIdx].workTime > heap.Nodes[rChildIdx].workTime {
		return rChildIdx
	}
	return lChildIdx
}

func (heap *JobMinHeap) HDelete() (deletedNode JobNode) {
	idx := 1
	deletedNode = heap.Nodes[idx]
	heap.Nodes[idx] = heap.Nodes[heap.NodeNum]
	heap.Nodes[heap.NodeNum] = JobNode{}
	heap.NodeNum--

	for {
		if idx * 2 > heap.NodeNum {
			break
		}
		childIdx := heap.GetMinChildIdx(idx)
		if heap.Nodes[childIdx].workTime >= heap.Nodes[idx].workTime {
			break
		}
		agentNode := heap.Nodes[idx]
		heap.Nodes[idx] = heap.Nodes[childIdx]
		heap.Nodes[childIdx] = agentNode
		idx = childIdx
	}
	return
}