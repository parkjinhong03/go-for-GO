package apiGatewayUsecase

import (
	"MSA.example.com/1/entities"
	"MSA.example.com/1/protocol"
	"MSA.example.com/1/tool/customError"
	natsEncoder "MSA.example.com/1/tool/encoder/nats"
	"MSA.example.com/1/tool/message"
	"MSA.example.com/1/usecase"
	"encoding/json"
	"github.com/eapache/go-resiliency/breaker"
	"github.com/eapache/go-resiliency/deadline"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"strings"
	"time"
)

type authServiceHandler struct {
	natsM    message.NatsMessage
	validate *validator.Validate
	natsE    natsEncoder.Encoder
	breakeR  *breaker.Breaker
}

func NewAuthServiceHandler(nastM message.NatsMessage, validator *validator.Validate,
	natsE natsEncoder.Encoder, breakeR *breaker.Breaker) *authServiceHandler {
	return &authServiceHandler{
		natsM:    nastM,
		validate: validator,
		natsE:    natsE,
		breakeR:  breakeR,
	}
}

func (h *authServiceHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	paths := strings.Split(r.URL.Path, "/")

	switch paths[0] {
	case "signup":
		toFunc := func(i <-chan struct{}) error {
			h.SignUpHandler(rw, r)
			return nil
		}
		brFunc := func() error {
			dl := deadline.New(time.Second * 3)
			return dl.Run(toFunc)
		}
		err := h.breakeR.Run(brFunc)
		switch err {
		case breaker.ErrBreakerOpen:
			rw.WriteHeader(http.StatusServiceUnavailable)
			return
		case deadline.ErrTimedOut:
			rw.WriteHeader(http.StatusRequestTimeout)
			return
		default:
			return
		}
	default:
		rw.WriteHeader(http.StatusNotFound)
	}
}

func (h *authServiceHandler) SignUpHandler(rw http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	request := entities.SignUpRequestEntities{}
	err := d.Decode(&request)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.validate.Struct(&request)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// 수신자 경합과 메시지 순서 유지 (인스턴스 복제 시)
	// 중복 메시지 처리 -> 메시지 추적과 중복 메시지 솎아
	// 트랜잭셔널 메시징 기능 추가 -> DB 테이블을 메시지 큐로 활용
	// 사가 트랜잭션에 대한 비격리 대책 -> 시맨틱 락 구현 -> 완료
	// user_id 중복 체크 -> req/resp | row 생성 -> pub/sub
	// 비동기 처리 전환 및 푸시 기능 구햔
	// 로깅 구현 (ElasticSearch, Logstash, Kibana)
	// 비번 암호화 util 추가
	enErr := h.natsE.Encode(protocol.AuthSignUpRequestProtocol{
		Required: protocol.RequiredProtocol{
			Usage:        "AuthSignUpRequest",
			InputChannel: "auth.signup",
		},
		RequestId:    r.Header.Get("X-Request-ID"),
		UserId:       request.UserId,
		UserPwd:      request.UserPwd,
		Name:         request.Name,
		PhoneNumber:  request.PhoneNumber,
		Introduction: request.Introduction,
		Email:        request.Email,
	}).(*customError.ProxyWriteError)

	if enErr.Err != nil {
		switch enErr.Err.(type) {
		case validator.ValidationErrors:
			rw.WriteHeader(http.StatusBadRequest)
		default:
			log.Printf("some error occurs while encoding to message, err: %v\n", enErr.Err)
			// 타임아웃 발생 시, 아래 구문에서 경고 로그가 뜸, 리펙토링 필요
			rw.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	p := protocol.ApiGatewaySignUpResponseProtocol{}
	if err := json.Unmarshal(enErr.ReturnMsg.Data, &p); err != nil {
		log.Printf("some error occurs while decoding message into struct, err: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := h.validate.Struct(&p); err != nil {
		log.Printf("the message recived from auth.signup is invalid, err: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if p.Success {
		rw.WriteHeader(http.StatusCreated)
		reply := entities.SignUpResponseEntities{
			StatusCode: http.StatusOK,
			ResultUser: p.ResultUser,
		}
		encoder := json.NewEncoder(rw)

		if encoder.Encode(reply) != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	// 에러 코드 상수화 추가
	switch p.ErrorCode {
	case usecase.UserIdDuplicateErrorCode:
		rw.WriteHeader(470)
		return
	case usecase.ParsingFailureErrorCode:
		log.Println("auth.signup fail to parse error code")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	default:
		log.Println("undefined error code's come from auth.signup")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}