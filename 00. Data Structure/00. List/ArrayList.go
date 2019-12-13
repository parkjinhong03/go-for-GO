package main

import (
	"errors"
	"fmt"
)

const ListLen = 100
type LData interface{}

// 배열기반 리스트를 정의한 구조체
type List struct {
	arr [ListLen]LData	// 리스트의 저장소인 배열
	numOfData int		// 저장된 데이터의 수
	curPosition int		// 데이터 참조위치를 기록
}

func NewList() *List {
	return &List{
		numOfData:	 0,
		curPosition: -1,
	}
}

func (plist *List) LInsert(data LData) {
	if plist.numOfData >= ListLen {
		fmt.Println("저장이 불가능합니다.")
		return
	}

	plist.arr[plist.numOfData] = data
	plist.numOfData++
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

func (plist *List) LRemove() (LData, error) {
	if plist.curPosition < 0 {
		return nil, errors.New("삭제가 불가능합니다")
	}
	rdata := plist.arr[plist.curPosition]

	for i:=plist.curPosition; i<plist.numOfData-1; i++ {
		plist.arr[i] = plist.arr[i+1]
	}

	if plist.curPosition != 0 {
		plist.curPosition--
	}
	plist.numOfData--
	return rdata, nil
}

func (plist List) LPrint() {
	var data LData

	fmt.Printf("현재 참조 위치: %d | 현재 데이터의 수: %d\n", plist.curPosition, plist.numOfData)

	if plist.LFirst(&data) {
		fmt.Print(data, " ")

		for {
			if nData := plist.LNext(&data); nData {
				fmt.Print(data, " ")
				continue
			}
			break
		}
		fmt.Println()
	}

}
