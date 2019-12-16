package main

import "fmt"

type BTData interface {}

type VisitFunc func(BTData)

// 데이터와 왼쪽 자식 노드, 오른쪽 자식 노드의 주소값을 가지고 있는 구조체
type BTreeNode struct {
	data BTData
	left *BTreeNode
	right *BTreeNode
}

// 모든 값이 null인 이진 트리 노드 생성 및 주소값 반환
func MakeBTreeNode() *BTreeNode {
	return &BTreeNode{
		data:  nil,
		left:  nil,
		right: nil,
	}
}

// 인자 값으로 받은 이진 트리 노드의 데이터 반환
func GetData(bt *BTreeNode) BTData {
	return bt.data
}

// 인자 값으로 받은 이진 트리 노드에 데이터를 저장시킴
func SetData(bt *BTreeNode, data BTData) {
	bt.data = data
}

// 인자 값으로 받은 이진 트리 노드의 왼쪽 자식 노드의 주소값을 반환
func GetLeftSubTree(bt *BTreeNode) *BTreeNode {
	return bt.left
}

// 인자 값으로 받은 이진 트리 노드의 오른쪽 자식 노드의 주소값을 반환
func GetRightSubTree(bt *BTreeNode) *BTreeNode {
	return bt.right
}

// 두 번째의 인자 값으로 받은 이진 트리 노드를 첫 번째의 인자 값으로 받은 이진 트리 노드의 왼쪽 서브 트리로 연결시킴
func MakeLeftSubTree(main *BTreeNode, sub *BTreeNode) {
	main.left = sub
}

// 두 번째의 인자 값으로 받은 이진 트리 노드를 첫 번째의 인자 값으로 받은 이진 트리 노드의 오른쪽 서브 트리로 연결시킴
func MakeRightSubTree(main *BTreeNode, sub *BTreeNode) {
	main.right = sub
}

func InorderTraverse(bt *BTreeNode, action VisitFunc) {
	if bt == nil {
		return
	}

	InorderTraverse(bt.left, action)
	action(bt.data)
	InorderTraverse(bt.right, action)
}

func PreorderTraverse(bt *BTreeNode, action VisitFunc) {
	if bt == nil {
		return
	}

	action(bt.data)
	PreorderTraverse(bt.left, action)
	PreorderTraverse(bt.right, action)
}

func PostorderTraverse(bt *BTreeNode, action VisitFunc) {
	if bt == nil {
		return
	}

	PostorderTraverse(bt.left, action)
	PostorderTraverse(bt.right, action)
	action(bt.data)
}

func main() {
	bt1 := MakeBTreeNode()
	bt2 := MakeBTreeNode()
	bt3 := MakeBTreeNode()
	bt4 := MakeBTreeNode()
	bt5 := MakeBTreeNode()
	bt6 := MakeBTreeNode()

	SetData(bt1, 1)
	SetData(bt2, 2)
	SetData(bt3, 3)
	SetData(bt4, 4)
	SetData(bt5, 5)
	SetData(bt6, 6)

	MakeLeftSubTree(bt1, bt2)
	MakeRightSubTree(bt1, bt3)
	MakeLeftSubTree(bt2, bt4)
	MakeRightSubTree(bt2, bt5)
	MakeRightSubTree(bt3, bt6)

	PreorderTraverse(bt1, ShowIntData); fmt.Println()
	// 1 2 4 5 3 6
	InorderTraverse(bt1, ShowIntData); fmt.Println()
	// 4 2 5 1 3 6
	PostorderTraverse(bt1, ShowIntData); fmt.Println()
	// 4 5 2 6 3 1
}

func ShowIntData(data BTData) {
	fmt.Print(data, " ")
}