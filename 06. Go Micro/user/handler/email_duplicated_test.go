package handler

import (
	"github.com/google/uuid"
	"log"
	proto "user/proto/user"
)

type emailDuplicatedTest struct {
	Email               string
	Authorization       string
	XRequestId          string
	ExpectCode          uint32
	ExpectMessage       string
	ExpectMethods       map[method]returns
	ExpectAuthorization string
}

func (e emailDuplicatedTest) createTestFromForm() (test emailDuplicatedTest) {
	test = e

	if e.Email == none 			{ test.Email = "" } 	 	else if e.Email == "" 		  { test.Email = defaultEmail }
	if e.Authorization == none  { test.Authorization = "" } else if e.Authorization == "" { test.Authorization = "" }
	if e.XRequestId == none 	{ test.XRequestId = "" } 	else if e.XRequestId == "" 	  { test.XRequestId = uuid.New().String() }

	return
}

func (e emailDuplicatedTest) setRequestContext(req *proto.EmailDuplicatedRequest) {
	req.Email = e.Email
	req.Authorization = e.Authorization
	req.XRequestId = e.XRequestId
	return
}

func (e emailDuplicatedTest) onExpectMethods() {
	for method, returns := range e.ExpectMethods {
		e.onMethod(method, returns)
	}
	return
}

func (e emailDuplicatedTest) onMethod(method method, returns returns) {
	switch method {
	case "CheckIfEmailExist":
		mockStore.On("CheckIfEmailExist", e.Email).Return(returns...)
	default:
		log.Fatalf("%s method cannot be on booked\n", method)
	}
	return
}