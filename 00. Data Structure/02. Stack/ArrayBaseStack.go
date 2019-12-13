package main

import (
	"fmt"
)

const StackLen = 5
type Data interface {}

// stack 구조체 정의
type Stack struct {
	stackArr [StackLen]Data
	topIndex int
}

// 새로운 Stack 객체를 생성 및 초기화하고 그 객체의 주소 값을 반환해줌
func New() *Stack {
	return &Stack{
		stackArr: [StackLen]Data{},
		topIndex: -1,
	}
}

// 전달받은 리시버 변수인 Stack이 빈 값인지 확인해주는 메서드
func (ps *Stack) SIsEmpty() bool {
	if ps.topIndex == -1 {
		return true
	} else {
		return false
	}
}

// 잔딜 빋은 리시버 변수인 Stack이 꽉 찼는지 확인해주는 메서드
func (ps *Stack) SIsFull() bool {
	if ps.topIndex == StackLen-1 {
		return true
	} else {
		return false
	}
}

// 전달 받은 리시버 변수인 Stack에 새 값을 넣어주는 메서드
func (ps *Stack) SPush(data Data) {
	pushErrorHandler(func(ps *Stack, data Data) {
		if ps.SIsFull() {
			panic("Memory is FULL!!")
		} else {
			ps.topIndex++
			ps.stackArr[ps.topIndex] = data
		}
	})(ps, data)
}

// 전달 받은 리시버 변수인 Stack에서 가장 최근에 넣은 값을 삭제 후 반환하는 메서드
func (ps *Stack) SPop() Data {
	return popErrorHandler(func(ps *Stack) Data {
		if ps.SIsEmpty() {
			panic("Memory is EMPTY!!")
			return nil
		} else {
			rIdx := ps.topIndex
			ps.topIndex -= 1
			return ps.stackArr[rIdx]
		}
	})(ps)
}

// SPush 함수에 에러 처리 코드를 추가한 함수 반환
func pushErrorHandler(handler func(ps *Stack, data Data)) func(ps *Stack, data Data) {
	return func(ps *Stack, data Data) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()

		handler(ps, data)
	}
}

// SPop 함수에 에러 처리 코드를 추가한 함수 반환
func popErrorHandler(handler func(ps *Stack) Data) func (ps *Stack) Data {
	return func(ps *Stack) Data {
		defer func() {
			if err:=recover(); err != nil{
				fmt.Println(err)
			}
		}()

		return handler(ps)
	}
}

// 절달 받은 리시버 변수인 Stack의 값들을 최근에 넣은 순서대로 출력해주는 메서드
func (ps *Stack) SPrint() {
	for i:=ps.topIndex; i>=0; i-- {
		fmt.Print(ps.stackArr[i], " ")
	}
}


func main() {
	stack := New()

	stack.SPush("a")
	stack.SPush("ab")
	stack.SPush("abc")
	stack.SPush(1)
	stack.SPush(2)
	stack.SPush(3)

	stack.SPop()
	stack.SPop()

	stack.SPrint()
	// 1 abc ab a
}
