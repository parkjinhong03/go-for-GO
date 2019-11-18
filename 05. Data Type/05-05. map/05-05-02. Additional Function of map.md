# **map의 추가적인 기능**
## **값 찾기**
- **인덱스 연산자**([])로 특정 키에 대한 값을 얻을 수 있다.

- 키가 없을 때는 값으로 **지정한 타입의 제로값**을 반환한다.
    ~~~go
    numberMap := map[string]int{}
    numberMap["zero"] = 0
    numberMap["one"] = 1
    numberMap["two"] = 2
    fmt.Println(numberMap["zero"])  // 0
    fmt.Println(numberMap["one"])   // 1
    fmt.Println(numberMap["two"])   // 2
    fmt.Println(numberMap["three"]) // 0
    ~~~
- 이 예제에서는 **키가 "tree"인 요소**의 값이 **실제로 0**인지, 아니면 **없어서 제로값**을 반환했는지 알 수없다.
- 따라서 []연산자의 **두 번째 매개변수인 bool 값**을 이용해서 확일할 수 있다. 
    ~~~go
    if v, ok := numberMap["three"]; ok {
        fmt.Println("'three' is in numberMap. value:", v)
    } else {
        fmt.Println("'three' is not in numberMap")
    }
    ~~~

<br>

---
## **요소 추가, 수정, 삭제**
- 맵에 새 **키/값을 추가**하는 구문은 **기존 키의 값을 수정**하는 구문과 **같다.**

- 둘 다 **[] 연산자**를 사용한다.(myM ap[key] = value)
    - 키에 해당하는 요소가 없으면 **새로운 요소로 추가**됨
    - 키에 해당하는 요소가 있으면 **기존 값을 수정**함

    ~~~go
    numberMap = map[int]string{}
    numberMap[1] = "one"
    numberMap[2] = "two"
    fmt.Println(numberMap)  // map[1:one 2:two]

    numberMap[3] = "three"
    fmt.Println(numberMap)  // map[1:one 2:two 3:three]
    numberMap[3] = "삼"
    fmt.Println(numberMap)  // map[1:one 2:two 3:삼]
    ~~~

- 요소를 제거할 때는 **delete() 함수**를 사용한다.
    ~~~go
    delete(numberMap, 3)
    fmt.Println(numberMap)  //map[1:one 2:two]
    ~~~