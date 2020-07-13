package user

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"user/tool/random"
)

func TestDefaultUserDAOInsertUser(t *testing.T) {
	setUpEnv()

	tests := []insertUserTest{
		{
			AuthId:       1,
			Name:         "박진홍",
			PhoneNumber:  "01088378347",
			Email:        "jinhong0719@naver.com",
			Introduction: "안녕하세요 저의 이름은 박진홍입니다.",
			ExpectError:  nil,
		}, {
			AuthId:       1,
			Name:         "박진홍",
			PhoneNumber:  "01088378347",
			Email:        "jinhong0719@hanmail.net",
			Introduction: "안녕하세요 저의 이름은 박진홍입니다.",
			ExpectError:  AuthIdDuplicatedError,
		}, {
			AuthId:       2,
			Name:         "박진홍",
			PhoneNumber:  "01088378347",
			Email:        "jinhong0719@naver.com",
			Introduction: "안녕하세요 저의 이름은 박진홍입니다.",
			ExpectError:  EmailDuplicatedError,
		}, {
			AuthId:       3,
			Name:         "박진홍박진홍",
			PhoneNumber:  "01088378347",
			Email:        "jinhong0719@hanmail.net",
			Introduction: "안녕하세요 저의 이름은 박진홍입니다.",
			ExpectError:  NameTooLongError,
		}, {
			AuthId:       4,
			Name:         "박진홍",
			PhoneNumber:  "01088378347123",
			Email:        "jinhong0719@hanmail.net",
			Introduction: "안녕하세요 저의 이름은 박진홍입니다.",
			ExpectError:  PhoneNumberTooLongError,
		}, {
			AuthId:       5,
			Name:         "박진홍",
			PhoneNumber:  "01088378347",
			Email:        "jinhong0719jinhong0719@hanmail.net",
			Introduction: "안녕하세요 저의 이름은 박진홍입니다.",
			ExpectError:  EmailTooLongError,
		},
	}

	for _, test := range tests {
		_, err := test.Exec()
		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v\n)", test)
	}

	ud.Rollback()
	_ = ud.db.Close()
}

func TestDefaultUserDAOCheckIfEmailExist(t *testing.T) {
	setUpEnv()

	inits := []insertUserTest{
		{
			AuthId:       1,
			Name:         "박진홍",
			PhoneNumber:  "01088378347",
			Email:        "jinhong0719@naver.com",
			Introduction: "안녕하세요 저의 이름은 박진홍입니다.",
			ExpectError:  nil,
		},
	}

	for _, init := range inits {
		_, err := init.Exec()
		if err != nil { log.Fatal(err) }
	}

	tests := []checkIfEmailExistTest{
		{
			Email: "jinhong0719@naver.com",
			ExpectExist: true,
			ExpectError: nil,
		}, {
			Email: "jinhong0719@gmail.com",
			ExpectExist: false,
			ExpectError: nil,
		},
	}

	for _, test := range tests {
		exist, err := test.Exec()
		assert.Equalf(t, test.ExpectExist, exist, "exist assertion error (test case: %v)\n", test)
		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)\n", test)
	}

	ud.Rollback()
	_ = ud.db.Close()
}

func TestDefaultUserDAOInsertMessage(t *testing.T) {
	setUpEnv()
	msgId := random.GenerateString(32)

	tests := []insertMessageTest{
		{
			MsgId:       msgId,
			ExpectError: nil,
		}, {
			MsgId:       msgId,
			ExpectError: MessageDuplicatedError,
		}, {
			MsgId:       random.GenerateString(33),
			ExpectError: MessageTooLongError,
		},
	}

	for _, test := range tests {
		_, err := test.Exec()
		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)\n", test)
	}
}