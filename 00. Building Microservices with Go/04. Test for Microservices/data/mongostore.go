package data

import "gopkg.in/mgo.v2"

// 모의 객체가 아닌, 실제 데이터 저장소인 MongoDB를 DataStore 인터페이스로 사용하기 위한 구조체 정의
type MongoStore struct {
	session *mgo.Session
}