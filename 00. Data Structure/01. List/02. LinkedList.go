package main

import (
	"fmt"
)

// 시작 Node와 끝 Node를 가지는 List 구조체
type List struct {
	head *Node
	tail *Node
}

// value와 다음 Node의 주솟값을 가지는 각 노드를 표현할 구조체
type Node struct {
	data int
	next *Node
}
