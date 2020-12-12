// https://programmers.co.kr/learn/courses/30/lessons/42586?language=go

package algorithm1

type PercentLinkedQueue struct {
	start *PercentNode
	current *PercentNode
	end *PercentNode
	count int
}

type PercentNode struct {
	percent int
	next *PercentNode
}

func (queue *PercentLinkedQueue) Enqueue(percent int) {
	newNode := PercentNode{
		percent: percent,
	}
	if queue.start == nil {
		queue.start = &newNode
	}
	if queue.end != nil {
		queue.end.next = &newNode
	}
	queue.end = &newNode
	queue.count++
}

func (queue *PercentLinkedQueue) Dequeue() {
	queue.start = queue.start.next;
	queue.count--
}

func (queue *PercentLinkedQueue) InitCurrentNode() {
	queue.current = nil
}

func (queue *PercentLinkedQueue) GetNextNode() (node *PercentNode) {
	if queue.current == nil {
		queue.current = queue.start
	}
	node = queue.current
	queue.current = queue.current.next
	return
}

func solution(progresses []int, speeds []int) (results []int) {
	queue := &PercentLinkedQueue{}
	for _, progress := range progresses {
		queue.Enqueue(progress)
	}

	for {
		for i:=0; i<queue.count; i++ {
			queue.GetNextNode().percent += speeds[i]
		}

		done := 0
		for {
			if queue.count == 0 {
				break
			} else if queue.GetNextNode().percent >= 100 {
				queue.Dequeue()
				done++
			} else {
				break
			}
		}
		queue.InitCurrentNode()

		if done != 0 {
			results = append(results, done)
		}
		if queue.count == 0 {
			break
		}
	}

	return
}
