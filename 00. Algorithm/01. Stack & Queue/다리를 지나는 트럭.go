// https://programmers.co.kr/learn/courses/30/lessons/42583?language=go

package algorithm1

// 대기 트럭용 큐
type TructLinkedQueue struct {
	start *TructNode
	end *TructNode
	count int
}

type TructNode struct {
	weight int
	next *TructNode
}

func (queue *TructLinkedQueue) Enqueue(weight int) {
	node := TructNode{weight: weight}
	if queue.start == nil { queue.start = &node }
	if queue.end != nil { queue.end.next = &node }
	queue.end = &node
	queue.count++
}

func (queue *TructLinkedQueue) Dequeue() (deletedNode *TructNode) {
	deletedNode = queue.start
	queue.start = queue.start.next
	queue.count--
	return
}

// 다리 건너는 트럭용 큐
type AcrossLinkedQueue struct {
	start *AcrossNode
	current *AcrossNode
	end *AcrossNode
	count int
	acrossTime int
	totalWeight int
}

type AcrossNode struct {
	weight int
	remainTime int
	next *AcrossNode
}

func (queue *AcrossLinkedQueue) Enqueue(weight int) {
	node := AcrossNode{
		weight: weight,
		remainTime: queue.acrossTime,
	}
	if queue.start == nil { queue.start = &node }
	if queue.end != nil { queue.end.next = &node }
	queue.end = &node
	queue.count++
	queue.totalWeight += weight
}

func (queue *AcrossLinkedQueue) Dequeue() (deletedNode *AcrossNode) {
	deletedNode = queue.start
	queue.totalWeight -= queue.start.weight
	queue.start = queue.start.next
	queue.count--
	return
}

func (queue *AcrossLinkedQueue) GetNextNode() (node *AcrossNode) {
	if queue.current == nil {
		queue.current = queue.start
	}
	node = queue.current
	queue.current = queue.current.next
	return
}

func solution2(bridge_length int, weight int, truck_weights []int) (time int) {
	truckQueue := TructLinkedQueue{}
	for _, truck_weight := range truck_weights {
		truckQueue.Enqueue(truck_weight)
	}

	time = 1
	crossedTrucks := []int{}
	acrossQueue := AcrossLinkedQueue{acrossTime: bridge_length}
	for {
		if truckQueue.start != nil && (weight - acrossQueue.totalWeight) >= truckQueue.start.weight {
			acrossQueue.Enqueue(truckQueue.start.weight)
			truckQueue.Dequeue()
		}

		time++
		for i:=0; i<acrossQueue.count; i++ {
			acrossQueue.GetNextNode().remainTime--
		}
		if acrossQueue.start.remainTime == 0 {
			acrossed := acrossQueue.Dequeue()
			crossedTrucks = append(crossedTrucks, acrossed.weight)
		}

		if len(crossedTrucks) == len(truck_weights) {
			break
		}
	}

	return time
}
