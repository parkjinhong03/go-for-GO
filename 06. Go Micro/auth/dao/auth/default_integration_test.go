package user

import (
	"auth/model"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var ud *defaultDAO

func connectDB () *gorm.DB {
	dbUser := os.Getenv("DB_USER")
	dbPwd := os.Getenv("DB_PASSWORD")
	addr := fmt.Sprintf("%s:%s@/AuthTestDB?charset=utf8&parseTime=True&loc=Local", dbUser, dbPwd)
	db, err := gorm.Open("mysql", addr)
	if err != nil { log.Fatal(err) }
	return db
}

func init() {
	db := connectDB()

	db.DropTableIfExists(model.Auth{}, model.ProcessedMessage{})
	db.CreateTable(model.Auth{}, model.ProcessedMessage{})
	db.AutoMigrate(model.Auth{}, model.ProcessedMessage{})
	db.LogMode(false)
	_ = db.Close()
}

func setUpEnv() {
	db := connectDB()

	db.LogMode(false)
	ud = NewDefaultDAO(db)
	ud.db = ud.db.Begin()
}

type insertAuthTest struct {
	UserId      string
	UserPw      string
	Status      string
	ExpectError error
}

func (ia insertAuthTest) Exec() (*model.Auth, error) {
	return ud.InsertAuth(&model.Auth{
		UserId: ia.UserId,
		UserPw: ia.UserPw,
		Status: ia.Status,
	})
}

type insertMessageTest struct {
	MsgId       string
	ExpectError error
}

func (im insertMessageTest) Exec() (*model.ProcessedMessage, error) {
	return ud.InsertMessage(&model.ProcessedMessage{
		MsgId: im.MsgId,
	})
}

type checkIfUserIdExist struct {
	UserId      string
	expectExist bool
	expectError error
}

func (c checkIfUserIdExist) Exec() (bool, error) {
	return ud.CheckIfUserIdExist(c.UserId)
}

type updateStatusTest struct {
	id          uint
	status      string
	expectError error
}

func (us updateStatusTest) Exec() error {
	return ud.UpdateStatus(us.id, us.status)
}

func TestDefaultAuthDAOInsertAuth(t *testing.T) {
	setUpEnv()

	tests := []insertAuthTest{
		{
			UserId: "jinhong0719",
			UserPw: "testPW",
			Status: CreatePending,
			ExpectError: nil,
		}, {
			UserId: "jinhong0719",
			UserPw: "testPW",
			Status: CreatePending,
			ExpectError: UserIdDuplicatedError,
		}, {
			UserId: "jin0719",
			UserPw: "testPW",
			Status: "InvalidStatus",
			ExpectError: InvalidStatusError,
		}, {
			UserId: "jinhong0719jinhong0719",
			UserPw: "testPW",
			Status: CreatePending,
			ExpectError: UserIdTooLongError,
		},
	}

	for _, test := range tests {
		_, err := test.Exec()
		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)\n", test)
	}

	ud.Rollback()
	_ = ud.db.Close()
}

func TestDefaultAuthDAOUpdateStatus(t *testing.T) {
	setUpEnv()

	inits := []insertAuthTest{
		{
			UserId: "jinhong0719",
			UserPw: "TestPw",
			Status: CreatePending,
		},
	}

	var pk uint
	for _, init := range inits {
		r, err := init.Exec()
		if err != nil { log.Fatal(err) }
		pk = r.ID
	}

	tests := []updateStatusTest{
		{
			id:          pk,
			status:      Created,
			expectError: nil,
		}, {
			id:          pk,
			status:      Created,
			expectError: nil,
		}, {
			id:          pk+1,
			status:      Created,
			expectError: NonexistentUserError,
		}, {
			id:          pk,
			status:      "ThisIsInvalidStatus",
			expectError: InvalidStatusError,
		},
	}

	for _, test := range tests {
		err := test.Exec()
		assert.Equalf(t, test.expectError, err, "error assertion error (test case: %v)\n", test)
	}

	ud.Rollback()
	_ = ud.db.Close()
}