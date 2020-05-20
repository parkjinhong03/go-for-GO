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

type AuthServiceHandler struct {
	natsM message.NatsMessage
	validate *validator.Validate
	natsE natsEncoder.Encoder
}

func (h *AuthServiceHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
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
		err = h.natsE.Encode(protocol.AuthSignUpProtocol{
			RequestId:     r.Header.Get("X-Request-ID"),
			UserId:        request.UserId,
			UserPwd:       request.UserPwd,
			Name:          request.Name,
			PhoneNumber:   request.PhoneNumber,
			Introduction:  request.Introduction,
			Email:         request.Email,
			ReturnChannel: "auth.signup.return",
			InputChannel:  "auth.signup",
		})

		enErr := err.(*customError.ProxyWriteError)
		if enErr.Err != nil {
			log.Printf("some error occurs while encoding to message, err: %v\n", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusOK)
		return

	default:
		rw.WriteHeader(http.StatusNotFound)
		return

	}
}