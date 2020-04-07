# **Greeter**

service 폴더에서 연습했던 것을 바탕으로 조금 더 클 틀을 가진 마이크로서비스 예제입니다.

## **Contents**

- **srv** - RPC greeter 서비스
- **cli** - 특정 서비스를 호출하는 RPC 클라이언트
- **api** - RPC API와 RESTful API 예제
- **web** - go-web을 이용하여 웹 서비스를 만드는 방법

<br>

## **Run Service (srv)**

**go.micro.srv.greeter** 실행
```shell
go run srv/main.go
```

<br>

## **Run Client(cli)**

go.micro.srv.greeter를 호출하는 **클라이언트** 실행
```shell
go run cli/main.go
```

client 디렉토리에서 **다른 언어**를 이용한 **클라이언트** 사용의 **예시**를 볼 수 있습니다.

<br>

## **Run API (api)**

**micro API**를 통해 **HTTP 기반 요청**을 처리할 수 있습니다. **Micro**는 논리적으로 백엔드 서비스로부터 API 서비스들을 **분리**합니다. 기본적으로 micro API는 **HTTP 요청을 수락**한 후 **\*api.Request**와 **\*api.Response** 타입으로 **변환**합니다. 해당 타입은 [**여기에서**](https://github.com/micro/micro/tree/master/api/proto)(micro/api/proto)에서 찾아볼 수 있습니다.

**go.micro.api.greeter** 실행
```shell
go run api/api.go 
```

**micro API** 실행
```shell
micro api --handler=api
```

API를 통한 **go.micro.api.greeter** 호출
```shell
curl http://localhost:8080/greeter/say/hello?name=John
```

API 디렉토리에서 다른 API handlers의 예시를 찾아 볼 수 있습니다.


<br>

## **Run Web (web)**

**micro web**은 마이크로 서비스로 **웹 앱**을 실행시키기 위한 **웹 대시보드**와 **리버스 프록시**입니다.

**go.micro.web.greete**r 실행
```
go run web/web.go 
```

**micro web** 실행
```shell
micro web
```