// https://programmers.co.kr/learn/courses/30/lessons/42579?language=go

package algorithm0

import (
	"sort"
)

type Count struct {
	Genre string
	Play int
}

type Count2 struct {
	Num int
	Play int
}

func solution2(genres []string, plays []int) (results []int) {
	// map에 정리
	countMap := map[string]int{}
	playMap := map[string][]int{}
	for index, play := range plays {
		if _, ok := countMap[genres[index]]; !ok {
			countMap[genres[index]] = 0
		}
		if _, ok := playMap[genres[index]]; !ok {
			playMap[genres[index]] = []int{}
		}
		countMap[genres[index]] += play
		playMap[genres[index]] = append(playMap[genres[index]], index)
	}

	// 정렬
	counts := []Count{}
	for key, value := range countMap {
		counts = append(counts, Count{key, value})
	}
	sort.Slice(counts, func(i int, j int) bool {
		return counts[i].Play >= counts[j].Play
	})

	// 추출
	for _, count := range counts {
		playCounts := []Count2{}
		for _, num := range playMap[count.Genre] {
			playCounts = append(playCounts, Count2{num, plays[num]})
		}
		sort.Slice(playCounts, func(i int, j int) bool {
			return playCounts[i].Play > playCounts[j].Play
		})
		for index, playCount := range playCounts {
			if index >= 2 { break }
			results = append(results, playCount.Num)
		}
	}

	return
}
