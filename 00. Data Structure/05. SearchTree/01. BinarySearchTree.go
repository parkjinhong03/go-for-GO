package main

import (
	BTree "../04. Tree"
	"fmt"
)

type BSTData BTree.BTTData

// BST 객체 생성 및 초기화
func BSTMakeAndInit(pRoot **BTree.BTTreeNode) {
	*pRoot = nil
}

// 노드에 저장된 데이터 반환
func BSTGetNodeData(bst *BTree.BTTreeNode) BSTData {
	return BTree.GetTData(bst)
}

// BST를 대상으로 데이터 저장 (노드의 생성과정 포함)
func BSTInsert(pRoot **BTree.BTTreeNode, data BSTData) {
	pNode := BTree.MakeBTTreeNode() // parent node
	cNode := *pRoot
	nNode := BTree.MakeBTTreeNode() // new node

	// 새로운 노드가 추가될 위치를 적절한 찾는다.
	for cNode!=nil {
		if data == BTree.GetTData(cNode) {
			return // 만약 키 값이 중복된다면 삽입을 종료시킴
		}

		pNode = cNode

		if BTree.GetTData(cNode).(int) > data.(int) {
			cNode = BTree.TGetLeftSubTree(cNode)
		} else {
			cNode = BTree.TGetRightSubTree(cNode)
		}
	}

	// pNode의 자식 노드로 추가할 새 노드의 값 설정
	BTree.SetTData(nNode, data)

	// pNode의 자식 노드로 추가할 새 노드를 추가
	if pNode!=nil { // 새 노드가 루트 노드가 아니라면
		if data.(int) < BTree.GetTData(pNode).(int) {
			BTree.TMakeLeftSubTree(pNode, nNode)
		} else {
			BTree.TMakeRightSubTree(pNode, nNode)
		}
	} else { // 새 노드가 루트 노드라면
		*pRoot = nNode
	}
}

// BST를 대상으로 데이터 탐색
func BSTSearch(bst *BTree.BTTreeNode, target BSTData) *BTree.BTTreeNode {

}
