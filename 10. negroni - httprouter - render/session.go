package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/goincremental/negroni-sessions"
)

const (
	currentUserKey = "oauth2_current_user"	// 세션에 저장되는 CurrentUser 키
	sessionDuration = time.Hour	// 로그인 세션 유지 시간
)

type User struct {
	Uid string `json:"uid"`
	Name string `json:"name"`
	Email string `json:"user"`
	AvatarUrl string `json:"avatar_url"`
	Exprired time.Time `json:"expired"`
}

func (u *User) Valid() bool {
	// 현재 시간 기준으로 만료 시간 확인
	return u.Exprired.Sub(time.Now()) > 0
}

func (u *User) Refresh() {
	// 만료 시간 연장
	u.Exprired = time.Now().Add(sessionDuration)
}

func GetCurrentUser(r *http.Request) *User {
	// 세선에서 CurrentUser 정보를 불러옴
	s := sessions.GetSession(r)

	if s.Get(currentUserKey) == nil {
		return nil
	}

	data := s.Get(currentUserKey).([]byte)
	var u User
	json.Unmarshal(data, &u)
	return &u
}

func SetCurrentUser(r *http.Request, u *User) {
	if u != nil {
		// CurrentUser 만료 시간 갱신
		u.Refresh()
	}

	s := sessions.GetSession(r)
	val, _ := json.Marshal(u)
	s.Set(currentUserKey, val)
}