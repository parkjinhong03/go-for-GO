// 입력 유효성 검사는 공격자가 방화벽을 우회하며 비공개 API에 접근하였을 때 해당 공격자를 힘들게 만들 수 있는 방법 중 하나이다.
// 이번 예제에서 입력 유효성을 검사하기 위해 gopkg.in/go-playground/validator 패키지를 사용할 것 이다.

package validation

import (
	"encoding/json"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

type validationHandler struct {}

type Request struct {
	Name string `json:"name"`
	Email string `json:"email" validate:"email"`
	URL string `json:"url" validate:"url"`
}

var validate = validator.New()

func (v *validationHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	request := Request{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(rw, "Invalid request object", http.StatusBadRequest)
		return
	}

	err = validate.Struct(&request)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func NewValidationHandler() *validationHandler {
	return &validationHandler{}
}