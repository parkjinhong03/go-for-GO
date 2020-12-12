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
