package main

import (
	"fmt"
)

const StackLen = 5
type Data interface {}

type Stack struct {
	stackArr [StackLen]Data
	topIndex int
}

func New() *Stack {
	return &Stack{
		stackArr: [StackLen]Data{},
		topIndex: -1,
	}
}

func (ps *Stack) SIsEmpty() bool {
	if ps.topIndex == -1 {
		return true
	} else {
		return false
	}
}

func (ps *Stack) SIsFull() bool {
	if ps.topIndex == StackLen-1 {
		return true
	} else {
		return false
	}
}

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
