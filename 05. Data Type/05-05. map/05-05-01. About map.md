# **map**
## **map이란?**
- **맵**(map)은 **키**(key)와 **값**(value)으로 이루어진 **테이블 형태의 컬렉션**이다.

- **키 타입**은 **비교 연산자**를 사용해서 비교할 수 있어야 한다.
    - **내장 타입**(int, float32, string), **배열**, **구조체**, **포인터** -> **키로 사용 가능**
    - **참조 타입**(슬라이스, 맵) -> 비교 연산자 사용 X -> **키로 사용 불가능**

- 하지만 **값**(value)에는 **모든 타입의 값**을 사용할 수 있다.

<br>

## **생성과 초기화**
- 맵은 **다음과 같은 방식**으로 생성한다.
    - **map[키타입]값타입{}**
    - **map[키타입]값타입{초깃값}**
    - **make(map[키타입]값타입, 초기용량)**
    - **make(map[키타입]값타입)**

- 다음 예제들은 **string 타입 키**와 **int 타입 값**을 갖는 **map**을 생성한다.
    ~~~go
    numberMap := map[string]int
    numberMap["one"] = 1
    numberMap["two"] = 2
    numberMap["three"] = 3
    fmt.Println(numberMap)
    ~~~
    ~~~go
    numberMap := map[string]int{
        "one": 1,
        "two": 2,
        "three": 3, // 요소를 여러 줄로 표기할 때 요소의 끝에 콤마(,)를 붙여야 함
    }
    fmt.Println(numberMap)
    ~~~
    ~~~go
    numberMap := make(map[string]int, 3) // 용량 생략 가능
    numberMap["one"] = 1
    numberMap["two"] = 2
    numberMap["three"] = 3
    fmt.Println(numberMap)
    ~~~
    ~~~
    실행 결과

    1, 3 번쨰 코드 -> map[one:1 two:2 three:3]

    2 번쨰 코드 -> map[three:3 one:1 two:2]
    ~~~

- **키로 값에 접근**할 때는 배열이나 슬라이스처럼 **인덱스 연산자**([])를 사용한다.

- **맵의 요소는 정렬되어 있지 않으므로** 각 요소에 어떤 순서로 접근할지 예측할 수 없다.
    ~~~go
    numberMap := make(map[string]int) // 용량 생략 가능
    numberMap["one"] = 1
    numberMap["two"] = 2
    numberMap["three"] = 3

    for k, v := range numberMap {
        fmt.Println(k, v)
    }
    ~~~
    ~~~
    실행 결과

    one 1
    three 3
    two 2
    ~~~

<br>

## *Note*
- **슬라이스**는 **참조 타입**이므로 맵의 키로 **사용할 수 없다.**
- 하지만 **[]byte** 타입과 **[]int32** 타입은 **string으로 변환**하면 **맵의 키**로 사용할 수 있다.
    ~~~go
    groupMap := make(map[string]string)

	group1 := []int32{1, 2, 3}
	group2 := []int32{'a', 'b', 'c'}
	group3 := []int32{7, 8, 'd'}

	groupMap[string(group1)] = "first"
	groupMap[string(group2)] = "second"
	groupMap[string(group3)] = "third"

	fmt.Println(groupMap[string(group1)])
	fmt.Println(groupMap[string(group2)])
	fmt.Println(groupMap[string(group3)])
    ~~~
    ~~~
    실행 결과

    first
    second
    third
    ~~~
- 문자열은 **유니코드 문자의 코드값**을 정수로 표현한 값의 **시퀀스**이므로, **[]int32 타입**을 **문자열**로 **변환**할 수 있다.