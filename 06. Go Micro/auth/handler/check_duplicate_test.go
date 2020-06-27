package handler

import (
	proto "auth/proto/auth"
	"github.com/stretchr/testify/mock"
	"log"
)

type UserIdExistTest struct {
	UserId        string
	Authorization string
	ExpectCode    int64
	ExpectMethods map[string][]interface{}
}

func (c UserIdExistTest) createTestFromForm() (test UserIdExistTest) {
	test.UserId = c.UserId
	test.Authorization = c.Authorization
	test.ExpectMethods = c.ExpectMethods
	test.ExpectCode = c.ExpectCode

	if c.UserId == None 		{ test.UserId = "" } 		else if c.UserId == "" 		  { test.UserId = DefaultUserId }
	if c.Authorization == None 	{ test.Authorization = "" } else if c.Authorization == "" { test.Authorization = DefaultUserPw }

	return
}

func (c UserIdExistTest) setRequestContext(req *proto.UserIdExistRequest) {
	req.UserId = c.UserId
	req.Authorization = c.Authorization
}

func (c UserIdExistTest) onExpectMethods() {
	for name, returns := range c.ExpectMethods {
		c.onMethod(name).Return(returns)
	}
}

func (c UserIdExistTest) onMethod(method string) (call *mock.Call) {
	switch method {
	case "CheckIfUserIdExists":
		call = mockStore.On("CheckIfUserIdExists", c.UserId)
	default:
		log.Fatalf("%s method cannot be on booked\n", method)
	}
	return
}
