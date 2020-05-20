package usecase

import (
	"MSA.example.com/1/entities"
	"MSA.example.com/1/protocol"
	"MSA.example.com/1/tool/customError"
	natsEncoder "MSA.example.com/1/tool/encoder/nats"
	"MSA.example.com/1/tool/message"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"strings"
)

type authServiceHandler struct {
	natsM 		message.NatsMessage
	validate 	*validator.Validate
	natsE 		natsEncoder.Encoder
}

func NewAuthServiceHandler(nastM message.NatsMessage, validator *validator.Validate, natsE natsEncoder.Encoder) *authServiceHandler {
	return &authServiceHandler{
		natsM:    nastM,
		validate: validator,
		natsE:    natsE,
	}
}

func (h *authServiceHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	paths := strings.Split(r.URL.Path, "/")

	d := json.NewDecoder(r.Body)
	switch paths[0] {
	case "signup", "signup/":
		request := entities.AuthSignUpEntities{}
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

		// 타임아웃, 회로 차단기 구현
		// auth_microservice 응답 코드 추가
		// usecase 트랜직셔널, 트랜잭셔널 메시징 기능 추가
		err = h.natsE.Encode(protocol.AuthSignUpRequestProtocol{
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
		})

		enErr := err.(*customError.ProxyWriteError)
		if enErr.Err != nil {
			log.Printf("some error occurs while encoding to message, err: %v\n", err)
			switch enErr.Err.(type) {
			case validator.ValidationErrors:
				rw.WriteHeader(http.StatusBadRequest)
			default:
				rw.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		// ApiGateSignUpResponseProtocol Unmarshal 및 처리 로직 추가
		rw.WriteHeader(http.StatusOK)
		return

	default:
		rw.WriteHeader(http.StatusNotFound)
		return

	}
}