# **sync.RWMutex**
## **sync.RWMutex란?**
- sync.RWMutex는 sync.Mutex와 동작 방식이 유사하다.
- 하지만 sync.RWMutex는 읽기 동작과 쓰기 동작을 나누어 잠금 처리할 수 있다.

<br>

---
## **sync.RWMutex의 기능**
- **읽기 잠금**
    - 읽기 잠금은 **읽기 작업** 한해서 공유 데이터가 **변하지 않음**을 보장해준다.
    - 따라서 **다른 고루틴**에서 데이터를 읽은 수는 있지만, **변경할 수는 없다.**

- **쓰기 잠금**
    - 쓰기 잠금은 **공유 데이터**에 쓰기 작업을 **보장**하기 위한 것이다.
    - 따라서 **다른 고루틴**에서는 읽기와 쓰기 작업 **모두 할 수 없다.**

<br>

---
## **sync.RWMutex의 메서드**
- **func (rw \*RWMutex) Lock()** -> 쓰기 잠금
- **func (rw \*RWMutex) Unlock()** -> 쓰기 잠금 해제
- **func (rw \*RWMutex) Rlock()** -> 읽기 잠금
- **func (rw \*RWMutex) Runlock()** -> 읽기 잠금 해제