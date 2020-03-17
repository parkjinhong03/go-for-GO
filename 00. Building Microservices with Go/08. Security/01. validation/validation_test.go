package validation

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getMockObjects(v interface{}) (*validationHandler, *httptest.ResponseRecorder, *http.Request) {
	h := NewValidationHandler()
	rw := httptest.NewRecorder()
	if v == nil {
		return h, rw, httptest.NewRequest("GET", "/examples/validator", nil)
	}
	b, _ := json.Marshal(v)
	return h, rw, httptest.NewRequest("GET", "/examples/validator", bytes.NewReader(b))
}

func TestReturnsBadRequestWhenEmailIsNotSent(t *testing.T) {
	request := Request{
		URL:   "https://github.com/parkjinhong03",
	}
	h, rw, r := getMockObjects(request)

	h.ServeHTTP(rw, r)

	assert.Equal(t, rw.Code, http.StatusBadRequest, "Should have raised an error")
}

func TestReturnsBadRequestWhenEmailIsInvalid(t *testing.T) {
	request := Request{
		Email: "invalidEmail.com",
		URL:   "https://github.com/parkjinhong03",
	}
	h, rw, r := getMockObjects(request)

	h.ServeHTTP(rw, r)

	assert.Equal(t, rw.Code, http.StatusBadRequest, "Should have raised an error")
}

func TestReturnsOKWhenNameIsNotSent(t *testing.T) {
	request := Request{
		Email: "jinhong0719@naver.com",
		URL:   "https://github.com/parkjinhong03",
	}
	h, rw, r := getMockObjects(request)

	h.ServeHTTP(rw, r)

	assert.Equal(t, rw.Code, http.StatusOK, "Shouldn't have raised an error")
}