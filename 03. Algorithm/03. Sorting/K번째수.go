package algorithm3

import "sort"

func solution(array []int, commands [][]int) (results []int) {
	for _, command := range commands {
		slicedArr := append([]int{}, array[command[0]-1:command[1]]...)
		sort.Ints(slicedArr)
		results = append(results, slicedArr[command[2]-1])
	}
	return
}
