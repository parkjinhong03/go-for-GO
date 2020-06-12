# **Go Micro**
## **시작하기 전에...**
- [**Asim Aslam**](https://medium.com/@asimaslam)님 외 네 분으로 구성되어있는 [**micro.mu**](https://micro.mu/)에서 만든 프레임워크로, 
- [**공식 문서**](https://dev.micro.mu/)와 [**Codemotion**](https://www.youtube.com/watch?v=OcjMi9cXItY), 그리고 [**Micro In Action**](https://itnext.io/micro-in-action-getting-started-a79916ae3cac)에서 많은 많은 도움을 받았습니다.

<br>

---
## **목차**
### [**1. go-micro와 micro??**](#go-micro와-micro??)
### [**2. go-micro의 아키텍처**](#go-micro의-아키텍처)
### [**3. 서비스 분해**](#서비스-분해)
### [**4. 적용할 패턴**](#적용할-패턴)
### [**5. 앞으로 배울 패턴**](#앞으로-배울-패턴)

<br>

---
## **go-micro와 micro??**
- [**go-micro**](https://github.com/micro/go-micro)는 실제 **Golang 코드로 microservice를 개발**할 때 사용되는 프레임워크로, MSA에서 자주 쓰이는 공통 패턴들을 사용하기 쉽도록 **추상화 형태**로 제공해줍니다. 이 서비스들의 일반적인 유형은 **gRPC**입니다. 
- [**micro**](https://github.com/micro/micro)는 **command line tool**로, go-micro를 이용하여 개발한 서비스를 실행시키기 위해 필수로 사용해야 되는것은 아니지만, 여러 서비스들을 **개발하고 관리하는데 매우 편리함을 제공**해줍니다. 예를 들어, 템플릿 프로젝트 생성, 런타임 상태 검사, 서비스 목록 및 호출 등등의 기능들이 있습니다. 참고로 이 도구는 **go-micro 기반**입니다.
- 참고로 [**go-plugin**](https://github.com/micro/go-plugins)이라는 것도 있는데, 이는 앞에서도 말씀드렸듯이, go-micro가 서비스들을 **추상화 형태**로 제공하기 때문에 가능한 **일련의 플러그인**입니다. Service Discovery, Async Messaging, Transport Protocol 등등의 기능들을 **사용자가 원하는 형태로 사용**할 수 있게 제공해줍니다.


<br>

---
## **go-micro의 아키텍처**
**go-micro의 Architecture**
![Architecture](architecture.png)

**go-micro**는 **마이크로 서비스 개발 및 분산 시스템 환경 구축**을 단순화하는 것에 초점을 두었습니다. 따라서 그 철학에 맞게 분산 시스템에서 항상 필요한 작업들을 모두 제공해줍니다. 하지만 여기서 더 중요한 것은 이러한 공통 패턴들을 모두 **인터페이스로 추상화**하였다는 것 입니다. 따라서 사용자가 세부 구현 정보를 알 필요 없이, **유연하고 강력한 시스템**을 매우 빠르게 구축할 수 있습니다. 중요한 몇몇 인터페이스들은 다음과 같습니다.

- ### **Service Discovery**
    **Service Discovery**는 동적으로 url이 바뀌는 서비스에 **동기 IPC**로 호출을 하기 위해 필수로 구현해야하는 기능입니다. 이 기능을 **추상화한 인터페이스**는 [**go-micro**](https://github.com/micro/go-micro/blob/master/registry/registry.go#L20)에서 보실 수 있습니다. 이 인터페이스를 구현한 모든 객체들은 서비스에서 **Service Discovery 기능의 수행 주체**로 사용될 수 있습니다. [**go-plugins**](https://github.com/micro/go-plugins/tree/master/registry)에서 etcd, consul, zookeeper등과 같은 많은 구현 객체들을 확인하실 수 있습니다. 참고로 설정을 하지 않을 경우, **multicase DNS(mdns)** 기반인 Service Discovery가 기본으로 사용됩니다.

- ### **Async Messaging**
    **Async Messaging**은 서비스들끼리의 **느슨한 결합과 강력한 시스템**을 구축하기 위한 핵심 기술입니다. 이 기능을 **추상화한 인터페이스** 또한 [**go-micro**](https://github.com/micro/go-micro/blob/master/broker/broker.go#L5)에서 보실 수 있습니다. 이 인터페이스 구현 객체도 [**go-plugins**](https://github.com/micro/go-plugins/tree/master/broker)에서 확인하실 수 있으며, 대표적으로 **RabbitMQ, Kafka, NSQ, NATS** 등등이 있습니다. 참고로 **HTTP 기반의 구현 객체**가 기본으로 사용되고, 추가 설정이 필요하지 않습니다.

- ### **Codec**
    **Codec**은 마이크로 서비스간의 통신을 위해서 메시지를 **인코딩 및 디코딩** 해주는 기술입니다. 이 기능을 **추상화한 인터페이스**는 [**go-micro**](https://github.com/micro/go-micro/blob/master/codec/codec.go#L30)에서 보실 수 있고, 구현 객체들은 [**go-micro**](https://github.com/micro/go-micro/tree/master/codec)와 [**go-plugins**](https://github.com/micro/go-plugins/tree/master/codec)에서 보실 수 있습니다. 현재 json, bson, msgpack등 많은 Codec들을 지원하고 있습니다.

- ### **others...**
    - **Server** - 마이크로서비스의 **서비스들을 정의**하기 위해 쓰입니다.
    - **Transport** - **전송 프로토콜**을 정의하기 위해 쓰입니다.
    - **Selector** - 로드 밸런싱 전략 구현을 위해 **서비스 선택 논리**를 추상화한 것 입니다.
    - **Wrapper** - 서버와 클라이언트의 요청을 래핑하는 **미들웨어**를 정의하기 위해 쓰입니다.

<br>

---
## **서비스 분해**
- **Api Gateway**(examples.blog.api.gateway) - 사용자의 요청 받아 **적절한 서비스로 분산**시키는 api
- **Auth Service**(examples.blog.service.auth) - 사용자의 **계정 정보**를 관리하는 서비스
- **User Service**(examples.blog.service.user) - 사용자의 **회원 정보**를 관리하는 서비스
- **Blog Service**(examples.blog.service.blog) - 사용자의 **블로그 정보**를 관리하는 서비스
- **Post Service**(examples.blog.service.post) - 블로그의 **게시글 정보**를 관리하는 서비스
- **Visitor Service**(examples.blog.service.visitor) - 블로그의 **방문자 정보**를 관리하는 서비스
- **Subscribe Service**(examples.blog.service.subscribe) - 블로그의 **구독 관련 정보**를 관리하는 서비스

<br>

---
## **적용할 패턴**
> 앞으로 나오는 **여러 패턴**들과 **결정 사항**들은 [**마이크로서비스 패턴**](https://www.aladin.co.kr/shop/wproduct.aspx?ItemId=228694618)책을 통해 알 수 있었으며, 이 책을 한번 **읽어보시길 추천**합니다!
- ### **통신 패턴**
    - **통신 스타일**
        - **비동기 메시징 -> 사용 O** *(메시지 순서 유지, 중복 메시지 처리 기능 구현 예정)*
        - **동기 IPC -> 사용 O** *(protocol buffer를 이용한 gRPC 사용 예정)*
    - **디스커버리**
        - **서버 사이드 디스커버리 -> 사용 O**
        - **클라이언트 사이드 디스커버리 -> 사용 X** *(why? 내가 직접 구현해보고 싶어서..ㅎㅎ)*
        - **서드파티 등록 -> 사용 O**
        - **자가 등록 -> 사용 X** *(why? Go Micro의 registra를 사용할 예정!)*
    - **신뢰성**
        - **회로 차단기 -> 사용 O** *(타임 아웃과 함께 사용할 예정)*
    - **트랜잭셔널 메시징**
        - **트랜잭션 로그 테일링 -> 사용 O** *(디비지움, Debezium을 이용할 예정!)*
        - **폴링 발행기 패턴 -> 사용 X** *(why? 자주 폴링할 경우 비용이 유발되기 때문)*

- ### **데이터 일관성 패턴**
    - **사가 트랜잭션 -> 사용 O** *(추후에 이벤트 소싱, 애그리거트 학습 후 사용해볼 예정)*
        - **코레오그래피 사가 -> 사용 O**
        - **오케스트레이션 사가 -> 사용 X** *(why? 코레오그래피 사가의 불편한 점을 아직 못느낌)*
        - **사가 트랜잭션은 ACD!!** *(시맨틱 락으로 비격리 문제 해결 예정)*
    - **분산 트랜잭션 -> 사용 X** *(why? 기술 선택의 제약성, 동기 통신으로 인한 가용성 저하)*

- ### **관측성 패턴**
    - **헬스 체크 API -> 사용 O** *(이 기능도 이번엔 꼭 구현해보고 싶은데, 어떻게 해야 할지 감이 안옵니다..ㅠㅠ)*
    - **로그 수집 -> 사용 O** *([ElasticSearch](https://www.elastic.co/kr/elasticsearch/service?elektra=home&storm=sub1), [Logstash](https://www.elastic.co/kr/logstash), [Kibana](https://www.elastic.co/kr/kibana)로 구축할 예정)*
    - **분산 추적 -> 사용 O** *(Golang 오픈소스 프로젝트인 [Jaeger](https://www.jaegertracing.io/)로 구축할 예정)*
    - **예외 추적 -> 사용 X** (why? 뭔지도 잘 모르겠고 굳이 사용해야 될 필요성를 느끼지 못해서..?)
    
<br>

---
## **앞으로 배울 패턴**
- ### **데이터 쿼리 패턴**
    - **API 조합**
    - **CQRS**

- ### **서비스 배포 패턴**
    - **호스트별 다중 서비스** *(WAR 파일 등 언어에 특정하게 패키징하여 배포했던 방식)*
    - **호스트별 단일 서비스** *(서비스 기술 스택을 캡슐화한 요즘 방식)*
        - **서버리스 배포**
        - **컨테이너별 서비스**
        - **VM별 서비스**

- ### **서비스 테스트 자동화 패턴**
    - **컨슈머 주도 계약 테스트** *(?)*
    - **컨슈머 쪽 계약 테스트** *(??)*
    - **서비스 컴포넌트 테스트** *(???)*
