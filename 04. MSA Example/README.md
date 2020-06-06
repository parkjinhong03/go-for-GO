# **MSA Example**

## **시작하기 전에...**
- **크리스 리처드슨**의 **[마이크로서비스 패턴](https://www.aladin.co.kr/shop/wproduct.aspx?ItemId=228694618)** 책을 읽고 많은 유용한 지식을 얻을 수 있었으며
- **jfeng45**님의 **[github 레포지토리](https://github.com/jfeng45/servicetmpl)** 와 **[medium 블로그](https://medium.com/@jfeng45/go-micro-service-with-clean-architecture-application-layout-e6216dbb835a)** 에서 많은 도움을 받았습니다.

<br>

---
## **프로젝트 계기**
**Building Microservices with Go** 책을 읽고 3강을 정리한 후에 **go언어에 대한 익숙함**과 **MSA에 대한 기본적인 이론**은 익혔다고 생각하였습니다.   
하지만 막상 MSA 기반의 개인 프로젝트를 진행하려고 하니 **전자인 익숙함은 익혔을지라도 후자인 이론은 너무 허실하고 부족**하여 이 프로젝트를 어떻게 시작해야할지 감이 전혀 안오고 **너무 막막했습니다.**   
따라서 제가 결정한게 위에서 말씀드린 **마이크로서비스 패턴**이라는 책을 사서, 실제 MSA 서비스를 위한 **탄탄한 기본 지식뿐만 아니라, 프로젝트에 도입할 수 있는 용기**를 쌓는 것이였습니다.  
그렇게 해서 만들게된 것이 **4강인 MSA Example**으로, 생각했던 시간이 코드를 짰던 시간보다 **10배는 많을 정도**로 대충 짠것이 절대 아닌, 온갖 스트레스와 프로젝트 구조에 대한 막막하고 두려운 감정을 **오직 해내고 싶다는 마음 하나로 이겨내면서 만든 것**임을 꼭 알아주셨으면 합니다.  
그러면 이제부터 제 프로젝트 **(매우 주관적인)** 에 도입된 **프로젝트 구조와 공통 패턴들**을 알려드리겠습나다. 

<br>

---
## 

<br>

---
## **프로젝트 구조**
<img src="./Project.png" width="270" height="300">

-  ### **cmd**
    - 모든 마이크로서비스 **프로그램들의 시작점**으로, **main.go, server, docker-compose.yml, DockerFile**로 구성되어있습니다.
    - **main.go** - 엔드포인트를 만드는 **main 함수**로, 각각의 계층에 **의존성을 주입**하여 **서비스를 실행**시킵니다.
    - **server** - main.go 파일을 **배포 환경**인 리눅스 환경에서 실행시키기 위해 **크로스 컴파일**하여 만든 바이너리 **실행 파일**입니다.
    - **DockerFile** - 배포환경에서 서버를 실행할 Docker Container를 만들기 위해 **Docker image를 생성하는 DockerFile**입니다.
    - **docker-compose.yml** - DockerFile로 만든 이미지를 이용하여 **쉽게 컨테이너를 만들수 있게 해주는** yml 파일입니다.

- ### **dataservice**
    - 오직 **model layer에만 의존**하여 도메인 계층의 데이터를 검색하고 수정하는 **persistance layer**입니다.
    - **db_service_interface.go** - go언어로 구현된 마이크로서비스에서 제공하는 모든 **data service를 위한 인터페이스**가 있습니다. 다른 패키지들은 해당 **인터페이스에만 의존**해야 하며, DB CRUD의 세부 구현 사항에 대해 알 필요가 없습니다.
    - **테이블 이름_data_service.go** - 위에서 말씀드린 인터페이스들을 구현하는 객체들로, 의존성으론 DB 연결에 대한 객체 하나만 포함하고 있습니다. 참고로 MSA 특성 상, 특정 테이블에 접근할 수 있는 서비스는 정해져 있으므로 data service도 마찬가지로 알맞은 서비스에서 사용해야 합니다.

- ### **entities**
    - **entities**는 사용자가 호출하는 서버-클라이언트 사이의 http 통신을 위해 **json과 객체 사이에 인코딩 및 디코딩을 도와주는 객체**들을 모와놓은 폴더입니다. 
    - 참고로 **유효성 타입 검사**를 진행하기 위해서 태그를 이용합니다.

- ### **middleware**
    - **middleware**는 사용자의 **요청을 처리하기 전에** 공통으로 **수행하는 작업들**을 따로 처리하기 위해 모와둔 폴더입니다.
    - 예를 들면, 로깅을 쉽게 하기 위해 api gateway에서 **각 요청마다 X-Request-ID 값을 설정**해줘야 하는데, 이 작업은 실제 **비즈니스 로직을 실행하기 전에 수행**되어야 하므로 이러한 작업들을 **middleware로 분리**합니다.

- ### **model**
    - **model**은 MVC(Model-View-Contorller)의 **Model과 비슷한 개념**으로, 도메인 구조체를 포함하고 있는 **domain model layer**입니다.
    - DB 관련 **도메인 모델을 정의**하는데 **[gorm 패키지](https://github.com/jinzhu/gorm)** 의 도움을 받았으며 **persistance layer** 또한 마찬가지 입니다.

- ### **protocol**
    - **protocol**또한 **json과 객체 사이에 인코딩 및 디코딩을 도와주는 객체**들을 모와놓은 폴더입니다.
    - 하지만 entities와는 다르게 protocol은 서비스끼리 **Message Queue**를 이용하여 **이벤트를 발행할 때의 통신 규약**을 정하기 위한 것 입니다.
    - protocol 또한 **유효성 타입 검사**를 진행하기 위해 태그를 사용하였습니다.

- ### **tool**
    - **tool**은 [**golang standard project layout**](https://github.com/golang-standards/project-layout)의 /tools와 같이 서비스를 개발하면서 **자주 쓰일것 같은 기능들**을 묶어놓은 디렉터리로, 저의 프로젝트에는 **다음과 같은 세부 디렉토리**가 있습니다.
    - **proxy** - **proxy**는 구조체가 마샬링된 바이너리 값을 받아서 용도에 따라 해당하는 **채널에 메시지를 발행해주는 서비스 통신의 중간 매체**를 하는 객체입니다. Write 메서드가 존재하여 **io.Writer 인터페이스**를 구현하고 있습니다.
    - **encoder** - **encoder**는 필드로 저장되어있는 proxy의 **Write 메서드를 호출하는 객체**로, **json.Encoder 객체를 임베딩**하고 있습니다. 그리고 encoder의 종류(jsonEncoder, xmlEncoder)에 따라 구조체를 마샬링하는 방법이 다릅니다.
    - **message** - **message**는 다양한 방법으로 **NATS서버에 연결하여 연결 객체**를 반환하는 함수들을 가지고 있고, 반환 객체는 **nats.Conn 객체를 임베딩**하여 일부 **메서드를 오버라이딩**하여 원하는 기능을 추가하였습니다.
    - **dbc** - **dbc**는 **데이터베이스 서버와 연결하여 연결 객체를 반환**하는 함수들을 포함하고 있는 디렉터리 입니다. [**gorm**](https://github.com/jinzhu/gorm) **패키지** 를 이용해 연결을 하였고, DB 연결 해제를 위해 **시그널 핸들 기능**을 추가하였습니다.

- ### **usecase**
    - 마지막으로 usecase는 **business layer**로, 여러 상황에 따른 로직을 처리하는 usecase 객체를 포함하고 있습니다. 
    - **nats 인코더, data service 객체(DAO), 유효성 검사 객체**에 대해 **의존성**을 가지고 있습니다.
    - 참고로 http 요청을 다루는 Api Gateway의 usecase를 제외하면 다른 모든 usecase들은 메시지를 handling합니다.