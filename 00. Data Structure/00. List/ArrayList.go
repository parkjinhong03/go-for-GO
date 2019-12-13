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

// 배열 기반 리스트인 List 객체 생성 및 초기화 후 반환
func NewList() *List {
	return &List{
		numOfData:	 0,
		curPosition: -1,
	}
}

// List 객체의 배열 arr 필드의 맨 뒤에 새로운 값 추가
func (plist *List) LInsert(data LData) {
	if plist.numOfData >= ListLen {
		fmt.Println("저장이 불가능합니다.")
		return
	}

	plist.arr[plist.numOfData] = data
	plist.numOfData++
}

// List 객체의 현재 참조 위치를 맨 앞으로 이동시킨 후 메개변수로 받은 pdata에 참조된 값 저장
func (plist *List) LFirst(pdata *LData) bool {
	if plist.numOfData == 0 {
		return false
	}

	plist.curPosition = 0
	*pdata = plist.arr[plist.curPosition]
	return true
}

// List 객체의 현재 참조 위치를 다음 요소로 변경시킨 후 매개변수로 받은 pdata에 참조된 값 저장
func (plist *List) LNext(pdata *LData) bool {
	if plist.curPosition >= plist.numOfData-1 {
		return false
	}

	plist.curPosition ++
	*pdata = plist.arr[plist.curPosition]
	return true
}

// List 객체의 현재 참조된 값 삭제 및 재정렬 후 삭제한 값 반환
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

// List 객체의 현재 참조 위치, 데이터 갯수, 데이터를 출력
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

func main() {
	list := NewList()
	data := new(LData)

	list.LInsert(1)
	list.LInsert(2)
	list.LInsert(3)

	list.LPrint()
	// 현재 참조 위치: -1 | 현재 데이터의 수: 3
	// 1 2 3

	if _, err:=list.LRemove(); err != nil {
		fmt.Println(err)
		// 삭제가 불가능합니다
	}

	list.LFirst(data)
	if _, err:=list.LRemove(); err != nil {
		fmt.Println(err)
	}

	list.LPrint()
	// 현재 참조 위치: 0 | 현재 데이터의 수: 2
	// 2 3
}
