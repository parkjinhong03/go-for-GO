// 쓰로틀링이란 서비스가 처리할 수 있는 연결의 수를 제한하고 이 임계치를 초과하면 HTTP 에러 코드를 리턴하는 패턴이다.
// 이번 예제에서는 01장에서 사용했던 Go의 미들웨어 패턴이 매우 유용하게 사용된다.

package throttling

import (
	"net/http"
)

// 핸들러가 호출되기 전에 서버가 요청을 처리할 수 있는지 확인하기 위한 미들웨서 객체 선언
type LimitMiddleware struct {
	// 현재 서버에 들어와있는 요청의 갯수를 세기 위한 채널로, 요청이 들어오면 버퍼를 추가하고 응답이 완료되면 버퍼를 제거한다.
	connections chan struct{}
	next		http.Handler
}

// 받은 인자값으로 LimitMiddleware의 초깃값을 설정하여 반환하는 함수로, connections 매개변수는 제한할 연결의 갯수를 의미한다.
func NewLimitMiddleware(connections int, next http.Handler) *LimitMiddleware {
	// make 함수를 이용하여 버퍼를 가지고 있는 버퍼링 채널을 생성할 수 있다.
	cons := make(chan struct{}, connections)

	return &LimitMiddleware{
		connections: cons,
		next:        next,
	}
}

func (l *LimitMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// select case 구문으로 여려 채널들을 핸들링할 수 있다.
	select {
	// l.connections 채널에 새로운 구조체를 송신할 수 있다면, 즉 해당 채널의 버퍼에 남아있는 칸이 있다면 핸들러를 수행한다.
	case l.connections <- struct{}{}:
		l.next.ServeHTTP(rw, r)
		// 핸들러를 수행 한 후, 요청이 완료돠었다는 표시를 하기 위해 채널에서 버퍼 하나를 뺴온다.
		<-l.connections
	// l.connections 채널에 새로운 구조체를 송신할 수 없다면, 즉 해당 채널의 버퍼에 남아있는 칸이 없다면 429(TooManyRequest)를 반환한다.
	default:
		http.Error(rw, "busy", http.StatusTooManyRequests)
	}
}
