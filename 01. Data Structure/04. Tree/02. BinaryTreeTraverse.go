package main

import "fmt"

type BTTData interface {}

type VisitFunc func(BTTData)

// 데이터와 왼쪽 자식 노드, 오른쪽 자식 노드의 주소값을 가지고 있는 구조체
type BTTreeNode struct {
	data BTTData
	left *BTTreeNode
	right *BTTreeNode
}

// 모든 값이 null인 이진 트리 노드 생성 및 주소값 반환
func MakeBTTreeNode() *BTTreeNode {
	return &BTTreeNode{
		data:  nil,
		left:  nil,
		right: nil,
	}
}

// 인자 값으로 받은 이진 트리 노드의 데이터 반환
func GetTData(bt *BTTreeNode) BTTData {
	return bt.data
}

// 인자 값으로 받은 이진 트리 노드에 데이터를 저장시킴
func SetTData(bt *BTTreeNode, data BTTData) {
	bt.data = data
}

// 인자 값으로 받은 이진 트리 노드의 왼쪽 자식 노드의 주소값을 반환
func TGetLeftSubTree(bt *BTTreeNode) *BTTreeNode {
	return bt.left
}

// 인자 값으로 받은 이진 트리 노드의 오른쪽 자식 노드의 주소값을 반환
func TGetRightSubTree(bt *BTTreeNode) *BTTreeNode {
	return bt.right
}

// 두 번째의 인자 값으로 받은 이진 트리 노드를 첫 번째의 인자 값으로 받은 이진 트리 노드의 왼쪽 서브 트리로 연결시킴
func TMakeLeftSubTree(main *BTTreeNode, sub *BTTreeNode) {
	main.left = sub
}

// 두 번째의 인자 값으로 받은 이진 트리 노드를 첫 번째의 인자 값으로 받은 이진 트리 노드의 오른쪽 서브 트리로 연결시킴
func TMakeRightSubTree(main *BTTreeNode, sub *BTTreeNode) {
	main.right = sub
}

// 이진 트리 전체를 중위 순회한 결과를 출력하는 함수
func InorderTraverse(bt *BTTreeNode, action VisitFunc) {
	if bt == nil {
		return
	}

	InorderTraverse(bt.left, action)
	action(bt.data)
	InorderTraverse(bt.right, action)
}

// 이진 트리 전체를 전위 순회한 결과를 출력하는 함수
func PreorderTraverse(bt *BTTreeNode, action VisitFunc) {
	if bt == nil {
		return
	}

	action(bt.data)
	PreorderTraverse(bt.left, action)
	PreorderTraverse(bt.right, action)
}

// 이진 트리 전체를 후위 순회한 결과를 출력하는 함수
func PostorderTraverse(bt *BTTreeNode, action VisitFunc) {
	if bt == nil {
		return
	}

	PostorderTraverse(bt.left, action)
	PostorderTraverse(bt.right, action)
	action(bt.data)
}

func main() {
	bt1 := MakeBTTreeNode()
	bt2 := MakeBTTreeNode()
	bt3 := MakeBTTreeNode()
	bt4 := MakeBTTreeNode()
	bt5 := MakeBTTreeNode()
	bt6 := MakeBTTreeNode()

	SetTData(bt1, 1)
	SetTData(bt2, 2)
	SetTData(bt3, 3)
	SetTData(bt4, 4)
	SetTData(bt5, 5)
	SetTData(bt6, 6)

	TMakeLeftSubTree(bt1, bt2)
	TMakeRightSubTree(bt1, bt3)
	TMakeLeftSubTree(bt2, bt4)
	TMakeRightSubTree(bt2, bt5)
	TMakeRightSubTree(bt3, bt6)

	PreorderTraverse(bt1, ShowIntData); fmt.Println()
	// 1 2 4 5 3 6
	InorderTraverse(bt1, ShowIntData); fmt.Println()
	// 4 2 5 1 3 6
	PostorderTraverse(bt1, ShowIntData); fmt.Println()
	// 4 5 2 6 3 1
}

func ShowIntData(data BTTData) {
	fmt.Print(data, " ")
}