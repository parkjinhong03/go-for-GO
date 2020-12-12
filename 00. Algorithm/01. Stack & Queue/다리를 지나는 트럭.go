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

