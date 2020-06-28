package subscriber

import (
	"auth/dao/user"
	"auth/model"
	proto "auth/proto/auth"
	"auth/tool/random"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"testing"
	"time"
)

type method string
type returns []interface{}

type createAuthTest struct {
	UserId        string
	UserPw        string
	Name          string
	PhoneNumber   string
	Email         string
	Introduction  string
	XRequestId    string
	ExpectMethods map[method]returns
}

func (c createAuthTest) createTestFromForm() (test createAuthTest) {
	test = c

	if c.UserId == none 		{ test.UserId = "" } 		else if c.UserId == "" 		  { test.UserId = defaultUserId }
	if c.UserPw == none 		{ test.UserPw = "" } 		else if c.UserPw == "" 		  { test.UserPw = defaultUserPw }
	if c.Name == none 		 	{ test.Name = "" } 			else if c.Name == "" 		  { test.Name = defaultName }
	if c.PhoneNumber == none  	{ test.PhoneNumber = "" }	else if c.PhoneNumber == ""   { test.PhoneNumber = defaultPN }
	if c.Introduction == none	{ test.Introduction = "" } 	else if c.Introduction == ""  { test.Introduction = "" }
	if c.Email == none 			{ test.Email = "" } 		else if c.Email == "" 		  { test.Email = defaultEmail }
	if c.XRequestId == none 	{ test.XRequestId = "" }	else if c.XRequestId == ""	  { test.XRequestId = uuid.New().String() }

	if _, ok := c.ExpectMethods["Insert"]; ok {
		test.setAuthContext(c.ExpectMethods["Insert"][0].(*model.Auth))
	}

	return
}

func (c createAuthTest) setAuthContext(auth *model.Auth) {
	auth.ID = authId
	auth.UserId = c.UserId
	auth.UserPw	= c.UserPw
	auth.Status = user.CreatePending
	auth.CreatedAt = time.Now()
	auth.UpdatedAt = time.Now()
	authId++
}

func (c createAuthTest) setRequestContext(req *proto.CreateAuthMessage) {
	req.UserId = c.UserId
	req.UserPw = c.UserPw
	req.Name = c.Name
	req.Email = c.Email
	req.PhoneNumber = c.PhoneNumber
	req.Introduction = c.Introduction
}

func (c createAuthTest) onExpectMethods() {
	for method, returns := range c.ExpectMethods {
		c.onMethod(method, returns)
	}
}

func (c createAuthTest) onMethod(method method, returns returns) {
	switch method {
	case "Insert":
		mockStore.On("Insert", &model.Auth{
			UserId: c.UserId,
			UserPw: c.UserPw,
			Status: user.CreatePending,
		}).Return(returns...)
	case "Commit":
		mockStore.On("Commit").Return(returns...)
	case "Rollback":
		mockStore.On("Rollback").Return(returns...)
	case "Ack":
		mockStore.On("Ack").Return(returns...)
	default:
		panic(fmt.Sprintf("%s method cannot be on booked\n", method))
	}
}

func TestCreateAuthValidMessage(t *testing.T) {
	setUp()
	msg := &proto.CreateAuthMessage{}
	var tests []createAuthTest

	forms := []createAuthTest{
		{
			ExpectMethods: map[method]returns{
				"Insert": {&model.Auth{}, nil},
				"Ack": {nil},
				"Commit": {&gorm.DB{}},
			},
		}, {
			ExpectMethods: map[method]returns{
				"Insert": {&model.Auth{}, errors.New("user id duplicated error")},
				"Rollback": {&gorm.DB{}},
			},
		}, {
			ExpectMethods: map[method]returns{
				"Insert": {&model.Auth{}, nil},
				"Ack": {errors.New("some error occurs while acknowledge message")},
				"Rollback": {&gorm.DB{}},
			},
		},
	}

	for _, form := range forms {
		tests = append(tests, form.createTestFromForm())
	}

	for _, test := range tests {
		test.setRequestContext(msg)
		test.onExpectMethods()

		header := make(map[string]string)
		header["XRequestId"] = test.XRequestId
		header["MessageId"] = random.GenerateString(32)
		header["Env"] = "Test"

		body, _ := json.Marshal(msg)
		_ = h.CreateAuth(NewCustomEvent(mockStore, header, body))
	}

	mockStore.AssertExpectations(t)
}
