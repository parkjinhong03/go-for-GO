package algorithm2

import "sort"

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

func solution(jobs [][]int) (avg int) {
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i][0] == jobs[j][0] && jobs[i][1] <= jobs[j][1]
	})

	heap := NewJobMinHeap(len(jobs) + 1)
	totalJobCnt := len(jobs)
	var currentRemainTime, totalTime, finishedCnt int
	var currentJobNode JobNode
	time := -1

	jobsMap := map[int][]int{}
	for _, job := range jobs {
		if _, ok := jobsMap[job[0]]; !ok {
			jobsMap[job[0]] = []int{}
		}
		jobsMap[job[0]] = append(jobsMap[job[0]], job[1])
	}

	for {
		time++

		if jobsArr, ok := jobsMap[time]; ok {
			for _, job := range jobsArr {
				heap.HInsert(time, job)
			}
			jobs = jobs[len(jobsArr):]
		}

		if currentRemainTime != 0 {
			currentRemainTime--
			if currentRemainTime == 0 {
				totalTime += time - currentJobNode.comeTime
				finishedCnt++
			}
		}

		if finishedCnt == totalJobCnt {
			break
		}

		if currentRemainTime == 0 && heap.NodeNum != 0 {
			execNode := heap.HDelete()
			currentRemainTime = execNode.workTime
			currentJobNode = execNode
		}
	}

	return totalTime/totalJobCnt
}
