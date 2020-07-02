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
