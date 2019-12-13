package main

import "fmt"

const ListLen = 100
type LData interface{}

// 배열기반 리스트를 정의한 구조체
type List struct {
	arr [ListLen]LData	// 리스트의 저장소인 배열
	numOfData int		// 저장된 데이터의 수
	curPosition int		// 데이터 참조위치를 기록
}

func (plist *List) NewList() *List {
	return &List{
		numOfData:	 0,
		curPosition: -1,
	}
}

func (plist *List) LInsert(data LData) {
	if plist.numOfData >= ListLen {
		fmt.Println("저장이 불가능합니다.")
	}

	plist.arr[plist.numOfData] = data
	plist.curPosition ++
}

func (plist *List) LFirst(pdata *LData) bool {
	if plist.numOfData == 0 {
		return false
	}

	plist.curPosition = 0
	*pdata = plist.arr[plist.curPosition]
	return true
}

func (plist *List) LNext(pdata *LData) bool {
	if plist.curPosition >= plist.numOfData-1 {
		return false
	}

	plist.curPosition ++
	*pdata = plist.arr[plist.curPosition]
	return true
}

func (plist *List) LRemove() LData {
	var rdata LData = plist.arr[plist.curPosition]

	for i:=plist.curPosition; i<plist.numOfData; i++ {
		plist.arr[i] = plist.arr[i+1]
	}

	plist.numOfData--
	plist.curPosition--
	return rdata
}

func (plist *List) LCount() int {
	return plist.numOfData
}
