# **Go Micro**
## **시작하기 전에...**
- 위대하신 [**Asim Aslam**](https://medium.com/@asimaslam)님께서 많은 기여를 하신 프레임워크로, [**Go Micro 공식 문서**](https://dev.micro.mu/)와 [**Codemotion의 강의 영상**](https://www.youtube.com/watch?v=OcjMi9cXItY), 그리고 Che Dan
님의 [**Micro In Action**](https://itnext.io/micro-in-action-part-5-message-broker-a3decf07f26a)에서
- 많은 **유용한 지식과 도움**을 받을 수 있었고, 덕분에 Go Micro를 프로젝트에 써볼 **용기**가 생겼습니다!

<br>

---
## **go-micro와 micro??**
- **go-micro**는 실제 **Golang 코드로 microservice를 개발**할 때 사용되는 프레임워크로, MSA에서 자주 쓰이는 공통 패턴들을 사용하기 쉽도록 **추상화 형태**로 제공해줍니다. 이 서비스들의 일반적인 유형은 **gRPC**입니다. 
- **micro**는 **command line tool**로, go-micro를 이용하여 개발한 서비스를 실행시키기 위해 필수로 사용해야 되는것은 아니지만, 여러 서비스들을 **개발하고 관리하는데 매우 편리함을 제공**해줍니다. 예를 들어, 템플릿 프로젝트 생성, 런타임 상태 검사, 서비스 목록 및 호출 등등의 기능들이 있습니다. 참고로 이 도구는 **go-micro 기반**입니다.

