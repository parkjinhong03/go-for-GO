## **시작하기 전에...**
- 해당 디렉터리 안에 있는 **코드**들 및 **README.md**에 대한 내용은 **다음 사이트**를 **참조**하였습니다.
- **https://github.com/micro/go-micro**
- **https://github.com/micro/example**
- **https://micro.mu/docs**

<br>


# **Go Micro** [![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/micro/go-micro?tab=doc) [![Travis CI](https://api.travis-ci.org/micro/go-micro.svg?branch=master)](https://travis-ci.org/micro/go-micro) [![Go Report Card](https://goreportcard.com/badge/micro/go-micro)](https://goreportcard.com/report/github.com/micro/go-micro)

[**Go Micro**](https://github.com/micro/go-micro)는 **microservices**를 개발하는데 용이한 **프레임워크**입니다.

## **개요**

**Go Micro**는 **RPC**와 **이벤트 기반 통신**을 포함한 **분산 시스템 개발**에 대한 핵심 요구 사항들을 제공합니다.   
**Go Micro**는 당신이 빠르게 시작할 수 있도록 **기본적인 것**을 제공하고, **pluggable architecture** 기반으로 구현된 [**go-plugin**](https://github.com/micro/go-plugins)을 이용해 모든것을 **손쉽게 교체**할 수 있습니다.  
또한 [**Twitter**](https://twitter.com/microhq)와 [**Slack**](https://micro.mu/slack)와 같은 여러 **커뮤니티**들이 활성화되어 있습니다.

<img src="https://micro.mu/docs/images/go-micro.svg" />


<br>

## **특징**

**Go Micro**는 분산 시스템의 세부 정보를 **추상화**하고, 주요 특징은 다음과 같습니다.

- **Service Discovery** - 서비스 실행 시 자동으로 **서비스가 등록**되고, 호출시 **이름을 확인**하는 기능을 가지고 있습니다. **Service discovery**는 microservice 개발에 있어서 아주 **중요한 기능**입니다. 기본 **검색 메커니즘**은 zero config 시스템인 **mDNS(Multicast DNS)** 입니다.

- **Load Balancing** - 서버의 과부하를 막기위해 **Client Side Load Balancing** 기능을 제공합니다. **Client Side Load Balancing**란 **클라이언트** 측에서 로직을 구현하여 중간에 **스위치**를 거치지 않고 **바로 서버로 요청**을 하는 **부하 분산**의 방식 중 하나입니다.

- **Message Encoding** - **Content-Type**에 기반한 **동적 메세지 인코딩**을 지원합니다. 다양한 메세지들은 **암호화** 되어 서로 다른 클라이언트에 보내지고, 클라이언트와 서버는 이 문제를 기본적으로 처리합니다. 메세지에는 구글의 **protobuf** 또는 **json** 등과 같은 유형이 있습니다.

- **Request/Response** - **양방향 스트리밍**을 지원하는 **RPC**를 기반 으로 응답/요청을 처리한다. 기본 통신 방식은 [**gRPC**](https://grpc.io/)이다.

- **Async Messaging** - **비동기 통신** 및 **이벤트 기반 아키텍쳐**를 위한 Pub/Sub 패턴이 내장되어있습니다. 기본 메시징 시스템은 내장되어있는 [**NATS**](https://nats.io/) 서버입니다.

- **Pluggable Interfaces** - 위에서 말했다 시피, Go Micro는 **Go의 인터페이스**를 이용하여 분산 시스템의 정보를 **추상화**하고 있습니다. 따라서 이러한 인터페이스들을 만족하는 여러 객체들을 [**여기**](https://github.com/micro/go-plugins)에서 찾아 **교체**할 수 있습니다.


<br>

## **시작하기**

- 이 [**문서**](https://micro.mu/docs/framework.html)에서 go-micro의 **아키텍쳐**, **설치** 그리고 **사용**에 대한 세부 정보를 확인할 수 있습니다.