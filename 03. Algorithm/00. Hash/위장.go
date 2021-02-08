// https://programmers.co.kr/learn/courses/30/lessons/42578?language=go

package algorithm0

func solution(clothes [][]string) (cases int) {
	myClothes := map[string][]string{}
	for _, cloth := range clothes {
		if _, ok := myClothes[cloth[1]]; !ok {
			myClothes[cloth[1]] = []string{}
		}
		myClothes[cloth[1]] = append(myClothes[cloth[1]], cloth[0])
	}

	cases = 1
	for _, clothes := range myClothes {
		cases *= len(clothes) + 1
	}
	cases--

	return
}
