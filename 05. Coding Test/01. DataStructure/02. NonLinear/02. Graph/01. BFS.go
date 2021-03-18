// BFS는 그래프 탐색(하나의 정점으로부터 시작하여 차례대로 모든 정점들을 한 번씩 방문하는 탐색) 중 한 방법
// 루트 노드 혹은 다른 임의의 노드에서 시작해서 인접한 노드들을 먼저 탐색함 (넓게 O, 깊게 X)
// 두 노드 사이의 최단 경로 혹은 임의의 경로를 찾으려할 때 사용한다.

// 방문할 노드들을 Stack이 아닌 Queue에 저장하는 이유는? 
//  -> 먼저 들어온 노드를 먼저 방문해야 노드의 단계들을 순차적으로 방문할 수 있기 때문

package main

import (
	"fmt"
)

func main() {
	// 노드 개수가 5개인 그래프 선언
	g := Graph(5)

	// from 번째 노드부터 to 번째 노드들 까지의 간선 추가
	for from, to := range map[int][]int {
		0: {1, 2, 4},
		1: {0, 2},
		2: {0, 1, 3, 4},
		3: {2},
		4: {0, 2},
	} {
		g.AddEdge(from, to)
	}
	
	fmt.Println(g.BFS(0)) // [0 1 2 4 3]
}

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

// from 번째의 노드에서 to 번째 노드들 사이의 간선 추가
func (g *graph) AddEdge(from int, to []int) {
	g.adj[from] = append(g.adj[from], to...)
}

// start 번째 노드를 시작으로 BFS 방식으로 탐색하여 탐색 순서(seq) 반환
func (g *graph) BFS(start int) (seq []int) {
	// 특정 노드를 이미 탐색하였는지 저장 및 확인하기 위한 배열
	// nodeNum: 3 -> visited 초깃값: [false false false]
	visited := make([]bool, g.nodeNum)
	
	// 다음으로 탐색할 노드들을 저장해두기 위한 Queue 객체 선언
	q := IntQueue()

	// 시작 노드를 확인했다는 것을 등록하고, Queue에 등록하여 다음 탐색 순서에 추가
	visited[start] = true
	q.Push(start)

	// Queue 요소가 비었을 때 까지 반복
	for q.Size() != 0 {
		// 이번 for문에서 인접 정점들을 탐색한 정점(노드)을 Queue에서 받음
		// 첫 for문이라면 위에서 삽입한 start 값을 받음
		current := q.Pop()

		// 해당 노드와 연결된 모든 정점들을 다음 탐색 순서에 등록하는 for문
		// 여기서 현재 정점과 연결된 같은 깊이의 정점들을 다음 탐색 순서에 등록함으로써 '너비 우선 탐색'이 되는거임
		for _, linked := range g.adj[current] {
			// 이미 처리한 인접 정점이라면 이번 for문 중단
			if visited[linked] {
				continue
			}

			// 방문하였다는 것을 기록하고 Queue를 통해 다음 탐색 순서에 등록
			visited[linked] = true
			q.Push(linked)
		}

		// 탐색 순서 리스트에 이번 for문에서 처리한 정점 추가
		seq = append(seq, current)
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
