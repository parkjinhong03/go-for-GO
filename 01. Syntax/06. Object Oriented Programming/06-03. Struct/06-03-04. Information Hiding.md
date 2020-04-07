# **정보 은닉**
## **Go에서의 접근 제어**
- **04-02장**에 나왔듯이, Go는 식별자의 **첫 번째 문자의 대소문자**로 **public**과 **private**을 구분한다.

- 마찬가지로 **구조체의 필드**도 **대소문자**로 구분한다.
    - **대문자로 시작하는 필드**(exported field)     -> **패키지 외부**에서 접근 가능
    - **소문자로 시작하는 필드**(non-exported field) -> **패키지 내부**에서만 접근 가능

- 대부분의 클래스 기반 객체 지향 언어에서는 **지정된 방식**으로 **내부 필드에 접근**하게 된다.
    - **getter**로 내부 정보를 **얻어**오고 **setter**로 내부 필드의 값을 **변경**한다.
    - 이렇게 하면 객체에 **잘못된 값을 할당**하거나 **비공개 정보가 노출**되는 것을 막을 수 있다.
    - 따라서 **Go**에서도 이와 같은 방식으로 **내부 정보를 보호한다.**

<br>

---
## **생성 함수**
- **구조체**를 생성할 때 초깃값을 지정하지 않으면 **제로값으로 초기화**한다.
- 하지만 종종 규칙에 따라 **초깃값**을 **원하는 값**으로 지정해야 할 때가 있다.
- **Go**의 구조체는 **생성자**를 지원하지 않지만, **구조체 생성을 위한 함수**를 만들어 대체할 수 있다.
    ~~~go
    type Item struct {
        name string
        price float64
        quantity int
    }

    func NewItem(name string, price float64, quantity int) *Item {
        if price <= 0 || quantity <= 0 || len(name) == 0 {
            return nil
        }
        return &Item{name, price, quantity}
    }
    ~~~

<br>

- **Item 타입**은 **public**로, **내부 필드**는 외부에서 값을 수정하지 못하는 **private**로 정의하였다.
- 그리고 **항상 유효한 *Item**이 생성되도록 **NewItem() 함수**를 제공한다.

<br>

### ***Go 코드 컨벤션***
- **생성 함수의 이름**은 **New**로 시작하도록 하고, 패키지명과 타입명이 **같을 때**는 **New()** 로 지정한다.

- **예시 코드**
    - **item 패키지**
        ~~~go
        package item

        type Item struct {
            name string
            price float64
            quantity int
        }

        func New(name string, price float64, quantity int) *Item {
            if price <= 0 || quantity <= 0 || len(name) == 0 {
                return nil
            }
            return &Item{name, price, quantity}
        }
        ~~~
    - **main 패키지**
        ~~~go
        package main

        import (
            "fmt"
            "/lib/item"
        )

        func main() {
            shirts := item.New("Men's Slim-Fit Shirts", 25000, 3)
            ...
        }
        ~~~

<br>

---
## **getter와 setter**
- **private** 필드에 **getter**와 **setter**를 제공하면 **항상 유효한 값**이 할당되게 할 수 있다.

<br>

**위의 item 패키지 코드에 추가**
~~~go

func (t *Item) Name() string {
	return t.name
}

func (t *Item) SetName(n string) {
	if len(n) != 0 {
		t.name = n
	}
}

func (t *Item) Price() float64 {
	return t.price
}

func (t *Item) SetPrice(p float64) {
	if p > 0 {
		t.price = p
	}
}

func (t *Item) Quantity() int {
	return t.quantity
}

func (t *Item) SetQuantity(q int) {
	if q > 0 {
		t.quantity = q
	}
}
~~~

**getter, setter 사용**
~~~go
shirts.SetPrice(30000)
shirts.SetQuantity(2)

fmt.Println("name:", shirts.Name())
fmt.Println("price:", shirts.Price())
fmt.Println("quantity:", shirts.Quantity())
~~~

<br>

### ***Go 코드 컨벤션***
- **getter** 메서드는 보통 **필드명**과 같은 이름으로 짓는다. (ex. name 필드 -> **Name()** 메서드)
- **setter** 메서드는 보통 **Set필드명**으로 짓는다. (ex. name 필드 -> **SetName(n string)** 메서드) 