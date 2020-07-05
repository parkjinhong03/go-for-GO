package subscriber

import (
	"fmt"
	"github.com/google/uuid"
	"time"
	"user/model"
	proto "user/proto/user"
	"user/tool/random"
)

type createUserTest struct {
	AuthId         uint32
	Name           string
	PhoneNumber    string
	Email          string
	Introduction   string
	XRequestID     string
	MessageID      string
	AfterMessageID string
	ExpectMethods  map[method]returns
	ExpectError    error
}

