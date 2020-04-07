# **내부 필드 접근**
## **내부 필드 접근 방식**
- 내부 필드에는 **. 연산자**로 접근한다
    ~~~go
    type Item struct {
        name string
        price float64
        quantity int
    }

    func (t Item) Cost() float64 {
        return t.price * float64(t.quantity)
    }

    func main() {
        var t Item
        t.name = "Men's Slip-Fit Shirt"
        t.price = 25000
        t.quantity = 3

        fmt.Println(t.name)     // Men's Slip-Fit Shirt 
        fmt.Println(t.price)    // 25000
        fmt.Println(t.quantity) // 3
        fmt.Println(t.Cost())   // 75000
    }
    ~~~

<br>

- **다른 구조체**를 **구조체의 내부 필드**로 지정하면 **지정한 필드명으로** 내부 구조체의 필드에 **접근할 수 있다.**
    ~~~go
    package main

    import (
        "fmt"
    )

    type dimension struct {
        width, height, length float64
    }

    type newItem struct {
        name string
        price float64
        quantity int
        packageDimension dimension
        itemDimension dimension
    }

    func main() {
        shoes := newItem{
            "Sports Shoes", 30000, 2,
            dimension{30, 270, 20},
            dimension{50, 300, 30},
        }

        // 내부 필드인 dimension 구조체의 값 출력
        fmt.Printf("%#v\n", shoes.itemDimension)
        fmt.Printf("%#v\n", shoes.packageDimension)

        // dimension 구조체의 내부 필드 값 출력
        fmt.Println(shoes.packageDimension.width)
        fmt.Println(shoes.packageDimension.height)
        fmt.Println(shoes.packageDimension.length)
    }
    ~~~
    ~~~
    실행 결과

    main.dimension{width:50, height:300, length:30}
    main.dimension{width:30, height:270, length:20}
    30
    270
    20
    ~~~
- **구조체 값을 출력**할 때 **필드명과 함께 출력**하려면 **%#v**를 사용하면 된다.

<br>

---
## **태그**
- 구조체 필드에 옵션으로 **태그(tag)를 정의**할 수 있다.
- 태그는 필드에 **추가된 문자열**이고, 필드에 **중요한 레이블**이나 **간단한 설명**을 추가할 때 유용하다.
- 태그는 **reflect.TypeOf() 함수**로 확인할 수 있다.
~~~go
package main

type Item struct {
    name        string  "상품 이름"
    price       float64 "상품 가격"
    quantity    int     "구매 수량"
}

func main() {
    tType := reflect.TypeOf(Item{})
    for i:=0; i<tType.NumField(); i++ {
        fmt.Println(tType.Field(i).Tag)
    }
}
~~~
~~~
실행 결과

상품 이름
상품 가격
구매 수량
~~~