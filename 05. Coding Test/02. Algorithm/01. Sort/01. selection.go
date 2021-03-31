// 제자리 정렬 (입력 배열 외에는 다른 메모리 요구 X)
// 해당 순서에 원소를 넣을 위치는 정해져 있고, 어떤 원소를 넣을지 선택하는 알고리즘
// 시간복잡도 -> (n-1) + (n-2) + … + 2 + 1 = n(n-1)/2 = O(n^2)

package main

import (
	"fmt"
)

func main() {
	fmt.Println(selectionSort([]int{9, 6, 7, 3, 5}))
}

func selectionSort(x []int) []int {
	for idx:=0; idx<len(x)-1; idx++ {
		var leastIdx = idx
		
		// 현재 인덱스 이후의 요소들 중 가장 작은 요소의 인덱스 검색
		for compareIdx:=idx+1; compareIdx<len(x); compareIdx++ {
			if x[compareIdx] < x[leastIdx] {
				leastIdx = compareIdx
			}
		}

		// 현재 인덱스와 가장 작은 요소의 인덱스가 다르다면 두 인덱스의 값 교체
		if leastIdx != idx {
			x[idx], x[leastIdx] = x[leastIdx], x[idx]
		}
	}

	return x
}
