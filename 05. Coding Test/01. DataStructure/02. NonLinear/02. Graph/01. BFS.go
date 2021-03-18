// BFS는 그래프 탐색(하나의 정점으로부터 시작하여 차례대로 모든 정점들을 한 번씩 방문하는 탐색) 중 한 방법
// 루트 노드 혹은 다른 임의의 노드에서 시작해서 인접한 노드들을 먼저 탐색함 (넓게 O, 깊게 X)
// 두 노드 사이의 최단 경로 혹은 임의의 경로를 찾으려할 때 사용한다.

// 방문할 노드들을 Stack이 아닌 Queue에 저장하는 이유는?
//  -> 먼저 들어온 노드를 먼저 방문해야 노드의 단계들을 순차적으로 방문할 수 있기 때문

package main

// 그래프 구현 객체
type graph struct {
	// 노드의 총 갯수
	nodeNum int

	// 노드 간의 간선을 표시하는 인접 리스트
	// Ex) nodeNum -> 3, adj -> [0: [1], 1: [0, 2], 2: [1]]
	// 0번째 노드: 1번째 노드와 연결
	// 1번째 노드: 0, 2번쨰 노드와 연결
	// 2번째 노드: 1번째 노드와 연결
	//     0
	//   / 
	//  1  -  2
	adj [][]int // adjacent list
}

func Graph(num int) (g *graph) {
	g = &graph{
		nodeNum: num,
		adj: [][]int{},
	}
	for i:=0; i<num; i++ {
		g.adj = append(g.adj, []int{})
	}
	return
}

// Int 타입 전용 Queue 구현 객체
type intQueue []int

func IntQueue() *intQueue {
	return &intQueue{}
}

// 요소 삽입
func (q *intQueue) Push(i int) {
	*q = append(*q, i)
}

// 첫 번째 요소 추출
func (q *intQueue) Pop() (i int) {
	i, *q = (*q)[0], (*q)[1:]
	return
}

// Queue 요소 갯수 반환
func (q *intQueue) Size() int {
	return len(*q)
}
