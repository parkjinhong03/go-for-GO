package algorithm2

import (
	"sort"
)

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
